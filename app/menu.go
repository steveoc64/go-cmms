package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func createMenu(Menu []shared.UserMenu) {
	w := dom.GetWindow()
	doc := w.Document()

	header := doc.GetElementByID("app-header")
	removeMenu()

	if len(Menu) > 1 {

		hamburger := doc.GetElementByID("hamburger").(*dom.HTMLAnchorElement)
		hamburger.Style().Set("display", "inline")

		// create a new menu
		ul := doc.CreateElement("ul").(*dom.HTMLUListElement)
		ul.SetClass("side-nav fixed")
		ul.SetID("sidebar-menu")

		logo := doc.CreateElement("li").(*dom.HTMLLIElement)
		logo.SetClass("brand-logo")
		// logo.SetTextContent("CMMS")
		img := doc.CreateElement("img").(*dom.HTMLImageElement)
		img.Src = "/img/logo.png"
		logo.AppendChild(img)

		ul.AppendChild(logo)

		for _, v := range Menu {
			li := doc.CreateElement("li").(*dom.HTMLLIElement)
			li.SetClass("bold")
			if v.Icon != "" {
				li.SetInnerHTML(`<i class="material-icons">` + v.Icon + `</i> ` + v.Title)
				li.AddEventListener("click", true, clickMenu)
			} else {
				li.SetTextContent(v.Title)
			}
			li.SetAttribute("href", v.URL)
			ul.AppendChild(li)
		}

		header.AppendChild(ul)
	}

	r.InterceptLinks()
}

func removeMenu() {
	w := dom.GetWindow()
	doc := w.Document()

	header := doc.GetElementByID("app-header")

	// remove the menu, if its already there
	olddiv := doc.GetElementByID("sidebar-menu")
	if olddiv != nil {
		header.RemoveChild(olddiv)
	}

}

func clickMenu(event dom.Event) {
	path := event.CurrentTarget().GetAttribute("href")
	event.PreventDefault()
	go r.Navigate(path)
}
