package main

import (
	"strconv"
	"time"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/gopherjs/gopherjs/js"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func adminReports(context *router.Context) {

	rep := shared.Report{}

	rep.DateFrom = time.Now()
	rep.DateTo = time.Now()
	BackURL := "/"
	title := "Reports"
	form := formulate.EditForm{}
	form.New("fa-bar-chart", title)

	// Layout the fields
	form.Row(1).
		AddCustom(1, "Selection", "Dates", "")

	form.Row(1).
		AddCustom(1, "Report", "Report", "")

		// AddDate(1, "From", "DateFrom").
		// AddDate(1, "To", "DateTo")

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate(BackURL)
	})

	w := dom.GetWindow()
	doc := w.Document()

	// All done, so render the form
	form.Render("edit-form", "main", &rep)

	sel := doc.QuerySelector("[name=Dates]")
	if sel != nil {
		sel.SetInnerHTML("")
		sg := doc.CreateElement("div")
		sg.Class().Add("site-grid")
		sel.AppendChild(sg)

		b1 := doc.CreateElement("div")
		b1.Class().Add("site__item")
		b2 := doc.CreateElement("div")
		b2.Class().Add("site__title")
		b2.SetInnerHTML("Year to Date")
		b1.AppendChild(b2)
		sg.AppendChild(b1)
		b1.AddEventListener("click", false, func(evt dom.Event) {
			print("clicked on year to date")
		})

		b1 = doc.CreateElement("div")
		b1.Class().Add("site__item")
		b2 = doc.CreateElement("div")
		b2.Class().Add("site__title")
		b2.SetInnerHTML("Current Quarter")
		b1.AppendChild(b2)
		sg.AppendChild(b1)
		b1.AddEventListener("click", false, func(evt dom.Event) {
			print("clicked on quarterly")
		})

		b1 = doc.CreateElement("div")
		b1.Class().Add("site__item")
		b2 = doc.CreateElement("div")
		b2.Class().Add("site__title")
		b2.SetInnerHTML("Current Month")
		b1.AppendChild(b2)
		sg.AppendChild(b1)
		b1.AddEventListener("click", false, func(evt dom.Event) {
			print("clicked on current month")
		})
	}

	// var chart = c3.generate({
	//     bindto: '#chart',
	//     data: {
	//       columns: [
	//         ['data1', 30, 200, 100, 400, 150, 250],
	//         ['data2', 50, 20, 10, 40, 15, 25]
	//       ]
	//     }
	// });
	c3 := js.Global.Get("c3")
	data1 := make([]interface{}, 7)
	data2 := make([]interface{}, 7)
	cols := make([]interface{}, 2)
	data1[0] = "data1"
	data1[1] = 30
	data1[2] = 200
	data1[3] = 100
	data1[4] = 400
	data1[5] = 150
	data1[6] = 250
	data2[0] = "data2"
	data2[1] = 50
	data2[2] = 20
	data2[3] = 10
	data2[4] = 40
	data2[5] = 15
	data2[6] = 25
	cols[0] = data1
	cols[1] = data2
	// data1 = {"data", 30, 200, 100, 400, 150, 250}
	// data2 :=	{"data2", 50, 20, 10, 40, 15, 25}
	params := js.M{
		"bindto": "[name=Report]",
		"data": js.M{
			"columns": cols,
		},
	}
	print("about to call c3 with params", params)
	c3.Call("generate", params)

}

func hashtagList(context *router.Context) {

	// gob.Register(shared.HashtagUpdateData{})
	Session.Subscribe("hashtag", _hashtagList)
	go _hashtagList("list", 0)
}

func _hashtagList(action string, id int) {
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
		Session.Navigate("/util")
	})

	form.NewRowEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate("/hashtag/add")
	})

	form.RowEvent(func(key string) {
		Session.Navigate("/hashtag/" + key)
	})

	form.Render("hashtag-list", "main", hashtags)
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
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&hashtag)
			go func() {
				newID := 0
				rpcClient.Call("TaskRPC.HashtagInsert", shared.HashtagRPCData{
					Channel: Session.Channel,
					Hashtag: &hashtag,
				}, &newID)
				print("added hashtag", newID)
				Session.Navigate(BackURL)
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
	Session.ID["hashtag"] = id

	Session.Subscribe("hashtag", _hashtagEdit)
	go _hashtagEdit("edit", id)
}

