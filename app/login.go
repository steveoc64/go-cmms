package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func Login(username string, passwd string, rem bool) {

	Session.Username = ""
	Session.UserRole = ""

	lc := &shared.LoginCredentials{
		Username:   username,
		Password:   passwd,
		RememberMe: rem,
		Channel:    Session.Channel,
	}
	// print("login params", lc)

	lr := &shared.LoginReply{}
	err := rpcClient.Call("LoginRPC.Login", lc, lr)
	if err != nil {
		print("RPC error", err.Error())
	}
	if lr.Result == "OK" {
		// createMenu(lr.Menu)
		loadRoutes(lr.Role, lr.Routes)
		Session.Username = lc.Username
		Session.UserRole = lr.Role
		Session.UserID = lr.ID
		hideLoginForm()
	} else {
		print("login failed")
		dom.GetWindow().Alert("Login Failed")
	}
}

func Logout() {
	print("log out")
	showLoginForm()
	initRouter() // reset all the routes to nothing
	// js.Global.Get("location").Set("hash", "")

	w := dom.GetWindow()
	loc := w.Location()
	print("loc=", loc)

	// loc.pathname =
	js.Global.Get("location").Call("replace", "/")

	// Session.Router.Navigate("/")
}

func hideLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLButtonElement)
	logoutBtn.Style().Set("display", "inline")
	// logoutBtn.AddEventListener("click", false, func(evt dom.Event) {
	// 	evt.PreventDefault()
	// 	Logout()
	// })

	username := fmt.Sprintf("%s - %s", Session.Username, Session.UserRole)
	userBtn := doc.GetElementByID("userbtn").(*dom.HTMLButtonElement)
	userBtn.SetTextContent(username)
	userBtn.Style().Set("display", "inline")
	// userBtn.AddEventListener("click", false, func(evt dom.Event) {
	// 	evt.PreventDefault()
	// 	userProfile()
	// })

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

	// loginBtn.AddEventListener("click", false, func(evt dom.Event) {
	// 	// print("clicked login btn")
	// 	evt.PreventDefault()

	// 	username := doc.GetElementByID("l-username").(*dom.HTMLInputElement).Value
	// 	passwd := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement).Value
	// 	rem := doc.GetElementByID("l-remember").(*dom.HTMLInputElement).Checked

	// 	go Login(username, passwd, rem)
	// })

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLButtonElement)
	logoutBtn.Style().Set("display", "none")
	userBtn := doc.GetElementByID("userbtn").(*dom.HTMLButtonElement)
	userBtn.Style().Set("display", "none")

	// removeMenu()
}

func initLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()

	// Attach events once
	loginBtn := doc.GetElementByID("l-loginbtn").(*dom.HTMLInputElement)
	loginBtn.AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		username := doc.GetElementByID("l-username").(*dom.HTMLInputElement).Value
		passwd := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement).Value
		rem := doc.GetElementByID("l-remember").(*dom.HTMLInputElement).Checked
		go Login(username, passwd, rem)
	})

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLButtonElement)
	logoutBtn.AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		Logout()
	})

	userBtn := doc.GetElementByID("userbtn").(*dom.HTMLButtonElement)
	userBtn.AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		userProfile()
	})

	doc.QuerySelector("#homepage").AddEventListener("click", false, func(evt dom.Event) {
		print("clicked on homepage thing")
		Session.Router.Navigate("/")
	})
}
