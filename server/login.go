package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"log"
	"time"
)

type LoginRPC struct{}

func (l *LoginRPC) Login(lc *shared.LoginCredentials, lr *shared.LoginReply) error {
	start := time.Now()

	// do some authentication here

	// send a login reply
	lr.Result = "RPC OK"
	lr.Token = "abc123toehunoehnoenuh"
	lr.Menu = []string{"RPC Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}

	// Get the connection we are on
	log.Println("channel is", lc.Channel)
	conn := Connections.Get(lc.Channel)
	log.Println("got conn", conn)
	if conn != nil {
		conn.Login(lc.Username, 34234)
	}
	Connections.Show("after upgrading connection with login")

	log.Printf(`RPC ->
    » login(%s,%s,%t,%d)
    « (%s,%s) %s\n`,
		lc.Username, lc.Password, lc.RememberMe, lc.Channel,
		lr.Result, lr.Token,
		time.Since(start))

	return nil
}
