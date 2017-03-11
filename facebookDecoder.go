package pony

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
)

// FacebookMessengerDecoder translates Facebook Messenger messages
type FacebookMessengerDecoder struct {
}

func (decoder *FacebookMessengerDecoder) receive(msgData []byte) ([]Message, error) {
	var call facebookMessengerWebhookMessageCallback
	if err := json.Unmarshal(msgData, &call); err != nil {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
		return nil, errors.New("Cannot parse")
	}

	messages := make([]Message, len(call.Entries))
	for _, page := range call.Entries {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s", page.PageID)
		for _, fbMsg := range page.Messages {

			messages = append(messages, message)
		}
	}
	return messages, nil
}

// Message models
type facebookMessengerReceivedMessage struct {
	facebookMessengerWebhookMessageCallbackMessageParties
	Time    facebookMessengerAPITime                               `json:"timestamp"`
	Message facebookMessengerWebhookMessageCallbackMessageRecieved `json:"message,omitempty"`
}

type facebookMessengerMessageParty struct {
	FacebookUserID string `json:"id"`
}

// MessageParty
func (p facebookMessengerMessageParty) ID() string {
	return p.FacebookUserID
}

// Webhook models

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
	Sender    facebookMessengerMessageParty `json:"sender"`
	Recipient facebookMessengerMessageParty `json:"recipient"`
}

type facebookMessengerWebhookMessageCallbackMessageRecieved struct {
	ID          string                                                     `json:"mid"`
	Sequence    int                                                        `json:"seq"`
	Text        string                                                     `json:"text"`
	Attachments []facebookMessengerWebhookMessageCallbackMessageAttachment `json:"attachment,omitempty"`
	QuickReply  facebookMessengerWebhookMessageCallbackMessageQuickReply   `json:"quick_reply,omitempty"`
}

type facebookMessengerWebhookMessageCallbackMessageAttachment struct {
	Type    facebookMessengerMessageAttachmentType `json:"type"`
	Payload interface{}                            `json:"payload"`
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
			Type    facebookMessengerfacebookMessengerMessageAttachmentType       `json:"type"`
			Payload facebookMessengerWebhookMessageCallbackMessageAttachmentMedia `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.Type = s.Type
		a.Payload = s.Payload
	case "location":
		var s struct {
			Type    facebookMessengerfacebookMessengerMessageAttachmentType          `json:"type"`
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
	Payload map[string]interface{} `json:"payload"`
}
