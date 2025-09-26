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
)

type Message struct {
	Tp      string `json:"tp"`
	Content string `json:"content"`
}

type SenderType struct {
	Cid      string `json:"cid"`
	Nik      string `json:"nik"`
	IsSender bool   `json:"issender"`
}

func decMessage(txt []byte) (msg *Message) {
	msg = new(Message)
	json.Unmarshal(txt, msg)
	return
}

func decSender(txt string) (st *SenderType) {
	st = new(SenderType)
	json.Unmarshal([]byte(txt), st)
	return
}
