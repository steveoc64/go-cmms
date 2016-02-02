package main

import (
	"fmt"
	// "github.com/gopherjs/gopherjs/js"
	"bufio"
	"encoding/gob"
	//	"errors"
	"github.com/gopherjs/websocket"
	//	"github.com/steveoc64/go-cmms/shared"
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

	// Call PingRPC to burn through the message with seq = 0
	/*	in := &shared.PingReq{
			Msg: "Use up the first msg",
		}
		out := &shared.PingRep{}*/
	//rpcClient.Call("PingRPC.Ping", in, out)

	// Now we can spawn a pinger against the backetd
	//go sendPings(55000)

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
	print("rpc ->", r.ServiceMethod, r.Seq)
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

type MsgPayload struct {
	Msg string
}

func (c *myClientCodec) ReadResponseHeader(r *rpc.Response) error {
	err := c.dec.Decode(r)
	//	print("rpc header <-", r)
	if err != nil {
		print("rpc error", err)
		if err != nil && err.Error() == "extra data in buffer" {
			err = c.dec.Decode(r)
		}
	}
	if r.Seq == 0 {
		print("Async update from server -", r.ServiceMethod)
		return nil
		//return errors.New("Async update from server")
	}
	return err
}

func (c *myClientCodec) ReadResponseBody(body interface{}) error {
	err := c.dec.Decode(body)
	//print("rpc <-", body)
	return err
}

func (c *myClientCodec) Close() error {
	print("calling close")
	return c.rwc.Close()
}
