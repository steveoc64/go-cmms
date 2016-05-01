package main

import (
	"github.com/steveoc64/go-cmms/shared"
	//"log"
)

type PingRPC struct{}

func (p *PingRPC) Ping(msg string, out *shared.AsyncMessage) error {
	out.Action = "Ping"
	out.ID = 0
	return nil
}
