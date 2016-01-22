package main

import (
	"fmt"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

func getWSBaseURL() string {
	document := dom.GetWindow().Document().(dom.HTMLDocument)
	location := document.Location()

	wsProtocol := "ws"
	if location.Protocol == "https:" {
		wsProtocol = "wss"
	}
	return fmt.Sprintf("%s://%s:%s/ws", wsProtocol, location.Hostname, location.Port)
}

func websocketInit() *websocket.Conn {
	wsBaseURL := getWSBaseURL()
	ws, err := websocket.Dial(wsBaseURL)
	if err != nil {
		print("failed to open websocket")
	}
	return ws
}
