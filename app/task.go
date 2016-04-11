package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/form"
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
	EditTask shared.SchedTaskEditData
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
			el := doc.QuerySelector("#popup-form")
			el.Class().Add("md-show")
			doc.QuerySelector("#descr").(*dom.HTMLTextAreaElement).Focus()
		})

		// Hit escape = close the modal dialog
		doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				print("hit escape here")
				evt.PreventDefault()
				doc.QuerySelector("#popup-form").Class().Remove("md-show")
				doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
			}
		})

		// On a change of Frequency, change which sub-edit is visible
		doc.QuerySelector("#freq").AddEventListener("change", false, func(evt dom.Event) {
			f := doc.QuerySelector("#freq").(*dom.HTMLSelectElement).SelectedIndex

			// Hide all
			doc.QuerySelector("#freq-0").Class().Remove("task-show")
			doc.QuerySelector("#freq-1").Class().Remove("task-show")
			doc.QuerySelector("#freq-2").Class().Remove("task-show")
			doc.QuerySelector("#freq-3").Class().Remove("task-show")
			doc.QuerySelector("#freq-4").Class().Remove("task-show")

			switch f {
			case 0:
				doc.QuerySelector("#freq-0").Class().Add("task-show")
			case 1:
				doc.QuerySelector("#freq-1").Class().Add("task-show")
			case 2:
				doc.QuerySelector("#freq-2").Class().Add("task-show")
			case 3:
				doc.QuerySelector("#freq-3").Class().Add("task-show")
			case 4:
				doc.QuerySelector("#freq-4").Class().Add("task-show")
			}
		})

		// Handler for clicking the close button
		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			doc.QuerySelector("#popup-form").Class().Remove("md-show")
			doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		})
		// Hit escape = close the modal dialog
		doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				evt.PreventDefault()
				doc.QuerySelector("#popup-form").Class().Remove("md-show")
				doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
			}
		})

		// Handler for submitting the edit form
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			doc.QuerySelector("#popup-form").Class().Remove("md-show")
			doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
			print("add the new item")

			// Parse the form element and get a form.Form object in return.
			f, err := form.Parse(doc.QuerySelector(".grid-form"))
			if err != nil {
				print("form parse error", err.Error())
				// return
			}
			if err := f.Bind(&data.EditTask.Task); err != nil {
				print("form bind error", err.Error())
				// return
			}

			// fill in the ID fields
			data.EditTask.Task.MachineID = id
			data.EditTask.Task.ToolID = 0
			data.EditTask.Task.Component = ""

			// manually get the textarea
			data.EditTask.Task.Descr = doc.GetElementByID("descr").(*dom.HTMLTextAreaElement).Value

			// manually get the selected freq
			freq := doc.GetElementByID("freq").(*dom.HTMLSelectElement)
			switch freq.SelectedIndex {
			case 0:
				data.EditTask.Task.Freq = "Weekly"
				for i, v := range doc.QuerySelectorAll("[name=week]") {
					week := v.(*dom.HTMLInputElement)
					if week.Checked {
						data.EditTask.Task.Week = i + 1
						// print("weekly on week =", data.EditTask.Week)
						break
					}
				}
			case 1:
				data.EditTask.Task.Freq = "Yearly"
				// print("yearly and startdate =", data.EditTask.StartDate)
			case 2:
				data.EditTask.Task.Freq = "Every N Days"
				// print("every N days =", data.EditTask.Days)
			case 3:
				data.EditTask.Task.Freq = "One Off"
				// print("once off at =", data.EditTask.OneOffDate)
			case 4:
				data.EditTask.Task.Freq = "Job Count"
				// print("every N jobs =", data.EditTask.Count)
			}

			// manually get the selected component
			comp := doc.GetElementByID("component").(*dom.HTMLSelectElement)
			print("comp = ", comp)

			go func() {
				data.EditTask.Channel = Session.Channel
				print("edit task =", data.EditTask)
				// retval := 0
				// rpcClient.Call("TaskRPC.Save", &data.EditTask, &retval)
			}()

		})
		//
		doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate(siteMachineEdit)
		})
	}()
}

///
func siteTasks(context *router.Context) {
	print("TODO - siteTasks")
}
