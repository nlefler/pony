package facebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type facebookMessengerWebhook struct {
	webhookPrefix   string
	validationToken string
	pageToken       string
	decoder         *facebookMessengerDecoder
	receiveOn       chan ReceivedMessage
}

func (wh *facebookMessengerWebhook) addRoutes(mux *http.ServeMux) {
	log.Printf("fb add routes %s\n", wh.webhookPrefix)
	makeHandler := func(wh *facebookMessengerWebhook,
		handler func(*facebookMessengerWebhook, http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			handler(wh, w, req)
		}
	}
	mux.HandleFunc(fmt.Sprintf("/%s/webhook", wh.webhookPrefix), makeHandler(wh, facebookWebhookDispatcher))
	mux.HandleFunc(fmt.Sprintf("/%s/authorize", wh.webhookPrefix), makeHandler(wh, facebookAuthorizeHandler))
}

func facebookAuthorizeHandler(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	log.Println("pony.facebook.authorize")
	w.WriteHeader(http.StatusOK)
}

func facebookWebhookDispatcher(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	// TODO(nl): verify signature
	switch req.Method {
	case "GET":
		facebookWebhookValidate(wh, w, req)
	case "POST":
		log.Printf("%s: POST", wh.webhookPrefix)
		jsonBytes, err := ioutil.ReadAll(req.Body)
		if len(jsonBytes) == 0 || err != nil {
			log.Printf("pony.facebook.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("%s: 400 %v", wh.webhookPrefix, err)
			return
		}
		messages, err := wh.decoder.receive(jsonBytes)
		if err != nil {
			log.Printf("%v", err)
			log.Printf("%s: 500 %v", wh.webhookPrefix, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, m := range messages {
			log.Printf("%s: Message %v", wh.webhookPrefix, m)
			wh.receiveOn <- m
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func facebookWebhookValidate(wh *facebookMessengerWebhook, w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("hub.mode")
	if mode != "subscribe" {
		log.Printf("pony.facebook.validate Failed, mode is %s", mode)
		w.WriteHeader(http.StatusOK)
		return
	}
	token := req.FormValue("hub.verify_token")
	if token != wh.validationToken {
		log.Printf("pony.facebook.validate Failed, token is %s", token)
		w.WriteHeader(http.StatusOK)
		return
	}

	challenge := req.FormValue("hub.challenge")
	if challenge == "" {
		log.Printf("pony.facebook.validate Failed, no challenge")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("pony.facebook.validate Validated")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
	return
}

/// Decoder

type facebookMessengerDecoder struct {
}

func (decoder *facebookMessengerDecoder) receive(msgData []byte) ([]ReceivedMessage, error) {
	var call facebookMessengerWebhookMessageCallback
	if err := json.Unmarshal(msgData, &call); err != nil {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %s", string(msgData))
		return nil, errors.New("Cannot parse")
	}

	messages := make([]ReceivedMessage, 0, len(call.Entries))
	for _, page := range call.Entries {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s", page.PageID)
		for _, fbMsg := range page.Messages {
			msg := ReceivedMessage{ID: fbMsg.Message.ID,
				Sender:      Sender{ID: fbMsg.Sender.ID},
				Text:        fbMsg.Message.Text,
				Attachments: fbMsg.Message.Attachments,
				QuickReply:  fbMsg.Message.QuickReply.Payload}
			log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s, message; %v", page.PageID, fbMsg)
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

type receivedParty struct {
	ID string `json:"id"`
}

type facebookMessengerReceivedMessage struct {
	Sender    receivedParty                                          `json:"sender"`
	Recipient receivedParty                                          `json:"recipient"`
	Time      facebookMessengerAPITime                               `json:"timestamp"`
	Message   facebookMessengerWebhookMessageCallbackMessageRecieved `json:"message,omitempty"`
}

/// Webhook models

// facebookMessengerAPITime aliases time.Time to add custom parsing of unix timestamps
type facebookMessengerAPITime struct {
	Time time.Time
}

// facebookMessengerWebhookMessageCallback represents the webhook callback
type facebookMessengerWebhookMessageCallback struct {
	object  string                                         // Always 'page', so not exposed `json:"object"`
	Entries []facebookMessengerWebhookMessageCallbackEntry `json:"entry"`
}

// facebookMessengerWebhookMessageCallbackEntry represents messages delivered for a particular page
type facebookMessengerWebhookMessageCallbackEntry struct {
	PageID   string                             `json:"id"`
	Time     facebookMessengerAPITime           `json:"time"`
	Messages []facebookMessengerReceivedMessage `json:"messaging"`
}

// UnmarshalJSON parses a unix time into APITime
func (t *facebookMessengerAPITime) UnmarshalJSON(data []byte) error {
	u, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(u), 0)
	return nil
}

type facebookMessengerWebhookMessageCallbackMessageRecieved struct {
	ID          string                                                   `json:"mid"`
	Sequence    int                                                      `json:"seq"`
	Text        string                                                   `json:"text"`
	Attachments []MessageAttachment                                      `json:"attachment,omitempty"`
	QuickReply  facebookMessengerWebhookMessageCallbackMessageQuickReply `json:"quick_reply,omitempty"`
}

func (a MessageAttachment) UnmarshalJSON(data []byte) error {
	var strMap map[string]string
	err := json.Unmarshal(data, &strMap)
	if err != nil {
		return err
	}
	payloadType, ok := strMap["type"]
	if !ok {
		return errors.New("Missing 'type'")
	}
	switch payloadType {
	case "image":
		fallthrough
	case "audio":
		fallthrough
	case "video":
		var s struct {
			Type    MessageAttachmentContentType                                  `json:"type"`
			Payload facebookMessengerWebhookMessageCallbackMessageAttachmentMedia `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.Type = s.Type
		a.Payload = s.Payload
	case "location":
		var s struct {
			Type    MessageAttachmentContentType                                     `json:"type"`
			Payload facebookMessengerWebhookMessageCallbackMessageAttachmentLocation `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.Type = s.Type
		a.Payload = s.Payload
	}
	if err != nil {
		return err
	}
	return nil
}

type facebookMessengerWebhookMessageCallbackMessageAttachmentMedia struct {
	URL string `json:"url"`
}

type facebookMessengerWebhookMessageCallbackMessageAttachmentLocation struct {
	Coordinate facebookMessengerWebhookMessageCallbackMessageAttachmentLocationCoordinate `json:"coordinate"`
}

type facebookMessengerWebhookMessageCallbackMessageAttachmentLocationCoordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"long"`
}

type facebookMessengerWebhookMessageCallbackMessageQuickReply struct {
	Payload string `json:"payload"`
}
