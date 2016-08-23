package main

import (
	"fmt"
	"strconv"
	"strings"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func showSchedPhotos(task shared.SchedTask) {
	// print("populate the photos", task)

	w := dom.GetWindow()
	doc := w.Document()
	div := doc.QuerySelector("[name=Photos]")
	div.SetInnerHTML("")

	for _, v := range task.Photos {
		// print(k, ":", v)
		// Create an image widget, and add it to the photos block
		i := doc.CreateElement("img").(*dom.HTMLImageElement)
		i.SetAttribute("photo-id", fmt.Sprintf("%d", v.ID))
		i.Class().SetString("photopreview")
		i.Src = v.Preview
		switch v.Type {
		case "PDF":
			// Is a PDF, so wrap the image with a box that includes the filename
			// and auto-break on each doc

			wspan := doc.CreateElement("div")
			wspan.AppendChild(i)
			p := doc.CreateElement("p")
			p.SetInnerHTML(v.Filename)
			wspan.AppendChild(p)
			div.AppendChild(wspan)
		case "Image":
			div.AppendChild(i)
		default:
			print("adding attachment of unknown type", v.Type, "dt", v.Datatype, "fn", v.Filename)
			print("v", v)
		}
		// print("attaching click event to i", i)
		i.AddEventListener("click", false, func(evt dom.Event) {
			print("click on attachment preview image")
			evt.PreventDefault()
			theID, _ := strconv.Atoi(evt.Target().GetAttribute("photo-id"))

			go func() {
				photo := shared.Photo{}
				rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					ID:      theID,
				}, &photo)
				flds := strings.SplitN(photo.Data, ",", 2)
				print("got full photo", flds[0])
				switch flds[0] {
				case "data:application/pdf;base64":
					w.Open(photo.Data, "", "")
				case "data:image/jpeg;base64", "data:image/png;base64", "data:image/gif;base64":
					if el2 := doc.QuerySelector("#photo-full").(*dom.HTMLImageElement); el2 != nil {
						doc.QuerySelector("#show-image").Class().Add("md-show")
						el2.Src = photo.Data
					}
				}

			}()
		})
	}
}

type MachineSchedListData struct {
	Machine shared.Machine
	Tasks   []shared.SchedTask
}

