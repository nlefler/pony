package pony

// MessageParty represents a party in a conversation
type MessageParty interface {
	ID() string
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
