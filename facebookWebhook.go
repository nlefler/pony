package pony

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type facebookMessengerWebhook struct {
	webhookPrefix   string
	validationToken string
	pageToken       string
	decoder         *FacebookMessengerDecoder
}

func NewFacebookMessengerWebhook(pageName string, validationToken string, pageToken string, decoder *FacebookMessengerDecoder) Webhook {
	return &facebookMessengerWebhook{pageName, validationToken, pageToken, decoder}
}

func (wh *facebookMessengerWebhook) addRoutes(mux *http.ServeMux) {
	makeHandler := func(wh *facebookMessengerWebhook,
		handler func(*facebookMessengerWebhook, http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			handler(wh, w, req)
		}
	}
	mux.HandleFunc(fmt.Sprintf("%s/webhook", wh.webhookPrefix), makeHandler(wh, facebookWebhookDispatcher))
	mux.HandleFunc(fmt.Sprintf("%s/authorize", wh.webhookPrefix), makeHandler(wh, facebookAuthorizeHandler))
}

func facebookAuthorizeHandler(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	log.Println("pony.pony.authorize")
	w.WriteHeader(http.StatusOK)
}

func facebookWebhookDispatcher(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	// TODO(nl): verify signature
	switch req.Method {
	case "GET":
		facebookWebhookValidate(wh, w, req)
	case "POST":
		jsonBytes, err := ioutil.ReadAll(req.Body)
		if len(jsonBytes) == 0 || err != nil {
			log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		wh.decoder.receive(jsonBytes)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func facebookWebhookValidate(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("hub.mode")
	if mode != "subscribe" {
		log.Printf("pony.pony.validate Failed, mode is %s", mode)
		w.WriteHeader(http.StatusOK)
		return
	}
	token := req.FormValue("hub.verify_token")
	if token != wh.validationToken {
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
