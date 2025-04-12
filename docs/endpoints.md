# Discord2HTTP API Endpoints

All endpoints should be appended to your host URL.

For example, if hosted locally: `http://localhost:8080/numberofmessages`

## Message Endpoints

| Endpoint | Description |
|----------|-------------|
| `/numberofmessages` | Returns the count of stored messages |
| `/{index}/message` | Retrieves message content by index |
| `/{index}/user` | Retrieves message author username by index |

## Event Endpoints

| Endpoint | Description |
|----------|-------------|
| `/numberofevents` | Returns the count of stored Discord events |
| `/event/{index}/eventname` | Gets event name |
| `/event/{index}/time` | Gets event start time in ISO 8601 format |
| `/event/{index}/location` | Gets event location |
| `/event/{index}/description` | Gets event description |
| `/event/{index}/bannerurl` | Gets event banner image URL |