func _hashtagEdit(action string, id int) {

	BackURL := "/hashtags"

	switch action {
	case "edit":
		print("manually edit")
	case "delete":
		if id != Session.ID["hashtag"] {
			return
		}
		print("current record has been deleted")
		Session.Navigate(BackURL)
		return
	default:
		if id != Session.ID["hashtag"] {
			return
		}
	}
	hashtag := shared.Hashtag{}
	rpcClient.Call("TaskRPC.HashtagGet", shared.HashtagRPCData{
		Channel: Session.Channel,
		ID:      id,
	}, &hashtag)

	title := "Edit Hashtag - #" + hashtag.Name

	form := formulate.EditForm{}
	form.New("fa-hashtag", title)

	// Layout the fields
	form.Row(1).
		AddInput(1, "Hashtag Name (without the # symbol)", "Name")

	form.Row(1).
		AddCustom(1, "Markup Rules", "Markup", "")
	form.Row(1).
		AddBigTextarea(1, "Expand to", "Descr")
	form.Row(1).
		AddCustom(1, "Expands to :", "Expand", "")

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate(BackURL)
	})

	form.DeleteEvent(func(evt dom.Event) {
		evt.PreventDefault()
		hashtag.ID = id
		go func() {
			done := false
			rpcClient.Call("TaskRPC.HashtagDelete", shared.HashtagRPCData{
				Channel: Session.Channel,
				Hashtag: &hashtag,
			}, &done)
			Session.Navigate(BackURL)
		}()
	})

	form.SaveEvent(func(evt dom.Event) {
		evt.PreventDefault()
		form.Bind(&hashtag)
		go func() {
			done := false
			rpcClient.Call("TaskRPC.HashtagUpdate", shared.HashtagRPCData{
				Channel: Session.Channel,
				Hashtag: &hashtag,
			}, &done)
			Session.Navigate(BackURL)
		}()
	})

	// All done, so render the form
	form.Render("edit-form", "main", &hashtag)
	setMarkupButtons("Descr")

	// Add an action grid
	form.ActionGrid("hash-action", "#action-grid", hashtag.ID, func(url string) {
		Session.Navigate(url)
	})
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
		Session.Navigate(BackURL)
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
				if w.Confirm("Generate Machine Parts ?") {
					rpcClient.Call("UtilRPC.Machine", Session.Channel, &retval)
				}
			case "part":
				if w.Confirm("Generate Part Class ?") {
					rpcClient.Call("UtilRPC.Parts", Session.Channel, &retval)
				}
			case "cat":
				if w.Confirm("Generate Category ?") {
					rpcClient.Call("UtilRPC.Cats", Session.Channel, &retval)
				}
			case "mtt":
				if w.Confirm("Generat MTT Links ?") {
					rpcClient.Call("UtilRPC.MTT", Session.Channel, &retval)
				}
			// case "photomove":
			// 	if w.Confirm("Generate Photos ?") {
			// 		rpcClient.Call("UtilRPC.PhotoMove", Session.Channel, &retval)
			// 	}
			case "sms":
				Session.Navigate("/sms")
				return
			case "hashtag":
				Session.Navigate("/hashtags")
				return
			case "editor":
				Session.Navigate("/testeditor")
			case "taskfigs":
				if w.Confirm("Generate Accumulators on all Tasks ?") {
					rpcClient.Call("UtilRPC.TaskFigs", Session.Channel, &retval)
				}
			case "phototest":
				Session.Navigate("/phototest")
			case "thumbnails":
				if w.Confirm("Generate New Thumbnails and Previews ?") {
					rpcClient.Call("UtilRPC.Thumbnails", Session.Channel, &retval)
				}
			case "users":
				Session.Navigate("/usersonline")
			default:
				print("ERROR - unknown utility", url)
				return
			}
			// el.Value = retval
			el.SetTextContent(retval)
		}()
	})

}
func testeditor(context *router.Context) {

	go func() {

		BackURL := "/util"
		form := formulate.EditForm{}
		form.New("fa-edit", "Test Markup Rules")

		// Layout the fields

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

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", nil)

		setMarkupButtons("Notes")
	}()

}

func showProgress(txt string) {

	// display the photo upload progress widget
	w := dom.GetWindow()
	doc := w.Document()

	if ee := doc.QuerySelector("#photoprogress"); ee != nil {
		ee.Class().Add("md-show")
		if pt := doc.QuerySelector("#progresstext"); pt != nil {
			pt.SetInnerHTML(txt)
		}
	}
}

func hideProgress() {

	// display the photo upload progress widget
	w := dom.GetWindow()
	doc := w.Document()

	if ee := doc.QuerySelector("#photoprogress"); ee != nil {
		ee.Class().Remove("md-show")
		if pt := doc.QuerySelector("#progresstext"); pt != nil {
			pt.SetInnerHTML("")
		}
	}
}
