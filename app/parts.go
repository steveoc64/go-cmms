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
		form.Column("Number of Parts", "Count")

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
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&partClass)
			go func() {
				data := shared.PartClassUpdateData{
					Channel:   Session.Channel,
					PartClass: &partClass,
				}
				newID := 0
				rpcClient.Call("PartRPC.InsertClass", data, &newID)
				print("added class ID", newID)
				Session.Router.Navigate(BackURL)
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
		req := shared.PartListReq{
			Channel: Session.Channel,
			Class:   partClass,
		}
		class := shared.PartClass{}
		rpcClient.Call("PartRPC.List", req, &data)
		rpcClient.Call("PartRPC.GetClass", partClass, &class)

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
							data := shared.PartClassUpdateData{
								Channel:   Session.Channel,
								PartClass: &class,
							}
							done := false
							rpcClient.Call("PartRPC.DeleteClass", data, &done)
						}()
						Session.Router.Navigate(BackURL)
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
			go func() {
				class.Name = doc.QuerySelector("#class-name").(*dom.HTMLInputElement).Value
				data := shared.PartClassUpdateData{
					Channel:   Session.Channel,
					PartClass: &class,
				}
				done := false
				rpcClient.Call("PartRPC.UpdateClass", data, &done)
			}()
		})
		doc.QuerySelector("#class-descr").AddEventListener("change", false, func(evt dom.Event) {
			print("TODO - Description has changed")
			go func() {
				class.Descr = doc.QuerySelector("#class-descr").(*dom.HTMLInputElement).Value
				data := shared.PartClassUpdateData{
					Channel:   Session.Channel,
					PartClass: &class,
				}
				done := false
				rpcClient.Call("PartRPC.UpdateClass", data, &done)
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

		rpcClient.Call("PartRPC.Get", id, &part)
		rpcClient.Call("PartRPC.ClassList", Session.Channel, &classes)
		rpcClient.Call("PartRPC.StockList", id, &stocks)
		rpcClient.Call("PartRPC.PriceList", id, &prices)

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
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				data := shared.PartUpdateData{
					Channel: Session.Channel,
					Part:    &part,
				}
				done := false
				rpcClient.Call("PartRPC.Delete", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&part)
			go func() {
				data := shared.PartUpdateData{
					Channel: Session.Channel,
					Part:    &part,
				}
				done := false
				rpcClient.Call("PartRPC.Update", data, &done)
				NewBackURL := ""
				if done {
					// Go back to parts list
					NewBackURL = fmt.Sprintf("/parts/%d", part.Class)
				} else {
					// refresh this screen
					NewBackURL = fmt.Sprintf("/part/%d", part.ID)
				}
				Session.Router.Navigate(NewBackURL)
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
		rpcClient.Call("PartRPC.GetClass", id, &class)
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
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&part)
			go func() {
				data := shared.PartUpdateData{
					Channel: Session.Channel,
					Part:    &part,
				}
				newID := 0
				rpcClient.Call("PartRPC.Insert", data, &newID)
				print("Added new part", newID)
				NewBackURL := fmt.Sprintf("/parts/%d", part.Class)
				Session.Router.Navigate(NewBackURL)
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
