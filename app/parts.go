package main

import (
	"fmt"
	"strconv"

	"itrak-cmms/shared"

	"github.com/gopherjs/gopherjs/js"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func addTree(tree []shared.Category, ul *dom.HTMLUListElement, depth int) {

	w := dom.GetWindow()
	doc := w.Document()
	// print("adding from ", tree, " to ", ul)
	// Add a LI for each category
	for _, tv := range tree {
		// print("Tree Value", i, tv)
		widgetID := fmt.Sprintf("category-%d", tv.ID)
		li := doc.CreateElement("li")
		li.SetID(widgetID)
		chek := doc.CreateElement("input").(*dom.HTMLInputElement)
		chek.Type = "checkbox"
		label := doc.CreateElement("label")
		label.SetAttribute("for", widgetID)
		label.SetInnerHTML(tv.Name)
		label.SetAttribute("data-type", "category")
		label.SetAttribute("data-id", fmt.Sprintf("%d", tv.ID))
		// label.Class().Add("category")
		label.SetID(widgetID + "-label")
		chek.SetAttribute("data-type", "category")
		chek.SetAttribute("data-id", fmt.Sprintf("%d", tv.ID))
		chek.SetID(widgetID + "-chek")
		li.AppendChild(label)
		li.AppendChild(chek)
		ul.AppendChild(li)

		if len(tv.Subcats) > 0 {
			ul2 := doc.CreateElement("ul").(*dom.HTMLUListElement)
			li.AppendChild(ul2)
			addTree(tv.Subcats, ul2, depth+1)
		} else {
			if depth == 0 {
				ulempty := doc.CreateElement("ul")
				li.AppendChild(ulempty)
				liempty := doc.CreateElement("li")
				liempty.SetInnerHTML("(no sub-categories)")
				ulempty.AppendChild(liempty)
			}
		}

		ul3 := doc.CreateElement("ul")
		li.AppendChild(ul3)
		if len(tv.Parts) > 0 {
			for _, part := range tv.Parts {
				partID := fmt.Sprintf("part-%d", part.ID)
				li2 := doc.CreateElement("li")
				li2.SetID(partID)
				li2.SetInnerHTML(fmt.Sprintf(`%s : %s`, part.StockCode, part.Name))
				li2.Class().Add("stock-item")
				li2.SetAttribute("data-type", "part")
				li2.SetAttribute("data-id", fmt.Sprintf("%d", part.ID))
				ul3.AppendChild(li2)
			}
		} else {
			if depth > 0 {
				li3 := doc.CreateElement("li")
				li3.SetInnerHTML("(no parts)")
				ul3.AppendChild(li3)
			}
		}
	}
}

func partsList(context *router.Context) {

	currentCat := 0
	currentPart := 0
	var lastSelectedClass *dom.TokenList

	go func() {
		tree := []shared.Category{}
		rpcClient.Call("PartRPC.GetTree", shared.PartTreeRPCData{
			Channel:    Session.Channel,
			CategoryID: 0,
		}, &tree)
		print("got tree", tree)

		machineTypes := []shared.MachineType{}
		rpcClient.Call("MachineRPC.MachineTypes", shared.MachineRPCData{
			Channel: Session.Channel,
		}, &machineTypes)
		print("machine types", machineTypes)

		// thePart := shared.Part{}
		// theCat := shared.Category{}

		machineTools := []shared.MachineTypeTool{}

		BackURL := "/"

		form := formulate.EditForm{}
		form.New("fa-puzzle-piece", "Parts List")

		// create the swapper panels
		swapper := formulate.Swapper{
			Name:     "Details",
			Selected: 1,
		}

		form.Row(2).
			AddCustom(1, "Parts Tree", "tree", "tree").
			AddSwapper(1, "Details", &swapper)

		catPanel := swapper.AddPanel("Category")
		catPanel.AddRow(1).AddInput(1, "Category Name", "CatName")
		catPanel.AddRow(1).AddInput(1, "Stock Code", "CatStockCode")
		catPanel.AddRow(1).AddInput(1, "Description", "CatDescr")
		catPanel.AddRow(2).
			AddSelect(1, "Machine Type", "MachineType",
				machineTypes, "ID", "Name", 0, 0).
			AddSelect(1, "Tool", "MachineTool",
				machineTools, "ID", "Name", 0, 0)

		catPanel.Row(4).
			AddButton(1, "Save", "SaveCat").
			AddButton(1, "+ Category", "AddCat").
			AddButton(1, "+ Part", "AddPart").
			AddButton(1, "- Delete", "DelCat")

		// Layout the fields for Parts
		partPanel := swapper.AddPanel("Part")

		partPanel.Row(2).
			AddInput(1, "Name", "Name").
			AddInput(1, "Stock Code", "StockCode")

		partPanel.Row(1).
			AddInput(1, "Description", "Descr")

		partPanel.Row(4).
			AddDecimal(1, "ReOrder Level", "ReorderStocklevel", 2, "1").
			AddDecimal(1, "ReOrder Qty", "ReorderQty", 2, "1").
			AddDecimal(1, "Current Stock", "CurrentStock", 2, "1").
			AddInput(1, "Qty Type", "QtyType")

		partPanel.Row(4).
			// AddDisplay(2, "Last Price Update", "LastPriceDateDisplay").
			AddInput(2, "Supplier Info", "SupplierInfo").
			AddDecimal(1, "Latest Price", "LatestPrice", 2, "1").
			AddDisplay(1, "Valuation", "ValuationString")

		partPanel.Row(1).
			AddTextarea(1, "Notes", "Notes")

		partPanel.Row(2).
			AddButton(1, "Save", "SavePart").
			AddButton(1, "- Delete", "DelPart")

		// ID                   int        `db:"id"`
		// 	Class                int        `db:"class"`
		// 	Category             int        `db:"category"`
		// 	Name                 string     `db:"name"`
		// 	Descr                string     `db:"descr"`
		// 	StockCode            string     `db:"stock_code"`
		// 	ReorderStocklevel    float64    `db:"reorder_stocklevel"`
		// 	ReorderQty           float64    `db:"reorder_qty"`
		// 	LatestPrice          float64    `db:"latest_price"`
		// 	LastPriceDate        *time.Time `db:"last_price_date"`
		// 	LastPriceDateDisplay string     `db:"last_price_date_display"`
		// 	CurrentStock         float64    `db:"current_stock"`
		// 	ValuationString      string     `db:"valuation_string"`
		// 	Valuation            float64    `db:"valuation"`
		// 	QtyType              string     `db:"qty_type"`
		// 	Picture              string     `db:"picture"`
		// 	Notes                string     `db:"notes"`

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.Render("parts-tree", "main", tree)

		// Fill in the custom field
		w := dom.GetWindow()
		doc := w.Document()

		t := doc.QuerySelector(`[name="tree"]`)
		t.SetInnerHTML("") // Init the tree panel

		tlcBtn := doc.CreateElement("input").(*dom.HTMLInputElement)
		tlcBtn.Type = "button"
		tlcBtn.Class().Add("button-primary")
		tlcBtn.Value = "New Category"
		tlcBtn.AddEventListener("click", false, func(evt dom.Event) {
			print("clicked on the add base level category btn", tree)

			go func() {

				// Create a fresh new parts category at the backend
				newCatID := 0
				rpcClient.Call("PartRPC.AddCategory", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      0,
				}, &newCatID)

				// Read the new category back
				newCat := shared.Category{}
				rpcClient.Call("PartRPC.GetCategory", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      newCatID,
				}, &newCat)

				// Append the new category to the tree widget
				tree = append(tree, newCat)

				// Manually create a new entry in the DOM for the newly created category
				ul := doc.QuerySelector(".treeview").(*dom.HTMLUListElement)
				widgetID := fmt.Sprintf("category-%d", newCatID)
				li := doc.CreateElement("li")
				li.SetID(widgetID)
				chek := doc.CreateElement("input").(*dom.HTMLInputElement)
				// chek.Class().Add("cheka")
				chek.Type = "checkbox"
				li.AppendChild(chek)
				label := doc.CreateElement("label")
				label.SetAttribute("for", widgetID)
				label.SetInnerHTML(newCat.Name)
				label.SetAttribute("data-type", "category")
				label.SetAttribute("data-id", fmt.Sprintf("%d", newCatID))
				// label.Class().Add("category")
				label.SetID(widgetID + "-label")
				chek.SetAttribute("data-type", "category")
				chek.SetAttribute("data-id", fmt.Sprintf("%d", newCatID))
				chek.SetID(widgetID + "-chek")
				li.AppendChild(label)
				ul.AppendChild(li)

				// add an empty list for the categories
				ul2 := doc.CreateElement("ul").(*dom.HTMLUListElement)
				li.AppendChild(ul2)

				// add an empty list for the parts
				ul3 := doc.CreateElement("ul")
				li.AppendChild(ul3)
				// ul.AppendChild(ul3)
				// li3 := doc.CreateElement("li")
				// li3.SetInnerHTML("(no parts)")
				// ul3.AppendChild(li3)

				currentCat = newCatID
			}()
		})

		// Create the Tree's UL element
		ul := doc.CreateElement("ul").(*dom.HTMLUListElement)
		ul.SetClass("treeview")

		// Recursively add elements to the tree
		addTree(tree, ul, 0)

		t.AppendChild(ul)

		t.AppendChild(tlcBtn)

		// Handlers for the various buttons
		btnAddPart := doc.QuerySelector(`[name=AddPart]`)
		btnAddCat := doc.QuerySelector(`[name=AddCat]`)
		btnDelCat := doc.QuerySelector(`[name=DelCat]`)
		btnSaveCat := doc.QuerySelector(`[name=SaveCat]`)
		btnSavePart := doc.QuerySelector(`[name=SavePart]`)
		btnDelPart := doc.QuerySelector(`[name=DelPart]`)

		btnAddCat.AddEventListener("click", false, func(evt dom.Event) {
			go func() {
				newCat := 0
				rpcClient.Call("PartRPC.AddCategory", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      currentCat,
				}, &newCat)
				// print("Add category ", newCat, "to current cat", currentCat)

				// Find the UL element for the current Cat, and add a new LI to it !
				theLI := doc.QuerySelector(fmt.Sprintf("#category-%d", currentCat)).(*dom.HTMLLIElement)
				// print("got ", theLI)
				// theUL := theLI.LastChild().(*dom.HTMLUListElement)
				cNodes := theLI.ChildNodes()
				print("cnodes =", cNodes)
				theUL := theLI.ChildNodes()[2].(*dom.HTMLUListElement)

				print("got ", theLI, theUL)

				widgetID := fmt.Sprintf("category-%d", newCat)
				li := doc.CreateElement("li")
				li.SetID(widgetID)
				chek := doc.CreateElement("input").(*dom.HTMLInputElement)
				chek.Type = "checkbox"
				li.AppendChild(chek)
				label := doc.CreateElement("label")
				label.SetAttribute("for", widgetID)
				label.SetInnerHTML("New Category")
				label.SetAttribute("data-type", "category")
				label.SetAttribute("data-id", fmt.Sprintf("%d", newCat))
				// label.Class().Add("category")
				label.SetID(widgetID + "-label")
				chek.SetAttribute("data-type", "category")
				chek.SetAttribute("data-id", fmt.Sprintf("%d", newCat))
				chek.SetID(widgetID + "-chek")
				li.AppendChild(label)
				theUL.AppendChild(li)

				// add an empty list for the categories
				ul2 := doc.CreateElement("ul").(*dom.HTMLUListElement)
				li.AppendChild(ul2)

				// Add an empty list for the parts
				ul3 := doc.CreateElement("ul")
				li.AppendChild(ul3)
				// theUL.AppendChild(ul3)
				// li3 := doc.CreateElement("li")
				// li3.SetInnerHTML("(no parts)")
				// ul3.AppendChild(li3)
			}()
		})

		btnAddPart.AddEventListener("click", false, func(evt dom.Event) {
			go func() {
				newPart := shared.Part{}
				rpcClient.Call("PartRPC.AddPart", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      currentCat,
				}, &newPart)
				// print("Add part ", newPart.ID, " to current cat", currentCat)

				// Find the UL element for the current Cat, and add a new LI to it !
				theLI := doc.QuerySelector(fmt.Sprintf("#category-%d", currentCat)).(*dom.HTMLLIElement)
				// print("got ", theLI)
				// theUL := theLI.LastChild().(*dom.HTMLUListElement)
				nodes := theLI.ChildNodes()
				print("nodes -", nodes)
				theUL := theLI.ChildNodes()[2].(*dom.HTMLUListElement)

				// print("got ", theLI, theUL)

				partID := fmt.Sprintf("part-%d", newPart.ID)
				li2 := doc.CreateElement("li")
				li2.SetID(partID)
				li2.SetInnerHTML(newPart.StockCode)
				li2.Class().Add("stock-item")
				li2.SetAttribute("data-type", "part")
				li2.SetAttribute("data-id", fmt.Sprintf("%d", newPart.ID))
				theUL.AppendChild(li2)

			}()
		})

		btnSavePart.AddEventListener("click", false, func(evt dom.Event) {
			go func() {
				// Read the part from the DOM
				thePart := shared.Part{}
				rpcClient.Call("PartRPC.Get", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      currentPart,
				}, &thePart)

				// print("The original Part ", thePart)
				swapper.Panels[1].Bind(&thePart)
				// print("The modified Part ", thePart)

				// Update the database
				done := false
				rpcClient.Call("PartRPC.Update", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      currentPart,
					Part:    &thePart,
				}, &done)

				// Find the LI element for the current Part, and redraw it
				theLI := doc.QuerySelector(fmt.Sprintf("#part-%d", currentPart)).(*dom.HTMLLIElement)
				// print("got ", theLI)
				theLI.SetInnerHTML(fmt.Sprintf(`%s : %s`, thePart.StockCode, thePart.Name))
				thePart.ValuationString = thePart.DisplayValuation()
				swapper.Panels[1].Paint(&thePart)
			}()
		})

		btnSaveCat.AddEventListener("click", false, func(evt dom.Event) {

			go func() {
				// Read the part from the DOM
				theCat := shared.Category{}
				rpcClient.Call("PartRPC.GetCategory", shared.PartRPCData{
					Channel: Session.Channel,
					ID:      currentCat,
				}, &theCat)

				print("The original Cat ", theCat)
				theCat.Name = doc.QuerySelector("[name=CatName]").(*dom.HTMLInputElement).Value
				theCat.Descr = doc.QuerySelector("[name=CatDescr]").(*dom.HTMLInputElement).Value
				theCat.StockCode = doc.QuerySelector("[name=CatStockCode]").(*dom.HTMLInputElement).Value
				print("The modified Cat ", theCat)

				// Update the database
				done := false
				rpcClient.Call("PartRPC.UpdateCategory", shared.CategoryRPCData{
					Channel:  Session.Channel,
					ID:       currentCat,
					Category: &theCat,
				}, &done)

				// Find the LI element for the current Part, and redraw it
				theLabel := doc.QuerySelector(fmt.Sprintf("#category-%d-label", currentCat)).(*dom.HTMLLabelElement)
				// print("got ", theLI)
				theLabel.SetInnerHTML(theCat.Name)
			}()
		})

		btnDelCat.AddEventListener("click", false, func(evt dom.Event) {
			doIt := js.Global.Call("confirm", "Delete this category ?")
			print("doIt", doIt)
			if doIt {

				go func() {
					parentCat := 0
					rpcClient.Call("PartRPC.DelCategory", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      currentCat,
					}, &parentCat)
					// print("Del Cat ", currentCat)

					// Find the LI element for the current Cat, and remove it
					theLI := doc.QuerySelector(fmt.Sprintf("#category-%d", currentCat)).(*dom.HTMLLIElement)
					// print("got ", theLI)

					theLI.ParentNode().RemoveChild(theLI)
					if parentCat != 0 {
						theCat := shared.Category{}
						rpcClient.Call("PartRPC.GetCategory", shared.PartRPCData{
							Channel: Session.Channel,
							ID:      parentCat,
						}, &theCat)
						// print("Cat", dataID, theCat)
						currentCat = theCat.ID
						doc.QuerySelector("[name=CatName]").(*dom.HTMLInputElement).Value = theCat.Name
						doc.QuerySelector("[name=CatDescr]").(*dom.HTMLInputElement).Value = theCat.Descr
						doc.QuerySelector("[name=CatStockCode]").(*dom.HTMLInputElement).Value = theCat.StockCode
						swapper.Select(0)
					}
				}()
			}
			print("Delete current cat", currentCat)
		})

		btnDelPart.AddEventListener("click", false, func(evt dom.Event) {
			doIt := js.Global.Call("confirm", "Delete this part ?")
			print("doIt", doIt)
			if doIt {

				go func() {
					rpcClient.Call("PartRPC.DelPart", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      currentPart,
					}, &currentCat)
					// print("Del part ", currentPart)

					// Find the LI element for the current Part, and remove it
					theLI := doc.QuerySelector(fmt.Sprintf("#part-%d", currentPart)).(*dom.HTMLLIElement)
					// print("got ", theLI)

					theLI.ParentNode().RemoveChild(theLI)
					theCat := shared.Category{}
					rpcClient.Call("PartRPC.GetCategory", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      currentCat,
					}, &theCat)
					// print("Cat", dataID, theCat)
					currentCat = theCat.ID
					doc.QuerySelector("[name=CatName]").(*dom.HTMLInputElement).Value = theCat.Name
					doc.QuerySelector("[name=CatDescr]").(*dom.HTMLInputElement).Value = theCat.Descr
					doc.QuerySelector("[name=CatStockCode]").(*dom.HTMLInputElement).Value = theCat.StockCode
					swapper.Select(0)
				}()
				print("Delete current part", currentPart)
			}
		})

		// Add functions on the tree
		// Handlers on the table itself
		ul.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			li := evt.Target()
			print("click event on the list", evt, "with tag", li.TagName())

			switch li.TagName() {
			case "LI":
				print("This could be a part or category, or a whole line")
				theType := li.GetAttribute("data-type")
				print("data-type", theType)
				switch theType {
				case "category", "part":
					print("valid LI, proceed")
				default:
					// check that it has an ID
					if li.ID() == "" {
						print("Clicked on some empty line - ignore the click")
						return
					}
					print("valid category header .. get the matching label and work from there")
					li = li.FirstChild().(dom.Element)
					print("li has morphed into", li)
				}
			case "LABEL":
				print("This must be a category")
			case "INPUT":
				print("clicking on checkboxes in the tree is totally broken at the moment due to CSS weirdness, so eat the event and do nothing for now")
				return
				print("Lets toggle the input for now")
				theInput := li.(*dom.HTMLInputElement)
				print("theInput", theInput)
				theInput.Checked = !theInput.Checked
				print("Clicked on the input, so lets find the label")
				li = li.ParentElement().FirstChild().(dom.Element)
				print("li has morphed into", li)
			default:
				print("dont know what to do about that object type - do nothing")
				return
			}

			if lastSelectedClass != nil {
				lastSelectedClass.Remove("listselected")
			}
			lastSelectedClass = li.Class()
			// print("LI class =", lastSelectedClass)

			// only add this if the target is a LI, otherwise it makes no sense to do this
			lastSelectedClass.Add("listselected")

			dataType := li.GetAttribute("data-type")
			dataID := li.GetAttribute("data-id")
			actualID, _ := strconv.Atoi(dataID)
			switch dataType {
			case "category":
				go func() {
					theCat := shared.Category{}
					rpcClient.Call("PartRPC.GetCategory", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      actualID,
					}, &theCat)
					// print("Cat", dataID, theCat)
					currentCat = theCat.ID
					doc.QuerySelector("[name=CatName]").(*dom.HTMLInputElement).Value = theCat.Name
					doc.QuerySelector("[name=CatDescr]").(*dom.HTMLInputElement).Value = theCat.Descr
					doc.QuerySelector("[name=CatStockCode]").(*dom.HTMLInputElement).Value = theCat.StockCode

					if theCat.NumParts > 0 || theCat.NumSubcats > 0 {
						// turn off the delete btn
						// print("hide the del btn")
						btnDelCat.Class().Add("hidden")
					} else {
						// turn on the delete btn
						btnDelCat.Class().Remove("hidden")
					}

					print("expand out the cat", currentCat)
					theCheka := doc.QuerySelector(fmt.Sprintf("#category-%d-chek", currentCat)).(*dom.HTMLInputElement)
					theCheka.Checked = !theCheka.Checked

					swapper.Select(0)
					doc.QuerySelector(`[name=CatName]`).(*dom.HTMLInputElement).Focus()
				}()
				// print("Category", dataID)
				swapper.Select(0)
			case "part":
				go func() {
					thePart := shared.Part{}
					rpcClient.Call("PartRPC.Get", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      actualID,
					}, &thePart)
					print("Part", dataID, thePart)
					currentPart = thePart.ID
					thePart.ValuationString = thePart.DisplayValuation()
					swapper.Panels[1].Paint(&thePart)
					doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = thePart.ValuationString
					swapper.Select(1)
					doc.QuerySelector(`[name=Name]`).(*dom.HTMLInputElement).Focus()

				}()
			}
		})

		// Auto calculate the valuation on change of fields
		doc.QuerySelector("[name=CurrentStock]").AddEventListener("change", false, func(evt dom.Event) {
			thePart := shared.Part{}
			print("current stock count changes")
			s := doc.QuerySelector("[name=CurrentStock]").(*dom.HTMLInputElement).Value
			p := doc.QuerySelector("[name=LatestPrice]").(*dom.HTMLInputElement).Value
			s1, _ := strconv.ParseFloat(s, 64)
			p1, _ := strconv.ParseFloat(p, 64)
			thePart.CurrentStock = s1
			thePart.LatestPrice = p1
			thePart.ValuationString = thePart.DisplayValuation()
			doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = thePart.ValuationString
		})
		doc.QuerySelector("[name=LatestPrice]").AddEventListener("change", false, func(evt dom.Event) {
			thePart := shared.Part{}
			print("latest price changes")
			s := doc.QuerySelector("[name=CurrentStock]").(*dom.HTMLInputElement).Value
			p := doc.QuerySelector("[name=LatestPrice]").(*dom.HTMLInputElement).Value
			s1, _ := strconv.ParseFloat(s, 64)
			p1, _ := strconv.ParseFloat(p, 64)
			thePart.CurrentStock = s1
			thePart.LatestPrice = p1
			thePart.ValuationString = thePart.DisplayValuation()
			doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = thePart.ValuationString
		})

	}()
}

