package main

import (
	"net/http"

	"github.com/nlefler/pony/pony"
)

func main() {
	pony := pony.New("")

	mux := http.NewServeMux()
	pony.AddRoutes(mux)

	http.ListenAndServe("0.0.0.0:8080", mux)
}
