package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Configuration variables (will be set via flags)
var (
	maxMessages    int    // Maximum number of messages to store
	maxEvents      int    // Maximum number of events to store
	httpServerAddr string // HTTP server address
	discordToken   string // Discord bot token
	discordChannel string // Discord channel ID used globally
	discordGuild   string // Discord guild ID used for events
	enableEvents   bool   // Enable Discord events support
	enableChannel  bool   // Enable Discord channel messages support
	eventRefresh   int    // Event refresh interval in seconds
)

// Message represents a simplified Discord message
type Message struct {
	Index   int       // Position in the message list (0 is newest)
	UserID  string    // Discord user ID
	User    string    // Discord username
	Content string    // Message content
	Time    time.Time // Message timestamp
}

// Event represents a Discord scheduled event
type Event struct {
	Index       int       // Position in the event list (0 is newest)
	ID          string    // Discord event ID
	Name        string    // Event name
	Description string    // Event description
	Location    string    // Event location
	StartTime   time.Time // Event start time (in UTC)
	EndTime     time.Time // Event end time (in UTC)
	GuildID     string    // Guild ID the event belongs to
	ImageURL    string    // URL to the event cover image
}

// Global variables
var (
	messages      []Message
	events        []Event
	messagesMutex sync.RWMutex
	eventsMutex   sync.RWMutex
	discord       *discordgo.Session
	channelID     string
)

func main() {
	// Define command line flags
	flag.IntVar(&maxMessages, "max-messages", 5, "Maximum number of messages to store")
	flag.IntVar(&maxEvents, "max-events", 10, "Maximum number of events to store")
	flag.StringVar(&httpServerAddr, "server-addr", "localhost:8080", "HTTP server address")
	flag.StringVar(&discordToken, "token", "", "Discord bot token")
	flag.StringVar(&discordChannel, "channel", "", "Discord channel ID")
	flag.StringVar(&discordGuild, "guild", "", "Discord guild/server ID (required if events enabled without channel)")
	flag.BoolVar(&enableEvents, "events", false, "Enable Discord events support")
	flag.BoolVar(&enableChannel, "enablechannel", true, "Enable Discord channel messages support")
	flag.IntVar(&eventRefresh, "event-refresh", 3600, "Event refresh interval in seconds (0 to disable auto-refresh)")
	flag.Parse()

	// Validate required flags
	if discordToken == "" {
		log.Fatal("Discord token required. Use --token flag")
	}

	// Only require channel ID if channel messages are enabled
	if enableChannel && discordChannel == "" {
		log.Fatal("Discord channel ID required. Use --channel flag")
	}

	// If events are enabled but channel is disabled, guild ID is required
	if enableEvents && !enableChannel && discordGuild == "" && discordChannel == "" {
		log.Fatal("Discord guild ID required when events are enabled without channel. Use --guild flag")
	}

	// Set global channel ID
	channelID = discordChannel

	// Initialize the messages and events slices
	messages = make([]Message, 0, maxMessages)
	events = make([]Event, 0, maxEvents)
	// Create a new Discord session with message content intent
	var err error
	discord, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Add required intents to read message content
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// Register message handler only if channel messages are enabled
	if enableChannel {
		discord.AddHandler(messageCreate)
	}

	// Open the websocket connection to Discord
	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection: %v", err)
	}
	defer discord.Close()

	// Fetch initial messages if channel is enabled
	if enableChannel {
		fetchInitialMessages()
	}

	// Fetch initial events if enabled
	if enableEvents {
		fetchInitialEvents()

		// Start event refresh loop if auto-refresh is enabled
		if eventRefresh > 0 {
			startEventRefreshLoop(time.Duration(eventRefresh) * time.Second)
		}
	}

	// Set up HTTP routes
	if enableChannel {
		http.HandleFunc("/numberofmessages", handleNumberOfMessages)
		// Default message handler
		http.HandleFunc("/", handleMessageRequest)
	}

	// Set up event routes if enabled
	if enableEvents {
		http.HandleFunc("/numberofevents", handleNumberOfEvents)
		http.HandleFunc("/event/", handleEventRequest)
	}

	// Start HTTP server
	log.Printf("Starting HTTP server on %s", httpServerAddr)
	if enableChannel {
		log.Printf("Monitoring Discord channel: %s", channelID)
	}
	if enableEvents {
		if discordGuild != "" {
			log.Printf("Discord events monitoring enabled for guild: %s", discordGuild)
		} else if channelID != "" {
			log.Printf("Discord events monitoring enabled (using guild from channel)")
		}
	}
	log.Fatal(http.ListenAndServe(httpServerAddr, nil))
}
