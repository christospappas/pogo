package server

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"sync"
)

const chanBufferSize = 100

type Client struct {
	id       string
	ws       *websocket.Conn
	server   *Server
	sockets  map[string]*Socket
	ch       chan *Message
	data     map[string]interface{}
	doneCh   chan bool
	sockLock *sync.Mutex
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	// check for ws / server existence

	return &Client{
		uuid.New(), ws, server,
		make(map[string]*Socket),
		make(chan *Message, chanBufferSize),
		make(map[string]interface{}),
		make(chan bool),
		&sync.Mutex{},
	}
}

func (c *Client) Connect(namespace string) {
	c.sockLock.Lock()
	defer c.sockLock.Unlock()

	log.Println("[pogo] Client Connected [" + namespace + "]: " + c.id)
	ns := c.server.Of(namespace)
	socket := ns.AddClient(c)
	c.sockets[ns.name] = socket
}

func (c *Client) Disconnect(namespace string) {
	c.sockLock.Lock()
	defer c.sockLock.Unlock()
	socket := c.sockets[namespace]
	socket.Disconnect()
	delete(c.sockets, namespace)
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

func (c *Client) onClose() {
	log.Println("Client Disconnected: " + c.id)
	for _, socket := range c.sockets {
		socket.Disconnect()
	}
}

// forward message to appropriate socket handler
func (c *Client) onData(msg *Message) {
	switch msg.Event {
	case "pogo:connect":
		c.Connect(msg.Namespace)
	case "pogo:disconnect":
		c.Disconnect(msg.Namespace)
	default:
		if msg.Namespace != "" {
			if socket, ok := c.sockets[msg.Namespace]; ok {
				socket.OnData(msg)
			} else {
				// namespace doesn't exist
			}
		} else {
			c.sockets["/"].OnData(msg)
		}
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
			c.onClose()
			c.doneCh <- true
			return
		}
	}

}

func (c *Client) reader() {
	for {
		select {
		case <-c.doneCh:
			c.onClose()
			c.doneCh <- true
			return
		default:

			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)

			if err == io.EOF {
				c.doneCh <- true
				return
			} else {
				c.onData(&msg)
			}
		}
	}
}
