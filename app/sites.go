package main

import (
	"github.com/go-humble/router"
	"honnef.co/go/js/dom"
)

func siteMap(context *router.Context) {
	print("in the siteMap function")

	t, err := GetTemplate("gridform")
	if err != nil {
		print("Failed to load gridform", err.Error())
		return
	}
	print("got template", t)
	w := dom.GetWindow()
	doc := w.Document()
	el := doc.QuerySelector("main")
	print(w, doc, el)
	if err := t.ExecuteEl(el, nil); err != nil {
		print("Failed to add template to body", err.Error())
	}
	print("seemed to work k")
}

func siteList(context *router.Context) {
	print("in the siteList function")
}
