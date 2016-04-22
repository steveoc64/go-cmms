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

func taskEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		task := shared.Task{}

		rpcClient.Call("TaskRPC.Get", id, &task)

		BackURL := "/tasks"
		title := fmt.Sprintf("Task Details - %d", id)
		form := formulate.EditForm{}
		form.New("fa-server", title)

		task.DisplayStartDate = task.StartDate.Format("Mon, Jan 2 2006")
		task.DisplayDueDate = task.DueDate.Format("Mon, Jan 2 2006")
		if task.Username == nil {
			task.DisplayUsername = "Unassigned"
		} else {
			task.DisplayUsername = *task.Username
		}

		// Layout the fields
		form.Row(3).
			AddDisplay(1, "User", "DisplayUsername").
			AddDisplay(1, "Start Date", "DisplayStartDate").
			AddDisplay(1, "Due Date", "DisplayDueDate")

		form.Row(3).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName").
			AddDisplay(1, "Component", "Component")

		form.Row(1).
			AddTextarea(1, "Description", "Descr")

		form.Row(2).
			AddDisplay(1, "Labour Est $", "LabourEst").
			AddDisplay(1, "Material Est $", "MaterialEst")

		form.Row(2).
			AddInput(1, "Actual Labour $", "LabourCost").
			AddInput(1, "Actual Material $", "MaterialCost")

		form.Row(1).
			AddTextarea(1, "Notes", "Log")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			task.ID = id
			go func() {
				data := shared.TaskUpdateData{
					Channel: Session.Channel,
					Task:    &task,
				}
				done := false
				rpcClient.Call("TaskRPC.Delete", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)
			data := shared.TaskUpdateData{
				Channel: Session.Channel,
				Task:    &task,
			}
			go func() {
				done := false
				rpcClient.Call("TaskRPC.Update", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &task)

		// And attach actions
		switch Session.UserRole {
		case "Admin":
			form.ActionGrid("task-admin-actions", "#action-grid", task.ID, func(url string) {
				// Session.Router.Navigate(url)
				if task.SchedID != 0 {
					c := &router.Context{
						Path:        url,
						InitialLoad: false,
						Params: map[string]string{
							"id":   fmt.Sprintf("%d", task.SchedID),
							"back": fmt.Sprintf("/task/%d", task.ID),
						},
					}
					schedEdit(c)
				}
			})
		default:
			form.ActionGrid("task-actions", "#action-grid", task.ID, func(url string) {
				Session.Router.Navigate(url)
			})
		}

	}()

}

// Show a list of all tasks
func taskList(context *router.Context) {

	go func() {
		tasks := []shared.Task{}
		rpcClient.Call("TaskRPC.List", Session.Channel, &tasks)

		form := formulate.ListForm{}
		form.New("fa-server", "Task List - All Active Tasks")

		// Define the layout
		switch Session.UserRole {
		case "Admin", "Site Manager":
			form.Column("User", "Username")
			form.Column("TaskID", "ID")
		}
		form.Column("Date", "GetStartDate")
		// form.Column("Due", "GetDueDate")
		form.Column("Site", "SiteName")
		form.Column("Machine", "MachineName")
		form.Column("Component", "Component")
		form.Column("Description", "Descr")
		form.Column("Duration", "DurationDays")
		form.Column("Completed", "CompletedDate")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/")
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/task/" + key)
		})

		form.Render("task-list", "main", tasks)

	}()
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
			Session.Router.Navigate("/sched/" + key)
		})

		form.Render("machine-sched-list", "main", tasks)
	}()
}

func schedEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
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

	weekdays := []formulate.SelectOption{
		{1, "Mon"},
		{2, "Tue"},
		{3, "Wed"},
		{4, "Thur"},
		{5, "Fri"},
	}

	go func() {
		machine := shared.Machine{}
		task := shared.SchedTask{}
		technicians := []shared.User{}
		rpcClient.Call("TaskRPC.GetSched", id, &task)
		rpcClient.Call("MachineRPC.Get", task.MachineID, &machine)
		rpcClient.Call("UserRPC.GetTechnicians", machine.SiteID, &technicians)
		type Context struct {
			// Params is the parameters from the url as a map of names to values.
			Params map[string]string
			// Path is the path that triggered this particular route. If the hash
			// fallback is being used, the value of path does not include the '#'
			// symbol.
			Path string
			// InitialLoad is true iff this route was triggered during the initial
			// page load. I.e. it is true if this is the first path that the browser
			// was visiting when the javascript finished loading.
			InitialLoad bool
		}

		BackURL := context.Params["back"]
		if BackURL == "" {
			BackURL = fmt.Sprintf("/machine/sched/%d", machine.ID)
		}
		plainTitle := fmt.Sprintf("Sched Maint Task for - %s - %s", machine.Name, *machine.SiteName)
		title := plainTitle
		if task.Paused {
			title += " (PAUSED)"
		}

		form := formulate.EditForm{}
		form.New("fa-wrench", title)

		// create the swapper panels
		swapper := formulate.Swapper{
			Name:     "freq",
			Selected: 1,
		}

		// Add a set of swappable panels for freq options
		theWeek := 1
		if task.Week != nil {
			theWeek = *task.Week
		}
		theWeekDay := 1
		if task.WeekDay != nil {
			theWeekDay = *task.WeekDay
		}
		swapper.AddPanel("week").AddRow(2).
			AddRadio(1, "Week of the Month", "Week", weeks, "ID", "Name", theWeek).
			AddRadio(1, "Weekday", "WeekDay", weekdays, "ID", "Name", theWeekDay)

		swapper.AddPanel("year").AddRow(1).AddDate(1, "Day of the Year", "StartDate")
		swapper.AddPanel("days").AddRow(1).AddNumber(1, "Number of Days", "Days", "1")
		swapper.AddPanel("oneoff").AddRow(1).AddDate(1, "One Off Date", "OneOffDate")
		swapper.AddPanel("job").AddRow(1).AddNumber(1, "Job Count", "Count", "1")

		// Layout the fields
		currentFreq := 0
		for _, f := range freqs {
			if f.Name == task.Freq {
				currentFreq = f.ID
				break
			}
		}
		form.Row(3).
			AddRadio(1, "Frequency", "Freq", freqs, "ID", "Name", currentFreq).
			AddSwapper(2, "Frequency Options:", &swapper)

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

		currentComp := 0
		switch task.CompType {
		case "T":
			for i, tool := range machine.Components {
				if tool.ID == task.ToolID {
					currentComp = i + 1
					break
				}
			}
		case "C":
			for i, c := range compOther {
				if c.Name == task.Component {
					currentComp = i + 100
				}
			}
		}
		form.Row(2).
			AddGroupedSelect(1,
				"Component", "Component",
				[]formulate.SelectGroup{
					{"", compGen},
					{"Tools", compTools},
					{"Other Components", compOther},
				},
				currentComp).
			AddSelect(1, "Assign To Technician", "UserID", technicians, "ID", "Username", 0, task.UserID)

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

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				data := shared.SchedTaskUpdateData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}
				done := false
				rpcClient.Call("TaskRPC.DeleteSched", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)
			task.MachineID = machine.ID

			// interpret the Component from the grouped options
			comp, _ := strconv.Atoi(task.Component)
			// print("comp = ", comp)
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
					break
				}
			}
			switch targetFreq {
			case 1:
				task.Days = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 2:
				task.Days = nil
				task.Count = nil
				task.Week = nil
				task.OneOffDate = nil
			case 3:
				task.Week = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 4:
				task.Days = nil
				task.Count = nil
				task.StartDate = nil
				task.Week = nil
			case 5:
				task.Days = nil
				task.Week = nil
				task.StartDate = nil
				task.OneOffDate = nil
			}

			go func() {
				data := shared.SchedTaskUpdateData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}
				done := false
				rpcClient.Call("TaskRPC.UpdateSched", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &task)

		// and set the swap panel to match the data
		for i, f := range freqs {
			if f.Name == task.Freq {
				swapper.Select(i)
				break
			}
		}

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

		// Add some action buttons for this schedule
		form.ActionGrid("sched-actions", "#action-grid", task, func(url string) {
			done := false
			switch url {
			case "play":
				go rpcClient.Call("TaskRPC.SchedPlay", task.ID, &done)
				task.Paused = false
				doc.QuerySelector("#playtask").Class().Add("action-hidden")
				doc.QuerySelector("#pausetask").Class().Remove("action-hidden")
				task.Paused = false
				title := plainTitle
				if task.Paused {
					title += " (PAUSED)"
				}
				form.SetTitle(title)
			default:

			case "pause":
				go rpcClient.Call("TaskRPC.SchedPause", task.ID, &done)
				task.Paused = true
				doc.QuerySelector("#pausetask").Class().Add("action-hidden")
				doc.QuerySelector("#playtask").Class().Remove("action-hidden")
				title := plainTitle
				if task.Paused {
					title += " (PAUSED)"
				}
				form.SetTitle(title)
			default:
				Session.Router.Navigate(url)
			}
		})

		// Set the initial vis of the action items
		print("paused =", task.Paused)
		if task.Paused {
			doc.QuerySelector("#playtask").Class().Remove("action-hidden")
		} else {
			doc.QuerySelector("#pausetask").Class().Remove("action-hidden")
		}

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

	weekdays := []formulate.SelectOption{
		{1, "Mon"},
		{2, "Tue"},
		{3, "Wed"},
		{4, "Thur"},
		{5, "Fri"},
	}

	go func() {
		machine := shared.Machine{}
		task := shared.SchedTask{}
		technicians := []shared.User{}
		rpcClient.Call("MachineRPC.Get", id, &machine)
		rpcClient.Call("UserRPC.GetTechnicians", machine.SiteID, &technicians)

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
		swapper.AddPanel("week").AddRow(2).
			AddRadio(1, "Week of the Month", "Week", weeks, "ID", "Name", 1).
			AddRadio(1, "Weekday", "WeekDay", weekdays, "ID", "Name", 1)
		swapper.AddPanel("year").AddRow(1).AddDate(1, "Day of the Year", "StartDate")
		swapper.AddPanel("days").AddRow(1).AddNumber(1, "Number of Days", "Days", "1")
		swapper.AddPanel("oneoff").AddRow(1).AddDate(1, "One Off Date", "OneOffDate")
		swapper.AddPanel("job").AddRow(1).AddNumber(1, "Job Count", "Count", "1")

		// Layout the fields
		form.Row(3).
			AddRadio(1, "Frequency", "Freq", freqs, "ID", "Name", 1).
			AddSwapper(2, "Frequency Options:", &swapper)

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
				0).
			AddSelect(1, "Assign To Technician", "UserID", technicians, "ID", "Username", 0, 0)

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
			// print("comp = ", comp)
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
					break
				}
			}
			switch targetFreq {
			case 1:
				task.Days = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 2:
				task.Days = nil
				task.Count = nil
				task.Week = nil
				task.OneOffDate = nil
			case 3:
				task.Week = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 4:
				task.Days = nil
				task.Count = nil
				task.StartDate = nil
				task.Week = nil
			case 5:
				task.Days = nil
				task.Week = nil
				task.StartDate = nil
				task.OneOffDate = nil
			}

			go func() {
				data := shared.SchedTaskUpdateData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}
				newID := 0
				rpcClient.Call("TaskRPC.InsertSched", data, &newID)
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

func siteTaskList(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		site := shared.Site{}
		tasks := []shared.Task{}
		rpcClient.Call("SiteRPC.Get", id, &site)
		rpcClient.Call("TaskRPC.SiteList", id, &tasks)

		BackURL := fmt.Sprintf("/site/%d", id)
		form := formulate.ListForm{}
		form.New("fa-server", "Active Tasks for "+site.Name)

		// Define the layout

		switch Session.UserRole {
		case "Admin", "Site Manager":
			form.Column("User", "Username")
		}

		form.Column("Date", "GetStartDate")
		// form.Column("Due", "GetDueDate")
		// form.Column("Site", "SiteName")
		form.Column("Machine", "MachineName")
		form.Column("Component", "Component")
		form.Column("Description", "Descr")
		form.Column("Duration", "DurationDays")
		form.Column("Completed", "CompletedDate")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/task/" + key)
		})

		form.Render("site-task-list", "main", tasks)

	}()
}
