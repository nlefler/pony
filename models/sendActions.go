package models

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

// senderAction is a non-message related action
type senderAction struct {
	Action    senderActionType `json:"sender_action"`
	Recipient MessageParty     `json:"recipient"`
}

// OutgoingMessage is a placeholder for one of OutgoingTextMessage
type OutgoingMessage interface{}

// OutgoingTextMessage holds a message to send
type OutgoingTextMessage struct {
	Text string `json:"text"`
}
