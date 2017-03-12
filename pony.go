package pony

// Pony receives webhook messages and delegates
type Pony struct {
	sender   MessageParty
	services map[string]Service
}

func NewPony(sender MessageParty) *Pony {
	return &Pony{sender, make(map[string]Service)}
}

func (p *Pony) AddService(service Service) {
	p.services[service.ID()] = service
}

// SendMessage sends a message
func (p *Pony) SendMessage(msg Message, serviceID string) {
	s, ok := p.services[serviceID]
	if ok {
		s.Send(msg)
	}
}
