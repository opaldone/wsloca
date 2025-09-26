// Package controllers
package controllers

import (
	"net/http"
	"strconv"

	"wsloca/serv"

	"github.com/julienschmidt/httprouter"
)

var hub *serv.Hub

func init() {
	hub = serv.NewHub()
	go hub.Run()
}

// Ws handler to create client
func Ws(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cid := ps.ByName("cid")
	sender := ps.ByName("sender")

	s, _ := strconv.Atoi(sender)
	issender := s > 0

	serv.ServeWs(cid, issender, hub, w, r)
}
