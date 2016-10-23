package main

import (
	"log"
	"net/http"

	"github.com/nlefler/pony/models"
	"github.com/nlefler/pony/pony"
)

func main() {
	pony := pony.New("")

	mux := http.NewServeMux()
	pony.AddRoutes(mux)

	received := make(chan models.ReceivedMessage, 100)
	pony.SetMessageReceived(received)

	go func(ch chan models.ReceivedMessage) {
		for {
			msg := <-ch
			log.Printf("example Got message %s", msg.Message.Text)
		}
	}(received)

	http.ListenAndServe("0.0.0.0:8080", mux)
}
