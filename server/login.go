package main

import (
	"encoding/gob"
	"github.com/steveoc64/go-cmms/shared"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func login(ws *websocket.Conn) {
	start := time.Now()

	// read the credentials
	dec := gob.NewDecoder(ws)
	lc := &shared.LoginCredentials{}
	if err := dec.Decode(lc); err != nil {
		log.Println(err.Error())
		return
	}

	// send a login reply
	lr := shared.LoginReply{
		Result: "OK",
		Token:  "abc123toehunoehnoenuh",
		Menu:   []string{"Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"},
	}
	enc := gob.NewEncoder(ws)
	ws.PayloadType = websocket.BinaryFrame
	if err := enc.Encode(lr); err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf(`» %s %s
    » login(%s,%s,%t)
    « (%s,%s)`,
		ws.Request().RemoteAddr,
		time.Since(start),
		lc.Username, lc.Password, lc.RememberMe,
		lr.Result, lr.Token)
}
