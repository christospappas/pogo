package client

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type Client struct {
	ws       *websocket.Conn
	messages chan string
}

func NewClient() *Client {
	return &Client{
		messages: make(chan string, 10),
	}
}

func (c *Client) Connect(server string, origin string) {

	var err error
	if c.ws, err = websocket.Dial(server, "", origin); err != nil {
		panic("Error connecting to websocket server")
	}

	go c.receive()
	go c.send()
}

func (c *Client) receive() {
	for {
		var rec string
		websocket.Message.Receive(c.ws, &rec)
		fmt.Println(rec)
	}
}

func (c *Client) send() {
	for message := range c.messages {
		if err := websocket.Message.Send(c.ws, message); err != nil {
			fmt.Printf("[pogo] error sending message: %s", err)
		}
	}
}

func (c *Client) Emit(msg string) {
	c.messages <- msg
}
