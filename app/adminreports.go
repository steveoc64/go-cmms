package main

import (
	"strconv"
	"strings"
	"time"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/gopherjs/gopherjs/js"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func adminReports(context *router.Context) {
	print("TODO - adminReports")
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
		AddTextarea(1, "Expand to", "Descr")

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
			case "photomove":
				if w.Confirm("Generate Photos ?") {
					rpcClient.Call("UtilRPC.PhotoMove", Session.Channel, &retval)
				}
			case "sms":
				Session.Navigate("/sms")
				return
			case "hashtag":
				Session.Navigate("/hashtags")
				return
			case "editor":
				Session.Navigate("/testeditor")
			case "phototest":
				Session.Navigate("/phototest")
			default:
				print("ERROR - unknown utility", url)
				return
			}
			// el.Value = retval
			el.SetTextContent(retval)
		}()
	})

}

func phototest(context *router.Context) {

	go func() {

		photos := []shared.Photo{}
		rpcClient.Call("UtilRPC.PhotoList", shared.PhotoRPCData{
			Channel: Session.Channel,
		}, &photos)

		form := formulate.ListForm{}
		form.New("fa-camera-retro", "Photo List")

		// Define the layout
		form.Column("Notes", "Notes")
		form.Column("Table", "Entity")
		form.Column("ID", "EntityID")
		form.ImgColumn("Thumbnail", "Thumb")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/phototest/add")
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/phototest/" + key)
		})

		form.Render("photo-list", "main", photos)

		// manually set thumbnails on the fields for now, until formulate is refactored
		// w := dom.GetWindow()
		// doc := w.Document()

		// for _, v := range photos {
		// 	ename := fmt.Sprintf(`[name=Thumbnail-%d]`, v.ID)
		// 	print("ename = ", ename, v.ID)
		// 	el := doc.QuerySelector(ename).(*dom.HTMLImageElement)
		// 	el.Src = v.Thumbnail
		// }
	}()
}

func phototestEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	print("phototest edit", id)

	go func() {
		photo := shared.Photo{}

		rpcClient.Call("UtilRPC.GetPhoto", shared.PhotoRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &photo)
		// print("got photo", photo)

		BackURL := "/phototest"
		form := formulate.EditForm{}
		form.New("fa-camera-retro", "Photo Upload Tester")

		// Layout the fields

		form.Row(1).
			AddInput(1, "Notes", "Notes")

		form.Row(2).
			AddInput(1, "Table", "Entity").
			AddNumber(1, "ID", "EntityID", "0")

		form.Row(1).
			AddPreview(1, "Preview", "Preview")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("update photo")
			form.Bind(&photo)
			// print("bind the photo gives", photo)
			go func() {
				done := false
				rpcClient.Call("UtilRPC.UpdatePhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					ID:      id,
					Photo:   &photo,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("delete photo")
			form.Bind(&photo)
			// print("bind the photo gives", photo)
			go func() {
				done := false
				rpcClient.Call("UtilRPC.DeletePhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					ID:      id,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &photo)

		w := dom.GetWindow()
		doc := w.Document()

		if el := doc.QuerySelector("[name=PreviewPreview]").(*dom.HTMLImageElement); el != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				print("clicked on the photo")
				evt.PreventDefault()

				showProgress("Loading Photo ...")

				go func() {
					rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
						Channel: Session.Channel,
						ID:      id,
					}, &photo)

					print("got fullsize image")
					el.Src = photo.Photo
					el.Class().Remove("photopreview")
					el.Class().Add("photofull")
					hideProgress()

					// restyle the preview to be full size
				}()

			})
		}

	}()

}

