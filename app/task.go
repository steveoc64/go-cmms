package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// Add in a pile of constants for lists and things
const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

func taskMaint(context *router.Context) {
	print("TODO - taskmaint")
}

func taskList(context *router.Context) {
	print("TODO - taskList")
}

type MachineSchedListData struct {
	Machine shared.Machine
	Tasks   []shared.SchedTask
}

// Show a list of all Scheduled Maint items for this machine
func machineSchedList(context *router.Context) {

	id, err := strconv.Atoi(context.Params["machine"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machine := shared.Machine{}
		tasks := []shared.SchedTask{}
		rpcClient.Call("MachineRPC.Get", id, &machine)
		rpcClient.Call("TaskRPC.ListMachineSched", id, &tasks)

		BackURL := fmt.Sprintf("/machine/%d", id)

		form := formulate.ListForm{}
		form.New("fa-wrench", "Sched Maint for - "+machine.Name+" - "+*machine.SiteName)

		// Define the layout
		form.Column("Tool / Component", "Component")
		form.Column("Frequency", "ShowFrequency")
		form.Column("Description", "Descr")
		form.Column("$ Labour", "LabourCost")
		form.Column("$ Materials", "MaterialCost")
		form.Column("Duration", "DurationDays")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(fmt.Sprintf("/machine/sched/add/%d", id))
		})

		form.RowEvent(func(key string) {
			print("TODO - edit sched")
		})

		// form.Render("machine-sched-list", "main", tasks)
		// customData := MachineSchedListData{
		// 	Machine: machine,
		// 	Tasks:   tasks,
		// }
		// form.RenderCustom("machine-sched-list-custom", "main", customData)
		form.Render("machine-sched-list", "main", tasks)

		// // Add a handler for clicking on a row
		// doc.GetElementByID("machine-sched-list").AddEventListener("click", false, func(evt dom.Event) {
		// 	td := evt.Target()
		// 	tr := td.ParentElement()
		// 	key := tr.GetAttribute("key")
		// 	// print("TODO - copy selected element into data.EditTask, then display the modal", key)
		// 	Session.Router.Navigate("/sched/" + key)
		// })

		// // Add a handler for clicking on the add button
		// doc.QuerySelector(".data-add-btn").AddEventListener("click", false, func(evt dom.Event) {
		// 	evt.CurrentTarget().(*dom.BasicHTMLElement).Style().Set("display", "none")
		// 	el := doc.QuerySelector("#popup-form")
		// 	el.Class().Add("md-show")
		// 	doc.QuerySelector("#descr").(*dom.HTMLTextAreaElement).Focus()
		// })

		// // Hit escape = close the modal dialog
		// doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
		// 	if evt.(*dom.KeyboardEvent).KeyCode == 27 {
		// 		print("hit escape here")
		// 		evt.PreventDefault()
		// 		doc.QuerySelector("#popup-form").Class().Remove("md-show")
		// 		doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		// 	}
		// })

		// // On a change of Frequency, change which sub-edit is visible
		// doc.QuerySelector("#freq").AddEventListener("change", false, func(evt dom.Event) {
		// 	f := doc.QuerySelector("#freq").(*dom.HTMLSelectElement).SelectedIndex

		// 	// Hide all
		// 	doc.QuerySelector("#freq-0").Class().Remove("task-show")
		// 	doc.QuerySelector("#freq-1").Class().Remove("task-show")
		// 	doc.QuerySelector("#freq-2").Class().Remove("task-show")
		// 	doc.QuerySelector("#freq-3").Class().Remove("task-show")
		// 	doc.QuerySelector("#freq-4").Class().Remove("task-show")

		// 	switch f {
		// 	case 0:
		// 		doc.QuerySelector("#freq-0").Class().Add("task-show")
		// 	case 1:
		// 		doc.QuerySelector("#freq-1").Class().Add("task-show")
		// 	case 2:
		// 		doc.QuerySelector("#freq-2").Class().Add("task-show")
		// 	case 3:
		// 		doc.QuerySelector("#freq-3").Class().Add("task-show")
		// 	case 4:
		// 		doc.QuerySelector("#freq-4").Class().Add("task-show")
		// 	}
		// })

		// // Handler for clicking the close button
		// doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
		// 	doc.QuerySelector("#popup-form").Class().Remove("md-show")
		// 	doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		// })
		// // Hit escape = close the modal dialog
		// doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
		// 	if evt.(*dom.KeyboardEvent).KeyCode == 27 {
		// 		evt.PreventDefault()
		// 		doc.QuerySelector("#popup-form").Class().Remove("md-show")
		// 		doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		// 	}
		// })

		// // Handler for submitting the edit form
		// doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
		// 	evt.PreventDefault()
		// 	doc.QuerySelector("#popup-form").Class().Remove("md-show")
		// 	doc.QuerySelector(".data-add-btn").(*dom.BasicHTMLElement).Style().Set("display", "inline")
		// 	print("add the new item")

		// 	// fill in the ID fields and init prior to binding
		// 	data.EditTask.Task.MachineID = id
		// 	data.EditTask.Task.ToolID = 0
		// 	data.EditTask.Task.Component = ""
		// 	data.EditTask.Task.Week = 0
		// 	data.EditTask.Task.Count = 0
		// 	data.EditTask.Task.Days = 0

		// 	// Parse the form element and get a form.Form object in return.
		// 	f, err := form.Parse(doc.QuerySelector(".grid-form"))
		// 	if err != nil {
		// 		print("form parse error", err.Error())
		// 		// return
		// 	}
		// 	if err := f.Bind(&data.EditTask.Task); err != nil {
		// 		print("form bind error", err.Error())
		// 		// return
		// 	}

		// 	// manually get the textarea
		// 	data.EditTask.Task.Descr = doc.GetElementByID("descr").(*dom.HTMLTextAreaElement).Value

		// 	// manually get the selected freq
		// 	freq := doc.GetElementByID("freq").(*dom.HTMLSelectElement)
		// 	switch freq.SelectedIndex {
		// 	case 0:
		// 		data.EditTask.Task.Freq = "Monthly"
		// 		data.EditTask.Task.StartDate = time.Unix(0, 0)
		// 		data.EditTask.Task.OneOffDate = time.Unix(0, 0)

		// 		for i, v := range doc.QuerySelectorAll("[name=week]") {
		// 			week := v.(*dom.HTMLInputElement)
		// 			if week.Checked {
		// 				data.EditTask.Task.Week = i + 1
		// 				// print("weekly on week =", data.EditTask.Week)
		// 				break
		// 			}
		// 		}
		// 	case 1:
		// 		data.EditTask.Task.Freq = "Yearly"
		// 		data.EditTask.Task.OneOffDate = time.Unix(0, 0)
		// 		// print("yearly and startdate =", data.EditTask.StartDate)
		// 	case 2:
		// 		data.EditTask.Task.Freq = "Every N Days"
		// 		data.EditTask.Task.StartDate = time.Unix(0, 0)
		// 		data.EditTask.Task.OneOffDate = time.Unix(0, 0)
		// 		// print("every N days =", data.EditTask.Days)
		// 	case 3:
		// 		data.EditTask.Task.Freq = "One Off"
		// 		data.EditTask.Task.StartDate = time.Unix(0, 0)
		// 		// print("once off at =", data.EditTask.OneOffDate)
		// 	case 4:
		// 		data.EditTask.Task.Freq = "Job Count"
		// 		data.EditTask.Task.StartDate = time.Unix(0, 0)
		// 		data.EditTask.Task.OneOffDate = time.Unix(0, 0)
		// 		// print("every N jobs =", data.EditTask.Count)
		// 	}

		// 	// manually get the selected component
		// 	comp := doc.GetElementByID("component").(*dom.HTMLSelectElement).SelectedIndex
		// 	if comp == 0 {
		// 		data.EditTask.Task.CompType = "A"
		// 	} else {
		// 		// If the selected item # is <= number of tools, then its a tool
		// 		if comp <= len(data.Machine.Components) {
		// 			data.EditTask.Task.CompType = "T"
		// 			data.EditTask.Task.ToolID = data.Machine.Components[comp-1].ID
		// 		} else {
		// 			// else it is one of the standard items
		// 			data.EditTask.Task.CompType = "C"
		// 			switch comp - len(data.Machine.Components) {
		// 			case 1:
		// 				data.EditTask.Task.Component = "RollBed"
		// 			case 2:
		// 				data.EditTask.Task.Component = "Uncoiler"
		// 			case 3:
		// 				data.EditTask.Task.Component = "Electrical"
		// 			case 4:
		// 				data.EditTask.Task.Component = "Hydraulic"
		// 			case 5:
		// 				data.EditTask.Task.Component = "Lube"
		// 			case 6:
		// 				data.EditTask.Task.Component = "Printer"
		// 			case 7:
		// 				data.EditTask.Task.Component = "Console"
		// 			}
		// 		}

		// 	}

		// 	go func() {
		// 		data.EditTask.Channel = Session.Channel
		// 		print("edit task =", data.EditTask)
		// 		retval := 0
		// 		rpcClient.Call("TaskRPC.Save", &data.EditTask, &retval)
		// 		// refresh and redraw the whole page
		// 		machineSchedList(context)
		// 	}()

		// })
		// //
		// doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
		// 	Session.Router.Navigate(siteMachineEdit)
		// })
	}()
}

