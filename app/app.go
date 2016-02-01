package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	window := dom.GetWindow()
	doc = Doc{window.Document()}
	appBody := doc.GetElementByID("app")

	layout := createLayout()
	nav := createNavBar()
	splash := createSplash()
	lf := createLoginForm()

	// Add the heirachy of parts
	layout.AppendChild(nav)
	layout.AppendChild(splash)
	layout.AppendChild(lf)
	appBody.AppendChild(layout)

	// Make a websocket connection
	websocketInit()

	// All Done !
}
