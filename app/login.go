package main

import (
	"encoding/gob"
	"github.com/steveoc64/go-cmms/shared"
)

func login() {
	// print("Clicked Login Btn", username.Value, pw.Value, rem.Value, rem.Checked)

	lc := shared.LoginCredentials{
		Username:   username.Value,
		Password:   pw.Value,
		RememberMe: rem.Checked,
	}
	// print(lc, ws)

	// Send login info
	ws.Send("login")
	enc := gob.NewEncoder(ws)
	if err := enc.Encode(lc); err != nil {
		print("failed to send login info", err.Error())
		return
	}

	// Read reply
	lr := &shared.LoginReply{}
	done := false
	for !done {
		dec := gob.NewDecoder(ws)
		if err := dec.Decode(lr); err != nil {
			if err.Error() == "extra data in buffer" {
				print("ignoring extra data in buffer")
			} else {
				print("rx login reply", err.Error())
				return
			}
		} else {
			print("got reply", lr)
			done = true
		}
	}
	// print(lr)
	hideLoginForm()
	createMenu(lr.Menu)
}
