# Discord2HTTP Command Line Flags

Configuration flags that can be used when launching Discord2HTTP.

## Basic Configuration

| Flag | Description |
|------|-------------|
| `--token` | Discord bot token (required) |
| `--channel` | Discord channel ID (required if enablechannel is true) |
| `--guild` | Discord guild/server ID (required if events enabled without channel) |
| `--server-addr` | HTTP server address (default: "localhost:8080") |

## Message Configuration

| Flag | Description |
|------|-------------|
| `--enablechannel` | Enable Discord channel messages support (default: true) |
| `--max-messages` | Maximum number of messages to store (default: 5) |

## Event Configuration

| Flag | Description |
|------|-------------|
| `--events` | Enable Discord events support (default: false) |
| `--max-events` | Maximum number of events to store (default: 10) |
| `--event-refresh` | Event refresh interval in seconds, 0 to disable auto-refresh (default: 3600) |
