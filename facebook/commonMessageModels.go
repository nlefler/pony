package facebook

// MessageParty represents a party in a conversation
type MessageParty interface {
	PonyFacebookMessengerID() string
}

// MessageAttachmentContentType represents the content type of a message attachment
type MessageAttachmentContentType string

const (
	// MessageAttachmentContentTypeImage represents an image attachment
	MessageAttachmentContentTypeImage MessageAttachmentContentType = "image"
	// MessageAttachmentContentTypeAudio represents an audio attachment
	MessageAttachmentContentTypeAudio MessageAttachmentContentType = "audio"
	// MessageAttachmentContentTypeVideo represents an video attachment
	MessageAttachmentContentTypeVideo MessageAttachmentContentType = "video"
	// MessageAttachmentContentTypeFile represents an file attachment
	MessageAttachmentContentTypeFile MessageAttachmentContentType = "file"
	// MessageAttachmentContentTypeLocation represents an location attachment
	MessageAttachmentContentTypeLocation MessageAttachmentContentType = "location"
)

// MessageAttachment is extra content in a message
type MessageAttachment interface {
	Type() MessageAttachmentContentType
	Payload() interface{}
}

// Message is a Message
type Message interface {
	ID() string
	Sender() MessageParty
	Recipients() []MessageParty
	Text() string
	Attachments() []MessageAttachment
}

// SenderActionType is a non-message related action
type senderActionType string

const (
	// MarkSeen a received message as read
	MarkSeen senderActionType = "mark_seen"
	// TypingOn shows the typing indicator
	TypingOn senderActionType = "typing_on"
	// TypingOff disables the typing indicator
	TypingOff senderActionType = "typing_off"
)

// SenderAction is a non-message related action
type SenderAction struct {
	Action    senderActionType `json:"sender_action"`
	Recipient MessageParty     `json:"recipient"`
}

type basicTextMessage struct {
	id         string
	sender     MessageParty
	recipients []MessageParty
	text       string
}

// NewBasicTextMessage is a convenience implementation of a Message with no attachments
func NewBasicTextMessage(sender MessageParty, recipients []MessageParty, text string) Message {
	return basicTextMessage{"", sender, recipients, text}
}

func (m basicTextMessage) ID() string {
	return m.id
}

func (m basicTextMessage) Sender() MessageParty {
	return m.sender
}

func (m basicTextMessage) Recipients() []MessageParty {
	return m.recipients
}

func (m basicTextMessage) Text() string {
	return m.text
}

func (m basicTextMessage) Attachments() []MessageAttachment {
	return []MessageAttachment{}
}
