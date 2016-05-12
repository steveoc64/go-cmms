package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func stoppageList(context *router.Context) {
	Session.Subscribe("event", _stoppageList)
	go _stoppageList("list", 0)
}

func _stoppageList(action string, id int) {

	events := []shared.Event{}
	rpcClient.Call("EventRPC.List", Session.Channel, &events)

	form := formulate.ListForm{}
	form.New("fa-pause-circle-o", "Current Stoppages")

	// Define the layout
	form.Column("Raised By", "Username")
	form.Column("Date", "GetStartDate")

	if Session.UserRole == "Admin" {
		form.Column("Completed", "GetCompleted")
	}

	form.Column("Site", "SiteName")
	form.Column("Machine", "MachineName")
	form.Column("Component", "ToolType")
	form.Column("Notes", "Notes")

	switch Session.UserRole {
	case "Admin", "Site Manager":
		form.Column("Status", "GetStatus")
	}

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate("/")
	})

	form.RowEvent(func(key string) {
		Session.Navigate("/stoppage/" + key)
	})

	form.PrintEvent(func(evt dom.Event) {
		dom.GetWindow().Print()
	})

	form.Render("stoppage-list", "main", events)

	// completed events

	if Session.UserRole == "Admin" {

		cevents := []shared.Event{}
		rpcClient.Call("EventRPC.ListCompleted", Session.Channel, &cevents)

		cform := formulate.ListForm{}
		cform.New("fa-pause-circle-o", "Completed Stoppages")

		// Define the layout
		cform.Column("Raised By", "Username")
		cform.Column("Date", "GetStartDate")
		cform.Column("Completed", "GetCompleted")
		cform.Column("Site", "SiteName")
		cform.Column("Machine", "MachineName")
		cform.Column("Component", "ToolType")
		cform.Column("Notes", "Notes")
		cform.Column("Status", "GetStatus")

		cform.RowEvent(func(key string) {
			Session.Navigate("/stoppage/" + key)
		})

		w := dom.GetWindow()
		doc := w.Document()

		// force a page break for printing
		div := doc.CreateElement("div")
		div.Class().Add("page-break")
		doc.QuerySelector("main").AppendChild(div)

		div = doc.CreateElement("div").(*dom.HTMLDivElement)
		div.SetID("cevent")
		doc.QuerySelector("main").AppendChild(div)
		cform.Render("stoppage-list", "#cevent", cevents)
	}

}

func stoppageEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	Session.ID["event"] = id

	Session.Subscribe("event", _stoppageEdit)
	go _stoppageEdit("edit", id)
}

