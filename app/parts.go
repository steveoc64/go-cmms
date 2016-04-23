package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

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

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		if Session.UserRole == "Admin" {
			form.NewRowEvent(func(evt dom.Event) {
				evt.PreventDefault()
				Session.Router.Navigate("/class/add")
			})
		}

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/parts/" + key)
		})

		form.Render("class-select", "main", data)

	}()
}

func classAdd(context *router.Context) {
	print("TODO - classAdd")
}

// Show a list of all parts for the given class
func partList(context *router.Context) {

	partClass, _ := strconv.Atoi(context.Params["id"])
	print("show parts of class", partClass)

	go func() {
		data := []shared.Part{}
		req := shared.PartListReq{
			Channel: Session.Channel,
			Class:   partClass,
		}
		class := shared.PartClass{}
		rpcClient.Call("PartRPC.List", req, &data)
		rpcClient.Call("PartRPC.GetClass", partClass, &class)

		// load a form for the class
		loadTemplate("class-edit", "main", class)

		BackURL := "/class/select"
		Title := fmt.Sprintf("Parts of type - %s", class.Name)

		form := formulate.ListForm{}
		form.New("fa-puzzle-piece", Title)

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Description", "Descr")
		form.Column("Stock Code", "StockCode")
		form.Column("Reorder Lvl/Qty", "ReorderDetails")
		form.Column("Qty", "QtyType")
		form.Column("Latest Price", "DisplayPrice")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		if Session.UserRole == "Admin" {
			form.NewRowEvent(func(evt dom.Event) {
				evt.PreventDefault()
				Session.Router.Navigate(fmt.Sprintf("/part/add/%d", class.ID))
			})
		}

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/part/" + key)
		})

		form.Render("parts-list", "#parts-list-goes-here", data)

		// Add an onChange callback to the class edit fields
		w := dom.GetWindow()
		doc := w.Document()

		doc.QuerySelector("#class-name").AddEventListener("change", false, func(evt dom.Event) {
			print("TODO - Name has changed")
		})
		doc.QuerySelector("#class-descr").AddEventListener("change", false, func(evt dom.Event) {
			print("TODO - Description has changed")
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
		rpcClient.Call("PartRPC.Get", id, &part)

		BackURL := fmt.Sprintf("/parts/%d", part.Class)
		title := fmt.Sprintf("Part Details - %s - %s", part.Name, part.StockCode)
		form := formulate.EditForm{}
		form.New("fa-puzzle-piece", title)

		// Layout the fields

		form.Row(2).
			AddInput(1, "Name", "Name").
			AddInput(1, "Stock Code", "StockCode")

		form.Row(1).
			AddInput(1, "Description", "Descr")

		form.Row(3).
			AddNumber(1, "ReOrder Level", "ReorderStocklevel", "1").
			AddNumber(1, "ReOrder Qty", "ReorderQty", "1").
			AddInput(1, "Qty Type", "QtyType")

		form.Row(1).
			AddTextarea(1, "Notes", "Notes")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("TODO - delete part")
			return

			// machine.ID = id
			// go func() {
			// 	data := shared.MachineUpdateData{
			// 		Channel: Session.Channel,
			// 		Machine: &machine,
			// 	}
			// 	done := false
			// 	rpcClient.Call("MachineRPC.Delete", data, &done)
			// 	Session.Router.Navigate(BackURL)
			// }()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("TODO - part save")
			return
			// form.Bind(&machine)
			// go func() {
			// 	data := shared.MachineUpdateData{
			// 		Channel: Session.Channel,
			// 		Machine: &machine,
			// 	}
			// 	done := false
			// 	rpcClient.Call("MachineRPC.Update", data, &done)
			// 	Session.Router.Navigate(BackURL)
			// }()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &part)

		// And attach actions
		form.ActionGrid("part-actions", "#action-grid", part.ID, func(url string) {
			Session.Router.Navigate(url)
		})

	}()
}

func partAdd(context *router.Context) {
	print("TODO partAdd")
}
