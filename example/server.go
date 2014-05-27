package main

import (
	"github.com/christospappas/pogo"
	"github.com/christospappas/pogo/server"
	"log"
)

func main() {

	pogo.Namespace("/analytics", func(ns *server.Namespace) {

		ns.On("track", func(msg *server.Message, c *server.Client) {
			// ns.Channel("something").Send("")
		})

		ns.On("version", func(msg *server.Message, c *server.Client) {
			log.Println("this is version 1234")
		})

		ns.On("sendAll", func(msg *server.Message, c *server.Client) {
			log.Println("OMG! Sending woohoo!...")
		})

	})

	// TODO: Pattern based channel & Event matching
	// pogo.Channel("/posts/:id").On("msg", func(msg *server.Message, c *server.Client) {

	// })

	// pogo.Resource("posts", postsController)

	// pogo.Create("/posts")

	// pogo.Update("/posts/{id}", func(msg *server.Message, c *server.Client) {

	// })

	// pogo.Retrieve("/posts/{id}", func(msg *server.Message, c *server.Client) {

	// })

	// pogo.List("/posts", func(msg *server.Message, c *server.Client) {

	// })

	pogo.On("track", func(msg *server.Message, c *server.Client) {
		log.Println("oooh we received a track event on / ")
	})

	pogo.On("sendAll", func(msg *server.Message, c *server.Client) {
		log.Println("Sending woohoo!...")
	})

	pogo.Listen()

}
