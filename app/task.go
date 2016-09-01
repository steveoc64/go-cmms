package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"itrak-cmms/shared"

	"github.com/gopherjs/gopherjs/js"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

// Add in a pile of constants for lists and things
const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

func taskMaint(context *router.Context) {
	print("TODO - taskmaint")
}

func calcAllDone(task shared.Task) bool {
	if task.CompletedDate != nil {
		return false
	}

	if task.LabourHrs == 0.0 {
		print("Needs some labour input")
		return false
	}
	for _, v := range task.Checks {
		if !v.Done {
			print("Task has incomplete checklist items")
			return false
		}
	}

	print("Task appears to be complete")
	return true
}

func setCheckboxes(task shared.Task) {
	w := dom.GetWindow()
	doc := w.Document()

	for _, v := range task.Checks {
		tc := fmt.Sprintf("#taskchek-%d", v.Seq)
		if c := doc.QuerySelector(tc); c != nil {
			cb := c.(*dom.HTMLInputElement)
			cb.Checked = v.Done
			cb.Disabled = v.Done
		} else {
			print("no such dom element", tc)
		}
	}
}

func taskEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	Session.ID["task"] = id

	Session.Subscribe("task", _taskEdit)
	go _taskEdit("edit", id)
}

func _taskEdit(action string, id int) {

	BackURL := "/tasks"

	switch action {
	case "edit":
		print("manually edit")
	case "delete":
		if id != Session.ID["task"] {
			return
		}
		print("current record has been deleted")
		Session.Navigate(BackURL)
		return
	default:
		if id != Session.ID["task"] {
			return
		}
	}

	task := shared.Task{}

	rpcClient.Call("TaskRPC.Get", shared.TaskRPCData{
		Channel: Session.Channel,
		ID:      id,
	}, &task)

	task.AllDone = calcAllDone(task)
	// print("task with parts and checks attached =", task)

	RefreshURL := fmt.Sprintf("/task/%d", id)
	title := fmt.Sprintf("Task Details - %06d", id)
	form := formulate.EditForm{}
	form.New("fa-server", title)

	if task.StartDate != nil {
		task.DisplayStartDate = task.StartDate.Format("Mon, Jan 2 2006")
	}
	if task.DueDate != nil {
		task.DisplayDueDate = task.DueDate.Format("Mon, Jan 2 2006")
	}
	if task.Username == nil {
		task.DisplayUsername = "Unassigned"
	} else {
		task.DisplayUsername = *task.Username
	}

	// Define a function to add the actions and set all the callbacks
	setActions := func(actionID int) {

		switch Session.UserRole {
		case "Admin", "Site Manager":
			form.ActionGrid("task-admin-actions", "#action-grid", task, func(url string) {
				if strings.HasPrefix(url, "/sched") {
					if task.SchedID != 0 {
						Session.Navigate(fmt.Sprintf("/sched/%d", task.SchedID))
						// 	c := &router.Context{
						// 		Path:        url,
						// 		InitialLoad: false,
						// 		Params: map[string]string{
						// 			"id":   fmt.Sprintf("%d", task.SchedID),
						// 			"back": fmt.Sprintf("/task/%d", task.ID),
						// 		},
						// 	}
						// 	schedEdit(c)
					}
				} else if strings.HasPrefix(url, "/stoppage") {
					if task.EventID != 0 {
						c := &router.Context{
							Path:        url,
							InitialLoad: false,
							Params: map[string]string{
								"id":   fmt.Sprintf("%d", task.EventID),
								"back": fmt.Sprintf("/task/%d", task.ID),
							},
						}
						stoppageEdit(c)
					}
				} else if strings.HasPrefix(url, "/retransmit") {
					print("re-transmit the SMS for this one")
					go func() {
						result := ""
						rpcClient.Call("TaskRPC.Retransmit", shared.TaskRPCData{
							Channel: Session.Channel,
							Task:    &task,
						}, &result)
						print(result)
						js.Global.Call("alert", result)
					}()
				} else {
					go func() {
						done := false
						rpcClient.Call("TaskRPC.Complete", shared.TaskRPCData{
							Channel: Session.Channel,
							Task:    &task,
						}, &done)
						Session.Navigate(BackURL)
					}()
				}
			})
		case "Technician":
			form.ActionGrid("task-actions", "#action-grid", task, func(url string) {
				go func() {
					done := false
					print("calling task.complete ???")
					rpcClient.Call("TaskRPC.Complete", shared.TaskRPCData{
						Channel: Session.Channel,
						Task:    &task,
					}, &done)
					Session.Navigate(BackURL)
				}()
			})
		}
	}

	// Layout the fields
	partsTitle := "Parts Used - record qty for each part used"

	useRole := Session.UserRole
	if Session.CanAllocate {
		useRole = "Admin"
	}

	switch useRole {
	case "Admin":

		techs := []shared.User{}
		rpcClient.Call("UserRPC.GetTechnicians", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      0,
		}, &techs)
		assignedTo := 0
		if task.AssignedTo != nil {
			assignedTo = *task.AssignedTo
		}

		form.Row(3).
			AddSelect(1, "User", "AssignedTo", techs, "ID", "Username", 0, assignedTo).
			AddDisplay(1, "Start Date", "DisplayStartDate").
			AddDisplay(1, "Due Date", "DisplayDueDate")

		form.Row(3).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName").
			AddDisplay(1, "Component", "Component")

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Attachments", "Photos", "")
			// AddPreview(1, "", "StoppagePreview").
			// AddPreview(1, "", "Preview1").
			// AddPreview(1, "", "Preview2").
			// AddPreview(1, "", "Preview3")

		form.Row(1).
			AddTextarea(1, "Notes", "Log")

		form.Row(1).
			AddCustom(1, "Notes and CheckLists", "CheckList", "")

		rateStr := ""
		if task.LabourHrs != 0 {
			rate := int(task.LabourCost / task.LabourHrs)
			rateStr = fmt.Sprintf(" @ $%d/hr", rate)
		}

		form.Row(4).
			AddDecimal(1, fmt.Sprintf("Hours%s", rateStr), "LabourHrs", 2, "0.5").
			// AddDecimal(1, "Actual Labour $", "LabourCost", 2, "1").
			AddDisplay(1, "Actual Labour $", "LabourCost").
			// AddDecimal(1, "Actual Material $", "MaterialCost", 2, "1")
			AddDisplay(1, "Actual Material $", "MaterialCost").
			AddDisplay(1, "Other Costs", "OtherCost")

		// form.Row(4).
		// 	AddDisplay(2, "Labour Est $", "LabourEst").
		// 	AddDecimal(1, fmt.Sprintf("Hours%s", rateStr), "LabourHrs", 2, "0.5").
		// 	AddDecimal(1, "Actual Labour $", "LabourCost", 2, "1")

		// form.Row(4).
		// 	AddDisplay(2, "Material Est $", "MaterialEst").
		// 	AddDecimal(2, "Actual Material $", "MaterialCost", 2, "1")

	case "Site Manager":
		form.Row(3).
			AddDisplay(1, "User", "DisplayUsername").
			AddDisplay(1, "Start Date", "DisplayStartDate").
			AddDisplay(1, "Due Date", "DisplayDueDate")

		form.Row(3).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName").
			AddDisplay(1, "Component", "Component")

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Attachments", "Photos", "")
			// AddPreview(1, "", "StoppagePreview").
			// AddPreview(1, "", "Preview1").
			// AddPreview(1, "", "Preview2").
			// AddPreview(1, "", "Preview3")

		form.Row(1).
			AddDisplayArea(1, "Notes", "Log")

		form.Row(1).
			AddCustom(1, "Notes and CheckLists", "CheckList", "")

		form.Row(4).
			AddDecimal(1, "Hours", "LabourHrs", 2, "0.5").
			// AddDecimal(1, "Actual Labour $", "LabourCost", 2, "1").
			AddDisplay(1, "Actual Labour $", "LabourCost").
			// AddDecimal(1, "Actual Material $", "MaterialCost", 2, "1")
			AddDisplay(1, "Actual Material $", "MaterialCost").
			AddDisplay(1, "Other Costs", "OtherCost")

		// form.Row(4).
		// 	AddDisplay(2, "Labour Est $", "LabourEst").
		// 	AddDecimal(1, "Hours", "LabourHrs", 2, "0.5").
		// 	AddDecimal(1, "Actual Labour $", "LabourCost", 2, "1")

		// form.Row(4).
		// 	AddDisplay(2, "Material Est $", "MaterialEst").
		// 	AddDecimal(2, "Actual Material $", "MaterialCost", 2, "1")

	case "Technician":
		row := form.Row(5).
			AddDisplay(1, "Start Date", "DisplayStartDate").
			AddDisplay(1, "Due Date", "DisplayDueDate").
			// AddDecimal(1, "Actual Material $", "MaterialCost", 2, "1").
			// AddDecimal(1, "Actual Labour $", "LabourCost", 2, "1").
			AddDisplay(1, "Actual Material $", "MaterialCost").
			AddDisplay(1, "Other Costs $", "OtherCost")
		// AddDisplay(1, "Actual Labour $", "LabourCost")

		if task.CompletedDate == nil {
			row.AddDecimal(1, "Hours", "LabourHrs", 2, "0.5")
		} else {
			row.AddDisplay(1, "Hours", "LabourHrs")
		}

		form.Row(3).
			AddDisplay(1, "Site", "SiteName").
			AddDisplay(1, "Machine", "MachineName").
			AddDisplay(1, "Component", "Component")

		form.Row(5).
			AddPhoto(1, "Add Photo", "NewPhoto").
			AddCustom(4, "Attachments", "Photos", "")

		if task.CompletedDate == nil {
			form.Row(1).
				AddTextarea(1, "Notes", "Log")
		} else {
			form.Row(1).
				AddDisplayArea(1, "Notes", "Log")
		}

		form.Row(1).
			AddCustom(1, "Notes and CheckLists", "CheckList", "")
	}

	// create the swapper panels
	swapper := formulate.Swapper{
		Name:     "PartDetails",
		Selected: 1,
	}

	form.Row(1).
		AddCustom(1, "Parts Used", "PartsUsed", "hidden")

	form.Row(2).
		AddCustom(1, partsTitle, "PartList", "").
		AddSwapper(1, "Part Details", &swapper)

	catPanel := swapper.AddPanel("Category")
	catPanel.BindWithForm = false

	catPanel.AddRow(1).AddDisplay(1, "Category Name", "CatName")
	catPanel.AddRow(1).AddDisplay(1, "Stock Code", "CatStockCode")
	catPanel.AddRow(1).AddDisplay(1, "Description", "CatDescr")

	// Layout the fields for Parts
	partPanel := swapper.AddPanel("Part")
	partPanel.BindWithForm = false

	partPanel.Row(1).
		AddDecimal(1, "Qty Used", "QtyUsed", 2, "1")

	partPanel.Row(2).
		AddDisplay(1, "Name", "Name").
		AddDisplay(1, "Stock Code", "StockCode")

	partPanel.Row(1).
		AddDisplay(1, "Description", "Descr")

	partPanel.Row(4).
		AddDisplay(1, "ReOrder Level", "ReorderStocklevel").
		AddDisplay(1, "ReOrder Qty", "ReorderQty").
		AddDisplay(1, "Current Stock", "CurrentStock").
		AddDisplay(1, "Qty Type", "QtyType")

	partPanel.Row(1).
		AddDisplay(1, "Supplier Details", "SupplierInfo")

	partPanel.Row(1).
		AddDisplay(1, "Notes", "Notes")

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate(BackURL)
	})

	form.PrintEvent(func(evt dom.Event) {
		dom.GetWindow().Print()
	})

	if useRole == "Admin" {
		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			task.ID = id
			go func() {
				done := false
				rpcClient.Call("TaskRPC.Delete", shared.TaskRPCData{
					Channel: Session.Channel,
					Task:    &task,
				}, &done)
				w := dom.GetWindow()
				w.Alert(fmt.Sprintf("Task %06d Deleted", task.ID))
				Session.Navigate(BackURL)
			}()
		})
	}

	print("useRole =", useRole)
	if useRole == "Admin" ||
		(Session.UserRole == "Technician" && task.CompletedDate == nil) {

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&task)

			if task.NewPhoto.Data != "" {
				showProgress("Updating Task ...")

				// If the uploaded data is a PDF, then use that data instead of the preview
				task.NewPhoto.Data = ImageCache.GetImage()
			}
			go func() {
				updatedTask := shared.Task{}
				rpcClient.Call("TaskRPC.Update", shared.TaskRPCData{
					Channel: Session.Channel,
					Task:    &task,
				}, &updatedTask)
				Session.Navigate(RefreshURL)
				hideProgress()
			}()
		})

	}

	// All done, so render the form
	form.Render("edit-form", "main", &task)
	setPhotoField("NewPhoto")
	showPartsButtons(id)
	showTaskPhotos(task)

	w := dom.GetWindow()
	doc := w.Document()

	print("rendering the checklist here ...")
	renderMarkup(doc.QuerySelector("[name=CheckList]").(*dom.HTMLDivElement), expandHashtags(task.Descr))
	setCheckboxes(task)

	// loadTemplate("task-check-list", "[name=CheckList]", task)
	loadTemplate("task-parts-tree", "[name=PartList]", task)

	// on change of the labour hrs, update the all done flag
	if useRole == "Admin" || task.CompletedDate == nil {
		lh := doc.QuerySelector("[name=LabourHrs]").(*dom.HTMLInputElement)
		lh.AddEventListener("change", false, func(evt dom.Event) {
			wasDone := task.AllDone
			task.LabourHrs, _ = strconv.ParseFloat(lh.Value, 64)
			// fire off the change to the backend
			form.Bind(&task)
			go func() {
				updatedTask := shared.Task{}
				rpcClient.Call("TaskRPC.Update", shared.TaskRPCData{
					Channel: Session.Channel,
					Task:    &task,
				}, &updatedTask)
				if Session.UserRole == "Admin" {
					w := dom.GetWindow()
					doc := w.Document()
					doc.QuerySelector("[name=LabourCost]").(*dom.HTMLInputElement).Value = fmt.Sprintf("%.2f", updatedTask.LabourCost)
					print("updated =", updatedTask.LabourCost)
				}
			}()

			task.AllDone = calcAllDone(task)
			if wasDone != task.AllDone {
				setActions(2)
			}
		})
	}

	var lastSelectedClass *dom.TokenList
	currentCat := 0
	currentPart := 0

	// click on a button in the partsused area
	if el := doc.QuerySelector("[name=PartsUsed]"); el != nil {
		el.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()

			switch evt.Target().TagName() {
			case "INPUT":
				btn := evt.Target().(*dom.HTMLInputElement)
				partID, _ := strconv.Atoi(btn.GetAttribute("part-id"))
				// print("clicked on a btn", btn, partID)

				// Now expand the tree to show the selected item
				li := doc.QuerySelector(fmt.Sprintf("#part-%d", partID))
				// print("got LI", li)

				// now walk up the tree from there, expanding the parent checkboxes
				expandAll(li)

				// deselect the currently selected part
				if lastSelectedClass != nil {
					lastSelectedClass.Remove("listselected")
				}
				currentPartLI := doc.QuerySelector(fmt.Sprintf("#part-%d", currentPart))
				if currentPartLI != nil {
					currentPartLI.Class().Remove("listselected")
				}
				lastSelectedClass = li.Class()
				lastSelectedClass.Add("listselected")
				go func() {
					thePart := shared.Part{}
					rpcClient.Call("PartRPC.Get", shared.PartRPCData{
						Channel: Session.Channel,
						ID:      partID,
					}, &thePart)
					// print("Part", dataID, thePart)
					currentPart = thePart.ID
					swapper.Panels[1].Paint(&thePart)
					// thePart.ValuationString = thePart.DisplayValuation()
					// doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = thePart.ValuationString

					// Now get the Qty used and populate that field
					qtyUsed := 0.0
					rpcClient.Call("TaskRPC.GetQtyUsed", shared.TaskRPCPartData{
						Channel: Session.Channel,
						ID:      id,
						Part:    thePart.ID,
					}, &qtyUsed)
					swapper.Select(1)
					qel := doc.QuerySelector("[name=QtyUsed]").(*dom.HTMLInputElement)
					qel.Value = fmt.Sprintf("%.2f", qtyUsed)
					qel.Focus()
				}()
			}
		})
	}

	print("show parts tree here")

	t := doc.QuerySelector(`[name="parts-tree-div"]`)
	t.SetInnerHTML("") // Init the tree panel

	// Create the Tree's UL element
	ul := doc.CreateElement("ul").(*dom.HTMLUListElement)
	ul.SetClass("treeview")

	// Fetch the complete parts tree from the backend
	go func() {
		tree := []shared.Category{}
		rpcClient.Call("PartRPC.GetTree", shared.PartTreeRPCData{
			Channel:    Session.Channel,
			CategoryID: 0,
		}, &tree)
		// print("got tree", tree)

		// Recursively add elements to the tree
		addTaskPartsTree(tree, ul, 0)

		t.AppendChild(ul)

		// Add a change function on the Qty field
		doc.QuerySelector("[name=QtyUsed]").AddEventListener("change", false, func(evt dom.Event) {
			// print("Change in value of QtyUsed on task", currentPart, id)
			qtyValue := doc.QuerySelector("[name=QtyUsed]").(*dom.HTMLInputElement).Value
			// print("qv", qtyValue)
			qty, _ := strconv.ParseFloat(qtyValue, 64)
			// print("new qty", qty)

			go func() {
				partsUsedUpdate := shared.PartsUsedUpdate{}

				rpcClient.Call("TaskRPC.AddParts", shared.TaskRPCPartData{
					Channel: Session.Channel,
					ID:      id,
					Part:    currentPart,
					Qty:     qty,
				}, &partsUsedUpdate)
				// print("parts used update", partsUsedUpdate)
				// print("new stock level after using", qty, "=", newStockLevel)
				doc.QuerySelector("[name=CurrentStock]").(*dom.HTMLInputElement).Value = fmt.Sprintf("%.2f", partsUsedUpdate.NewStockOnHand)
				doc.QuerySelector("[name=MaterialCost]").(*dom.HTMLInputElement).Value = fmt.Sprintf("%.2f", partsUsedUpdate.TotalMaterialCost)

				// populate the parts button display
				showPartsButtons(id)
			}()

		})

		// Add functions on the tree
		ul.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			li := evt.Target()
			// print("click event on the list", evt, "with tag", li.TagName())

			switch li.TagName() {
			case "LI":
				// print("This could be a part or category, or a whole line")
				theType := li.GetAttribute("data-type")
				// print("data-type", theType)
				switch theType {
				case "category", "part":
					// print("valid LI, proceed")
				default:
					// check that it has an ID
					if li.ID() == "" {
						// print("Clicked on some empty line - ignore the click")
						return
					}
					// print("valid category header .. get the matching label and work from there")
					li = li.FirstChild().(dom.Element)
					// print("li has morphed into", li)
				}
			case "LABEL":
				// print("This must be a category")
			case "INPUT":
				// print("clicking on checkboxes in the tree is totally broken at the moment due to CSS weirdness, so eat the event and do nothing for now")
				return
				// print("Lets toggle the input for now")
				theInput := li.(*dom.HTMLInputElement)
				// print("theInput", theInput)
				theInput.Checked = !theInput.Checked
				// print("Clicked on the input, so lets find the label")
				li = li.ParentElement().FirstChild().(dom.Element)
				// print("li has morphed into", li)
			default:
				// print("dont know what to do about that object type - do nothing")
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

					// print("expand out the cat", currentCat)
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
					// print("Part", dataID, thePart)
					currentPart = thePart.ID
					swapper.Panels[1].Paint(&thePart)
					// thePart.ValuationString = thePart.DisplayValuation()
					// doc.QuerySelector("[name=ValuationString]").(*dom.HTMLInputElement).Value = thePart.ValuationString

					// Now get the Qty used and populate that field
					qtyUsed := 0.0
					rpcClient.Call("TaskRPC.GetQtyUsed", shared.TaskRPCPartData{
						Channel: Session.Channel,
						ID:      id,
						Part:    thePart.ID,
					}, &qtyUsed)
					swapper.Select(1)
					qel := doc.QuerySelector("[name=QtyUsed]").(*dom.HTMLInputElement)
					qel.Value = fmt.Sprintf("%.2f", qtyUsed)
					qel.Focus()

				}()
			}
		})

		// }
	}()

	// click on the parts button, expand the div to show a tree
	if el := doc.QuerySelector("[name=parts-button]"); el != nil {
		el.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			// loadPartsTree()

			// }()
		})
	}

	if el := doc.QuerySelector("[name=CheckList]"); el != nil {

		// First pass - set each checkbox to checked based on whether that item is complete or not

		el.AddEventListener("click", false, func(evt dom.Event) {
			// evt.PreventDefault()
			clickedOn := evt.Target()
			switch clickedOn.TagName() {
			default:
				evt.PreventDefault()
			case "INPUT":
				ie := clickedOn.(*dom.HTMLInputElement)
				key, _ := strconv.Atoi(ie.GetAttribute("key"))
				if key < 1 {
					w.Alert("ERROR: Checkbox has an invalid ID .. please pass this on to IT support, thanks.")
				} else {

					// print("clicked on key", key)

					task.Checks[key-1].Done = true
					now := time.Now()
					task.Checks[key-1].DoneDate = &now
					task.AllDone = calcAllDone(task)

					// save to the backend
					go func() {
						done := false
						rpcClient.Call("TaskRPC.Check", shared.TaskCheckUpdate{
							Channel:   Session.Channel,
							TaskCheck: &task.Checks[key-1],
						}, &done)
					}()
					// loadTemplate("task-check-list", "[name=CheckList]", task)
					setActions(3)
				}
			}

		})
	}

	// And attach actions

	setActions(1)
}

