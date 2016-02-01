package main

import (
	"github.com/steveoc64/go-cmms/shared"
	//"log"
)

type PingRPC struct{}

func (p *PingRPC) Ping(in *shared.PingReq, out *shared.PingRep) error {

	//log.Println("Ping ", in.PingName, in.Seq)
	out.Seq = in.Seq
	out.Result = "Got " + in.PingName
	return nil
}
