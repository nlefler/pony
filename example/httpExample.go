package main

import (
	"github.com/nlefler/pony/facebook"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fb := facebook.NewFacebookMessenger("glucloser_test", "pasta", "EAAFuSRpTlbQBAMobErudrFZCeDOqN8f4ZBlpZAvwJESjZCHygbC3kIij8nv7PL0RAN6a4hA87KKDfcc6gZB5TglBMtGNAZAMAmcVkJUFtKRTOXd3LZBXHWKRVTcrlWGkM0zCidYHTWErnwYUp8kguQB8yaIkAr1B8zxk1FPRXwZBSAZDZD")
	fb.Setup(mux)

	go func(ch <-chan facebook.Message) {
		for msg := range ch {
			log.Printf("%v\n", msg)

			s := msg.Sender()
			r := msg.Recipients()[0]
			resp := message{user{r.FacebookMessengerID()}, user{s.FacebookMessengerID()}, msg.Text()}
			log.Printf("Sending from %s to %s", r.FacebookMessengerID(), s.FacebookMessengerID())
			fb.Send(resp)
		}
	}(fb.ReceiveOn())

	//	go func (mux *http.ServeMux) {
	http.ListenAndServe("0.0.0.0:8080", mux)
	//	}(mux)
}

type user struct {
	id string
}

func (u user) FacebookMessengerID() string {
	return u.id
}

type message struct {
	sender    user
	recipient user
	text      string
}

func (m message) ID() string {
	return "x"
}

func (m message) Sender() facebook.MessageParty {
	return m.sender
}

func (m message) Recipients() []facebook.MessageParty {
	return []facebook.MessageParty{m.recipient}
}

func (m message) Text() string {
	return m.text
}

func (m message) Attachments() []facebook.MessageAttachment {
	return []facebook.MessageAttachment{}
}
