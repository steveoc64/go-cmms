package main

import (
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func adminReports(context *router.Context) {
	print("TODO - adminReports")

}

func hashtagList(context *router.Context) {

	go func() {
		hashtags := []shared.Hashtag{}
		rpcClient.Call("TaskRPC.HashtagList", Session.Channel, &hashtags)

		form := formulate.ListForm{}
		form.New("fa-hashtag", "Hashtag List")

		// Define the layout
		form.Column("Name", "HashName")
		form.Column("Expands To", "Descr")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/util")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/hashtag/add")
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/hashtag/" + key)
		})

		form.Render("hashtag-list", "main", hashtags)
	}()
}

func hashtagAdd(context *router.Context) {
	go func() {
		hashtag := shared.Hashtag{}

		BackURL := "/hashtags"
		title := "Add New Hashtag"
		form := formulate.EditForm{}
		form.New("fa-hashtag", title)

		// Layout the fields
		form.Row(1).
			AddInput(1, "Hashtag Name (without the # symbol)", "Name")

		form.Row(1).
			AddTextarea(1, "Expand to", "Descr")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&hashtag)
			data := shared.HashtagUpdateData{
				Channel: Session.Channel,
				Hashtag: &hashtag,
			}
			go func() {
				newID := 0
				rpcClient.Call("TaskRPC.HashtagInsert", data, &newID)
				print("added hashtag", newID)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &hashtag)

	}()

}

func hashtagEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		hashtag := shared.Hashtag{}
		rpcClient.Call("TaskRPC.HashtagGet", id, &hashtag)

		BackURL := "/hashtags"
		title := "Edit Hashtag - #" + hashtag.Name

		form := formulate.EditForm{}
		form.New("fa-hashtag", title)

		// Layout the fields
		form.Row(1).
			AddInput(1, "Hashtag Name (without the # symbol)", "Name")

		form.Row(1).
			AddTextarea(1, "Expand to", "Descr")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			hashtag.ID = id
			go func() {
				data := shared.HashtagUpdateData{
					Channel: Session.Channel,
					Hashtag: &hashtag,
				}
				done := false
				rpcClient.Call("TaskRPC.HashtagDelete", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&hashtag)
			data := shared.HashtagUpdateData{
				Channel: Session.Channel,
				Hashtag: &hashtag,
			}
			go func() {
				done := false
				rpcClient.Call("TaskRPC.HashtagUpdate", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &hashtag)

		// Add an action grid
		form.ActionGrid("hash-action", "#action-grid", hashtag.ID, func(url string) {
			Session.Router.Navigate(url)
		})

	}()

}

func adminUtils(context *router.Context) {

	BackURL := "/"
	title := "Admin Utilities"
	form := formulate.EditForm{}
	form.New("fa-gear", title)

	form.Row(1).AddCodeBlock(1, "Results", "Results")

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Router.Navigate(BackURL)
	})

	// All done, so render the form
	form.Render("edit-form", "main", nil)

	w := dom.GetWindow()
	doc := w.Document()

	// el := doc.QuerySelector("[name=Results]").(*dom.HTMLTextAreaElement)
	el := doc.QuerySelector("[name=Results]")

	isSteve := false
	if Session.UserRole == "Admin" && Session.Username == "steve" {
		isSteve = true
	}
	form.ActionGrid("admin-util-actions", "#action-grid", isSteve, func(url string) {
		retval := ""
		go func() {
			switch url {
			case "backup":
				rpcClient.Call("UtilRPC.Backup", Session.Channel, &retval)
			case "top":
				rpcClient.Call("UtilRPC.Top", Session.Channel, &retval)
			case "logs":
				rpcClient.Call("UtilRPC.Logs", Session.Channel, &retval)
			case "machine":
				rpcClient.Call("UtilRPC.Machine", Session.Channel, &retval)
			case "part":
				rpcClient.Call("UtilRPC.Parts", Session.Channel, &retval)
			case "hashtag":
				Session.Router.Navigate("/hashtags")
				return
			default:
				print("ERROR - unknown utility", url)
				return
			}
			// el.Value = retval
			el.SetTextContent(retval)
		}()
	})

}
