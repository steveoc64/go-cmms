package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func taskMaint(context *router.Context) {
	print("TODO - taskmaint")
}

func taskList(context *router.Context) {
	print("TODO - taskList")
}

type MachineSchedListData struct {
	Machine  shared.Machine
	EditTask shared.SchedTask
	Tasks    []shared.SchedTask
}

// Show a list of all Scheduled Maint items for this machine
func machineSchedList(context *router.Context) {

	id, err := strconv.Atoi(context.Params["machine"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		siteMachineEdit := fmt.Sprintf("/machine/%d", id)

		data := MachineSchedListData{}
		rpcClient.Call("MachineRPC.Get", id, &data.Machine)
		rpcClient.Call("TaskRPC.MachineList", id, &data.Tasks)
		loadTemplate("machine-sched-list", "main", data)

		// Add a handler for clicking on a row
		doc.GetElementByID("machine-sched-list").AddEventListener("click", false, func(evt dom.Event) {
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			print("TODO - copy selected element into data.EditTask, then display the modal", key)
		})

		// Add a handler for clicking on the add button
		doc.QuerySelector(".data-add-btn").AddEventListener("click", false, func(evt dom.Event) {
			evt.CurrentTarget().(*dom.BasicHTMLElement).Style().Set("display", "none")
			print("add new task for this machine")
			el := doc.QuerySelector("#popup-form")
			el.Class().Add("md-show")
			// doc.QuerySelector("#focusme").(*dom.HTMLInputElement).Focus()
		})

		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			print("clicked on close button")
			doc.QuerySelector("#popup-form").Class().Remove("md-show")
			doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		})
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			print("clicked on save button")
			doc.QuerySelector("#popup-form").Class().Remove("md-show")
			doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		})

		doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate(siteMachineEdit)
		})
	}()
}

func siteTasks(context *router.Context) {
	print("TODO - siteTasks")
}