func machineSchedAdd(context *router.Context) {
	id, err := strconv.Atoi(context.Params["machine"])
	if err != nil {
		print(err.Error())
		return
	}

	freqs := []formulate.SelectOption{
		{1, "Monthly"},
		{2, "Yearly"},
		{3, "Every N Days"},
		{4, "One Off"},
		{5, "Job Count"},
	}

	weeks := []formulate.SelectOption{
		{1, "1st Week"},
		{2, "2nd Week"},
		{3, "3rd Week"},
		{4, "4th Week"},
	}

	go func() {
		machine := shared.Machine{}
		task := shared.SchedTask{}
		rpcClient.Call("MachineRPC.Get", id, &machine)

		BackURL := fmt.Sprintf("/machine/sched/%d", machine.ID)
		title := fmt.Sprintf("Add Sched Maint Task for - %s - %s", machine.Name, *machine.SiteName)

		form := formulate.EditForm{}
		form.New("fa-wrench", title)

		// create the swapper panels
		swapper := formulate.Swapper{
			Name:     "freq",
			Selected: 1,
		}

		// Add a set of swappable panels for freq options
		swapper.AddPanel("week").AddRow(1).AddRadio(1, "Week of the Month", "Week", weeks, "ID", "Name", 1)
		swapper.AddPanel("year").AddRow(1).AddDate(1, "Day of the Year", "StartDate")
		swapper.AddPanel("days").AddRow(1).AddNumber(1, "Number of Days", "Days", "1")
		swapper.AddPanel("oneoff").AddRow(1).AddDate(1, "One Off Date", "OneOffDate")
		swapper.AddPanel("job").AddRow(1).AddNumber(1, "Job Count", "Count", "1")

		// Layout the fields
		form.Row(2).
			AddRadio(1, "Frequency", "Freq", freqs, "ID", "Name", 1).
			AddSwapper(1, "Frequency Options:", &swapper)

		compGen := []formulate.SelectOption{
			{0, "General Maintenance"},
		}

		compTools := []formulate.SelectOption{}
		for i, comp := range machine.Components {
			newOpt := formulate.SelectOption{i + 1, comp.Name}
			compTools = append(compTools, newOpt)
		}

		compOther := []formulate.SelectOption{
			{100, "RollBed"},
			{101, "Uncoiler"},
			{102, "Electrical"},
			{103, "Hydraulic"},
			{104, "Lube"},
			{105, "Printer"},
			{106, "Console"},
		}

		form.Row(2).
			AddGroupedSelect(1,
				"Component", "Component",
				[]formulate.SelectGroup{
					{"", compGen},
					{"Tools", compTools},
					{"Other Components", compOther},
				},
				0)

		form.Row(1).
			AddTextarea(1, "Task Description", "Descr")

		form.Row(3).
			AddNumber(1, "Labour Cost", "LabourCost", "1").
			AddNumber(1, "Material Cost", "MaterialCost", "1").
			AddNumber(1, "Duration (days)", "DurationDays", "1")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)
			task.MachineID = machine.ID

			// interpret the Component from the grouped options
			comp, _ := strconv.Atoi(task.Component)
			print("comp = ", comp)
			if comp == 0 {
				task.CompType = "A"
				task.Component = compGen[0].Name
			} else if comp < len(machine.Components) {
				task.CompType = "T"
				task.ToolID = machine.Components[comp-1].ID
				task.Component = compTools[comp-1].Name
			} else {
				task.CompType = "C"
				offset := comp - len(machine.Components)
				task.Component = compOther[offset-1].Name
			}

			// convert the selected freq into a meaningful string
			targetFreq, _ := strconv.Atoi(task.Freq)
			for _, f := range freqs {
				if f.ID == targetFreq {
					task.Freq = f.Name
				}
			}

			go func() {
				data := shared.SchedTaskUpdateData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}
				newID := 0
				rpcClient.Call("TaskRPC.Insert", data, &newID)
				print("added task ID", newID)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &task)
		swapper.SelectByName("week")

		// Setup a callback on the freq selector
		w := dom.GetWindow()
		doc := w.Document()

		doc.QuerySelector("[name=radio-Freq]").AddEventListener("click", false, func(evt dom.Event) {
			clickedOn := evt.Target()
			switch clickedOn.TagName() {
			case "INPUT":
				ie := clickedOn.(*dom.HTMLInputElement)
				key, _ := strconv.Atoi(ie.Value)
				switch key {
				case 1:
					swapper.SelectByName("week")
				case 2:
					swapper.SelectByName("year")
				case 3:
					swapper.SelectByName("days")
				case 4:
					swapper.SelectByName("oneoff")
				case 5:
					swapper.SelectByName("job")
				}
			}
		})

	}()

}

