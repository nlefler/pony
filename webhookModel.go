package pony

import "net/http"

// Webhook is a common interface for all webhooks which receive messages
type Webhook interface {
	addRoutes(mux *http.ServeMux)
}
