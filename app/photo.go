package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"

	"honnef.co/go/js/dom"
)

func setPhotoField(f string) {

	w := dom.GetWindow()
	doc := w.Document()

	// add a handler on the photo field
	if el := doc.QuerySelector(fmt.Sprintf("[name=%s", f)).(*dom.HTMLInputElement); el != nil {
		el.AddEventListener("change", false, func(evt dom.Event) {
			files := el.Files()
			fileReader := js.Global.Get("FileReader").New()
			fileReader.Set("onload", func(e *js.Object) {
				target := e.Get("target")
				imgData := target.Get("result").String()
				//print("imgdata =", imgData)
				imgEl := doc.QuerySelector(fmt.Sprintf("[name=%sPreview]", f)).(*dom.HTMLImageElement)
				imgEl.Src = imgData
				imgEl.Class().Remove("hidden")
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