func _stoppageEdit(action string, id int) {

	BackURL := "/stoppages"

	switch action {
	case "edit":
		print("manually edit")
	case "delete":
		if id != Session.ID["event"] {
			return
		}
		print("current record has been deleted")
		Session.Navigate(BackURL)
		return
	default:
		if id != Session.ID["event"] {
			return
		}
	}

	event := shared.Event{}

	rpcClient.Call("EventRPC.Get", id, &event)

	title := fmt.Sprintf("Stoppage Details - %06d", id)
	form := formulate.EditForm{}
	form.New("fa-pause-circle-o", title)

	// print("and the startdate is ", event.StartDate)
	// print("and the startdate is ", event.StartDate.String())
	// event.DisplayDate = event.StartDate.String()
	event.DisplayDate = event.StartDate.Format("Mon, Jan 2 2006 15:04:05")

	// Layout the fields
	switch Session.UserRole {
	case "Admin":
		form.Row(2).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName")

		form.Row(3).
			AddDisplay(1, "Component", "ToolType").
			AddDisplay(1, "StartDate", "DisplayDate").
			AddDisplay(1, "Raised By", "Username")

		form.Row(1).
			AddBigTextarea(1, "Notes", "Notes")

		form.Row(1).
			AddCustom(1, "Assigned To", "AssignedTo", "")

	case "Site Manager":
		form.Row(2).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName")

		form.Row(3).
			AddDisplay(1, "Component", "ToolType").
			AddDisplay(1, "StartDate", "DisplayDate").
			AddDisplay(1, "Raised By", "Username")

		form.Row(1).
			AddDisplayArea(1, "Notes", "Notes")

		form.Row(1).
			AddCustom(1, "Assigned To", "AssignedTo", "")
	}

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate(BackURL)
	})

	form.PrintEvent(func(evt dom.Event) {
		dom.GetWindow().Print()
	})

	// Only Admin has the power to delete, update, or dig deeper on an event
	if Session.UserRole == "Admin" {
		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			event.ID = id
			go func() {
				data := shared.EventUpdateData{
					Channel: Session.Channel,
					Event:   &event,
				}
				done := false
				rpcClient.Call("EventRPC.Delete", data, &done)
				Session.Navigate(BackURL)
			}()
		})

		if Session.UserRole == "Admin" {
			form.SaveEvent(func(evt dom.Event) {
				evt.PreventDefault()
				form.Bind(&event)
				data := shared.EventUpdateData{
					Channel: Session.Channel,
					Event:   &event,
				}
				go func() {
					done := false
					rpcClient.Call("EventRPC.Update", data, &done)
					Session.Navigate(BackURL)
				}()
			})
		}
	}

	// All done, so render the form
	form.Render("edit-form", "main", &event)

	// and show the assignments
	loadTemplate("stoppage-assigned-to", "[name=AssignedTo]", event)

	// And attach actions
	switch Session.UserRole {
	case "Admin":
		form.ActionGrid("event-actions", "#action-grid", event, func(url string) {
			Session.Navigate(url)
		})
	case "Site Manager":
		form.ActionGrid("event-sm-actions", "#action-grid", event.ID, func(url string) {
			Session.Navigate(url)
		})
	}

}

func stoppageComplete(context *router.Context) {

	if Session.UserRole != "Admin" {
		return
	}

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		event := shared.Event{}
		event.ID = id
		data := shared.EventUpdateData{
			Channel: Session.Channel,
			Event:   &event,
		}
		done := false
		rpcClient.Call("EventRPC.Complete", data, &done)
		print("Completed Event", id)
		Session.Navigate("/stoppages")
	}()

}

func stoppageNewTask(context *router.Context) {
	if Session.UserRole != "Admin" {
		return
	}

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		event := shared.Event{}
		techs := []shared.User{}

		rpcClient.Call("EventRPC.Get", id, &event)
		rpcClient.Call("UserRPC.GetTechnicians", 0, &techs)

		now1 := time.Now()
		now2 := time.Now()

		assign := shared.AssignEvent{
			Channel:     Session.Channel,
			SiteName:    event.SiteName,
			MachineName: event.MachineName,
			ToolType:    event.ToolType,
			DisplayDate: event.StartDate.Format("Mon, Jan 2 2006 15:04:05"),
			Username:    event.Username,
			Event:       &event,
			StartDate:   &now1,
			DueDate:     &now2,
			Notes:       event.Notes,
		}

		BackURL := fmt.Sprintf("/stoppage/%d", id)
		title := fmt.Sprintf("Raise Task for Stoppage - %06d", id)
		form := formulate.EditForm{}
		form.New("fa-sign-in", title)

		// print("and the startdate is ", event.StartDate)
		// print("and the startdate is ", event.StartDate.String())
		event.DisplayDate = event.StartDate.String()
		event.DisplayDate = event.StartDate.Format("Mon, Jan 2 2006 15:04:05")

		// Layout the fields
		form.Row(2).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName")

		form.Row(3).
			AddDisplay(1, "Component", "ToolType").
			AddDisplay(1, "StartDate", "DisplayDate").
			AddDisplay(1, "Raised By", "Username")

		form.Row(3).
			AddSelect(1, "Assign To", "AssignTo", techs, "ID", "Username", 1, 0).
			AddDate(1, "Workorder Start Date", "StartDate").
			AddDate(1, "Workorder Due Date", "DueDate")

		form.Row(1).
			AddBigTextarea(1, "Notes", "Notes")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&assign)

			go func() {
				newID := 0
				rpcClient.Call("EventRPC.Workorder", assign, &newID)
				print("new Task raised", newID)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &assign)

	}()

}
