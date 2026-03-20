// Package serv
package serv

import (
	"encoding/json"
	"sort"
	"sync"
)

// Hub with clients
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	lockHub    sync.RWMutex
}

type ClientDebType struct {
	Cid      string
	Nik      string
	Issender bool
}

// NewHub create new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) addClient(cl *Client) {
	h.lockHub.Lock()
	h.clients[cl.cid] = cl
	h.lockHub.Unlock()
}

func (h *Hub) removeClient(cl *Client) {
	h.lockHub.Lock()
	if _, ok := h.clients[cl.cid]; ok {
		delete(h.clients, cl.cid)
		close(cl.send)
	}
	h.lockHub.Unlock()

	if cl.issender {
		h.senderStopped(cl)
	}
}

// Run hub
func (h *Hub) Run() {
	for {
		select {
		case uqcl := <-h.register:
			h.addClient(uqcl)
		case uqcl := <-h.unregister:
			h.removeClient(uqcl)
		}
	}
}

func (h *Hub) listSenders() (res string) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	var lis []Message

	for _, cl := range h.clients {
		lis = append(lis, Message{
			Cid:      cl.cid,
			Nik:      cl.nik,
			IsSender: cl.issender,
		})
	}

	bont, _ := json.Marshal(lis)
	res = string(bont)

	return
}

func (h *Hub) hiReceivers(msg *Message) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	bts, _ := json.Marshal(msg)

	for _, cl := range h.clients {
		if cl.issender {
			continue
		}

		cl.send <- bts
	}
}

func (h *Hub) senderStopped(cl *Client) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	msg := new(Message)
	msg.Tp = SENDERST
	msg.Cid = cl.cid

	bts, _ := json.Marshal(msg)

	for _, cl := range h.clients {
		if cl.issender {
			continue
		}

		cl.send <- bts
	}
}

func (h *Hub) reqLoca(msg *Message) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	clin, ex := h.clients[msg.Cid]

	if !ex {
		return
	}

	bts, _ := json.Marshal(msg)
	clin.send <- bts
}

func (h *Hub) ansLoca(msg *Message) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	bts, _ := json.Marshal(msg)

	for _, cl := range h.clients {
		if cl.issender {
			continue
		}

		cl.send <- bts
	}
}

func (h *Hub) reqChat(msg *Message) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	clin, ex := h.clients[msg.Cid]

	if !ex {
		return
	}

	bts, _ := json.Marshal(msg)
	clin.send <- bts
}

func (h *Hub) GetShowClients() (list []ClientDebType) {
	h.lockHub.RLock()
	defer h.lockHub.RUnlock()

	for _, cl := range h.clients {
		item := ClientDebType{
			Cid:      cl.cid,
			Nik:      cl.nik,
			Issender: cl.issender,
		}

		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Cid < list[j].Cid
	})

	return
}
