package main

import (
	"encoding/gob"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"net/rpc"
	"time"
)

type ConnectionsList struct {
	sockets []*websocket.Conn
}

func (c *ConnectionsList) Add(ws *websocket.Conn) *ConnectionsList {
	c.sockets = append(c.sockets, ws)
	return c
}

func (c *ConnectionsList) Drop(ws *websocket.Conn) *ConnectionsList {
	c.sockets = append(c.sockets, ws)
	return c
}

func (c *ConnectionsList) Show(header string) *ConnectionsList {
	fmt.Println("==================================")
	fmt.Println(header)
	for i, ws := range c.sockets {
		fmt.Printf("  %d:", i+1)
		fmt.Println(ws.Request().RemoteAddr)
	}
	fmt.Println("==================================")
	return c
}

var Connections *ConnectionsList

func webSocket(c *echo.Context) error {

	ws := c.Socket()
	ws.PayloadType = websocket.BinaryFrame

	Connections.Add(ws).Show("Connections Grows To:")
	go sendPings(ws, 50000)
	rpc.ServeConn(ws)
	return nil
}

type MsgPayload struct {
	Msg string
}

// Constantly Ping the Backend
func sendPings(ws *websocket.Conn, ms time.Duration) {

	ticker := time.NewTicker(time.Millisecond * ms)
	enc := gob.NewEncoder(ws)
	r := rpc.Response{}
	for _ = range ticker.C {
		r.ServiceMethod = "Ping"
		log.Println("sending ping to client", r)
		if err := enc.Encode(&r); err != nil {
			log.Println("Some sort of error sending Ping header", err.Error())
			return
		}
		payload := &MsgPayload{}
		enc.Encode(payload)
	}
}
