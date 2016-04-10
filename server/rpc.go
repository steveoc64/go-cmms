package main

import (
	"log"
	"net/rpc"
)

func registerRPC() {

	log.Println("Registering RPC services:")

	if err := rpc.Register(new(LoginRPC)); err != nil {
		log.Fatal(err)
	}
	log.Println("» Login")

	if err := rpc.Register(new(SiteRPC)); err != nil {
		log.Fatal(err)
	}
	log.Println("» Site")

	if err := rpc.Register(new(MachineRPC)); err != nil {
		log.Fatal(err)
	}
	log.Println("» Machine")

	if err := rpc.Register(new(UserRPC)); err != nil {
		log.Fatal(err)
	}
	log.Println("» User")
}
