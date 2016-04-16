package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// Show a list of all parts
func partList(context *router.Context) {

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := []shared.Part{}
		rpcClient.Call("PartRPC.List", Session.Channel, &data)
		loadTemplate("part-list", "main", data)

		// Add a handler for clicking on a row
		doc.GetElementByID("part-list").AddEventListener("click", false, func(evt dom.Event) {
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			Session.Router.Navigate("/part/" + key)
		})

		// Add a handler for clicking on the add butto
		doc.QuerySelector(".data-add-btn").AddEventListener("click", false, func(evt dom.Event) {
			print("add new part")
		})
	}()
}

func partEdit(context *router.Context) {
	print("TODO partEdit")
}

func partAdd(context *router.Context) {
	print("TODO partAdd")
}
