#!/bin/bash
# Discord2HTTP Easy Linux start script
# Edit this file to change settings

# Required values
TOKEN="your_discord_token_here"

# Required for channel reading, set to false if you only want event support
# Delete or comment out the CHANNEL line if you are not using the channel ID
ENABLE_CHANNEL="true"
CHANNEL="your_channel_id_here"

# Required to be set to true if you want to enable event reading
# GUILD must be set if there is no channel ID specified
# Remove or comment out the GUILD line if you are using a channel ID
ENABLE_EVENTS="false"
GUILD="your_discord_server_id_here"

# Optional parameters (these are default values)
MAX_MESSAGES=5
MAX_EVENTS=10
SERVER_ADDR="localhost:8080"
EVENT_REFRESH=3600

# Run the application with parameters
echo "Starting Discord2HTTP..."
./Discord2HTTP-linux-amd64 \
  --token="${TOKEN}" \
  --channel="${CHANNEL}" \
  --max-messages="${MAX_MESSAGES}" \
  --max-events="${MAX_EVENTS}" \
  --server-addr="${SERVER_ADDR}" \
  --guild="${GUILD}" \
  --events="${ENABLE_EVENTS}" \
  --enablechannel="${ENABLE_CHANNEL}" \
  --event-refresh="${EVENT_REFRESH}"