package main

import (
	"log"
	"net/rpc"
)

func registerRPC() {

	log.Println("Registering RPC services")

	l := new(LoginRPC)
	if err := rpc.Register(l); err != nil {
		log.Fatal(err)
	}

	p := new(PingRPC)
	if err := rpc.Register(p); err != nil {
		log.Fatal(err)
	}

	// r := new(RouteRPC)
	// if err := rpc.Register(r); err != nil {
	// 	log.Fatal(err)
	// }
}
