package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/form"
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// "github.com/steveoc64/go-cmms/shared"
// "honnef.co/go/js/dom"

func machineList(context *router.Context) {
	print("TODO machineList")
}

type MachineEditData struct {
	Machine shared.Machine
}

func machineEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := MachineEditData{}
		rpcClient.Call("MachineRPC.Get", id, &data.Machine)
		loadTemplate("machine-edit", "main", data)

		siteMachineList := fmt.Sprintf("/site/machine/%d", data.Machine.SiteId)

		// Add handlers for this form
		doc.QuerySelector("legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate(siteMachineList)
		})
		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(siteMachineList)
		})
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			// Parse the form element and get a form.Form object in return.
			f, err := form.Parse(doc.QuerySelector(".grid-form"))
			if err != nil {
				print("form parse error", err.Error())
				return
			}
			if err := f.Bind(&data.Machine); err != nil {
				print("form bind error", err.Error())
				return
			}
			// manually get the textarea
			data.Machine.Notes = doc.GetElementByID("notes").(*dom.HTMLTextAreaElement).Value

			updateData := &shared.MachineUpdateData{
				Channel: Session.Channel,
				Machine: &data.Machine,
			}
			go func() {
				retval := 0
				print("calling MachineRPC.Save")
				rpcClient.Call("MachineRPC.Save", updateData, &retval)
			}()
			Session.Router.Navigate(siteMachineList)
		})

		// Add an Action Grid
		loadTemplate("machine-actions", "#action-grid", id)
		for _, ai := range doc.QuerySelectorAll(".action__item") {
			url := ai.(*dom.HTMLDivElement).GetAttribute("url")
			if url != "" {
				ai.AddEventListener("click", false, func(evt dom.Event) {
					url := evt.CurrentTarget().GetAttribute("url")
					Session.Router.Navigate(url)
				})
			}
		}

	}()

}
