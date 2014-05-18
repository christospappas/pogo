package main

import (
	"fmt"
	"github.com/christospappas/pogo/server"
	"net/http"
)

func main() {
	// TODO: Allow multiple endpoints on the same server
	s := pogo.NewServer("/")

	// TODO: bind to channel events
	// s.Channel("/posts/:id").On("update", func(..
	s.On("someEvent", func(msg *pogo.Message, c *pogo.Client) { // figure out how to inject
		// TODO: s.Channel("blah").Broadcast(msg)
		fmt.Println("someEvent recieved!")
	})

	s.Listen()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
