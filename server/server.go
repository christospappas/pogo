package pogo

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"sync"
)

// Server maintains the connected clients, channels and handlers
type Server struct {
	pattern  string
	clients  map[string]*Client
	channels map[string]Channel
	handlers map[string]Handler
	chanLock *sync.Mutex
}

// Handler can be any callable function.
type Handler func(msg *Message, c *Client)

// Creates new server for specified pattern.
func NewServer(pattern string) *Server {

	log.Println("Starting pogo " + Version())
	clients := make(map[string]*Client)
	channels := make(map[string]Channel)
	handlers := make(map[string]Handler)
	return &Server{pattern, clients, channels, handlers, &sync.Mutex{}}
}

func (s *Server) On(event string, handler Handler) {
	s.handlers[event] = handler
}

func (s *Server) HandleSubscribe(chanName string, c *Client) {
	s.chanLock.Lock()
	defer s.chanLock.Unlock()

	if ch, ok := s.channels[chanName]; ok {
		log.Println("[pogo] Adding Client to Channel: " + chanName)
		ch.Subscribe(c)
	} else {
		log.Println("[pogo] Creating new Channel: " + chanName)
		ch := NewChannel(chanName)
		s.channels[chanName] = ch
		ch.Subscribe(c)
	}

}

func (s *Server) HandleUnsubscribe(chanName string, c *Client) {

	if ch, ok := s.channels[chanName]; ok {
		ch.Unsubscribe(c)

		// delete channel if no more subscribers
		if ch.SubscriberCount() == 0 {
			delete(s.channels, chanName)
		}

	} else {
		log.Println("[pogo] Channel not found: " + chanName)
	}

}

func (s *Server) HandleMessage(msg *Message, c *Client) {

	log.Println("[pogo] Command: " + msg.Event)

	switch msg.Event {
	case "channel:subscribe":
		c.server.HandleSubscribe(msg.Channel, c)
	case "channel:unsubscribe":
		c.server.HandleUnsubscribe(msg.Channel, c)
	default:
		if fn, ok := s.handlers[msg.Event]; ok {
			fn(msg, c)
		} else {
			log.Println("[pogo] No Event found: " + msg.Event)
		}
	}

}

// Disconnects a client from all channels and the server
func (s *Server) HandleDisconnect(c *Client) {
	for _, ch := range c.subscriptions {
		ch.Unsubscribe(c)
	}
	delete(s.clients, c.id)
	log.Println("[pogo] Client Disconnected: " + c.id)
}

func (s *Server) Listen() {
	log.Println("[pogo] Listening...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				// do something
			}
		}()

		client := NewClient(ws, s)

		if client != nil {
			s.clients[client.id] = client
			client.Listen()
		}

	}

	http.Handle(s.pattern, websocket.Handler(onConnected))

}
