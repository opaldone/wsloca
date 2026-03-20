package serv

import (
	"encoding/json"
)

const (
	RLOCA    = "rloca"
	ALOCA    = "aloca"
	CLIST    = "clist"
	SENDERHI = "sender_hi"
	SENDERST = "sender_stop"
	GOCHAT   = "go_chat"
)

type Message struct {
	Tp       string `json:"tp"`
	Cid      string `json:"cid"`
	Nik      string `json:"nik"`
	IsSender bool   `json:"issender"`
	Roomid   string `json:"roomid"`
	Content  string `json:"content"`
}

func decMessage(txt []byte) (msg *Message) {
	msg = new(Message)
	json.Unmarshal(txt, msg)
	return
}
