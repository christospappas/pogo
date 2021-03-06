/*
Package server is a websocket microframework inspired by Socket.IO.

Usage example:

  m := server.New()

*/
package server

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"sync"
)

// Handler can be any callable function.
type Handler func(msg *Message, c *Client)

type Server struct {
	namespaces map[string]*Namespace
	nsLock     *sync.Mutex
}

func New() *Server {
	log.Println("Starting pogo " + Version())
	return &Server{make(map[string]*Namespace), &sync.Mutex{}}
}

func (s *Server) On(event string, handler Handler) {
	ns := s.Of("/")
	ns.On(event, handler)
}

func (s *Server) Of(name string) *Namespace {
	s.nsLock.Lock()
	s.nsLock.Unlock()

	if ns, ok := s.namespaces[name]; ok {
		return ns
	} else {
		ns := NewNamespace(name)
		s.namespaces[name] = ns
		return ns
	}
}

func (s *Server) OnConnect(ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		if err != nil {
			// do something
		}
	}()

	client := NewClient(ws, s)

	if client != nil {
		client.Connect("/")
		client.Listen()
	}
}
