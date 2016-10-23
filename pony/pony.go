package pony

import (
	"net/http"

	"github.com/nlefler/pony/message"
)

// Pony receives webhook messages and delegates
type Pony struct {
	validationToken string
}

// New constructs a new Pony
func New(validationToken string) *Pony {
	return &Pony{validationToken}
}

// AddRoutes adds webhook routes to the provided ServeMux
func (p *Pony) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/webhook", makeHandler(p, webhookDispatcher))
	mux.HandleFunc("/authorize", makeHandler(p, authorizeHandler))
	mux.HandleFunc("/up", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("UP")) })
}

func makeHandler(p *Pony, handler func(*Pony, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(p, w, req)
	}
}

func webhookDispatcher(p *Pony, w http.ResponseWriter, req *http.Request) {
	// TODO(nl): verify signature
	switch req.Method {
	case "GET":
		webhookValidate(p, w, req)
	case "POST":
		handler := message.ReceiptHandler{}
		handler.ServeHTTP(w, req)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func authorizeHandler(p *Pony, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func webhookValidate(p *Pony, w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("hub.mode")
	if mode != "subscribe" {
		w.WriteHeader(http.StatusOK)
		return
	}
	token := req.FormValue("hub.verify_token")
	if token != p.validationToken {
		w.WriteHeader(http.StatusOK)
		return
	}

	challenge := req.FormValue("hub.challenge")
	if challenge == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
	return
}
