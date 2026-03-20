package serv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"wsloca/tools"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1048576
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	cid        string
	nik        string
	issender   bool
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	lockClient sync.RWMutex
}

func (c *Client) sendMeListClients() {
	str := c.hub.listSenders()

	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = CLIST
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) setMe(msg *Message) {
	c.lockClient.Lock()
	defer c.lockClient.Unlock()

	if len(msg.Nik) > 0 {
		c.nik = msg.Nik
	}
}

func (c *Client) processMessage(msg *Message) {
	switch msg.Tp {
	case CLIST:
		c.sendMeListClients()
	case SENDERHI:
		c.setMe(msg)
		msg.Cid = c.cid
		c.hub.hiReceivers(msg)
	case RLOCA:
		c.hub.reqLoca(msg)
	case ALOCA:
		msg.Cid = c.cid
		c.hub.ansLoca(msg)
	case GOCHAT:
		c.hub.reqChat(msg)
	}
}

// Gets some message
func (c *Client) readPump() {
	defer func() {
		c.stopClient()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				tools.Danger("readPump", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		msg := decMessage(message)

		c.processMessage(msg)
	}
}

// Message to the outside
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for range n {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) stopClient() {
	c.hub.unregister <- c
	c.conn.Close()
}

// ServeWs handles websocket requests from the peer.
func ServeWs(issenderIn bool,
	hubIn *Hub, win http.ResponseWriter, rin *http.Request,
) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connIn, err := upgrader.Upgrade(win, rin, nil)
	if err != nil {
		tools.Danger("ServeWs", err)
		return
	}

	px := "re-"
	if issenderIn {
		px = "se-"
	}

	cid := fmt.Sprintf("%s%s", px, tools.RandUID())

	nc := &Client{
		cid:      cid,
		issender: issenderIn,
		hub:      hubIn,
		conn:     connIn,
		send:     make(chan []byte, 256),
	}

	nc.hub.register <- nc
	go nc.writePump()
	go nc.readPump()
}
