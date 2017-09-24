package facebook

import (
	"fmt"
	"net/http"
)

/// Service

type FacebookMessenger struct {
	id      string
	webhook *facebookMessengerWebhook
	sender  *facebookMessengerSender
}

func NewFacebookMessenger(pageName string, validationToken string, pageToken string) *FacebookMessenger {
	id := fmt.Sprintf("com.pony.facebook.messenger.%s", pageName)
	webhook := &facebookMessengerWebhook{pageName, validationToken, pageToken,
		&facebookMessengerDecoder{}, make(chan ReceivedMessage, 100)}
	sender := newFacebookMessengerSender(pageToken)
	return &FacebookMessenger{id, webhook, &sender}
}

func (fb *FacebookMessenger) Setup(mux *http.ServeMux) {
	fb.webhook.addRoutes(mux)
}

func (fb *FacebookMessenger) ID() string {
	return fb.id
}

func (fb *FacebookMessenger) Send(msg ContentMessage) {
	fb.sender.send(msg)
}

func (fb *FacebookMessenger) SendAction(action Action) {
	fb.sender.send(action)
}

func (fb *FacebookMessenger) ReceiveOn() <-chan ReceivedMessage {
	return fb.webhook.receiveOn
}
