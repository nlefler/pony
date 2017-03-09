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

func makeHandler(p *Pony, handler func(*Pony, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(p, w, req)
	}
}

func webhookDispatcher(p *Pony, w http.ResponseWriter, req *http.Request) {
	// TODO(nl): verify signature
	switch req.Method {
	case "GET":
		webhookValidate(p, w, req)
	case "POST":
		p.receiptHandler.Handle(w, req)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func authorizeHandler(p *Pony, w http.ResponseWriter, req *http.Request) {
	log.Println("pony.pony.authorize")
	w.WriteHeader(http.StatusOK)
}

func webhookValidate(p *Pony, w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("hub.mode")
	if mode != "subscribe" {
		log.Printf("pony.pony.validate Failed, mode is %s", mode)
		w.WriteHeader(http.StatusOK)
		return
	}
	token := req.FormValue("hub.verify_token")
	if token != p.validationToken {
		log.Printf("pony.pony.validate Failed, token is %s", token)
		w.WriteHeader(http.StatusOK)
		return
	}

	challenge := req.FormValue("hub.challenge")
	if challenge == "" {
		log.Printf("pony.pony.validate Failed, no challenge")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("pony.pony.validate Validated")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
	return
}
