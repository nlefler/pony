package pony

// Pony receives webhook messages and delegates
type Pony struct {
	sender   *Sender
	services []Service
}

func NewPony(sender *Sender) *Pony {
	return &Pony{sender, make([]Service, 1)}
}

// SendMessage sends a message
func (p *Pony) SendMessage(recipient MessageParty, msg OutgoingMessage, serviceID string) {
	p.sender.Send(recipient, msg)
}
