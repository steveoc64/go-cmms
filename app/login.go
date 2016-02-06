package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"log"
)

func rpc_login() {

	lc := &shared.LoginCredentials{
		Username:   username.Value,
		Password:   pw.Value,
		RememberMe: rem.Checked,
		Channel:    channelID,
	}
	print("login params", lc)
	lr := &shared.LoginReply{}
	err := rpcClient.Call("LoginRPC.Login", lc, lr)
	if err != nil {
		log.Println(err.Error())
	}
	hideLoginForm()
	createMenu(lr.Menu)
}
