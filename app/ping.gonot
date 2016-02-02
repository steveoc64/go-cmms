package main

import (
	"github.com/gopherjs/websocket"
	"github.com/steveoc64/go-cmms/shared"
	"net/rpc"
	"time"
)

// Constantly Ping the Backend
func sendPings(ms time.Duration) {

	in := &shared.PingReq{
		PingName: "Test",
		Seq:      1,
	}
	out := &shared.PingRep{}

	ticker := time.NewTicker(time.Millisecond * ms)
	for _ = range ticker.C {
		err := rpcClient.Call("PingRPC.Ping", in, out)
		if err != nil {
			print(err.Error())
		} else {
			//print("Ping", t, in, out)
			in.Seq++
		}
	}
}

type PingRPC struct{}

// Run a server at this end to receive Pings
func PingServer(ws *websocket.Conn) {

	print("Starting up the Ping RPC server")
	p := new(PingRPC)
	if err := rpc.Register(p); err != nil {
		print(err.Error())
	} else {
		print("registered OK")
	}
	//ws.PayloadType = websocket.BinaryFrame
	print("serveconn")
	rpc.ServeConn(ws)
	print("done serving the connection")
}

func (p *PingRPC) ClientPing(in *shared.PingReq, out *shared.PingRep) error {

	//log.Println("Ping ", in.PingName, in.Seq)
	out.Seq = in.Seq
	out.Result = "Got " + in.PingName
	return nil
}
