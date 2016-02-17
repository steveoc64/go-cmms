package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
	"log"
	"strings"
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
		log.Println(err.Error())
	}
	if lr.Result == "OK" {
		hideLoginForm()
		createMenu(lr.Menu)
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

		// Go to the home route
		print("jump to route", lr.Home)
		rr := &shared.RouteReq{
			Channel: channelID,
			Name:    lr.Home,
		}
		rres := &shared.RouteResponse{}
		err := rpcClient.Call("RouteRPC.Get", rr, rres)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("Got template", rres.Template)
		createContent(rres.Template)

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

	sidebar := doc.GetElementByID("sidebar-menu")
	if sidebar != nil {
		// sidebar.Style().Set("display", "none")
		print("removing sidebar", sidebar)
		sidebar.ParentNode().RemoveChild(sidebar)
	}

	uname := doc.GetElementByID("d-username").(*dom.HTMLLIElement)
	uname.SetTextContent("")
	urole := doc.GetElementByID("d-role").(*dom.HTMLLIElement)
	urole.SetTextContent("")
	usite := doc.GetElementByID("d-site").(*dom.HTMLLIElement)
	usite.SetTextContent("")
}

func createMenu(menu []string) {
	w := dom.GetWindow()
	doc := w.Document()

	if len(menu) < 1 {
		return
	}

	layout := doc.GetElementByID("layout")
	div := doc.CreateElement("div").(*dom.HTMLDivElement)
	div.SetID("sidebar-menu")
	div.SetClass(mainColor + " col s1 m2 lighten-1 text-white sidebar")

	ul := doc.CreateElement("ul").(*dom.HTMLUListElement)
	div.AppendChild(ul)

	for _, v := range menu {
		li := doc.CreateElement("li").(*dom.HTMLLIElement)
		a := doc.CreateElement("a").(*dom.HTMLAnchorElement)
		a.URLUtils.Href = "#"
		a.SetTextContent(v)
		li.AppendChild(a)
		ul.AppendChild(li)
	}

	layout.AppendChild(div)
	print("created menu", div)
}

func createContent(template string) {
	w := dom.GetWindow()
	doc := w.Document()

	layout := doc.GetElementByID("layout")
	oldcontent := doc.GetElementByID("content")
	if oldcontent != nil {
		layout.RemoveChild(oldcontent)
	}
	div := doc.CreateElement("div").(*dom.HTMLDivElement)
	div.SetID("content")
	div.SetClass("col s11 m10")
	div.SetInnerHTML(template)

	layout.AppendChild(div)
	print("created content", div)
}
