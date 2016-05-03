package main

import (
	"fmt"
	"time"
	// "github.com/gopherjs/gopherjs/js"
	"bufio"
	"encoding/gob"
	//	"errors"
	"io"
	"net"
	"net/rpc"

	"github.com/gopherjs/websocket"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

var ws net.Conn
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

func websocketInit() net.Conn {
	Session.Channel = 0

	wsBaseURL := getWSBaseURL()
	print("init websocket", wsBaseURL)
	wss, err := websocket.Dial(wsBaseURL)
	if err != nil {
		print("failed to open websocket", wsBaseURL, err.Error())
	}
	ws = wss

	encBuf := bufio.NewWriter(ws)
	client := &myClientCodec{
		rwc:    ws,
		dec:    gob.NewDecoder(ws),
		enc:    gob.NewEncoder(encBuf),
		encBuf: encBuf,
	}
	rpcClient = rpc.NewClientWithCodec(client)

	// Call PingRPC to burn through the message with seq = 0
	out := &shared.AsyncMessage{}
	rpcClient.Call("PingRPC.Ping", "init channel", out)

	return wss
}

// codec to encode requests and decode responses.''
type myClientCodec struct {
	rwc           io.ReadWriteCloser
	dec           *gob.Decoder
	enc           *gob.Encoder
	encBuf        *bufio.Writer
	serviceMethod string
	async         bool
}

func (c *myClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	// print("rpc ->", r.ServiceMethod)
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
	c.async = false
	c.serviceMethod = ""
	err := c.dec.Decode(r)
	// print("rpc header <-", r)
	if err != nil {

		if err != nil && err.Error() == "extra data in buffer" {
			err = c.dec.Decode(r)
		} else {
			print("decode error", err.Error())
		}
		if err != nil {
			print("rpc error", err)
			// force application reload
			go autoReload()
		}
	}
	if err == nil && r.Seq == 0 {
		// print("Async update from server -", r.ServiceMethod)
		c.async = true
		c.serviceMethod = r.ServiceMethod
		return nil
		//return errors.New("Async update from server")
	}
	return err
}

func autoReload() {

	print("Connection has expired !!")
	print("Logging out in ... 3")

	go func() {
		time.Sleep(time.Second)
		print("................ 2")
		time.Sleep(time.Second)
		print("........... 1")
		time.Sleep(time.Second)
		print(" !! BYE !!")
		Logout()
	}()
}

func (c *myClientCodec) ReadResponseBody(body interface{}) error {

	if c.async {
		// Read the response body into a string
		msg := shared.AsyncMessage{}
		// var b string
		// c.dec.Decode(&b)
		err := c.dec.Decode(&msg)
		if err != nil {
			print("decode error", err.Error())
		}
		// print("appear to be async with body of ", body)
		// processAsync(c.serviceMethod, body)
		processAsync(c.serviceMethod, msg)
		return nil
	}

	err := c.dec.Decode(body)
	return err
}

func (c *myClientCodec) Close() error {
	print("calling close")
	return c.rwc.Close()
}

// type AsyncMessage struct {
// 	Action string
// 	Data   interface{}
// }

func processAsync(method string, msg shared.AsyncMessage) {

	// print("processing async with method =", method)
	// print("and body =", body)

	switch method {
	case "Ping":
		Session.Channel = msg.ID
		print("Set channel to", Session.Channel)
	case "PingRPC.Ping":
		print("Keepalive")
	default:
		print("Msg:", method, "Action:", msg.Action, "ID:", msg.ID)
		if Session.Subscribe == method {
			go Session.SFn(&msg)
		}
	}
}

func Subscribe(name string, f func(*shared.AsyncMessage)) int {
	// print("subscribing to ", name)
	Session.Subscribe = name
	Session.SFn = f
	return 0
}
