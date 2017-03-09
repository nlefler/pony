package pony

import (
	"log"
	"net/http"

	"github.com/nlefler/pony/message"
	"github.com/nlefler/pony/models"
)

// Pony receives webhook messages and delegates
type Pony struct {
	receiptHandler  *ReceiptHandler
	sender          *Sender
}


// SetMessageReceived replaces the channel received messages will be sent to
func (p *Pony) SetMessageReceived(ch chan ReceivedMessage) {
	p.receiptHandler.Received = ch
}

// SendMessage sends a message
func (p *Pony) SendMessage(recipient MessageParty, msg OutgoingMessage) {
	p.sender.Send(recipient, msg)
}

// AddRoutes adds webhook routes to the provided ServeMux
func (p *Pony) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/up", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("UP")) })
}
