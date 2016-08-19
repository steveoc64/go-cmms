package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

var PDFImage string
var PDFData string
var isPDF bool

func GetPDFImage() {
	go func() {
		rpcClient.Call("UtilRPC.GetPDF", Session.Channel, &PDFImage)
		print("Cached standard PDF image", PDFImage[:44])
	}()
}

func setPhotoOnlyField(f string) {
	setPhotoUploadField(f, false)
}

func setPhotoField(f string) {
	setPhotoUploadField(f, true)
}

func setPhotoUploadField(f string, allowPDF bool) {

	w := dom.GetWindow()
	doc := w.Document()

	// add a handler on the photo field
	if el := doc.QuerySelector(fmt.Sprintf("[name=%s", f)).(*dom.HTMLInputElement); el != nil {

		// Set the attribute to say what types of data the field can accept
		el.AddEventListener("change", false, func(evt dom.Event) {
			print("in the change event and allowPDF = ", allowPDF)
			files := el.Files()
			fileReader := js.Global.Get("FileReader").New()
			fileReader.Set("onload", func(e *js.Object) {
				target := e.Get("target")
				imgData := target.Get("result").String()
				print("imgdata =", imgData[:80])
				flds := strings.Split(imgData, ";")
				imgEl := doc.QuerySelector(fmt.Sprintf("[name=%sPreview]", f)).(*dom.HTMLImageElement)

				switch flds[0] {
				case "data:application/pdf":
					// if is pdf, then load the standard preview into the field
					if allowPDF {

						imgEl.Src = PDFImage
						imgEl.Class().Remove("hidden")
						PDFData = imgData
						isPDF = true
					} else {
						w.Alert("ERROR: This screen only allows photos, not PDF files.")
					}
				case "data:image/jpeg", "data:image/png":
					// if is image, then load the image into the preview
					imgEl.Src = imgData
					imgEl.Class().Remove("hidden")
					isPDF = false
				}
			})
			fileReader.Set("onerror", func(e *js.Object) {
				err := e.Get("target").Get("error")
				print("Error reading file", err)
			})
			if len(files) > 0 {
				fileReader.Call("readAsDataURL", files[0])
			}
		})

	}
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
	// print("phototest edit", id)

	go func() {
		photo := shared.Photo{}

		rpcClient.Call("UtilRPC.GetPhoto", shared.PhotoRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &photo)
		// print("got photo", photo)

		BackURL := "/phototest"
		form := formulate.EditForm{}
		form.New("fa-camera-retro", "Photo Edit Tester")

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
			// print("delete photo")
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
				// print("clicked on the photo")
				evt.PreventDefault()

				showProgress("Loading File ...")

				go func() {
					rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
						Channel: Session.Channel,
						ID:      id,
					}, &photo)

					flds := strings.SplitN(photo.Photo, ",", 2)
					print("got full photo", flds[0])
					switch flds[0] {
					case "data:application/pdf;base64":
						w.Open(photo.Photo, "", "")
					case "data:image/jpeg;base64", "data:image/png;base64":
						// print("got fullsize image")
						el.Src = photo.Photo
						el.Class().Remove("photopreview")
						el.Class().Add("photofull")
					}

					hideProgress()

					// restyle the preview to be full size
				}()

			})
		}

	}()

}

func phototestAdd(context *router.Context) {
	// print("phototest add")

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

			showProgress("Uploading File ...")

			evt.PreventDefault()
			form.Bind(&photo)

			// If the uploaded data is a PDF, then use that data instead of the preview
			if isPDF {
				photo.Photo = PDFData
			}

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
		setPhotoField("Photo")

	}()

}
