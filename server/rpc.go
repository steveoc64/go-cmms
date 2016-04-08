package main

import (
	"log"
	"net/rpc"
)

func registerRPC() {

	log.Println("Registering RPC services")

	_r := new(LoginRPC)
	if err := rpc.Register(_r); err != nil {
		log.Fatal(err)
	}
	log.Println("» Login")

	_r := new(SiteRPC)
	if err := rpc.Register(_r); err != nil {
		log.Fatal(err)
	}
	log.Println("» Site")

	_r := new(MachineRPC)
	if err := rpc.Register(_r); err != nil {
		log.Fatal(err)
	}
	log.Println("» Site")

	_r := new(UserProfileRPC)
	if err := rpc.Register(_r); err != nil {
		log.Fatal(err)
	}
	log.Println("» UserProfile")
}
