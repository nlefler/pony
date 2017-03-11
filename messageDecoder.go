package pony

type MessageDecoder interface {
	decode(msgData []byte) Message
}
