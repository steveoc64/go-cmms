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
	log.Println("» Login")

	s := new(SiteRPC)
	if err := rpc.Register(s); err != nil {
		log.Fatal(err)
	}
	log.Println("» Site")

	up := new(UserProfileRPC)
	if err := rpc.Register(up); err != nil {
		log.Fatal(err)
	}
	log.Println("» UserProfile")
}
