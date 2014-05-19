package pogo

// Message represents a decoded JSON object
// sent by the websocket client.
type Message struct {
	Id      int                    `json:id`
	Event   string                 `json:event`
	Channel string                 `json:channel`
	Data    map[string]interface{} `json:data`
}