func showTaskPhotos(task shared.Task) {
	// print("populate the photos", task)

	w := dom.GetWindow()
	doc := w.Document()
	div := doc.QuerySelector("[name=Photos]")
	div.SetInnerHTML("")

	for _, v := range task.Photos {
		// print(k, ":", v)
		// Create an image widget, and add it to the photos block
		i := doc.CreateElement("img").(*dom.HTMLImageElement)
		i.SetAttribute("photo-id", fmt.Sprintf("%d", v.ID))
		i.Class().SetString("photopreview")
		i.Src = v.Preview
		switch v.Type {
		case "PDF":
			// Is a PDF, so wrap the image with a box that includes the filename
			// and auto-break on each doc

			wspan := doc.CreateElement("div")
			wspan.AppendChild(i)
			p := doc.CreateElement("p")
			p.SetInnerHTML(v.Filename)
			wspan.AppendChild(p)
			div.AppendChild(wspan)
		case "Data":
			wspan := doc.CreateElement("div")
			wspan.AppendChild(i)
			p := doc.CreateElement("p")
			p.SetInnerHTML(v.Filename)
			wspan.AppendChild(p)
			div.AppendChild(wspan)
		case "Image":
			div.AppendChild(i)
		case "photo":
			print("WARNING: This is an old format attachment of type photo")
			print("v", v)
			print("Please run  \"update photo set type='Image' where type='photo';\" on database ....")
		default:
			print("adding attachment of unknown type", v.Type, "dt", v.Datatype, "fn", v.Filename)
			print("v", v)
		}
		// print("attaching click event to i", i)
		i.AddEventListener("click", false, func(evt dom.Event) {
			print("click on attachment preview image")
			evt.PreventDefault()
			theID, _ := strconv.Atoi(evt.Target().GetAttribute("photo-id"))

			go func() {
				photo := shared.Photo{}
				rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					ID:      theID,
				}, &photo)
				flds := strings.SplitN(photo.Data, ",", 2)
				print("got full photo", flds[0])
				switch flds[0] {
				// case "data:application/pdf;base64":
				default:
					w.Open(photo.Data, "", "")
				case "data:image/jpeg;base64", "data:image/png;base64", "data:image/gif;base64":
					if el2 := doc.QuerySelector("#photo-full").(*dom.HTMLImageElement); el2 != nil {
						doc.QuerySelector("#show-image").Class().Add("md-show")
						el2.Src = photo.Data
					}
				}

			}()
		})
	}
}

