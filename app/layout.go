package main

import (
	"honnef.co/go/js/dom"
)

var mainColor = "blue-grey"

// Create a gridded layout
func createLayout() *dom.HTMLDivElement {

	// Create the basic layout
	layout := doc.addDiv()
	layout.SetID("layout")
	layout.SetClass("row")
	return layout
}

func createNavBar() *dom.BasicHTMLElement {

	nav := doc.addElement("nav")
	nav.SetClass(mainColor + " darken-2")

	// Wrapper
	navwrapper := doc.addDiv()

	// Logo
	logo := doc.addDiv()
	logo.SetClass("brand-logo")
	logo.SetTextContent("CMMS")
	navwrapper.AppendChild(logo)

	// Top level options
	ul := doc.addUL()
	ul.SetClass("right hide-on-med-and-down")
	login := doc.addLI("Login", "#")
	ul.AppendChild(login)
	login.AddEventListener("click", false, func(event dom.Event) {
		// print("clicked login btn")
		showLoginForm()
	})

	navwrapper.AppendChild(ul)

	nav.AppendChild(navwrapper)

	return nav
}

// Create a Splash screen
func createSplash() *dom.HTMLDivElement {

	div := doc.addDiv()
	div.SetID("splash")
	img := doc.addIMG("/img/sbs01.jpg")
	banner := doc.addH("h3", "CMMS Facilities Management")
	div.AppendChild(img)
	div.AppendChild(banner)
	return div
}

var username, pw, rem *dom.HTMLInputElement
var udiv, pdiv, remdiv *dom.HTMLDivElement

// Create a login form
func createLoginForm() *dom.HTMLDivElement {

	// basic container and framework for the form
	div := doc.addDiv()
	div.SetClass("container")
	div.SetID("loginform")

	row := doc.addDiv()
	row.SetClass("row")

	col := doc.addDiv()
	col.SetClass("col s6 offset-s3")

	h3 := doc.addH("h3", "Login")
	h3.SetClass("center-align")

	// username
	ruser := doc.addRow()
	udiv, username = doc.addInputField("l-username", "text", "User Name", "s12")
	ruser.AppendChild(udiv)

	// passwd
	rpass := doc.addRow()
	pdiv, pw = doc.addInputField("l-passwd", "password", "PassWord", "s12")
	rpass.AppendChild(pdiv)

	// remember me
	rrem := doc.addRow()
	remdiv, rem = doc.addCheckbox("l-remember", "Remember Me ?", "s6")
	sub := doc.addSubmit("Login", "s6", rpc_login)
	rrem.AppendChild(remdiv)
	rrem.AppendChild(sub)

	// submit

	col.AppendChild(h3)
	col.AppendChild(ruser)
	col.AppendChild(rpass)
	col.AppendChild(rrem)

	row.AppendChild(col)
	div.AppendChild(row)
	div.Style().SetProperty("display", "none", "")

	return div
}

func showLoginForm() {

	// hide splash and show login
	sp := doc.GetElementByID("splash").(*dom.HTMLDivElement)
	sp.Style().SetProperty("display", "none", "")
	lf := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	lf.Style().SetProperty("display", "inline", "")
	lb := doc.GetElementByID("nav-login-btn")
	lb.SetTextContent("Login")

	// It the menu exists, remove it
	sb := doc.GetElementByID("sidebar-menu")
	if sb != nil {
		layout := doc.GetElementByID("layout")
		layout.RemoveChild(sb)
		print("removing menu")
	} else {
		print("no menu to remove")
	}
}

func hideLoginForm() {

	// hide login, and change to logout btn
	lf := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	lf.Style().SetProperty("display", "none", "")
	lb := doc.GetElementByID("nav-login-btn")
	lb.SetTextContent("Logout")
}

func createMenu(menu []string) {

	layout := doc.GetElementByID("layout")
	div := doc.addDiv()
	div.SetID("sidebar-menu")
	div.SetClass(mainColor + " col s1 m2 lighten-1 text-white sidebar")

	ul := doc.addUL()
	div.AppendChild(ul)

	for _, v := range menu {
		li := doc.addLI(v, "#")
		ul.AppendChild(li)
	}

	layout.AppendChild(div)
	print("created menu", div)
}