// List all scheduled maint tasks that include this hashtag
func hashtagUsed(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		hashtag := shared.Hashtag{}
		tasks := []shared.SchedTask{}
		data := shared.HashtagRPCData{
			Channel: Session.Channel,
			ID:      id,
		}
		rpcClient.Call("TaskRPC.HashtagGet", data, &hashtag)
		rpcClient.Call("TaskRPC.ListHashSched", data, &tasks)

		BackURL := fmt.Sprintf("/hashtag/%d", id)

		form := formulate.ListForm{}
		form.New("fa-wrench", "Sched Maint that includes #"+hashtag.Name)

		// Define the layout
		form.Column("Tool / Component", "Component")
		form.Column("Frequency", "ShowFrequency")
		form.Column("Description", "Descr")
		form.MultiImgColumn("Documents", "Photos", "Thumb")
		form.Column("$ Labour", "LabourCost")
		form.Column("$ Materials", "MaterialCost")
		form.Column("Duration", "DurationDays")
		form.Column("Job Status", "ShowPaused")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/sched/" + key)
		})

		form.Render("hash-sched-list", "main", tasks)
	}()
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
		data := shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}
		rpcClient.Call("MachineRPC.Get", data, &machine)
		rpcClient.Call("TaskRPC.ListMachineSched", data, &tasks)

		BackURL := fmt.Sprintf("/machine/%d", id)

		form := formulate.ListForm{}
		form.New("fa-wrench", "Sched Maint for - "+machine.Name+" - "+*machine.SiteName)

		// Define the layout
		form.Column("Tool / Component", "Component")
		form.Column("Frequency", "ShowFrequency")
		form.Column("Description", "Descr")
		form.MultiImgColumn("Documents", "Photos", "Thumb")
		// form.Column("$ Labour", "LabourCost")
		// form.Column("$ Materials", "MaterialCost")
		form.Column("Duration", "DurationDays")
		form.Column("Job Status", "ShowPaused")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(fmt.Sprintf("/machine/sched/add/%d", id))
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/sched/" + key)
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
		{4, "Every N Months"},
		{5, "One Off"},
		{6, "Job Count"},
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

		rpcClient.Call("TaskRPC.GetSched", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &task)

		rpcClient.Call("MachineRPC.Get", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      task.MachineID,
		}, &machine)

		rpcClient.Call("UserRPC.GetTechnicians", shared.UserRPCData{
			Channel: Session.Channel,
			ID:      machine.SiteID,
		}, &technicians)

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
		swapper.AddPanel("months").AddRow(1).AddNumber(1, "Every N Months", "Months", "1")
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
			newOpt := formulate.SelectOption{ID: i + 1, Name: comp.Name}
			compTools = append(compTools, newOpt)
		}

		compOther := []formulate.SelectOption{
			{101, "Uncoiler"},
			{102, "RollBed"},
			{103, "Conveyor"},
			{104, "Electrical"},
			{105, "Hydraulic"},
			{106, "Pnuematic"},
			{107, "Lube"},
			{108, "Printer"},
			{109, "Console"},
			{110, "Encoder"},
			{111, "StripGuide"},
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
			AddCustom(1, "Markup Rules", "Markup", "")
		form.Row(1).
			AddBigTextarea(1, "Task Description", "Descr")
		form.Row(1).
			AddCustom(1, "Expands to :", "Expand", "")

		form.Row(3).
			AddDecimal(1, "Labour Cost", "LabourCost", 2, "1").
			AddDecimal(1, "Material Cost", "MaterialCost", 2, "1").
			AddNumber(1, "Duration (days)", "DurationDays", "1")

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Photos", "Photos", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				done := false
				rpcClient.Call("TaskRPC.DeleteSched", shared.SchedTaskRPCData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)
			print("post bind", task)
			task.MachineID = machine.ID

			// interpret the Component from the grouped options
			comp, _ := strconv.Atoi(task.Component)
			print("comp2 = ", comp)
			// print("len comp = ", len(machine.Components))
			if comp == 0 {
				task.CompType = "A"
				task.Component = compGen[0].Name
			} else if comp <= len(machine.Components) {
				task.CompType = "T"
				task.ToolID = machine.Components[comp-1].ID
				task.Component = compTools[comp-1].Name
			} else {
				task.CompType = "C"
				offset := comp - len(machine.Components)
				print("edit offset = ", offset)
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
				task.Months = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 2:
				task.Days = nil
				task.Months = nil
				task.Count = nil
				task.Week = nil
				task.OneOffDate = nil
			case 3:
				task.Week = nil
				task.Count = nil
				task.Months = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 4:
				task.Days = nil
				task.Week = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 5:
				task.Days = nil
				task.Months = nil
				task.Count = nil
				task.StartDate = nil
				task.Week = nil
			case 6:
				task.Days = nil
				task.Months = nil
				task.Week = nil
				task.StartDate = nil
				task.OneOffDate = nil
			}

			// If the uploaded data is a PDF, then use that data instead of the preview
			if task.NewPhoto.Data != "" && isPDF {
				task.NewPhoto.Data = PDFData
				isPDF = false
			}

			go func() {
				done := false
				rpcClient.Call("TaskRPC.UpdateSched", shared.SchedTaskRPCData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}, &done)
				// Session.Navigate(BackURL)
				Session.Reload(context)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &task)
		setMarkupButtons("Descr")
		setPhotoField("NewPhoto")
		showSchedPhotos(task)

		// Setup a callback on the freq selector
		w := dom.GetWindow()
		doc := w.Document()

		// and set the swap panel to match the data
		for i, f := range freqs {
			if f.Name == task.Freq {
				swapper.Select(i)
				break
			}
		}

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
					swapper.SelectByName("months")
				case 5:
					swapper.SelectByName("oneoff")
				case 6:
					swapper.SelectByName("job")
				}
			}
		})

		// Add some action buttons for this schedule
		form.ActionGrid("sched-actions", "#action-grid", task, func(url string) {
			done := false
			switch url {
			case "play":
				go rpcClient.Call("TaskRPC.SchedPlay", shared.SchedTaskRPCData{
					Channel: Session.Channel,
					ID:      task.ID,
				}, &done)
				task.Paused = false
				doc.QuerySelector("#playtask").Class().Add("action-hidden")
				doc.QuerySelector("#pausetask").Class().Remove("action-hidden")
				task.Paused = false
				title := plainTitle
				if task.Paused {
					title += " (PAUSED)"
				}
				form.SetTitle(title)

			case "pause":
				go rpcClient.Call("TaskRPC.SchedPause", shared.SchedTaskRPCData{
					Channel: Session.Channel,
					ID:      task.ID,
				}, &done)
				task.Paused = true
				doc.QuerySelector("#pausetask").Class().Add("action-hidden")
				doc.QuerySelector("#playtask").Class().Remove("action-hidden")
				title := plainTitle
				if task.Paused {
					title += " (PAUSED)"
				}
				form.SetTitle(title)
			default:
				Session.Navigate(url)
			}
		})

		// Set the initial vis of the action items
		// print("paused =", task.Paused)
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
		{4, "Every N Months"},
		{5, "One Off"},
		{6, "Job Count"},
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
		rpcClient.Call("MachineRPC.Get", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machine)
		rpcClient.Call("UserRPC.GetTechnicians", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      machine.SiteID,
		}, &technicians)

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
		swapper.AddPanel("months").AddRow(1).AddNumber(1, "Every N Months", "Months", "1")
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
			newOpt := formulate.SelectOption{ID: i + 1, Name: comp.Name}
			compTools = append(compTools, newOpt)
		}

		compOther := []formulate.SelectOption{
			{101, "Uncoiler"},
			{102, "RollBed"},
			{103, "Conveyor"},
			{104, "Electrical"},
			{105, "Hydraulic"},
			{106, "Pnuematic"},
			{107, "Lube"},
			{108, "Printer"},
			{109, "Console"},
			{110, "Encoder"},
			{111, "StripGuide"},
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
			AddCustom(1, "Markup Rules", "Markup", "")
		form.Row(1).
			AddBigTextarea(1, "Task Description", "Descr")
		form.Row(1).
			AddCustom(1, "Expands to :", "Expand", "")

		form.Row(3).
			AddDecimal(1, "Labour Cost", "LabourCost", 2, "1").
			AddDecimal(1, "Material Cost", "MaterialCost", 2, "1").
			AddNumber(1, "Duration (days)", "DurationDays", "1")

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Photos", "Photos", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)
			task.MachineID = machine.ID

			// interpret the Component from the grouped options
			comp, _ := strconv.Atoi(task.Component)
			// print("comp1 = ", comp)
			// print("len comp", len(machine.Components))
			if comp == 0 {
				task.CompType = "A"
				task.Component = compGen[0].Name
			} else if comp <= len(machine.Components) {
				task.CompType = "T"
				task.ToolID = machine.Components[comp-1].ID
				task.Component = compTools[comp-1].Name
			} else {
				task.CompType = "C"
				offset := comp - len(machine.Components)
				print("offset = ", comp)
				task.Component = compOther[offset].Name
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
				task.Months = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 2:
				task.Days = nil
				task.Months = nil
				task.Count = nil
				task.Week = nil
				task.OneOffDate = nil
			case 3:
				task.Week = nil
				task.Months = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 4:
				task.Week = nil
				task.Days = nil
				task.Count = nil
				task.StartDate = nil
				task.OneOffDate = nil
			case 5:
				task.Days = nil
				task.Months = nil
				task.Count = nil
				task.StartDate = nil
				task.Week = nil
			case 6:
				task.Days = nil
				task.Months = nil
				task.Week = nil
				task.StartDate = nil
				task.OneOffDate = nil
			}

			go func() {
				newID := 0
				rpcClient.Call("TaskRPC.InsertSched", shared.SchedTaskRPCData{
					Channel:   Session.Channel,
					SchedTask: &task,
				}, &newID)
				// print("added task ID", newID)
				Session.Navigate(fmt.Sprintf("/sched/%d", newID))
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &task)
		setMarkupButtons("Descr")
		setPhotoField("NewPhoto")

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
					swapper.SelectByName("months")
				case 5:
					swapper.SelectByName("oneoff")
				case 6:
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
		rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &site)
		rpcClient.Call("TaskRPC.SiteList", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &tasks)

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
		// form.Column("Machine", "MachineName")
		// form.Column("Component", "Component")
		form.Column("Component", "GetComponent")
		form.Column("Description", "Descr")
		form.Column("Duration", "DurationDays")
		form.Column("Completed", "CompletedDate")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/task/" + key)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.Render("site-task-list", "main", tasks)

	}()
}

func schedTaskList(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		sched := shared.SchedTask{}
		tasks := []shared.Task{}
		machine := shared.Machine{}

		rpcClient.Call("TaskRPC.GetSched", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &sched)

		rpcClient.Call("TaskRPC.SchedList", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &tasks)

		rpcClient.Call("MachineRPC.Get", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      sched.MachineID,
		}, &machine)

		BackURL := fmt.Sprintf("/sched/task/%d", id)
		form := formulate.ListForm{}

		Title := fmt.Sprintf("Generated Tasks for - %s - %s", machine.Name, *machine.SiteName)
		form.New("fa-server", Title)
		// Define the layout

		switch Session.UserRole {
		case "Admin", "Site Manager":
			form.Column("User", "Username")
		}

		form.Column("Date", "GetStartDate")
		// form.Column("Due", "GetDueDate")
		// form.Column("Site", "SiteName")
		// form.Column("Machine", "MachineName")
		// form.Column("Component", "Component")
		form.Column("Component", "GetComponent")
		form.Column("Description", "Descr")
		form.Column("Duration", "DurationDays")
		form.Column("Completed", "CompletedDate")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/task/" + key)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.Render("site-task-list", "main", tasks)
	}()

}
