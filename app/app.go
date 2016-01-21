package main

import (
	"honnef.co/go/js/dom"
	"strconv"
)

func inputKeyUp(event dom.Event) {
	input := event.Target().(*dom.HTMLInputElement)

	span := dom.GetWindow().Document().GetElementByID("inputvalue")
	span.SetInnerHTML(input.Value)
}

func checkParent(target, elm dom.Element) bool {
	for target.ParentElement() != nil {
		if target.IsEqualNode(elm) {
			return true
		}
		target = target.ParentElement()
	}
	return false
}

func main() {
	w := dom.GetWindow()
	print("w =", w)
	d := w.Document()
	print("d =", d)
	h := d.GetElementByID("foo")
	print("h =", h)
	k := d.GetElementByID("keycode")
	print("k =", k)
	i := d.GetElementByID("input").(*dom.HTMLInputElement)
	vv := d.GetElementByID("vv").(*dom.HTMLDivElement)
	hidex := d.GetElementByID("hidex").(*dom.HTMLDivElement)
	createx := d.GetElementByID("createx").(*dom.HTMLDivElement)

	print(vv.Dataset()["what"])
	h.AddEventListener("click", false, func(event dom.Event) {
		event.PreventDefault()
		h.SetInnerHTML("I am Clicked")
		println("This message is printed in browser console")
	})

	w.AddEventListener("keydown", false, func(event dom.Event) {
		ke := event.(*dom.KeyboardEvent)
		k.SetInnerHTML(strconv.Itoa(ke.KeyCode))
	})

	i.Focus()
	i.AddEventListener("keyup", false, inputKeyUp)

	hidex.AddEventListener("click", false, func(event dom.Event) {
		hidex.Style().SetProperty("display", "none", "")
	})

	createx.AddEventListener("click", false, func(event dom.Event) {
		div := d.CreateElement("div").(*dom.HTMLDivElement)
		div.Style().SetProperty("color", "red", "")
		div.SetTextContent("I am new div")
		createx.AppendChild(div)
	})

	audbtn := d.GetElementByID("audbtn").(*dom.HTMLButtonElement)
	audbtn.AddEventListener("click", false, func(event dom.Event) {
		a := d.GetElementByID("audio").(*dom.HTMLAudioElement)
		if a.Paused {
			a.Play()
			audbtn.SetTextContent("Click Me to Pause Sound")
		} else {
			a.Pause()
			audbtn.SetTextContent("Click Me to Play Sound")
		}
	})

	toggle := d.GetElementByID("menu-dropdown").(*dom.HTMLDivElement)
	menu := d.GetElementByID("menuDiv-dropdown").(*dom.HTMLDivElement)

	d.AddEventListener("click", false, func(event dom.Event) {
		if !checkParent(event.Target(), menu) {
			// click NOT on the menu
			if checkParent(event.Target(), toggle) {
				// click on the link
				menu.Class().Toggle("invisible")
			} else {
				// click both outside link and outside menu, hide menu
				menu.Class().Add("invisible")
			}
		}
	})
}
