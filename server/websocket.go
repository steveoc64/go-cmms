package main

import (
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"net/rpc"
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
	rpc.ServeConn(ws)
	return nil
}
