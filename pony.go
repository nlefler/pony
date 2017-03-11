package pony

import (
	"net/http"
)

// Pony receives webhook messages and delegates
type Pony struct {
	sender *Sender
}

func NewPony() *Pony {
	return &Pony{}
}

// SetMessageReceived replaces the channel received messages will be sent to
func (p *Pony) SetMessageReceived(ch chan Message) {
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
