package pony

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nlefler/pony/models"
)

const (
	sendURLFormat = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

// Sender holds state for sending messages
type Sender struct {
	pageToken string
	sendURL   string
}

// NewSender makes a Sender
func NewSender(pageToken string) *Sender {
	return &Sender{pageToken, fmt.Sprintf(sendURLFormat, pageToken)}
}

// Send sends a message
func (s *Sender) Send(recipient MessageParty, message OutgoingMessage) {
	payload := OutgoingMessagePayload{Recipient: recipient, Message: message}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("message.sendHandler.Send Error marshaling %v", err)
		return
	}
	payloadReader := bytes.NewBuffer(payloadData)
	resp, err := http.Post(s.sendURL, "application/json", payloadReader)
	if err != nil {
		log.Printf("message.sendHandler.Send Error sending %v", err)
		return
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("message.sendHandler.Send Error reading response %v", err)
		return
	}
	log.Printf("message.sendHandler.Send Got response %v", string(respData))
}