func classAdd(context *router.Context) {

	go func() {
		partClass := shared.PartClass{}
		BackURL := "/class/select"
		title := "Add Machine Type for Parts List"
		form := formulate.EditForm{}
		form.New("fa-puzzle-piece", title)

		// Layout the fields

		form.Row(2).
			AddInput(1, "Name", "Name")

		form.Row(1).
			AddInput(1, "Description", "Descr")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&partClass)
			go func() {
				newID := 0
				rpcClient.Call("PartRPC.InsertClass", shared.PartClassRPCData{
					Channel:   Session.Channel,
					PartClass: &partClass,
				}, &newID)
				print("added class ID", newID)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &partClass)

	}()

}

// Show a list of all parts for the given class
func partList(context *router.Context) {

	partClass, _ := strconv.Atoi(context.Params["id"])
	// print("show parts of class", partClass)

	go func() {
		data := []shared.Part{}
		class := shared.PartClass{}
		rpcClient.Call("PartRPC.List", shared.PartRPCData{
			Channel: Session.Channel,
			ID:      partClass,
		}, &data)
		rpcClient.Call("PartRPC.GetClass", shared.PartClassRPCData{
			Channel: Session.Channel,
			ID:      partClass,
		}, &class)

		tree := []shared.Category{}
		rpcClient.Call("PartRPC.GetTree", shared.PartTreeRPCData{
			Channel:    Session.Channel,
			CategoryID: 76,
		}, &tree)
		print("got tree", tree)

		for i, t := range tree {
			print("tree", i, t)
			for i, p := range t.Parts {
				print("  part", i, p)
			}
			for i, c := range t.Subcats {
				print("  subcat", i, c)
			}
		}

		BackURL := "/class/select"
		Title := fmt.Sprintf("Parts of type - %s", class.Name)

		// load a form for the class
		if partClass == 0 {
			loadTemplate("class-display", "main", class)
		} else {
			loadTemplate("class-edit", "main", class)
			w := dom.GetWindow()
			doc := w.Document()

			if el := doc.QuerySelector(".data-del-btn"); el != nil {

				if el := doc.QuerySelector(".md-confirm-del"); el != nil {
					el.AddEventListener("click", false, func(evt dom.Event) {
						go func() {
							done := false
							rpcClient.Call("PartRPC.DeleteClass", shared.PartClassRPCData{
								Channel:   Session.Channel,
								PartClass: &class,
							}, &done)
						}()
						Session.Navigate(BackURL)
					})
				}

				el.AddEventListener("click", false, func(evt dom.Event) {
					doc.QuerySelector("#confirm-delete").Class().Add("md-show")
				})

				if el := doc.QuerySelector(".md-close-del"); el != nil {
					el.AddEventListener("click", false, func(evt dom.Event) {
						doc.QuerySelector("#confirm-delete").Class().Remove("md-show")
					})
				}

				if el := doc.QuerySelector("#confirm-delete"); el != nil {
					el.AddEventListener("keyup", false, func(evt dom.Event) {
						if evt.(*dom.KeyboardEvent).KeyCode == 27 {
							evt.PreventDefault()
							doc.QuerySelector("#confirm-delete").Class().Remove("md-show")
						}
					})
				}
			}
		}

		form := formulate.ListForm{}
		form.New("fa-puzzle-piece", Title)

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Description", "Descr")
		form.Column("Stock Code", "StockCode")
		form.Column("Reorder Lvl/Qty", "ReorderDetails")
		form.Column("Stock", "CurrentStock")
		form.Column("Qty", "QtyType")
		form.Column("Latest Price", "DisplayPrice")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		if Session.UserRole == "Admin" {
			form.NewRowEvent(func(evt dom.Event) {
				evt.PreventDefault()
				Session.Navigate(fmt.Sprintf("/part/add/%d", class.ID))
			})
		}

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/part/" + key)
		})

		form.Render("parts-list", "#parts-list-goes-here", data)

		// Add an onChange callback to the class edit fields
		w := dom.GetWindow()
		doc := w.Document()

		doc.QuerySelector("#class-name").AddEventListener("change", false, func(evt dom.Event) {
			print("TODO - Name has changed")
			go func() {
				class.Name = doc.QuerySelector("#class-name").(*dom.HTMLInputElement).Value
				done := false
				rpcClient.Call("PartRPC.UpdateClass", shared.PartClassRPCData{
					Channel:   Session.Channel,
					PartClass: &class,
				}, &done)
			}()
		})
		doc.QuerySelector("#class-descr").AddEventListener("change", false, func(evt dom.Event) {
			print("TODO - Description has changed")
			go func() {
				class.Descr = doc.QuerySelector("#class-descr").(*dom.HTMLInputElement).Value
				done := false
				rpcClient.Call("PartRPC.UpdateClass", shared.PartClassRPCData{
					Channel:   Session.Channel,
					PartClass: &class,
				}, &done)
			}()
		})

	}()
}

func partEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		part := shared.Part{}
		classes := []shared.PartClass{}
		stocks := []shared.PartStock{}
		prices := []shared.PartPrice{}
		data := shared.PartRPCData{
			Channel: Session.Channel,
			ID:      id,
		}

		rpcClient.Call("PartRPC.Get", data, &part)
		rpcClient.Call("PartRPC.ClassList", data, &classes)
		rpcClient.Call("PartRPC.StockList", data, &stocks)
		rpcClient.Call("PartRPC.PriceList", data, &prices)

		BackURL := fmt.Sprintf("/parts/%d", part.Class)
		title := fmt.Sprintf("Part Details - %s - %s", part.Name, part.StockCode)
		form := formulate.EditForm{}
		form.New("fa-puzzle-piece", title)

		// convert the last_price_date into a display field
		part.LastPriceDateDisplay = ""
		if part.LastPriceDate != nil {
			part.LastPriceDateDisplay = part.LastPriceDate.Format("Mon, Jan 2 2006")
		}
		part.ValuationString = part.DisplayValuation()

		// Layout the fields

		form.Row(1).
			AddSelect(1, "For Machine Type", "Class", classes, "ID", "Name", 1, part.Class)

		form.Row(2).
			AddInput(1, "Name", "Name").
			AddInput(1, "Stock Code", "StockCode")

		form.Row(1).
			AddInput(1, "Description", "Descr")

		form.Row(4).
			AddDecimal(1, "ReOrder Level", "ReorderStocklevel", 2, "1").
			AddDecimal(1, "ReOrder Qty", "ReorderQty", 2, "1").
			AddDecimal(1, "Current Stock", "CurrentStock", 2, "1").
			AddInput(1, "Qty Type", "QtyType")

		form.Row(4).
			AddDisplay(2, "Last Price Update", "LastPriceDateDisplay").
			AddDecimal(1, "Latest Price", "LatestPrice", 2, "1").
			AddDisplay(1, "Valuation", "ValuationString")

		form.Row(1).
			AddTextarea(1, "Notes", "Notes")

		form.Row(1).
			AddCustom(1, "Stock History", "StockList", "")

		form.Row(1).
			AddCustom(1, "Price History", "PriceList", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				done := false
				rpcClient.Call("PartRPC.Delete", shared.PartRPCData{
					Channel: Session.Channel,
					Part:    &part,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&part)
			go func() {
				done := false
				rpcClient.Call("PartRPC.Update", shared.PartRPCData{
					Channel: Session.Channel,
					Part:    &part,
				}, &done)
				NewBackURL := ""
				if done {
					// Go back to parts list
					NewBackURL = fmt.Sprintf("/parts/%d", part.Class)
				} else {
					// refresh this screen
					NewBackURL = fmt.Sprintf("/part/%d", part.ID)
				}
				Session.Navigate(NewBackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &part)

		// Inject the StockLevel list
		stocklist := formulate.ListForm{}
		stocklist.New("", "")
		stocklist.ColumnFormat("Date", "DateFromDisplay", `width="30%"`)
		stocklist.ColumnFormat("Description", "Descr", `width="50%" text-align="right"`)
		stocklist.ColumnFormat("Stock", "StockLevel", `width="20%" text-align="right"`)
		stocklist.Render("part-stock-list", "[name=StockList]", stocks)

		// Inject the Price list
		pricelist := formulate.ListForm{}
		pricelist.New("", "")
		pricelist.ColumnFormat("Date", "DateFromDisplay", `width="30%"`)
		pricelist.ColumnFormat("Description", "Descr", `width="30%"`)
		pricelist.ColumnFormat("Price", "PriceDisplay", `width="20%" text-align="right"`)
		pricelist.Render("part-price-list", "[name=PriceList]", prices)

		// Auto calculate the valuation on change of fields
		w := dom.GetWindow()
		doc := w.Document()
		doc.QuerySelector("[name=CurrentStock]").AddEventListener("change", false, func(evt dom.Event) {
			s := doc.QuerySelector("[name=CurrentStock]").(*dom.HTMLInputElement).Value
			p := doc.QuerySelector("[name=LatestPrice]").(*dom.HTMLInputElement).Value
			s1, _ := strconv.ParseFloat(s, 64)
			p1, _ := strconv.ParseFloat(p, 64)
			part.CurrentStock = s1
			part.LatestPrice = p1
			part.ValuationString = part.DisplayValuation()
			doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = part.ValuationString
		})
		doc.QuerySelector("[name=LatestPrice]").AddEventListener("change", false, func(evt dom.Event) {
			s := doc.QuerySelector("[name=CurrentStock]").(*dom.HTMLInputElement).Value
			p := doc.QuerySelector("[name=LatestPrice]").(*dom.HTMLInputElement).Value
			s1, _ := strconv.ParseFloat(s, 64)
			p1, _ := strconv.ParseFloat(p, 64)
			part.CurrentStock = s1
			part.LatestPrice = p1
			part.ValuationString = part.DisplayValuation()
			doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = part.ValuationString
		})

		// // And attach actions
		// form.ActionGrid("part-actions", "#action-grid", part.ID, func(url string) {
		// 	print("clicked on url", url)
		// })

	}()
}

func partAdd(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		part := shared.Part{}
		part.Class = id
		classes := []shared.PartClass{}
		class := shared.PartClass{}
		rpcClient.Call("PartRPC.GetClass", shared.PartClassRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &class)
		rpcClient.Call("PartRPC.ClassList", Session.Channel, &classes)

		BackURL := fmt.Sprintf("/parts/%d", part.Class)
		title := fmt.Sprintf("Add Part for Machine Type - %s - %s", class.Name, class.Descr)
		form := formulate.EditForm{}
		form.New("fa-puzzle-piece", title)

		// Layout the fields

		form.Row(1).
			AddSelect(1, "For Machine Type", "Class", classes, "ID", "Name", 1, part.Class)

		form.Row(2).
			AddInput(1, "Name", "Name").
			AddInput(1, "Stock Code", "StockCode")

		form.Row(1).
			AddInput(1, "Description", "Descr")

		form.Row(3).
			AddDecimal(1, "ReOrder Level", "ReorderStocklevel", 2, "1").
			AddDecimal(1, "ReOrder Qty", "ReorderQty", 2, "1").
			AddInput(1, "Qty Type", "QtyType")

		form.Row(2).
			AddDecimal(1, "Latest Price", "LatestPrice", 2, "1").
			AddDecimal(1, "Current Stock", "CurrentStock", 2, "1")

		form.Row(1).
			AddTextarea(1, "Notes", "Notes")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&part)
			go func() {
				newID := 0
				rpcClient.Call("PartRPC.Insert", shared.PartRPCData{
					Channel: Session.Channel,
					Part:    &part,
				}, &newID)
				print("Added new part", newID)
				NewBackURL := fmt.Sprintf("/parts/%d", part.Class)
				Session.Navigate(NewBackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &part)

	}()
}

func partPriceList(context *router.Context) {
	print("TODO - partPriceList")
}

func partStockList(context *router.Context) {
	print("TODO - partStockList")
}

// No longer required functions

// Show a list of machine classes, select one to show the parts for that class
func classSelect(context *router.Context) {

	go func() {
		data := []shared.PartClass{}
		rpcClient.Call("PartRPC.ClassList", Session.Channel, &data)

		BackURL := "/"

		form := formulate.ListForm{}
		form.New("fa-puzzle-piece", "Select Machine Type for Parts List")

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Description", "Descr")
		form.Column("Number of Parts", "Count")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		if Session.UserRole == "Admin" {
			form.NewRowEvent(func(evt dom.Event) {
				evt.PreventDefault()
				Session.Navigate("/class/add")
			})
		}

		form.RowEvent(func(key string) {
			Session.Navigate("/parts/" + key)
		})

		form.Render("class-select", "main", data)

	}()
}
