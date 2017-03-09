package pony

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

// OutgoingMessage is a placeholder for one of OutgoingTextMessage
type OutgoingMessage interface{}

// OutgoingMessagePayload wraps a message to send
type OutgoingMessagePayload struct {
	Recipient MessageParty    `json:"recipient"`
	Message   OutgoingMessage `json:"message"`
}

// OutgoingTextMessage holds a message to send
type OutgoingTextMessage struct {
	Text string `json:"text"`
}
