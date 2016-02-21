package main

import (
	"strings"

	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func login(username string, passwd string, rem bool) {

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
		w := dom.GetWindow()
		doc := w.Document()

		// Fill in the details on the nav bar
		uuname := strings.ToUpper(username[0:1]) + username[1:]
		uname := doc.GetElementByID("d-username").(*dom.HTMLLIElement)
		uname.SetTextContent(uuname)
		urole := doc.GetElementByID("d-role").(*dom.HTMLLIElement)
		urole.SetTextContent(lr.Role)
		usite := doc.GetElementByID("d-site").(*dom.HTMLLIElement)
		usite.SetTextContent(lr.Site)
	} else {
		print("login failed")
	}
}

func hideLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()
	loginForm := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	// loginForm.Class().SetString("hidden")
	loginForm.Style().Set("display", "none")

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLAnchorElement)
	// logoutBtn.Class().SetString("visible")
	logoutBtn.Style().Set("display", "inline")
}

func showLoginForm() {
	w := dom.GetWindow()
	doc := w.Document()
	loginForm := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	// loginForm.Class().SetString("visible")
	loginForm.Style().Set("display", "inline")

	logoutBtn := doc.GetElementByID("logoutbtn").(*dom.HTMLAnchorElement)
	// logoutBtn.Class().SetString("hidden")
	logoutBtn.Style().Set("display", "none")

	removeMenu()

	uname := doc.GetElementByID("d-username").(*dom.HTMLLIElement)
	uname.SetTextContent("")
	urole := doc.GetElementByID("d-role").(*dom.HTMLLIElement)
	urole.SetTextContent("")
	usite := doc.GetElementByID("d-site").(*dom.HTMLLIElement)
	usite.SetTextContent("")
}
