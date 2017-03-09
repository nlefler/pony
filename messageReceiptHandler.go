package pony

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nlefler/pony/models"
)

// ReceiptHandler handles message events
type ReceiptHandler struct {
	Received chan ReceivedMessage
}

// Handle handles a reeived message
func (rh *ReceiptHandler) Handle(w http.ResponseWriter, req *http.Request) {

}

func dispatch(rh *ReceiptHandler, msg ReceivedMessage) {
	if msg.Message.ID != "" {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Message %s", msg.Message.Text)
	}
	rh.Received <- msg
}