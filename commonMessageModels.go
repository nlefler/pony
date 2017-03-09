package pony

// MessageParty represents a party in a conversation
type MessageParty interface {
	Id() string
	FirstName() string
	LastName() string
}
