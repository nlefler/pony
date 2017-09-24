package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	facebookMessengerSendURLFormat = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

type facebookMessengerSender struct {
	pageToken string
	sendURL   string
}

func newFacebookMessengerSender(pageToken string) facebookMessengerSender {
	return facebookMessengerSender{pageToken,
		fmt.Sprintf(facebookMessengerSendURLFormat,
			pageToken)}
}

type outgoingMessage struct {
	Recipient Recipient              `json:"recipient"`
	Message   outgoingContentMessage `json:"message,omitempty"`
	Action    SenderActionType       `json:"sender_action,omitempty"`
}

type outgoingContentMessage struct {
	Text         string             `json:"text"`
	Attachment   *MessageAttachment `json:"attachment,omitempty"`
	QuickReplies []QuickReplies     `json:"quick_replies,omitempty"`
}

// Send sends a message
func (s *facebookMessengerSender) send(message interface{}) {
	outMessage := outgoingMessage{}
	switch message := message.(type) {
	case ContentMessage:
		outMessage.Recipient = message.Recipient
		outMessage.Message = outgoingContentMessage{}
		outMessage.Message.Text = message.Text
		if message.Attachment.Type != "" {
			outMessage.Message.Attachment = &message.Attachment
		}
		outMessage.Message.QuickReplies = message.QuickReplies
	case Action:
		outMessage.Recipient = message.Recipient
		outMessage.Action = message.Action
	}
	payloadData, err := json.Marshal(outMessage)
	log.Println(string(payloadData))
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
