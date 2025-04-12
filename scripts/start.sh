#!/bin/sh

# Start script for Discord2HTTP
# This script converts environment variables to command-line flags

# Set default values for optional parameters
MAX_MESSAGES="${MAX_MESSAGES:-5}"
MAX_EVENTS="${MAX_EVENTS:-10}"
SERVER_ADDR="${SERVER_ADDR:-localhost:8080}"
EVENT_REFRESH="${EVENT_REFRESH:-3600}"
ENABLE_EVENTS="${ENABLE_EVENTS:-false}"
ENABLE_CHANNEL="${ENABLE_CHANNEL:-true}"

# Check for required token
if [ -z "${DISCORD_TOKEN}" ]; then
  echo "Error: DISCORD_TOKEN environment variable must be set"
  exit 1
fi

# Build base command with required token
CMD="/app/discord2http --token=${DISCORD_TOKEN}"

# Add optional parameters - only add channel and guild if they are provided
[ -n "${DISCORD_CHANNEL}" ] && CMD="${CMD} --channel=${DISCORD_CHANNEL}"
[ -n "${DISCORD_GUILD}" ] && CMD="${CMD} --guild=${DISCORD_GUILD}"

# Add other configuration parameters with their default values
CMD="${CMD} --max-messages=${MAX_MESSAGES}"
CMD="${CMD} --max-events=${MAX_EVENTS}"
CMD="${CMD} --server-addr=${SERVER_ADDR}"
CMD="${CMD} --event-refresh=${EVENT_REFRESH}"

# Handle boolean flags correctly
[ "${ENABLE_EVENTS}" = "true" ] && CMD="${CMD} --events"
[ "${ENABLE_CHANNEL}" = "false" ] && CMD="${CMD} --enablechannel=false"

# Print the command for debugging (omit the token for security)
echo "Starting Discord2HTTP with configuration:"
echo "$CMD" | sed "s/--token=[^ ]*/--token=******/"

# Execute the command
exec $CMD
