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

type CachedImages struct {
	ImageData string

	PDFImage string
	PDFData  string
	isPDF    bool

	RawDataImage string
	RawData      string
	isRawData    bool
}

func (c *CachedImages) Clear() {
	c.isPDF = false
	c.isRawData = false
	c.PDFData = ""
	c.RawData = ""
	c.ImageData = ""
}

func (c *CachedImages) SetPDF(data string) {
	c.PDFData = data
	c.isPDF = true
	c.isRawData = false
}

func (c *CachedImages) SetRawData(data string) {
	c.RawData = data
	c.isRawData = true
	c.isPDF = false
}

func (c *CachedImages) SetImage(data string) {
	c.ImageData = data
	c.isRawData = false
	c.isPDF = false
}

func (c *CachedImages) GetImage() string {
	if c.isPDF {
		return c.PDFData
	}
	if c.isRawData {
		return c.RawData
	}
	return c.ImageData
}

var ImageCache CachedImages

func GetPDFImage() {
	go func() {
		rpcClient.Call("UtilRPC.GetPDFImage", Session.Channel, &ImageCache.PDFImage)
		// print("Cached standard PDF image", ImageCache.PDFImage[:44])

		rpcClient.Call("UtilRPC.GetRawDataImage", Session.Channel, &ImageCache.RawDataImage)
		// print("Cached raw data image", ImageCache.RawDataImage[:44])
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
			// print("filename may =", el.Value)
			lastSlash := strings.LastIndex(el.Value, `\`)
			fileName := el.Value
			if lastSlash > -1 {
				fileName = fileName[lastSlash+1:]
			}
			// print("filename computed to", fileName)
			// print("in the change event and allowPDF = ", allowPDF)
			files := el.Files()
			fileReader := js.Global.Get("FileReader").New()
			fileReader.Set("onload", func(e *js.Object) {
				target := e.Get("target")
				imgData := target.Get("result").String()
				// print("imgdata =", imgData[:80])
				flds := strings.Split(imgData, ";")
				// print("attachment type", flds[0])
				imgEl := doc.QuerySelector(fmt.Sprintf("[name=%sPreview]", f)).(*dom.HTMLImageElement)

				ImageCache.Clear()
				switch flds[0] {
				case "data:application/pdf":
					// if is pdf, then load the standard preview into the field
					if allowPDF {

						imgEl.Src = ImageCache.PDFImage
						imgEl.Class().Remove("hidden")
						ImageCache.SetPDF(imgData)
						// print("photo changed and looks like a PDF")
					} else {
						w.Alert("ERROR: This screen only allows photos, not PDF files.")
					}
				case "data:image/jpeg", "data:image/png", "data:image/gif":
					// if is image, then load the image into the preview
					imgEl.Src = imgData
					ImageCache.SetImage(imgData)
					imgEl.Class().Remove("hidden")
				default:
					// print("Adding data of unknown type", flds[0])
					if allowPDF {
						imgEl.Src = ImageCache.RawDataImage
						imgEl.Class().Remove("hidden")
						ImageCache.SetRawData(imgData)
					} else {
						w.Alert("ERROR: This screen only allows photos, please try again")
					}
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
		form.Column("Filename", "Filename")

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
		print("got photo", photo)

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

		form.Row(2).
			AddInput(1, "Type", "Type").
			AddInput(1, "DataType", "Datatype")

		form.Row(6).
			AddInput(3, "Filename", "Filename").
			AddDisplay(1, "Full File Size", "Length").
			AddDisplay(1, "Preview Size", "LengthP").
			AddDisplay(1, "Thumbnail Size", "LengthT")

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

		// manually fill in the filename field
		if el := doc.QuerySelector("[name=PreviewFilename]").(*dom.HTMLSpanElement); el != nil {
			if photo.Filename != "" {
				el.SetInnerHTML(photo.Filename)
				el.Class().Remove("hidden")
			}
		}
		if el := doc.QuerySelector("[name=PreviewPreview]").(*dom.HTMLImageElement); el != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				// print("clicked on the photo")
				evt.PreventDefault()

				go func() {
					// rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
					// 	Channel: Session.Channel,
					// 	ID:      id,
					// }, &photo)

					// flds := strings.SplitN(photo.Data, ",", 2)
					// print("got full photo", flds[0])
					// switch flds[0] {
					switch photo.Type {
					case "PDF", "Data":
						print("open file in new window")
						if photo.Data == "" {
							print("get copy of the full data")
							showProgress("Loading File ...")
							rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
								Channel: Session.Channel,
								ID:      id,
							}, &photo)
						}
						w.Open(photo.Data, "", "")
					case "Image":
						// case "data:image/jpeg;base64", "data:image/png;base64", "data:image/gif;base64":
						// print("got fullsize image")

						// toggle the state
						cl := el.Class()
						if cl.Contains("photopreview") {
							cl.Remove("photopreview")
							cl.Add("photofull")
							// go fullscreen
							if photo.Data == "" {
								print("getting a copy of the full image")
								showProgress("Loading File ...")
								rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
									Channel: Session.Channel,
									ID:      id,
								}, &photo)
							}
							el.Src = photo.Data
						} else {
							// reduce back to preview
							el.Src = photo.Preview
							cl.Add("photopreview")
							cl.Remove("photofull")
						}

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
		photo := shared.Phototest{}

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
			print("binding into", photo)
			form.Bind(&photo)
			print("post bind into", photo)

			// If the uploaded data is a PDF, then use that data instead of the preview
			photo.Photo.Data = ImageCache.GetImage()

			go func() {
				newID := 0
				rpcClient.Call("UtilRPC.AddPhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					Photo: &shared.Photo{
						Data:     photo.Photo.Data,
						Filename: photo.Photo.Filename,
						Notes:    photo.Notes,
					},
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
