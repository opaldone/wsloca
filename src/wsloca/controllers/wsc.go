// Package controllers
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"wsloca/serv"
	"wsloca/tools"

	"github.com/julienschmidt/httprouter"
)

var hub *serv.Hub

func init() {
	hub = serv.NewHub()
	go hub.Run()
}

// Ws handler to create client
func Ws(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sender := ps.ByName("sender")

	s, _ := strconv.Atoi(sender)
	issender := s > 0

	serv.ServeWs(issender, hub, w, r)
}

func Di(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uq := ps.ByName("uq")

	if uq != "shpa" && uq != tools.GetKeyCSRF() {
		fmt.Printf("\ncsrf\t\t%s\n",
			tools.GetKeyCSRF(),
		)

		return
	}

	deb := hub.GetShowClients()

	GenerateHTMLEmp(w, r, deb, "stru/dix")
}
