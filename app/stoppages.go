package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func stoppageList(context *router.Context) {
	Session.Subscribe("event", _stoppageList)
	go _stoppageList("list", 0)
}

func _stoppageList(action string, id int) {

	events := []shared.Event{}
	rpcClient.Call("EventRPC.List", Session.Channel, &events)

	// print("events =", events)

	form := formulate.ListForm{}
	form.New("fa-pause-circle-o", "Current Stoppages")

	// Define the layout
	form.Column("ID/User", "GetUserNameID")
	form.Column("Date", "GetStartDate")

	// if Session.UserRole == "Admin" {
	// 	form.Column("Completed", "GetCompleted")
	// }

	form.ColumnFormat("Site", "SiteName", "GetSiteClass")
	// form.Column("Machine", "MachineName")
	// form.Column("Component", "ToolType")
	form.Column("Component", "GetComponent")
	form.MultiImgColumn("Photo", "Photos", "Thumb")
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

	// manually display the images, until formulate is refactored
	// w := dom.GetWindow()
	// doc := w.Document()
	// for _, v := range events {
	// 	if v.PhotoThumbnail != "" {
	// 		ename := fmt.Sprintf(`[name=PhotoThumbnail-%d]`, v.ID)
	// 		print("ename = ", ename)
	// 		el := doc.QuerySelector(ename).(*dom.HTMLImageElement)
	// 		el.Src = v.PhotoThumbnail
	// 	}
	// }

	// completed events

	if Session.UserRole == "Admin" {

		cevents := []shared.Event{}
		rpcClient.Call("EventRPC.ListCompleted", Session.Channel, &cevents)

		cform := formulate.ListForm{}
		cform.New("fa-pause-circle-o", "Completed Stoppages (Last 30 Days)")

		// Define the layout
		cform.Column("Raised By", "Username")
		cform.Column("Date", "GetStartDate")
		cform.Column("Completed", "GetCompleted")
		cform.Column("Site", "SiteName")
		// cform.Column("Machine", "MachineName")
		// cform.Column("Component", "ToolType")
		cform.Column("Component", "GetComponent")
		cform.MultiImgColumn("Photo", "Photos", "Thumb")
		cform.Column("Notes", "Notes")
		// cform.Column("Status", "GetStatus")

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
		cform.Render("cstoppage-list", "#cevent", cevents)

		// manually display the images, until formulate is refactored
		// for _, v := range events {
		// 	if v.PhotoThumbnail != "" {
		// 		ename := fmt.Sprintf(`[name=PhotoThumbnail-%d]`, v.ID)
		// 		print("ename = ", ename)
		// 		el := doc.QuerySelector(ename).(*dom.HTMLImageElement)
		// 		el.Src = v.PhotoThumbnail
		// 	}
		// }

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
	RefreshURL := fmt.Sprintf("/stoppage/%d", id)

	switch action {
	case "edit":
		// print("manually edit")
	case "delete":
		if id != Session.ID["event"] {
			return
		}
		// print("current record has been deleted")
		Session.Navigate(BackURL)
		return
	default:
		if id != Session.ID["event"] {
			return
		}
	}

	event := shared.Event{}

	rpcClient.Call("EventRPC.Get", shared.EventRPCData{
		Channel: Session.Channel,
		ID:      id,
	}, &event)

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

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Attachments", "Photos", "")

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

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Attachments", "Photos", "")

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
				done := false
				rpcClient.Call("EventRPC.Delete", shared.EventRPCData{
					Channel: Session.Channel,
					Event:   &event,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		if Session.UserRole == "Admin" {
			form.SaveEvent(func(evt dom.Event) {
				evt.PreventDefault()
				form.Bind(&event)

				if event.NewPhoto.Data != "" {
					showProgress("Updating Event ...")

					// If the uploaded data is a PDF, then use that data instead of the preview
					event.NewPhoto.Data = ImageCache.GetImage()
				}

				go func() {
					done := false
					rpcClient.Call("EventRPC.Update", shared.EventRPCData{
						Channel: Session.Channel,
						Event:   &event,
					}, &done)
					hideProgress()
					// Session.Navigate(BackURL)
					// print("go to", RefreshURL)
					Session.Navigate(RefreshURL)
				}()
			})
		}
	}

	// All done, so render the form
	form.Render("edit-form", "main", &event)
	setPhotoField("NewPhoto")
	showEventPhotos(event)

	// and show the assignments
	loadTemplate("stoppage-assigned-to", "[name=AssignedTo]", event)

	// And attach actions
	switch Session.UserRole {
	case "Admin":
		form.ActionGrid("event-actions", "#action-grid", event, func(url string) {
			Session.Navigate(url)
		})
	case "Site Manager":
		// lookup the current user, if can allocate, then give them the admin options
		u := shared.User{}
		rpcClient.Call("UserRPC.Get", shared.UserRPCData{
			Channel: Session.Channel,
			ID:      Session.UserID,
		}, &u)
		if u.CanAllocate {
			form.ActionGrid("event-actions", "#action-grid", event, func(url string) {
				Session.Navigate(url)
			})
		} else {
			form.ActionGrid("event-sm-actions", "#action-grid", event.ID, func(url string) {
				Session.Navigate(url)
			})
		}
	}

}

func showEventPhotos(event shared.Event) {
	// print("populate the photos", event)

	w := dom.GetWindow()
	doc := w.Document()
	div := doc.QuerySelector("[name=Photos]")
	div.SetInnerHTML("")

	for _, v := range event.Photos {
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
		case "photo":
			print("WARNING: This is an old format attachment of type photo")
			print("v", v)
			print("Please run  \"update photo set type='Image' where type='photo';\" on database ....")
		default:
			wspan := doc.CreateElement("div")
			wspan.AppendChild(i)
			p := doc.CreateElement("p")
			p.SetInnerHTML(v.Filename)
			wspan.AppendChild(p)
			div.AppendChild(wspan)
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
				// case "data:application/pdf;base64":
				default:
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

func stoppageComplete(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	BackURL := fmt.Sprintf("/stoppage/%d", id)

	if Session.UserRole != "Admin" {

		go func() {

			// read the can allocate flag to see if this is allowed or not
			u := shared.User{}
			rpcClient.Call("UserRPC.Get", shared.UserRPCData{
				Channel: Session.Channel,
				ID:      Session.UserID,
			}, &u)
			if !u.CanAllocate {
				print("Not allowed to allocate tasks")
				Session.Navigate(BackURL)
			}
		}()
	}

	go func() {
		event := shared.Event{}
		event.ID = id
		done := false
		rpcClient.Call("EventRPC.Complete", shared.EventRPCData{
			Channel: Session.Channel,
			Event:   &event,
		}, &done)
		print("Completed Event", id)
		Session.Navigate("/stoppages")
	}()

}

func stoppageNewTask(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	BackURL := fmt.Sprintf("/stoppage/%d", id)

	if Session.UserRole != "Admin" {

		// read the can allocate flag to see if this is allowed or not
		go func() {

			u := shared.User{}
			rpcClient.Call("UserRPC.Get", shared.UserRPCData{
				Channel: Session.Channel,
				ID:      Session.UserID,
			}, &u)
			if !u.CanAllocate {
				print("Not allowed to allocate tasks")
				Session.Navigate(BackURL)
			}
		}()
	}

	go func() {
		event := shared.Event{}
		techs := []shared.User{}

		rpcClient.Call("EventRPC.Get", shared.EventRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &event)
		rpcClient.Call("UserRPC.GetTechnicians", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      0,
		}, &techs)

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
			Photo:       event.NewPhoto,
		}

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

		form.Row(2).
			AddPhoto(1, "Photos", "Photo").
			AddPreview(1, "", "Preview")

		form.Row(1).
			AddCustom(1, "Markup Rules", "Markup", "")
		form.Row(1).
			AddBigTextarea(1, "Notes", "Notes")
		form.Row(1).
			AddCustom(1, "Expands to :", "Expand", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&assign)
			print("post bind of the assign record", assign)
			Session.Navigate(BackURL)
			showProgress("Creating Task ...")

			go func() {
				newID := 0
				rpcClient.Call("EventRPC.Workorder", assign, &newID)
				print("new Task raised", newID)
				hideProgress()
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &assign)
		setMarkupButtons("Notes")
		setPhotoField("Photo")
	}()

}
