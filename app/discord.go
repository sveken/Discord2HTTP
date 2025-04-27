package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp" // Added regexp package
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Regular expression to find user mentions like <@1234567890>
var userMentionRegex = regexp.MustCompile(`<@(\d+)>`)

// resolveMentions finds user mentions (<@USER_ID>) in content and replaces them with @username
func resolveMentions(content string) string {
	return userMentionRegex.ReplaceAllStringFunc(content, func(mention string) string {
		// Extract user ID from the mention string (e.g., "<@1234567890>" -> "1234567890")
		matches := userMentionRegex.FindStringSubmatch(mention)
		if len(matches) < 2 {
			return mention // Should not happen with the defined regex, but safety first
		}
		userID := matches[1]

		// Fetch user information from Discord
		user, err := discord.User(userID)
		if err != nil {
			// Log the error and return the original mention if user fetch fails
			log.Printf("Error fetching user %s: %v", userID, err)
			return mention
		}

		// Return the formatted username mention
		return fmt.Sprintf("@%s", user.Username)
	})
}

// fetchInitialMessages retrieves the most recent messages from the Discord channel
func fetchInitialMessages() {
	msgs, err := discord.ChannelMessages(channelID, maxMessages, "", "", "")
	if err != nil {
		log.Printf("Error fetching initial messages: %v", err)
		return
	}

	// Discord returns messages from newest to oldest, so we need to add them in reverse
	messagesMutex.Lock()
	defer messagesMutex.Unlock()
	messages = make([]Message, 0, len(msgs))
	for i, msg := range msgs {
		messages = append(messages, Message{
			Index:   i,
			UserID:  msg.Author.ID,
			User:    msg.Author.Username,
			Content: msg.Content,
			Time:    msg.Timestamp.UTC(), // Convert timestamp to UTC
		})
	}

	log.Printf("Loaded %d initial messages", len(messages))
}

// messageCreate is called when a new message is created in Discord
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from other channels
	if m.ChannelID != channelID {
		return
	}

	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Add the new message to our list
	messagesMutex.Lock()
	defer messagesMutex.Unlock()
	// Create new message
	newMsg := Message{
		Index:   0, // New message is always at index 0
		UserID:  m.Author.ID,
		User:    m.Author.Username,
		Content: m.Content,
		Time:    m.Timestamp.UTC(), // Convert timestamp to UTC
	}

	// Shift all existing message indices
	for i := range messages {
		messages[i].Index++
	}

	// Add new message at the beginning
	if len(messages) < maxMessages {
		messages = append([]Message{newMsg}, messages...)
	} else {
		// Remove the oldest message
		messages = append([]Message{newMsg}, messages[:maxMessages-1]...)
	}
}

// handleNumberOfMessages returns the current number of stored messages
func handleNumberOfMessages(w http.ResponseWriter, r *http.Request) {
	messagesMutex.RLock()
	defer messagesMutex.RUnlock()

	fmt.Fprintf(w, "%d", len(messages))
}

// handleMessageRequest handles requests for message data by index
func handleMessageRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// If the path is empty or invalid, return an error
	if len(parts) != 2 || (parts[1] != "user" && parts[1] != "message") {
		http.Error(w, "Invalid request format. Use /{index}/user or /{index}/message", http.StatusBadRequest)
		return
	}

	// Parse message index
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid message index", http.StatusBadRequest)
		return
	}

	messagesMutex.RLock()
	defer messagesMutex.RUnlock()

	// Check if the index is valid
	if index < 0 || index >= len(messages) {
		http.Error(w, "Message index out of range", http.StatusNotFound)
		return
	}
	// Return the requested information
	if parts[1] == "user" {
		fmt.Fprint(w, messages[index].User)
	} else if parts[1] == "message" {
		if messages[index].Content == "" {
			fmt.Fprint(w, "[Empty message content]")
		} else {
			// Resolve mentions before sending the content
			resolvedContent := resolveMentions(messages[index].Content)
			fmt.Fprint(w, resolvedContent)
		}
	}
}