type SchedEditData struct {
	Channel int
	Machine shared.Machine
	Task    shared.SchedTask
}

func schedEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := SchedEditData{}
		data.Channel = Session.Channel
		print("getting sched task", id)
		rpcClient.Call("TaskRPC.GetSched", id, &data.Task)
		rpcClient.Call("MachineRPC.Get", data.Task.MachineID, &data.Machine)
		loadTemplate("sched-edit", "main", data)
		// doc.QuerySelector("#focusme").(*dom.HTMLInputElement).Focus()
		print("passed through data", data.Task.Freq, data.Task.StartDate)

		machineSched := fmt.Sprintf("/machine/schedules/%d", data.Machine.ID)

		// Back to the list
		doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate(machineSched)
		})

		// now setup which parts of the form are visible
		switch data.Task.Freq {
		case "Monthly":
			doc.GetElementByID("freq-0").Class().Add("task-show")
		case "Yearly":
			doc.GetElementByID("freq-1").Class().Add("task-show")
		case "Every N Days":
			doc.GetElementByID("freq-2").Class().Add("task-show")
		case "One Off":
			doc.GetElementByID("freq-3").Class().Add("task-show")
		case "Job Count":
			doc.GetElementByID("freq-4").Class().Add("task-show")
		}

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

		// Handler for submitting the edit form
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			print("update the scheduled task")
			Session.Router.Navigate(machineSched)
		})

	}()

}

///
func siteTaskList(context *router.Context) {
	print("TODO - siteTaskList")
}
