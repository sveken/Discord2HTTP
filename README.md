# Discord2HTTP
This project is for turning a Discord channel, events or both into simple HTTP strings that can be parsed and formatted in games to make, for example, an event panel that auto-updates and stays in sync with the Discord server. 

This was designed around Resonite for use as a way to automatically sync in-world panels to a Discord server.

## In Resonite
Place info on prebuilt panels here with pictures. After working with the wolfinator (Big Red Wolfy)


#### The Web Endpoints
A list of the web endpoints that are usable and parsable in games. [available here.](docs/endpoints.md)
#### List of Flags
A list of all the flags that can be used against the binary. [available here](docs/flags.md)
## Setup
You will be required to make a Discord bot that has `Message content intent` enabled in the bot settings on the Discord developers portal. 
Copy the bot token from your created application from the Discord developers portal.
If you only need events, you can use the server/guild ID found by enabling developer mode in Discord, then right clicking the server icon and copying the ID. 
Otherwise if no Guild/Server ID is set it will fallback to the channel ID to identify the server for events.

If you are only watching a channel you only need the channel ID. 

## Hosting

### Release Binaries.
[Download the latest release](https://github.com/sveken/Discord2HTTP/releases) for the platform you want. For windows and linux hosts there is an example start scripts under /scripts directory in this repo that can be edited with the required values and used for easy startup if you don't want to use docker.

### Docker Compose

```
services:
  discord2http:
    image: ghcr.io/sveken/discord2http:latest
    container_name: discord2http
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DISCORD_TOKEN=your_discord_bot_token_here
      - DISCORD_CHANNEL=your_discord_channel_id_here
      - DISCORD_GUILD=your_discord_guild_id_here
      - MAX_MESSAGES=5
      - MAX_EVENTS=10
      - ENABLE_EVENTS=false
      - ENABLE_CHANNEL=true
      - EVENT_REFRESH=3600
```


## Note. This program will make the channel you configure or the events if enabled public/accessible to anyone with access to this server/program.
This is required for use in game, for example a Resonite headless server to access this information. 

If this is a problem one example to get around this is if you have your headless server in docker, you can add this programs docker image to the same stack, meaning it will only be available for the headless server directly. 

Example Compose file for this. In resonite the url the headless would use would be discord2http:8080

```
services:
  headless:
    container_name: resonite-headless
    image: ghcr.io/voxelbonecloud/headless-docker:main 
    env_file: .env
    environment:
      CONFIG_FILE: Config.json
      ENABLE_MODS: false
      #ADDITIONAL_ARGUMENTS:
    tty: true
    stdin_open: true
    user: "1000:1000"
    volumes:
      - "/etc/localtime:/etc/localtime:ro"
      - ./Headless_Configs:/Config
      - ./Headless_Logs:/Logs
      - ./RML:/RML
    restart: on-failure:5
  discord2http:
    image: ghcr.io/sveken/discord2http:latest
    container_name: discord2http
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DISCORD_TOKEN=your_discord_bot_token_here
      - DISCORD_CHANNEL=your_discord_channel_id_here
      - DISCORD_GUILD=your_discord_guild_id_here
      - MAX_MESSAGES=5
      - MAX_EVENTS=10
      - ENABLE_EVENTS=false
      - ENABLE_CHANNEL=true
      - EVENT_REFRESH=3600
```
## Credits
This uses the DiscordGo Package from https://github.com/bwmarrin/discordgo