func phototestAdd(context *router.Context) {
	print("phototest add")

	go func() {
		photo := shared.Photo{}

		BackURL := "/phototest"
		form := formulate.EditForm{}
		form.New("fa-camera-retro", "Photo Upload Tester")

		// Layout the fields

		form.Row(1).
			AddInput(1, "Name", "Notes")

		form.Row(1).
			AddPhoto(1, "Photo", "Photo")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.SaveEvent(func(evt dom.Event) {

			showProgress("Uploading Photo ...")

			evt.PreventDefault()
			form.Bind(&photo)
			go func() {
				newID := 0
				rpcClient.Call("UtilRPC.AddPhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					Photo:   &photo,
				}, &newID)

				// sleep 1
				Session.Navigate(BackURL)
				hideProgress()
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &photo)

		// add a handler on the photo field
		w := dom.GetWindow()
		doc := w.Document()
		if el := doc.QuerySelector("[name=Photo]").(*dom.HTMLInputElement); el != nil {
			el.AddEventListener("change", false, func(evt dom.Event) {
				files := el.Files()
				fileReader := js.Global.Get("FileReader").New()
				fileReader.Set("onload", func(e *js.Object) {
					target := e.Get("target")
					imgData := target.Get("result").String()
					imgEl := doc.QuerySelector(`.photouppreview`).(*dom.HTMLImageElement)
					imgEl.Src = imgData
					imgEl.SetAttribute("src", imgData)
					imgEl.Class().Remove("hidden")
				})
				fileReader.Set("onerror", func(e *js.Object) {
					err := e.Get("target").Get("error")
					print("Error reading file", err)
				})
				fileReader.Call("readAsDataURL", files[0])
			})

		}
	}()

}

func testeditor(context *router.Context) {
	print("markdown editor test")

	go func() {

		BackURL := "/util"
		form := formulate.EditForm{}
		form.New("fa-edit", "Test Markdown Editor")

		// Layout the fields

		form.Row(1).
			AddCustom(1, "Markup Rules", "Markup", "")

		form.Row(1).
			AddBigTextarea(1, "Notes", "Notes")

		form.Row(1).
			AddCustom(1, "Expands to :", "Expand", "")

		clock := time.NewTimer(3 * time.Second)
		defer clock.Stop()

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

		// Add a change event on the big textarea
		w := dom.GetWindow()
		doc := w.Document()

		el := doc.QuerySelector("[name=Markup]").(*dom.HTMLDivElement)
		el.SetInnerHTML(`
Use the following markup to format the text :
<ul>
<li> <b>Bold Text</b>  Wrap the ^Bold Text^ using the ^ symbol.
<li> <u>Underline Text</u>  Wrap the _Underline Text_ using the _ symbol.
<li> <span class=redtext>Red Text</span>  Wrap the {Red Text} using the {} symbols.
<li> Start a line with  -  to create a list

</ul>

`)

		t := doc.QuerySelector("[name=Notes]").(*dom.HTMLTextAreaElement)
		et := t.Value

		exp := doc.QuerySelector("[name=Expand]").(*dom.HTMLDivElement)
		exp.SetInnerHTML("...")

		// set a timer to scan the contents
		if false {
			go func() {
				print("run a timer func to scan the contents of the notes")

				for now := range clock.C {
					print("got timer", now)

					if t.Value == et {
						print("no change")
					} else {
						print("notes has changed")
						et = t.Value
						renderMarkdown(exp, et)
					}
				}

			}()
		}

		t.AddEventListener("change", false, func(evt dom.Event) {
			notes := evt.Target().(*dom.HTMLTextAreaElement)
			print("notes has changed to", notes.Value)

			renderMarkdown(exp, notes.Value)
		})

	}()

}

func renderMarkdown(el *dom.HTMLDivElement, text string) {

	w := dom.GetWindow()
	doc := w.Document()

	el.SetInnerHTML("")

	// split the input into lines

	lines := strings.Split(text, "\n")

	para := ""
	for k, v := range lines {
		l := strings.Trim(v, " ")
		print("Line", k+1, ":", l)

		// if blank, then complete the paragraph
		if l == "" && len(para) > 0 {
			div := doc.CreateElement("div").(*dom.HTMLDivElement)
			div.SetInnerHTML(para)
			el.AppendChild(div)
			para = ""
		} else {
			// append this to the existing paragraph.
			if len(para) > 0 {
				para += " "
			}
			para += l
		}
	}
	if len(para) > 0 {
		div := doc.CreateElement("div").(*dom.HTMLDivElement)
		div.SetInnerHTML(para)
		el.AppendChild(div)
		para = ""
	}

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
