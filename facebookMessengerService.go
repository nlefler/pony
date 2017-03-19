package pony

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/// Service

type facebookMessenger struct {
	id      string
	webhook *facebookMessengerWebhook
	sender  *facebookMessengerSender
}

func NewFacebookMessenger(pageName string, validationToken string, pageToken string) Service {
	id := fmt.Sprintf("com.pony.facebook.messenger.%s", pageName)
	webhook := &facebookMessengerWebhook{pageName, validationToken, pageToken, &facebookMessengerDecoder{}, make(chan Message, 100)}
	sender := newFacebookMessengerSender(pageToken)
	return &facebookMessenger{id, webhook, &sender}
}

func (fb *facebookMessenger) Setup(mux *http.ServeMux) {
	fb.webhook.addRoutes(mux)
}

func (fb *facebookMessenger) ID() string {
	return fb.id
}

func (fb *facebookMessenger) Send(msg Message) {

}

func (fb *facebookMessenger) ReceiveOn() <-chan Message {
	return fb.webhook.receiveOn
}

/// Webhook

type facebookMessengerWebhook struct {
	webhookPrefix   string
	validationToken string
	pageToken       string
	decoder         *facebookMessengerDecoder
	receiveOn       chan Message
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
	log.Println("pony.pony.authorize")
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
			log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
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

/// Decoder

type facebookMessengerDecoder struct {
}

func (decoder *facebookMessengerDecoder) receive(msgData []byte) ([]Message, error) {
	var call facebookMessengerWebhookMessageCallback
	if err := json.Unmarshal(msgData, &call); err != nil {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
		return nil, errors.New("Cannot parse")
	}

	messages := make([]Message, 0, len(call.Entries))
	for _, page := range call.Entries {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s", page.PageID)
		for _, fbMsg := range page.Messages {
			log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s, message; %v", page.PageID, fbMsg)
			messages = append(messages, fbMsg)
		}
	}
	return messages, nil
}

/// Message models
// Message
type facebookMessengerReceivedMessage struct {
	facebookMessengerWebhookMessageCallbackMessageParties
	Time    facebookMessengerAPITime                               `json:"timestamp"`
	Message facebookMessengerWebhookMessageCallbackMessageRecieved `json:"message,omitempty"`
}

func (m facebookMessengerReceivedMessage) ID() string {
	return m.Message.ID
}

func (m facebookMessengerReceivedMessage) Sender() MessageParty {
	return m.FBSender
}

func (m facebookMessengerReceivedMessage) Recipients() []MessageParty {
	return []MessageParty{m.FBRecipient}
}

func (m facebookMessengerReceivedMessage) Text() string {
	return m.Message.Text
}

func (m facebookMessengerReceivedMessage) Attachments() []MessageAttachment {
	messages := make([]MessageAttachment, len(m.Message.Attachments))
	for _, a := range m.Message.Attachments {
		messages = append(messages, a)
	}
	return messages
}

type facebookMessengerWebhookMessageCallbackMessageAttachment struct {
	AttachmentType    facebookMessengerMessageAttachmentType `json:"type"`
	AttachmentPayload interface{}                            `json:"payload"`
}

func (a facebookMessengerWebhookMessageCallbackMessageAttachment) Type() MessageAttachmentContentType {
	switch a.AttachmentType {
	case facebookMessengerAttachmentImage:
		return MessageAttachmentContentTypeImage
	case facebookMessengerAttachmentAudio:
		return MessageAttachmentContentTypeAudio
	case facebookMessengerAttachmentVideo:
		return MessageAttachmentContentTypeVideo
	case facebookMessengerAttachmentFile:
		return MessageAttachmentContentTypeFile
	case facebookMessengerAttachmentLocation:
		return MessageAttachmentContentTypeLocation
	default:
		return ""
	}
}

func (a facebookMessengerWebhookMessageCallbackMessageAttachment) Payload() interface{} {
	return a.AttachmentPayload
}

type FacebookMessengerParty interface {
	FacebookMessengerID() string
}

type facebookMessengerMessageParty struct {
	FacebookUserID string `json:"id"`
}

// MessageParty
func (p facebookMessengerMessageParty) FacebookMessengerID() string {
	return p.FacebookUserID
}

/// Webhook models

// messageAttachmentType is the type of a media attachment
type facebookMessengerMessageAttachmentType string

// facebookMessengerAPITime aliases time.Time to add custom parsing of unix timestamps
type facebookMessengerAPITime struct {
	Time time.Time
}

const (
	// facebookMessengerAttachmentImage represents an image attachment
	facebookMessengerAttachmentImage facebookMessengerMessageAttachmentType = "image"
	// facebookMessengerAttachmentAudio represents an audio attachment
	facebookMessengerAttachmentAudio facebookMessengerMessageAttachmentType = "audio"
	// facebookMessengerAttachmentVideo represents an video attachment
	facebookMessengerAttachmentVideo facebookMessengerMessageAttachmentType = "video"
	// facebookMessengerAttachmentFile represents an file attachment
	facebookMessengerAttachmentFile facebookMessengerMessageAttachmentType = "file"
	// facebookMessengerAttachmentLocation represents an location attachment
	facebookMessengerAttachmentLocation facebookMessengerMessageAttachmentType = "location"
)

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

type facebookMessengerWebhookMessageCallbackMessageParties struct {
	FBSender    facebookMessengerMessageParty `json:"sender"`
	FBRecipient facebookMessengerMessageParty `json:"recipient"`
}

type facebookMessengerWebhookMessageCallbackMessageRecieved struct {
	ID          string                                                     `json:"mid"`
	Sequence    int                                                        `json:"seq"`
	Text        string                                                     `json:"text"`
	Attachments []facebookMessengerWebhookMessageCallbackMessageAttachment `json:"attachment,omitempty"`
	QuickReply  facebookMessengerWebhookMessageCallbackMessageQuickReply   `json:"quick_reply,omitempty"`
}

func (a facebookMessengerWebhookMessageCallbackMessageAttachment) UnmarshalJSON(data []byte) error {
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
			Type    facebookMessengerMessageAttachmentType                        `json:"type"`
			Payload facebookMessengerWebhookMessageCallbackMessageAttachmentMedia `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.AttachmentType = s.Type
		a.AttachmentPayload = s.Payload
	case "location":
		var s struct {
			Type    facebookMessengerMessageAttachmentType                           `json:"type"`
			Payload facebookMessengerWebhookMessageCallbackMessageAttachmentLocation `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.AttachmentType = s.Type
		a.AttachmentPayload = s.Payload
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
	Payload map[string]interface{} `json:"payload"`
}

/// Sender

const (
	facebookMessengerSendURLFormat = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

type facebookMessengerSender struct {
	pageToken string
	sendURL   string
}

func newFacebookMessengerSender(pageToken string) facebookMessengerSender {
	return facebookMessengerSender{pageToken, fmt.Sprintf(facebookMessengerSendURLFormat, pageToken)}
}

// Send sends a message
func (s *facebookMessengerSender) send(recipient MessageParty, message Message) {
	// TODO(nl): Recipients
	textMessage := outgoingTextMessage{message.Text()}
	payload := outgoingMessagePayload{message.Recipients()[0], textMessage}
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

type outgoingMessage interface{}

type outgoingMessagePayload struct {
	Recipient MessageParty    `json:"recipient"`
	Message   outgoingMessage `json:"message"`
}

type outgoingTextMessage struct {
	Text string `json:"text"`
}
