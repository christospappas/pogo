package server

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
)

const chanBufferSize = 100

type Client struct {
	id      string
	ws      *websocket.Conn
	server  *Server
	sockets map[string]*Socket
	ch      chan *Message
	data    map[string]interface{}
	doneCh  chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	// check for ws / server existence

	return &Client{
		uuid.New(), ws, server,
		make(map[string]*Socket),
		make(chan *Message, chanBufferSize),
		make(map[string]interface{}),
		make(chan bool),
	}
}

func (c *Client) Connect(namespace string) {
	log.Println("[pogo] Client Connected [" + namespace + "]: " + c.id)
	ns := c.server.Of(namespace)
	ns.Connect(c)
}

func (c *Client) Listen() {
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

// Disconnects a client from all channels and the server
func (c *Client) OnClose() {
	// for _, ch := range c.subscriptions {
	// 	ch.Unsubscribe(c)
	// }
	log.Println("[pogo] Client Disconnected: " + c.id)
}

// forward message to appropriate socket handler
func (c *Client) OnData(msg *Message) {
	if msg.Event != "pogo:connect" && msg.Namespace != "" {
		if socket, ok := c.sockets[msg.Namespace]; ok {
			socket.OnData(msg)
		}
	} else {
		c.sockets["/"].OnData(msg)
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
			c.OnClose()
			c.doneCh <- true
			return
		}
	}

}

func (c *Client) reader() {
	for {
		select {
		case <-c.doneCh:
			c.OnClose()
			c.doneCh <- true
			return
		default:

			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)

			if err == io.EOF {
				c.doneCh <- true
				return
			} else {
				c.OnData(&msg)
			}
		}
	}
}
