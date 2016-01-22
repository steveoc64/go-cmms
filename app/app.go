package main

import (
	"honnef.co/go/js/dom"
	// "strconv"
)

type Doc struct {
	dom.Document
}

func main() {
	window := dom.GetWindow()
	doc := Doc{window.Document()}
	appBody := doc.GetElementByID("app")

	layout := doc.createLayout()
	nav := doc.createNavBar()
	splash := doc.createSplash()
	lf := doc.createLoginForm()

	// Add the heirachy of parts
	layout.AppendChild(nav)
	layout.AppendChild(splash)
	layout.AppendChild(lf)
	appBody.AppendChild(layout)

	// Make a websocket connection
	ws := websocketInit()
	ws.Send("Hello")

	// All Done !
}
