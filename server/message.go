package server

// Message represents a decoded JSON object
// sent by the websocket client.
type Message struct {
	Id        int                    `json:id`
	Event     string                 `json:event`
	Namespace string                 `json:namespace`
	Channel   string                 `json:channel`
	Data      map[string]interface{} `json:data`
}
