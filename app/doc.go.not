package main

import (
	"honnef.co/go/js/dom"
)

type Doc struct {
	dom.Document
}

var doc Doc

// Create a new DIV element
func (doc *Doc) addDiv() *dom.HTMLDivElement {
	return doc.Document.CreateElement("div").(*dom.HTMLDivElement)
}

// Create a new DIV row element
func (doc *Doc) addRow() *dom.HTMLDivElement {
	row := doc.Document.CreateElement("div").(*dom.HTMLDivElement)
	row.SetClass("row")
	return row
}

// Create a new Basic Element
func (doc *Doc) addElement(name string) *dom.BasicHTMLElement {
	return doc.Document.CreateElement(name).(*dom.BasicHTMLElement)
}

// Create a new UL Element
func (doc *Doc) addUL() *dom.HTMLUListElement {
	return doc.Document.CreateElement("ul").(*dom.HTMLUListElement)
}

// Create a new List Item Element, with a link
func (doc *Doc) addLI(name string, link string) *dom.HTMLLIElement {
	li := doc.Document.CreateElement("li").(*dom.HTMLLIElement)
	a := doc.Document.CreateElement("a").(*dom.HTMLAnchorElement)
	a.URLUtils.Href = link
	// print("a=", a)
	a.SetTextContent(name)
	a.SetID("nav-login-btn")
	li.AppendChild(a)
	return li
}

// Create a Heading element
func (doc *Doc) addH(level string, content string) *dom.HTMLHeadingElement {
	h := doc.Document.CreateElement(level).(*dom.HTMLHeadingElement)
	h.SetTextContent(content)
	return h
}

// Create an IMG element
func (doc *Doc) addIMG(image string) *dom.HTMLImageElement {
	i := doc.Document.CreateElement("img").(*dom.HTMLImageElement)
	i.Src = image
	return i
}

// Create a Label element
func (doc *Doc) addLabel(id string, content string) *dom.HTMLLabelElement {
	l := doc.Document.CreateElement("label").(*dom.HTMLLabelElement)
	l.For = id
	l.SetTextContent(content)
	return l
}

// Create a simple input field
func (doc *Doc) addInput(id string, tp string, pl string) *dom.HTMLInputElement {
	i := doc.Document.CreateElement("input").(*dom.HTMLInputElement)
	i.Placeholder = pl
	i.Type = tp
	return i
}

// Create an InputField, with label
func (doc *Doc) addInputField(id string, tp string, pl string, cols string) (*dom.HTMLDivElement, *dom.HTMLInputElement) {

	div := doc.addDiv()
	div.SetClass("input-field col " + cols)

	input := doc.addInput(id, tp, pl)
	input.SetClass("validate")
	label := doc.addLabel(id, pl)

	div.AppendChild(input)
	div.AppendChild(label)

	return div, input
}

// Create an Checkbox, with label
func (doc *Doc) addCheckbox(id string, pl string, cols string) (*dom.HTMLDivElement, *dom.HTMLInputElement) {

	div := doc.addDiv()
	div.SetClass("input-field col " + cols)

	ch := doc.Document.CreateElement("input").(*dom.HTMLInputElement)
	ch.Type = "checkbox"
	ch.SetID(id)

	label := doc.addLabel(id, pl)

	div.AppendChild(ch)
	div.AppendChild(label)

	return div, ch
}

// Create a Submit Button
func (doc *Doc) addSubmit(name string, cols string, fn func()) *dom.HTMLDivElement {

	div := doc.addDiv()
	div.SetClass("col " + cols)

	btn := doc.Document.CreateElement("button").(*dom.HTMLButtonElement)
	btn.SetClass("btn btn-large waves-effect waves-light")
	btn.Name = name
	btn.Type = "button"
	btn.SetTextContent(name)
	btn.AddEventListener("click", false, func(event dom.Event) {
		// print("clicked submit btn")
		go fn()
	})

	div.AppendChild(btn)
	return div
}
