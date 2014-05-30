package main

import (
	"bufio"
	"fmt"
	"github.com/christospappas/pogo/client"
	"os"
	"strings"
)

func main() {
	client := client.NewClient()

	client.Connect("ws://localhost:8080/", "http://localhost/")

	read := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := read.ReadString('\n')

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}

		switch line {
		case "exit":
			return
		default:
			client.Emit(line)
		}
	}

}
