/*
Package pogo provides an out-of-box websocket server with a few defaults.

Usage Example:

*/
package pogo

import (
	"code.google.com/p/go.net/websocket"
	"github.com/christospappas/pogo/server"
	"log"
	"net/http"
)

var defaultServer *server.Server

func init() {
	defaultServer = server.New()
}

func On(event string, handler server.Handler) {
	defaultServer.On(event, handler)
}

func Of(ns string) *server.Namespace {
	return defaultServer.Of(ns)
}

func Namespace(name string, fn func(ns *server.Namespace)) {
	ns := defaultServer.Of(name)
	fn(ns)
}

func Listen() {
	log.Println("[pogo] Listening on port 8080")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				// do something
			}
		}()

		client := server.NewClient(ws, defaultServer)

		if client != nil {
			client.Connect("/")
			client.Listen()
		}

	}

	http.Handle("/", websocket.Handler(onConnected))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
