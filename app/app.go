package main

import (
	"honnef.co/go/js/dom"
)

var mainColor = "blue-grey"

func main() {
	// setup some vars
	w := dom.GetWindow()
	doc := w.Document()

	// Make a websocket connection
	websocketInit()

	// hide the splash screen
	splash := doc.GetElementByID("splash").(*dom.HTMLDivElement)
	splash.Style().Set("display", "none")

	// Activate the login form, and get focus on the username
	loginForm := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	loginForm.Style().Set("display", "inline")
	doc.GetElementByID("l-username").(*dom.HTMLInputElement).Focus()

	// loginForm.Class().SetString("visible")
	//	loginForm.Style().SetProperty("visibility", "visible", "")
	loginBtn := doc.GetElementByID("l-loginbtn").(*dom.HTMLButtonElement)
	loginBtn.AddEventListener("click", false, func(evt dom.Event) {
		print("clicked login btn")

		username := doc.GetElementByID("l-username").(*dom.HTMLInputElement).Value
		passwd := doc.GetElementByID("l-passwd").(*dom.HTMLInputElement).Value
		rem := doc.GetElementByID("l-remember").(*dom.HTMLInputElement).Checked

		go login(username, passwd, rem)
	})

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLAnchorElement)
	logoutBtn.AddEventListener("click", false, func(evt dom.Event) {
		print("clicked logout btn")
		showLoginForm()
	})

	// All Done !
}
