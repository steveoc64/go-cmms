package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"log"
)

func rpc_login() {

	print("called rpc login")
	lc := &shared.LoginCredentials{
		Username:   username.Value,
		Password:   pw.Value,
		RememberMe: rem.Checked,
	}
	lr := &shared.LoginReply{}
	err := rpcClient.Call("LoginRPC.Login", lc, lr)
	print("finished the rpc call")
	if err != nil {
		log.Println(err.Error())
	}
	hideLoginForm()
	createMenu(lr.Menu)
}
