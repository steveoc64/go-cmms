package main

import (
	"honnef.co/go/js/dom"
)

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
	print("a=", a)
	a.SetTextContent(name)
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
func (doc *Doc) addInputField(id string, tp string, pl string, cols string) *dom.HTMLDivElement {

	div := doc.addDiv()
	div.SetClass("input-field col " + cols)

	input := doc.addInput(id, tp, pl)
	input.SetClass("validate")
	label := doc.addLabel(id, pl)

	div.AppendChild(input)
	div.AppendChild(label)

	return div
}

// Create an Checkbox, with label
func (doc *Doc) addCheckbox(id string, pl string, cols string) *dom.HTMLDivElement {

	div := doc.addDiv()
	div.SetClass("input-field col " + cols)

	ch := doc.Document.CreateElement("input").(*dom.HTMLInputElement)
	ch.Type = "checkbox"
	ch.SetID(id)

	label := doc.addLabel(id, pl)

	div.AppendChild(ch)
	div.AppendChild(label)

	return div
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
		print("clicked submit btn")
		login()
	})

	div.AppendChild(btn)
	return div
}

// Create a gridded layout
func (doc *Doc) createLayout() *dom.HTMLDivElement {

	// Create the basic layout
	layout := doc.addDiv()
	layout.SetClass("row")
	return layout
}

// <nav>
//     <div class="nav-wrapper">
//       <a href="#" class="brand-logo">Logo</a>
//       <ul id="nav-mobile" class="right hide-on-med-and-down">
//         <li><a href="sass.html">Sass</a></li>
//         <li><a href="badges.html">Components</a></li>
//         <li><a href="collapsible.html">JavaScript</a></li>
//       </ul>
//     </div>
// </nav>
func (doc *Doc) createNavBar() *dom.BasicHTMLElement {

	nav := doc.addElement("nav")
	nav.SetClass("indigo")

	// Wrapper
	navwrapper := doc.addDiv()

	// Logo
	logo := doc.addDiv()
	logo.SetClass("brand-logo")
	logo.SetTextContent("CMMS")
	navwrapper.AppendChild(logo)

	// Top level options
	ul := doc.addUL()
	ul.SetClass("right hide-on-med-and-down")
	login := doc.addLI("Login", "#")
	ul.AppendChild(login)
	login.AddEventListener("click", false, func(event dom.Event) {
		print("clicked login btn")
		doc.showLoginForm()
	})

	navwrapper.AppendChild(ul)

	nav.AppendChild(navwrapper)

	return nav
}

// Create a Splash screen
func (doc *Doc) createSplash() *dom.HTMLDivElement {

	div := doc.addDiv()
	div.SetID("splash")
	img := doc.addIMG("/img/sbs01.jpg")
	banner := doc.addH("h3", "CMMS Facilities Management")
	div.AppendChild(img)
	div.AppendChild(banner)
	return div
}

// Create a login form
func (doc *Doc) createLoginForm() *dom.HTMLDivElement {

	// basic container and framework for the form
	div := doc.addDiv()
	div.SetClass("container")
	div.SetID("loginform")

	row := doc.addDiv()
	row.SetClass("row")

	col := doc.addDiv()
	col.SetClass("col s6 offset-s3")

	h3 := doc.addH("h3", "Login")
	h3.SetClass("center-align")

	// username
	ruser := doc.addRow()
	username := doc.addInputField("l-username", "text", "User Name", "s12")
	ruser.AppendChild(username)

	// passwd
	rpass := doc.addRow()
	pw := doc.addInputField("l-passwd", "password", "PassWord", "s12")
	rpass.AppendChild(pw)

	// remember me
	rrem := doc.addRow()
	rem := doc.addCheckbox("l-remember", "Remember Me ?", "s6")
	sub := doc.addSubmit("Login", "s6", login)
	rrem.AppendChild(rem)
	rrem.AppendChild(sub)

	// submit

	col.AppendChild(h3)
	col.AppendChild(ruser)
	col.AppendChild(rpass)
	col.AppendChild(rrem)

	row.AppendChild(col)
	div.AppendChild(row)
	div.Style().SetProperty("display", "none", "")

	return div
}

func (doc *Doc) showLoginForm() {

	// hide splash and show login
	sp := doc.GetElementByID("splash").(*dom.HTMLDivElement)
	sp.Style().SetProperty("display", "none", "")
	lf := doc.GetElementByID("loginform").(*dom.HTMLDivElement)
	lf.Style().SetProperty("display", "inline", "")
}

func login() {
	print("Here we are !!")
}
