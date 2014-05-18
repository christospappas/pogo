package pogo

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"reflect"
	"sync"
)

type Server struct {
	pattern  string
	clients  map[string]*Client
	channels map[string]Channel
	handlers map[string]Handler
	chanLock *sync.Mutex
}

// Handler can be any callable function.
type Handler func(msg *Message, c *Client)

func NewServer(pattern string) *Server {
	log.Println("Starting pogo " + Version())
	clients := make(map[string]*Client)
	channels := make(map[string]Channel)
	handlers := make(map[string]Handler)
	return &Server{pattern, clients, channels, handlers, &sync.Mutex{}}
}

func (s *Server) On(event string, handler Handler) {
	validateHandler(handler)
	s.handlers[event] = handler
}

func (s *Server) HandleSubscribe(chanName string, c *Client) {

	if ch, ok := s.channels[chanName]; ok {
		log.Println("[pogo] Adding Client to Channel: " + chanName)
		ch.Subscribe(c)
	} else {
		s.chanLock.Lock()
		defer s.chanLock.Unlock()

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
	if fn, ok := s.handlers[msg.Event]; ok {
		fn(msg, c)
	} else {
		log.Println("[pogo] No Event found: " + msg.Event)
	}
}

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

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("[pogo] handler must be a callable func")
	}
}
