package main

import (
	"encoding/gob"
	"github.com/steveoc64/go-cmms/shared"
	"golang.org/x/net/websocket"
	"log"
	"net/rpc"
	"time"
)

type LoginRPC struct{}

func loginServer(ws *websocket.Conn) {

	log.Println("Starting up the login RPC server on", ws.Request().RemoteAddr)
	l := new(LoginRPC)
	if err := rpc.Register(l); err != nil {
		log.Fatal(err)
	}
	ws.PayloadType = websocket.BinaryFrame
	rpc.ServeConn(ws)
	log.Println("done serving the connection")
}

func (l *LoginRPC) Login(lc *shared.LoginCredentials, lr *shared.LoginReply) error {
	start := time.Now()

	// do some authentication here

	// send a login reply
	lr.Result = "RPC OK"
	lr.Token = "abc123toehunoehnoenuh"
	lr.Menu = []string{"RPC Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}

	log.Printf(`RPC ->
    » login(%s,%s,%t)
    « (%s,%s) %s\n`,
		lc.Username, lc.Password, lc.RememberMe,
		lr.Result, lr.Token,
		time.Since(start))

	return nil
}

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
