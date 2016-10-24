package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nlefler/pony/models"
	"github.com/nlefler/pony/pony"
)

const (
	hostName = ""
	host     = "0.0.0.0"
	port     = 443
)

func main() {
	pony := pony.New("", "")

	mux := http.NewServeMux()
	pony.AddRoutes(mux)

	received := make(chan models.ReceivedMessage, 100)
	pony.SetMessageReceived(received)

	go func(ch chan models.ReceivedMessage) {
		for {
			msg := <-ch
			log.Printf("example Got message %s", msg.Message.Text)
			pony.SendMessage(msg.Sender, models.OutgoingTextMessage{Text: msg.Message.Text})
		}
	}(received)

	log.Println(fmt.Sprintf("Starting, serving at port %v", port))
	err := http.ListenAndServeTLS(fmt.Sprintf("%v:%v", host, port), "", "", mux)
	if err != nil {
		log.Fatal("ListenAndServeTLS: " + err.Error())
	}
}
