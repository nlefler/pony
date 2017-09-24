package facebook

type ReceivedMessage struct {
	ID          string
	Sender      Sender
	Text        string
	Attachments []MessageAttachment
	QuickReply  string
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
type MessageAttachment struct {
	Type    MessageAttachmentContentType `json:"type"`
	Payload interface{}                  `json:"payload"`
}

// Message is a Message
type Message struct {
	ID        string    `json:"-"`
	Recipient Recipient `json:"recipient"`
}

// Recipient is the Facebook user who will receive the message
type Recipient struct {
	ID string `json:"id"`
}

type Sender struct {
	ID string `json:"id"`
}

// Action is a sender action. It has no message content.
type Action struct {
	Message
	Action SenderActionType `json:"sender_action"`
}

// ContentMessage is a message with content: text, image, actions, etc
type ContentMessage struct {
	Message
	Text         string            `json:"text"`
	Attachment   MessageAttachment `json:"attachment,omitempty"`
	QuickReplies []QuickReplies    `json:"quick_replies,omitempty"`
}

// SenderActionType is a non-message related action
type SenderActionType string

const (
	// MarkSeen a received message as read
	MarkSeen SenderActionType = "mark_seen"
	// TypingOn shows the typing indicator
	TypingOn SenderActionType = "typing_on"
	// TypingOff disables the typing indicator
	TypingOff SenderActionType = "typing_off"
)

// QuickReplies represents a quick reply included with a sent message
type QuickReplies struct {
	Type    string `json:"content_type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
