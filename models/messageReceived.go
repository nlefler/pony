package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// MessageAttachmentType is the type of a media attachment
type MessageAttachmentType int

// APITime aliases time.Time to add custom parsing of unix timestamps
type APITime struct {
	Time time.Time
}

const (
	// Image represents an image attachment
	Image MessageAttachmentType = iota
	// Audio represents an audio attachment
	Audio MessageAttachmentType = iota
	// Video represents an video attachment
	Video MessageAttachmentType = iota
	// File represents an file attachment
	File MessageAttachmentType = iota
	// Location represents an location attachment
	Location MessageAttachmentType = iota
)

// WebhookMessageCallback represents the webhook callback
type WebhookMessageCallback struct {
	object  string                        // Always 'page', so not exposed `json:"object"`
	Entries []WebhookMessageCallbackEntry `json:"entry"`
}

// WebhookMessageCallbackEntry represents messages delivered for a particular page
type WebhookMessageCallbackEntry struct {
	PageID   string `json:"id"`
	APITime  `json:"time"`
	Messages []WebhookMessageCallbackMessage `json:"messaging"`
}

// UnmarshalJSON parses a unix time into APITime
func (t *APITime) UnmarshalJSON(data []byte) error {
	u, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(u), 0)
	return nil
}

// WebhookMessageCallbackMessage exposes message information on any message type
type WebhookMessageCallbackMessage struct {
	webhookMessageCallbackMessageParties
	Message WebhookMessageCallbackMessageRecieved `json:"message,omitempty"`
}

// WebhookMessageCallbackMessageParty represents a party in a conversation
type WebhookMessageCallbackMessageParty struct {
	ID string `json:"id"`
}

type webhookMessageCallbackMessageParties struct {
	Sender    WebhookMessageCallbackMessageParty `json:"sender"`
	Recipient WebhookMessageCallbackMessageParty `json:"recipient"`
}

type WebhookMessageCallbackMessageRecieved struct {
	ID          string                                    `json:"mid"`
	Sequence    int                                       `json:"seq"`
	Text        string                                    `json:"text"`
	Attachments []WebhookMessageCallbackMessageAttachment `json:"attachment"`
	QuickReply  WebhookMessageCallbackMessageQuickReply   `json:"quick_reply"`
}

type WebhookMessageCallbackMessageAttachment struct {
	Type    MessageAttachmentType `json:"type"`
	Payload interface{}           `json:"payload"`
}

func (a WebhookMessageCallbackMessageAttachment) UnmarshalJSON(data []byte) error {
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
			Type    MessageAttachmentType                        `json:"type"`
			Payload WebhookMessageCallbackMessageAttachmentMedia `json:"payload"`
		}
		err = json.Unmarshal(data, &s)
		a.Type = s.Type
		a.Payload = s.Payload
	case "location":
		var s struct {
			Type    MessageAttachmentType                           `json:"type"`
			Payload WebhookMessageCallbackMessageAttachmentLocation `json:"payload"`
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

type WebhookMessageCallbackMessageAttachmentMedia struct {
	URL string `json:"url"`
}

type WebhookMessageCallbackMessageAttachmentLocation struct {
	Coordinate WebhookMessageCallbackMessageAttachmentLocationCoordinate `json:"coordinate"`
}

type WebhookMessageCallbackMessageAttachmentLocationCoordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"long"`
}

type WebhookMessageCallbackMessageQuickReply struct {
	Payload map[string]interface{} `json:"payload"`
}