func showTaskPhotosOld(task shared.Task) {
	// print("populate the photos", task)

	w := dom.GetWindow()
	doc := w.Document()
	div := doc.QuerySelector("[name=Photos]")
	div.SetInnerHTML("")

	for _, v := range task.Photos {
		// print(k, ":", v)
		// Create an image widget, and add it to the photos block
		i := doc.CreateElement("img").(*dom.HTMLImageElement)
		i.SetAttribute("photo-id", fmt.Sprintf("%d", v.ID))
		i.Class().SetString("photopreview")
		i.Src = v.Preview
		div.AppendChild(i)

		i.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			theID, _ := strconv.Atoi(evt.Target().GetAttribute("photo-id"))
			// print("clicksed on photo ", theID)
			go func() {
				myPhoto := shared.Photo{}
				rpcClient.Call("UtilRPC.GetFullPhoto", shared.PhotoRPCData{
					Channel: Session.Channel,
					ID:      theID,
				}, &myPhoto)
				if myPhoto.Data != "" {
					if el2 := doc.QuerySelector("#photo-full").(*dom.HTMLImageElement); el2 != nil {
						doc.QuerySelector("#show-image").Class().Add("md-show")
						el2.Src = myPhoto.Data
					}
				}
			}()
		})
	}
}

func showPartsButtons(id int) {
	// print("populate the parts buttons")
	w := dom.GetWindow()
	doc := w.Document()
	if el := doc.QuerySelector("[name=PartsUsed]"); el != nil {
		// print("here with el", el)
		// el.Class().Add("hidden")
		div := el.(*dom.HTMLDivElement)
		// print("div", div)
		go func() {
			parts := []shared.TaskPart{}
			rpcClient.Call("TaskRPC.GetParts", shared.TaskRPCData{
				Channel: Session.Channel,
				ID:      id,
			}, &parts)
			// print("got parts", parts)
			el.SetInnerHTML("") //  clear out the div

			for _, v := range parts {
				// print("adding btn for", v)
				btn := doc.CreateElement("input").(*dom.HTMLInputElement)
				btn.Class().SetString("button button-outline")
				// btn.Class().SetString("button")
				btn.Type = "button"
				// btn.Value = fmt.Sprintf("%.1f x %s %s", v.QtyUsed, v.StockCode, v.PartName)
				btn.Value = fmt.Sprintf("%.1f x %s", v.QtyUsed, v.PartName)
				btn.SetAttribute("part-id", fmt.Sprintf("%d", v.PartID))
				div.AppendChild(btn)
				// print("added btn", btn)
			}
			el.Class().Remove("hidden")
		}()
	}
}

