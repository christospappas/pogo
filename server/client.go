package pogo

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
)

const chanBufferSize = 100

type Client struct {
	id            string
	ws            *websocket.Conn
	server        *Server
	ch            chan *Message
	subscriptions map[string]*Channel
	data          map[string]interface{}
	doneCh        chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	// check for ws / server existence

	return &Client{
		uuid.New(), ws, server,
		make(chan *Message, chanBufferSize),
		make(map[string]*Channel),
		make(map[string]interface{}),
		make(chan bool),
	}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Listen() {
	log.Println("[pogo] Client Connected: " + c.id)
	go c.writer()
	c.reader()
}

func (c *Client) Send(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		log.Println("Client is disconnected: " + c.id)
	}

}

func (c *Client) writer() {
	for {
		select {
		case msg := <-c.ch:
			log.Println("[pogo] Sending: ", msg)
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				log.Println("[pogo] Error sending message")
			}
		case <-c.doneCh:
			c.server.HandleDisconnect(c)
			c.doneCh <- true
			return
		}
	}

}

func (c *Client) reader() {
	for {
		select {
		case <-c.doneCh:
			c.server.HandleDisconnect(c)
			c.doneCh <- true
			return
		default:

			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)

			if err == io.EOF {
				c.doneCh <- true
				return
			} else {
				c.processMessage(&msg)
			}
		}
	}
}

// Need to handle errors here
func (c *Client) processMessage(msg *Message) {
	log.Println("[pogo] Command: " + msg.Event)
	switch msg.Event {
	case "channel:subscribe":
		chanName := msg.Data["channel"].(string)
		c.server.HandleSubscribe(chanName, c)
	case "channel:unsubscribe":
		chanName := msg.Data["channel"].(string)
		c.server.HandleUnsubscribe(chanName, c)
	default:
		c.server.HandleMessage(msg, c)
	}

	return
}
