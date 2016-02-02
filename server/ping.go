package main

import (
	"github.com/steveoc64/go-cmms/shared"
	//"log"
)

type PingRPC struct{}

func (p *PingRPC) Ping(in *shared.PingReq, out *shared.PingRep) error {
	out.Result = "Got " + in.Msg
	return nil
}