func expandAll(li dom.Element) {
	print("expanding", li)
	ul := li.ParentElement()
	print("parent UL", ul)
	if ul == nil {
		return
	}
	prevInput := ul.PreviousSibling()
	print("prevInput", prevInput)

	// Generally looking for the parent UL of the line clicked, which should be preceeded by
	// an input/checkbox field that defines whether or not the containing category is
	// expanded or not.  Due to some edge cases in data, we can end up with multiple UL
	// elements that follow the checkbox ... so the safe thing is to keep looking at the
	// prev sibling till we find the correct INPUT field
	for prevInput != nil && prevInput.NodeName() != "INPUT" {
		prevInput = prevInput.PreviousSibling()
	}

	if prevInput == nil {
		return
	}
	// got the parent checkbox, so expand it out
	print("prevInput type is ", prevInput.NodeType(), " name ", prevInput.NodeName())
	prevInput.(*dom.HTMLInputElement).Checked = true

	parentLI := ul.ParentElement()
	// print("parent LI", parentLI)
	if parentLI == nil {
		return
	}
	// recursion
	expandAll(parentLI)
}

func taskList(context *router.Context) {
	Session.Subscribe("task", _taskList)
	go _taskList("list", 0)
}

// Show a list of all tasks
func _taskList(action string, id int) {

	tasks := []shared.Task{}
	rpcClient.Call("TaskRPC.List", Session.Channel, &tasks)

	form := formulate.ListForm{}
	form.New("fa-server", "Task List - All Active Tasks")

	// Define the layout
	switch Session.UserRole {
	case "Admin", "Site Manager":
		form.Column("TaskID", "GetUserNameID")
		// form.Column("User", "Username")
		form.BoolColumn("Read", "IsRead")
	default:
		form.Column("ID", "GetID")
	}
	// form.Column("Src", "GetSource")
	form.MultiImgColumn("Photos", "Photos", "Thumb")
	form.Column("Date", "GetStartDate")
	// form.Column("Due", "GetDueDate")
	form.ColumnFormat("Site", "SiteName", "GetSiteClass")
	// form.Column("Machine", "MachineName")
	// form.Column("Component", "Component")
	form.Column("Component", "GetComponent")
	form.Column("Description", "Descr")
	form.Column("Duration", "DurationDays")

	switch Session.UserRole {
	case "Admin", "Site Mananger":
		form.Column("$$", "TotalCost")
	}

	// Add event handlers
	form.CancelEvent(func(evt dom.Event) {
		evt.PreventDefault()
		Session.Navigate("/")
	})

	form.RowEvent(func(key string) {
		Session.Navigate("/task/" + key)
	})

	form.PrintEvent(func(evt dom.Event) {
		dom.GetWindow().Print()
	})

	form.Render("task-list", "main", tasks)

	ctasks := []shared.Task{}
	rpcClient.Call("TaskRPC.ListCompleted", Session.Channel, &ctasks)

	cform := formulate.ListForm{}
	cform.New("fa-server", "Completed Tasks (last 30 days)")

	// Define the layout
	switch Session.UserRole {
	case "Admin", "Site Manager":
		cform.Column("User", "GetUserNameID")
	default:
		cform.Column("TaskID", "GetID")
	}
	// cform.Column("Src", "GetSource")
	cform.MultiImgColumn("Photos", "Photos", "Thumb")
	cform.Column("Date", "GetStartDate")
	// form.Column("Due", "GetDueDate")
	// cform.Column("Site", "SiteName")
	cform.ColumnFormat("Site", "SiteName", "GetSiteClass")
	// cform.Column("Machine", "MachineName")
	// cform.Column("Component", "Component")
	cform.Column("Component", "GetComponent")
	cform.Column("Description", "Descr")
	cform.Column("Duration", "DurationDays")

	switch Session.UserRole {
	case "Admin", "Site Mananger":
		cform.Column("$$", "TotalCost")
	}

	// cform.Column("Completed", "CompletedDate")
	cform.Column("Completed", "GetCompletedDate")

	cform.RowEvent(func(key string) {
		Session.Navigate("/task/" + key)
	})

	w := dom.GetWindow()
	doc := w.Document()

	// force a page break for printing
	div := doc.CreateElement("div")
	div.Class().Add("page-break")
	doc.QuerySelector("main").AppendChild(div)

	div = doc.CreateElement("div").(*dom.HTMLDivElement)
	div.SetID("ctasks")
	doc.QuerySelector("main").AppendChild(div)

	cform.Render("task-clist", "#ctasks", ctasks)

}

