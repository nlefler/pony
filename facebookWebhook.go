package pony

import "net/http"

type FacebookMessengerWebhook struct {
	validationToken string
	pageToken string
}

func NewFacebookMessengerWebhook(validationToken string, pageToken string) *FacebookMessengerWebhook {
	return &FacebookMessengerWebhook{validationToken, pageToken}
}

func (wh *FacebookMessengerWebhook) receive(rmsg ReceivedMessage) {

}

func (wh *FacebookMessengerWebhook) addRoutes(mux *http.ServeMux) {

}
