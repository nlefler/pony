package pony

import "net/http"

type Service interface {
	ID() string
	Send(msg Message)
	ReceiveOn() <-chan Message

	Setup(mux *http.ServeMux)
}
