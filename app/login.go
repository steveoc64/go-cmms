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
	// print("login params", lc)

	lr := &shared.LoginReply{}
	err := rpcClient.Call("LoginRPC.Login", lc, lr)
	if err != nil {
		print("RPC error", err.Error())
	}
	if lr.Result == "OK" {
		hideLoginForm(lc.Username)
		createMenu(lr.Menu)
		loadRoutes(lr.Role, lr.Routes)
	} else {
		print("login failed")
		dom.GetWindow().Alert("Login Failed")
	}
}

func Logout() {
	showLoginForm()
	initRouter() // reset all the routes to nothing
	js.Global.Get("location").Set("hash", "")
	r.Navigate("/")
}

func hideLoginForm(username string) {
	w := dom.GetWindow()
	doc := w.Document()

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLButtonElement)
	logoutBtn.Style().Set("display", "inline")
	logoutBtn.AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		Logout()
	})

	userBtn := doc.GetElementByID("userbtn").(*dom.HTMLButtonElement)
	userBtn.SetTextContent(username)
	userBtn.Style().Set("display", "inline")
	userBtn.AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		userProfile()
	})

	ldiv := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	ldiv.Style().Set("display", "none")
}

func showLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()

	// Destroy whateven is in main
	doc.QuerySelector("main").SetInnerHTML("")

	// Activate the login form, with an outlined loginbtn, and get focus on the username

	loginBtn := doc.GetElementByID("l-loginbtn").(*dom.HTMLInputElement)
	loginBtn.Class().Remove("button-outline")

	ldiv := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	ldiv.Style().Set("display", "block")
	doc.GetElementByID("l-username").(*dom.HTMLInputElement).Focus()

	loginBtn.AddEventListener("click", false, func(evt dom.Event) {
		// print("clicked login btn")
		evt.PreventDefault()

		username := doc.GetElementByID("l-username").(*dom.HTMLInputElement).Value
		passwd := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement).Value
		rem := doc.GetElementByID("l-remember").(*dom.HTMLInputElement).Checked

		go Login(username, passwd, rem)
	})

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLButtonElement)
	logoutBtn.Style().Set("display", "none")
	userBtn := doc.GetElementByID("userbtn").(*dom.HTMLButtonElement)
	userBtn.Style().Set("display", "none")

	// when the password has been entered, then set the button to full
	// pw := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement)
	// pw.AddEventListener("keypress", false, func(evt dom.Event) {
	// 	loginBtn.Class().Remove("button-outline")
	// })
	// pw.AddEventListener("change", false, func(evt dom.Event) {
	// 	loginBtn.Class().Remove("button-outline")
	// })

	removeMenu()
}