// fetchInitialEvents retrieves scheduled events from Discord server
func fetchInitialEvents() {
	var guildID string

	// If guild ID is directly provided, use it
	if discordGuild != "" {
		guildID = discordGuild
	} else if channelID != "" {
		// Otherwise try to get the guild ID from the channel
		channel, err := discord.Channel(channelID)
		if err != nil {
			log.Printf("Error fetching channel info: %v", err)
			return
		}
		guildID = channel.GuildID
	} else {
		log.Printf("Cannot fetch events: no guild ID or channel ID provided")
		return
	}

	// Fetch all scheduled events for the guild
	scheduledEvents, err := discord.GuildScheduledEvents(guildID, false)
	if err != nil {
		log.Printf("Error fetching events: %v", err)
		return
	}

	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	// Clear existing events and add the current ones
	events = make([]Event, 0, len(scheduledEvents))
	for i, evt := range scheduledEvents {
		// Convert times to UTC
		startTime := evt.ScheduledStartTime.UTC()
		var endTime time.Time
		if evt.ScheduledEndTime != nil {
			endTime = evt.ScheduledEndTime.UTC()
		}

		// Get location (entity type can be STAGE, VOICE or EXTERNAL)
		location := "Unknown"
		if evt.EntityType == discordgo.GuildScheduledEventEntityTypeExternal {
			location = evt.EntityMetadata.Location
		} else if evt.EntityType == discordgo.GuildScheduledEventEntityTypeStageInstance ||
			evt.EntityType == discordgo.GuildScheduledEventEntityTypeVoice {
			// Try to get channel name for stage/voice events
			if evt.ChannelID != "" {
				if ch, err := discord.Channel(evt.ChannelID); err == nil {
					location = ch.Name
				}
			}
		}
		// Construct the image URL if event has an image
		var imageURL string
		if evt.Image != "" {
			imageURL = fmt.Sprintf("https://cdn.discordapp.com/guild-events/%s/%s.png", evt.ID, evt.Image)
		}

		events = append(events, Event{
			Index:       i,
			ID:          evt.ID,
			Name:        evt.Name,
			Description: evt.Description,
			Location:    location,
			StartTime:   startTime,
			EndTime:     endTime,
			GuildID:     evt.GuildID,
			ImageURL:    imageURL,
		})
	}

	log.Printf("Loaded %d events", len(events))
}

// startEventRefreshLoop starts a goroutine that periodically refreshes events
func startEventRefreshLoop(refreshInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("Refreshing Discord events...")
			fetchInitialEvents()
		}
	}()
	log.Printf("Event refresh loop started with interval: %v", refreshInterval)
}

// handleNumberOfEvents returns the current number of stored events
func handleNumberOfEvents(w http.ResponseWriter, r *http.Request) {
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	fmt.Fprintf(w, "%d", len(events))
}

// handleEventRequest handles requests for event data by index
func handleEventRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Check if the path format is valid
	if len(parts) != 3 || parts[0] != "event" {
		http.Error(w, "Invalid request format. Use /event/{index}/{field}", http.StatusBadRequest)
		return
	}

	// Parse event index
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid event index", http.StatusBadRequest)
		return
	}

	// Field to retrieve is in parts[2]
	field := parts[2]

	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	// Check if the index is valid
	if index < 0 || index >= len(events) {
		http.Error(w, "Event index out of range", http.StatusNotFound)
		return
	} // Return the requested information based on the field
	switch field {
	case "eventname":
		fmt.Fprint(w, events[index].Name)
	case "time":
		// Return time in ISO 8601 format in UTC
		fmt.Fprint(w, events[index].StartTime.Format(time.RFC3339))
	case "location":
		fmt.Fprint(w, events[index].Location)
	case "description":
		fmt.Fprint(w, events[index].Description)
	case "bannerurl":
		// Return the CDN URL for the event image
		if events[index].ImageURL == "" {
			http.Error(w, "This event has no image", http.StatusNotFound)
			return
		}
		fmt.Fprint(w, events[index].ImageURL)
	default:
		http.Error(w, "Invalid field requested. Use eventname, time, location, description, or bannerurl", http.StatusBadRequest)
	}
}
