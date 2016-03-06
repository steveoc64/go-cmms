package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
	// "strings"
)

func Login(username string, passwd string, rem bool) {

	lc := &shared.LoginCredentials{
		Username:   username,
		Password:   passwd,
		RememberMe: rem,
		Channel:    channelID,
	}
	print("login params", lc)

	lr := &shared.LoginReply{}
	err := rpcClient.Call("LoginRPC.Login", lc, lr)
	if err != nil {
		print("RPC error", err.Error())
	}
	if lr.Result == "OK" {
		hideLoginForm()
		createMenu(lr.Menu)
		loadRoutes(lr.Role, lr.Routes)
	} else {
		print("login failed")
	}
}

func Logout() {
	showLoginForm()
	initRouter() // reset all the routes to nothing
	js.Global.Get("location").Set("hash", "")
	r.Navigate("/")
}

func hideLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLAnchorElement)
	logoutBtn.Style().Set("display", "inline")
	logoutBtn.AddEventListener("click", false, func(evt dom.Event) {
		print("clicked logout btn")
		evt.PreventDefault()
		Logout()
	})

}

func showLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()

	// Activate the login form, and get focus on the username
	loadTemplate("login", nil)
	doc.GetElementByID("l-username").(*dom.HTMLInputElement).Focus()

	loginBtn := doc.GetElementByID("l-loginbtn").(*dom.HTMLInputElement)
	loginBtn.AddEventListener("click", false, func(evt dom.Event) {
		print("clicked login btn")
		evt.PreventDefault()

		username := doc.GetElementByID("l-username").(*dom.HTMLInputElement).Value
		passwd := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement).Value
		rem := doc.GetElementByID("l-remember").(*dom.HTMLInputElement).Checked

		go Login(username, passwd, rem)
	})

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLAnchorElement)
	logoutBtn.Style().Set("display", "none")

	removeMenu()
}
