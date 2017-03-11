package pony

// MessageParty represents a party in a conversation
type MessageParty interface {
	ID() string
}

// MessageContentType represents the content type of a message
type MessageContentType string

const (
	// MessageContentTypeText represents a text attachment
	MessageContentTypeText MessageContentType = "text"
	// MessageContentTypeImage represents an image attachment
	MessageContentTypeImage MessageContentType = "image"
	// MessageContentTypeAudio represents an audio attachment
	MessageContentTypeAudio MessageContentType = "audio"
	// MessageContentTypeVideo represents an video attachment
	MessageContentTypeVideo MessageContentType = "video"
	// MessageContentTypeFile represents an file attachment
	MessageContentTypeFile MessageContentType = "file"
	// MessageContentTypeLocation represents an location attachment
	MessageContentTypeLocation MessageContentType = "location"
)

// Message is a Message
type Message interface {
	Id() string
	Sender() MessageParty
	Recipients() []MessageParty
	Content() interface{}
	ContentType() MessageContentType
}
