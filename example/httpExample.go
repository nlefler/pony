package main

import "github.com/nlefler/pony"

func main() {
	p := pony.NewPony(user{"example"})
	fb := pony.NewFacebookMessenger("glucloser_test", "", "")
	p.AddService(fb)

	msg := message{user{"sender"}, user{"recipient"}, "hello"}
	p.SendMessage(msg, fb.ID())
}

type user struct {
	id string
}

func (u user) ID() string {
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

func (m message) Sender() pony.MessageParty {
	return m.sender
}

func (m message) Recipients() []pony.MessageParty {
	return []pony.MessageParty{m.recipient}
}

func (m message) Text() string {
	return m.text
}

func (m message) Attachments() []pony.MessageAttachment {
	return []pony.MessageAttachment{}
}