func taskPartList(context *router.Context) {
	print("TODO - taksPartList")
}

// Show a list of all tasks related to a stoppage
func stoppageTaskList(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		tasks := []shared.Task{}
		event := shared.Event{}
		rpcClient.Call("EventRPC.Get", shared.EventRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &event)
		rpcClient.Call("TaskRPC.StoppageList", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &tasks)

		BackURL := fmt.Sprintf("/stoppage/%d", id)

		form := formulate.ListForm{}
		form.New("fa-server", fmt.Sprintf("Task List for Stoppage - %06d", id))

		form.Column("User", "Username")
		form.Column("TaskID", "ID")
		form.MultiImgColumn("Photos", "Photos", "Thumb")
		form.Column("Date", "GetStartDate")
		// form.Column("Due", "GetDueDate")
		form.Column("Site", "SiteName")
		// form.Column("Machine", "MachineName")
		// form.Column("Component", "Component")
		form.Column("Component", "GetComponent")
		form.Column("Description", "Descr")
		form.Column("Duration", "DurationDays")
		form.Column("Completed", "CompletedDate")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/task/" + key)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.Render("event-task-list", "main", tasks)

	}()
}

func technicianDiary(context *router.Context) {
	go func() {
		done := false
		rpcClient.Call("TaskRPC.Diary", Session.Channel, &done)
	}()
}

func addTaskPartsTree(tree []shared.Category, ul *dom.HTMLUListElement, depth int) {

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
