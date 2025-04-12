@echo off
REM Discord2HTTP Easy windows start script thing.
REM Edit this file to change settings,


REM Required values
set TOKEN=your_discord_token_here

REM Required for channel reading, set to fales if you only want event support. Delete or comment out the channel ID line if you are not using the channel ID.
set ENABLE_CHANNEL=true
set CHANNEL=your_channel_id_here

REM Required to to be set to true if you want to enable event reading. GUILD must be set if there is no channel ID specified. Remove or comment out the GUUILD line if you are using a channel ID.
set ENABLE_EVENTS=false
set GUILD=your_discord_server_id_here

REM Optional parameters (these are default values)
set MAX_MESSAGES=5
set MAX_EVENTS=10
set SERVER_ADDR=localhost:8080
set EVENT_REFRESH=3600



REM Run the application with parameters
echo Starting Discord2HTTP...
..\discord2http.exe ^
  --token=%TOKEN% ^
  --channel=%CHANNEL% ^
  --max-messages=%MAX_MESSAGES% ^
  --max-events=%MAX_EVENTS% ^
  --server-addr=%SERVER_ADDR% ^
  --guild=%GUILD% ^
  --events=%ENABLE_EVENTS% ^
  --enablechannel=%ENABLE_CHANNEL% ^
  --event-refresh=%EVENT_REFRESH%

REM Pause to see any errors
pause
