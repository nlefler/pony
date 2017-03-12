package pony

import "net/http"

type Service interface {
	ID() string
	Send(to MessageParty, msg Message)
	ReceiveOn() <-chan *Message

	Setup(mux *http.ServeMux)
}
