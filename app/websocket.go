package main

import (
	"fmt"
	// "github.com/gopherjs/gopherjs/js"
	"bufio"
	"encoding/gob"
	// "errors"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
	"io"
	"net/rpc"
)

var ws *websocket.Conn
var rpcClient *rpc.Client

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
	wss, err := websocket.Dial(wsBaseURL)
	if err != nil {
		print("failed to open websocket")
	}
	ws = wss

	encBuf := bufio.NewWriter(ws)
	client := &myClientCodec{ws, gob.NewDecoder(ws), gob.NewEncoder(encBuf), encBuf}
	rpcClient = rpc.NewClientWithCodec(client)

	// Now we can spawn a pinger against the backetd
	go sendPings(55000)

	//go PingServer(ws)

	return wss
}

// codec to encode requests and decode responses.''

type myClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

func (c *myClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	print("rpc writes a new header", r)
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *myClientCodec) ReadResponseHeader(r *rpc.Response) error {
	err := c.dec.Decode(r)
	if err != nil {
		print("rpc error", err.Error())
		return err
	}
	if r.ServiceMethod[:1] == "*" {
		print("Async update from server", r.ServiceMethod)
		return nil
		// return errors.New("Async update from server - NOT an RPC call")
	}
	return err
}

func (c *myClientCodec) ReadResponseBody(body interface{}) error {
	print("rpc reads message body")
	err := c.dec.Decode(body)
	print("rpc gets body", body)
	return err
}

func (c *myClientCodec) Close() error {
	print("calling close")
	return c.rwc.Close()
}
