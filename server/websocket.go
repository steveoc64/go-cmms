package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

var subscribers []*websocket.Conn

func webSocket(c *echo.Context) error {

	ws := c.Socket()
	msg := ""
	for {
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			return c.String(http.StatusOK, "Rx ws")
		}
		// log.Println(ws.Request().RemoteAddr, "Rx:", msg)
		// big routing table here based on the
		// incoming request
		switch msg {
		case "Hello":
			hello(ws)
		case "login":
			login(ws)
		}

	}
}

func hello(ws *websocket.Conn) {
	start := time.Now()
	subscribers = append(subscribers, ws)
	log.Printf("» %s » Hello %s", ws.Request().RemoteAddr, time.Since(start))
	showSubscriberPool("Connetion Pool Grows To:")
}

type socketMsg struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func publishSocket(event string, data interface{}) {
	var myEvent = &socketMsg{
		Event: event,
		Data:  data,
	}
	sendData, err := json.Marshal(myEvent)
	if err != nil {
		log.Println("Error constructing socket data", err.Error())
		return
	}

	gotKills := false
	var newSubs []*websocket.Conn
	fmt.Println("Publish event", event, data)
	for _, wss := range subscribers {
		err := websocket.Message.Send(wss, string(sendData))
		if err != nil {
			//			log.Println("Writing to connection", wss, "got error", err.Error(), "Removing connection from pool")
			// remove this connection from the ppool
			gotKills = true
		} else {
			newSubs = append(newSubs, wss)
		}
	}
	if gotKills {
		subscribers = newSubs
		showSubscriberPool("Pool Shrinks To:")
	}
}

func showSubscriberPool(header string) {

	fmt.Println("==================================")
	fmt.Println(header)
	for i, ws := range subscribers {
		fmt.Printf("  %d:", i+1)
		fmt.Println(ws.Request().RemoteAddr)
	}
	fmt.Println("==================================")
}

func pingSockets() {

	gotKills := false
	var newSubs []*websocket.Conn
	for _, wss := range subscribers {
		err := websocket.Message.Send(wss, ``)
		if err != nil {
			log.Println("Writing to connection", wss, "got error", err.Error(), "Removing connection from pool")
			gotKills = true
		} else {
			newSubs = append(newSubs, wss)
		}
	}
	if gotKills {
		subscribers = newSubs
		showSubscriberPool("Connection Pool Shrinks To:")
	}
}

func pinger() {
	// ticker := time.NewTicker(time.Second * 50) // just under the 1 min mark for nginx default timeouts
	ticker := time.NewTicker(time.Second * 5) // just under the 1 min mark for nginx default timeouts
	for range ticker.C {
		pingSockets()
	}
}
