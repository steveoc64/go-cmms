package main

import (
	"fmt"
	"itrak-cmms/shared"
	"strings"

	"honnef.co/go/js/dom"
)

func renderMarkup(el *dom.HTMLDivElement, text string) {

	w := dom.GetWindow()
	doc := w.Document()

	el.SetInnerHTML("")

	// split the input into lines

	lines := strings.Split(text, "\n")
	// print("lines =", lines)

	para := ""
	for _, v := range lines {
		l := strings.TrimRight(v, " ")
		l = strings.Replace(l, "<", "&lt;", -1)
		l = strings.Replace(l, ">", "&gt;", -1)
		// print("Line", k+1, ":", l)

		// if blank, then complete the paragraph
		if l == "" && len(para) > 0 {
			div := doc.CreateElement("div").(*dom.HTMLDivElement)
			div.SetInnerHTML(parsePara(para, true))
			el.AppendChild(div)
			para = ""
		} else {
			// append this to the existing paragraph.
			if len(para) > 0 {
				para += "\n"
			}
			para += l
		}
	}
	if len(para) > 0 {
		div := doc.CreateElement("div").(*dom.HTMLDivElement)
		div.SetInnerHTML(parsePara(para, true))
		el.AppendChild(div)
	}

}

// Parse a paragraph
func parsePara(para string, addbr bool) string {

	// print("parsing", para)
	retval := ""
	listmode := false

	for _, line := range strings.Split(para, "\n") {

		// Indentation
		if strings.HasPrefix(line, " ") {
			// println("add space")
			retval += "&nbsp;"
			retval += parsePara(line[1:], addbr)
			continue
		}

		// Add horizontal line
		if strings.HasPrefix(line, "---") {
			retval += "<hr>\n"
			continue
		}

		// Generate ordered list
		if strings.HasPrefix(line, "-") {
			if !listmode {
				listmode = true
				retval += "<ol>\n"
			}
			retval += fmt.Sprintf("<li>%s\n", parsePara(line[1:], false))
			continue
		}

		// Big Header
		if strings.HasPrefix(line, "!!!") {
			retval += fmt.Sprintf("<h1>%s</h1>\n", line[3:])
			continue
		}

		// Medium Header
		if strings.HasPrefix(line, "!!") {
			retval += fmt.Sprintf("<h2>%s</h2>\n", line[2:])
			continue
		}

		// Small Header
		if strings.HasPrefix(line, "!") {
			retval += fmt.Sprintf("<h3>%s</h3>\n", line[1:])
			continue
		}

		// Bold Text
		if x := strings.Index(line, "^"); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "^"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := fmt.Sprintf("%s<b>%s</b>%s", line[:x], line[x+1:x2], line[x2+1:])
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// Underscore Text
		if x := strings.Index(line, "_"); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "_"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := fmt.Sprintf("%s<u>%s</u>%s", line[:x], line[x+1:x2], line[x2+1:])
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// Italic Text
		if x := strings.Index(line, "~"); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "~"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := fmt.Sprintf("%s<i>%s</i>%s", line[:x], line[x+1:x2], line[x2+1:])
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// RED Text
		if x := strings.Index(line, "@"); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "@"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := fmt.Sprintf("%s<span class=redtext>%s</span>%s", line[:x], line[x+1:x2], line[x2+1:])
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// Green Text
		if x := strings.Index(line, "{"); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "}"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := fmt.Sprintf("%s<span class=greentext>%s</span>%s", line[:x], line[x+1:x2], line[x2+1:])
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// Checkbox
		if x := strings.Index(line, "["); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "]"); x2 > -1 {
				x2 += x + 1
				// println("x2 = ", x2)
				embolden := ""
				if addbr {
					embolden = fmt.Sprintf("%s\n<br><input type=checkbox><label class=label-inline>%s</label>\n<br>%s\n",
						line[:x], line[x+1:x2], line[x2+1:])
				} else {
					embolden = fmt.Sprintf("%s <input type=checkbox><label class=label-inline>%s</label> %s\n",
						line[:x], line[x+1:x2], line[x2+1:])
				}
				retval += parsePara(embolden, addbr)
				continue
			}
		}

		// Cleanup and exit
		if listmode {
			retval += "</ol>\n"
		}
		retval += line
		if addbr {
			retval += "<br>"
		}
	}

	return retval
}

func expandHashtags(s string) string {

	hashes := []shared.Hashtag{}
	rpcClient.Call("TaskRPC.HashtagListByLen", Session.Channel, &hashes)
	// print("hashes", hashes)
	if strings.Contains(s, "#") {
		// Keep looping through doing text conversions until there is
		// nothing left to expand
		stillLooking := true
		for stillLooking {
			stillLooking = false
			for _, v := range hashes {
				theHash := "#" + v.Name
				if strings.Contains(s, theHash) {
					s = strings.Replace(s, theHash, v.Descr, -1)
					stillLooking = true
				}
			}
		}
	} // contains hashes
	return s
}

func setMarkupButtons(txtfld string) {

	w := dom.GetWindow()
	doc := w.Document()

	el := doc.QuerySelector("[name=Markup]").(*dom.HTMLDivElement)
	el.SetInnerHTML(MarkupHelp)

	doc.QuerySelector("[name=helptext]").AddEventListener("click", false, func(evt dom.Event) {
		doc.QuerySelector("[name=helptext]").Class().Add("hidden")
		doc.QuerySelector("[name=helpbtn]").Class().Remove("hidden")
	})

	doc.QuerySelector("[name=helpbtn").AddEventListener("click", false, func(evt dom.Event) {
		evt.Target().Class().Add("hidden")
		doc.QuerySelector("[name=helptext]").Class().Remove("hidden")
	})

	exp := doc.QuerySelector("[name=Expand]").(*dom.HTMLDivElement)
	exp.SetInnerHTML(MarkupExpand)

	doc.QuerySelector("[name=expandbtn").AddEventListener("click", false, func(evt dom.Event) {
		el := doc.QuerySelector("[name=expanded-text]").(*dom.HTMLDivElement)
		el.SetInnerHTML("... expanding")
		notes := doc.QuerySelector(fmt.Sprintf("[name=%s]", txtfld)).(*dom.HTMLTextAreaElement)
		go func() {
			renderMarkup(el, expandHashtags(notes.Value))
		}()
	})
}

const MarkupHelp string = `
<input type=button class=button-primary name=helpbtn value=Help>
<div name=helptext class=hidden>

<h3>! Small Heading</h3>
<h2>!! Medium Heading</h2>
<h1>!!! Large Heading</h1>

<ul>
<li> <hr> Use 3 or more dashes (---) in a row to create a line break like the one above
<li> <b>Bold Text</b>  Wrap the ^Bold Text^ using the ^ symbol.
<li> <u>Underline Text</u>  Wrap the _Underline Text_ using the _ symbol.
<li> <i>Italic Text</i>  Wrap the ~Italic Text~ using the ~ symbol.
<li> <span class=redtext>Red Text</span>  Wrap the @Red Text@ using the @ symbol.
<li> <span class=greentext>Green Text</span>  Wrap the {Green Text} using the {} symbols.
<li> Start a line with  -  to add to an auto-numbered list
</ul>

<input type=checkbox id=testbox>
<label for=testbox class=label-inline>
[Enter a paragraph of text inside square brackets to associate a checkbox with the whole paragraph]
</label>
</div>
`
const MarkupExpand string = `
<input type=button class=button-primary name=expandbtn value=Test>
<div name=expanded-text>
</div>
`
