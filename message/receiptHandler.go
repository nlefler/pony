package message

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nlefler/pony/models"
)

// ReceiptHandler handles message events
type ReceiptHandler struct {
	Received chan models.ReceivedMessage
}

// Handle handles a reeived message
func (rh *ReceiptHandler) Handle(w http.ResponseWriter, req *http.Request) {
	jsonBytes, err := ioutil.ReadAll(req.Body)
	if len(jsonBytes) == 0 || err != nil {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	var call models.WebhookMessageCallback
	if err := json.Unmarshal(jsonBytes, &call); err != nil {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Can't parse request %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, page := range call.Entries {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Handling page %s", page.PageID)
		for _, msg := range page.Messages {
			go dispatch(rh, msg)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func dispatch(rh *ReceiptHandler, msg models.ReceivedMessage) {
	if msg.Message.ID != "" {
		log.Printf("message.receiptHandler.ReceiptHandler.ServeHTTP: Message %s", msg.Message.Text)
	}
	rh.Received <- msg
}
