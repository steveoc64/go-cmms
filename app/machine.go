package main

import (
	"fmt"
	"strconv"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

// "itrak-cmms/shared"
// "honnef.co/go/js/dom"

// Show a list of all machines for the given site
func siteMachineList(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		site := shared.Site{}
		machines := []shared.Machine{}

		rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &site)
		rpcClient.Call("SiteRPC.MachineList", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machines)

		form := formulate.ListForm{}
		form.New("fa-cogs", "Machine List for - "+site.Name)

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Description", "Descr")
		form.Column("Status", "Status")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(fmt.Sprintf("/site/%d", id))
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(fmt.Sprintf("/site/machine/add/%d", id))
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/machine/" + key)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.Render("site-machine-list", "main", machines)
		// form.Render("site-machine-list", "main", data)

	}()
}

func machineEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machine := shared.Machine{}
		users := []shared.User{}
		technicians := []shared.User{}
		// classes := []shared.PartClass{}
		machineTypes := []shared.MachineType{}

		rpcClient.Call("MachineRPC.Get", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machine)
		rpcClient.Call("UserRPC.GetManagers", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      machine.SiteID,
		}, &users)
		rpcClient.Call("UserRPC.GetTechnicians", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      machine.SiteID,
		}, &technicians)
		// rpcClient.Call("PartRPC.ClassList", Session.Channel, &classes)
		rpcClient.Call("MachineRPC.MachineTypes", shared.MachineRPCData{
			Channel: Session.Channel,
		}, &machineTypes)

		BackURL := fmt.Sprintf("/site/machine/%d", machine.SiteID)
		title := fmt.Sprintf("Machine Details - %s - %s", machine.Name, *machine.SiteName)
		form := formulate.EditForm{}
		form.New("fa-cogs", title)

		// Layout the fields

		form.Row(3).
			AddInput(1, "Name", "Name").
			AddInput(1, "Serial #", "Serialnum").
			Add(1, "Status", "text", "Status", "disabled")

		form.Row(1).
			AddSelect(1, "Machine Type", "MachineType",
				machineTypes, "ID", "Name",
				1, machine.MachineType)

		form.Row(1).
			AddCustom(1, "Diagram", "Diag", "")

		form.Row(1).
			AddInput(1, "Descrpition", "Descr")

		form.Row(2).
			AddSelect(1, "Stoppage Alerts To", "AlertsTo", users, "ID", "Name", 0, machine.AlertsTo).
			AddSelect(1, "Send Scheduled Tasks To", "TasksTo", technicians, "ID", "Name", 0, machine.TasksTo)

		form.Row(1).
			Add(1, "Notes", "textarea", "Notes", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			machine.ID = id
			go func() {
				done := false
				rpcClient.Call("MachineRPC.Delete", shared.MachineRPCData{
					Channel: Session.Channel,
					Machine: &machine,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machine)
			go func() {
				done := false
				rpcClient.Call("MachineRPC.Update", shared.MachineRPCData{
					Channel: Session.Channel,
					Machine: &machine,
				}, &done)
				// Session.Navigate(BackURL)
				Session.Reload(context)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machine)

		// render the machine diagram

		loadTemplate("machine-diag", "[name=Diag]", &shared.RaiseIssue{Machine: &machine})

		// on change the machine type, save and refresh
		w := dom.GetWindow()
		doc := w.Document()
		el := doc.QuerySelector("[name=MachineType]")
		if el != nil {
			doc.QuerySelector("[name=MachineType]").AddEventListener("change", false, func(evt dom.Event) {
				form.Bind(&machine)
				go func() {
					done := false
					rpcClient.Call("MachineRPC.Update", shared.MachineRPCData{
						Channel: Session.Channel,
						Machine: &machine,
					}, &done)
					// Session.Navigate(BackURL)
					Session.Reload(context)
				}()
			})
		}

		// And attach actions
		form.ActionGrid("machine-actions", "#action-grid", machine.ID, func(url string) {
			Session.Navigate(url)
		})

	}()

}

func siteMachineAdd(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machine := shared.Machine{}
		site := shared.Site{}
		rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &site)
		BackURL := fmt.Sprintf("/site/machine/%d", site.ID)
		title := fmt.Sprintf("Add Machine for Site - %s", site.Name)
		form := formulate.EditForm{}
		form.New("fa-cogs", title)

		// Layout the fields

		form.Row(3).
			Add(1, "Name", "text", "Name", `id="focusme"`).
			Add(1, "Serial #", "text", "Serialnum", "").
			Add(1, "Status", "text", "Status", "disabled")

		form.Row(1).
			Add(1, "Description", "text", "Descr", "")

		form.Row(1).
			Add(1, "Notes", "textarea", "Notes", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machine)
			machine.SiteID = site.ID
			machine.Status = "Running"
			go func() {
				newID := 0
				rpcClient.Call("MachineRPC.Insert", shared.MachineRPCData{
					Channel: Session.Channel,
					Machine: &machine,
				}, &newID)
				print("added machine ID", newID)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machine)

	}()

}

func machineReports(context *router.Context) {
	print("TODO - machineReports")
}

func machineStoppageList(context *router.Context) {
	print("TODO - machineStoppageList")
}

type SiteMachineListData struct {
	Site     shared.Site
	Machines []shared.Machine
}

func machineTypes(context *router.Context) {

	go func() {
		data := []shared.MachineType{}
		rpcClient.Call("MachineRPC.MachineTypes", shared.MachineRPCData{
			Channel: Session.Channel,
		}, &data)

		// print("got machine types", data)
		BackURL := "/"

		form := formulate.ListForm{}
		form.New("fa-cubes", "Machine Types")

		// Define the layout
		form.Column("Name", "Name")
		form.ImgColumn("Photo", "PhotoThumbnail")
		form.BoolColumn("Elec", "Electrical")
		form.BoolColumn("Hyd", "Hydraulic")
		form.BoolColumn("Pnue", "Pnuematic")
		form.BoolColumn("Lube", "Lube")
		form.BoolColumn("Prnt", "Printer")
		form.BoolColumn("Cnsl", "Console")
		form.BoolColumn("Unclr", "Uncoiler")
		form.BoolColumn("Rlbd", "Rollbed")
		form.BoolColumn("Cnvyr", "Conveyor")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		if Session.UserRole == "Admin" {
			form.NewRowEvent(func(evt dom.Event) {
				evt.PreventDefault()
				Session.Navigate("/machinetype/add")
			})
		}

		form.RowEvent(func(key string) {
			Session.Navigate("/machinetype/" + key)
		})

		form.Render("machine-type", "main", data)

		w := dom.GetWindow()
		doc := w.Document()

		el := doc.QuerySelector(".grid-form")
		print("el", el)
		if el != nil {

			doc.QuerySelector(".grid-form").AddEventListener("change", false, func(evt dom.Event) {
				print("something changed")

			})
		}
		// for _, v := range data {
		// 	if v.PhotoThumbnail != "" {
		// 		ename := fmt.Sprintf(`[name=PhotoThumbnail-%d]`, v.ID)
		// 		el := doc.QuerySelector(ename).(*dom.HTMLImageElement)
		// 		// print("img = ", el.Src)
		// 		el.Src = v.PhotoThumbnail
		// 	}
		// }
		// // print("fill in the img src inside the listform automatically")
	}()
}

func machineTypeAdd(context *router.Context) {
	print("TODO - add machine type")

	go func() {
		machineType := shared.MachineType{}

		BackURL := "/machinetypes"
		title := "Add New Machine Type"
		form := formulate.EditForm{}
		form.New("fa-cubes", title)

		// Layout the fields
		form.Row(1).
			AddInput(1, "Name", "Name")

		form.Row(3).
			AddCheck(1, "Electrical", "Electrical").
			AddCheck(1, "Hydraulic", "Hydraulic").
			AddCheck(1, "Pnuematic", "Pnuematic")

		form.Row(3).
			AddCheck(1, "Console", "Console").
			AddCheck(1, "Printer", "Printer").
			AddCheck(1, "Lube", "Lube")

		form.Row(3).
			AddCheck(1, "UnCoiler", "Uncoiler").
			AddCheck(1, "RollBed", "Rollbed").
			AddCheck(1, "Conveyor", "Conveyor")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machineType)
			go func() {
				newID := 0
				rpcClient.Call("MachineRPC.InsertMachineType", shared.MachineTypeRPCData{
					Channel:     Session.Channel,
					MachineType: &machineType,
				}, &newID)
				print("added machine type", newID)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machineType)

	}()

}

func machineTypeEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machineType := shared.MachineType{}

		rpcClient.Call("MachineRPC.GetMachineType", shared.MachineTypeRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machineType)

		// print("got machine type", machineType)

		BackURL := "/machinetypes"
		title := fmt.Sprintf("Machine Type Details - %s", machineType.Name)
		form := formulate.EditForm{}
		form.New("fa-cubes", title)

		// Layout the fields
		form.Row(1).
			AddInput(1, "Name", "Name")

		form.Row(1).
			AddCustom(1, "Diagram", "Diag", "")

		form.Row(1).
			AddPreview(1, "Photo", "PhotoPreview")

		form.Row(3).
			AddCheck(1, "Electrical", "Electrical").
			AddCheck(1, "Hydraulic", "Hydraulic").
			AddCheck(1, "Pnuematic", "Pnuematic")

		form.Row(3).
			AddCheck(1, "Console", "Console").
			AddCheck(1, "Printer", "Printer").
			AddCheck(1, "Lube", "Lube")

		form.Row(3).
			AddCheck(1, "UnCoiler", "Uncoiler").
			AddCheck(1, "RollBed", "Rollbed").
			AddCheck(1, "Conveyor", "Conveyor")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			machineType.ID = id
			go func() {
				done := false
				rpcClient.Call("MachineRPC.DeleteMachineType", shared.MachineTypeRPCData{
					Channel: Session.Channel,
					ID:      id,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machineType)
			go func() {
				done := false
				rpcClient.Call("MachineRPC.UpdateMachineType", shared.MachineTypeRPCData{
					Channel:     Session.Channel,
					ID:          id,
					MachineType: &machineType,
				}, &done)
				Session.Reload(context)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machineType)

		// Add a machine diag to the form
		// print("mt =", machineType)
		loadTemplate("machine-type-diag", "[name=Diag]", &machineType)

		// And attach actions
		form.ActionGrid("machine-type-actions", "#action-grid", id, func(url string) {
			Session.Navigate(fmt.Sprintf("/machinetype/%d/%s", id, url))
		})

	}()

}

func machineTypeTools(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		data := []shared.MachineTypeTool{}
		rpcClient.Call("MachineRPC.MachineTypeTools", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &data)
		machineType := shared.MachineType{}
		rpcClient.Call("MachineRPC.GetMachineType", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machineType)

		print("got machine type tools", data)
		BackURL := fmt.Sprintf("/machinetype/%d", id)

		form := formulate.ListForm{
			Draggable: true,
		}
		form.New("fa-cubes", fmt.Sprintf("Tools - %s", machineType.Name))
		// form.KeyField = "MachineID"

		// Define the layout
		form.Column("Position / Name", "GetName")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(fmt.Sprintf("/machinetype/%d/tool/add", id))
		})

		form.RowEvent(func(key string) {
			print("clicked on key", key)
			Session.Navigate(fmt.Sprintf("/machinetype/%d/tool/%s", id, key))
		})

		form.Render("machine-type-tools", "main", data)

		// // And attach actions
		// form.ActionGrid("machine-type-actions", "#action-grid", id, func(url string) {
		// 	Session.Navigate(fmt.Sprintf("/machinetype/%d/%s", id, url))
		// })

	}()
}

func machineTypeToolAdd(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {

		machineType := shared.MachineType{}
		rpcClient.Call("MachineRPC.GetMachineType", shared.MachineTypeRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machineType)

		machineTypeTool := shared.MachineTypeTool{
			MachineID: id,
			ID:        machineType.NumTools + 1,
		}

		BackURL := fmt.Sprintf("/machinetype/%d/tools", id)

		title := fmt.Sprintf("Add New Tool - %s", machineType.Name)

		form := formulate.EditForm{}
		form.New("fa-cubes", title)

		// Layout the fields
		form.Row(3).
			AddNumber(1, "Position", "ID", "1").
			AddInput(2, "Name", "Name")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machineTypeTool)
			go func() {
				newID := 0
				// print("saving", machineTypeTool)
				rpcClient.Call("MachineRPC.InsertMachineTypeTool", shared.MachineTypeToolRPCData{
					Channel:         Session.Channel,
					MachineTypeTool: &machineTypeTool,
				}, &newID)
				print("added tool", machineTypeTool.ID)
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machineTypeTool)

	}()

}

func machineTypeToolEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	tool, err := strconv.Atoi(context.Params["tool"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machineType := shared.MachineType{}

		rpcClient.Call("MachineRPC.GetMachineType", shared.MachineTypeRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machineType)

		// print("got machine type", machineType)

		machineTypeTool := shared.MachineTypeTool{}

		rpcClient.Call("MachineRPC.GetMachineTypeTool", shared.MachineTypeToolRPCData{
			Channel:   Session.Channel,
			MachineID: id,
			ID:        tool,
		}, &machineTypeTool)

		print("got machine type tool", machineTypeTool)

		BackURL := fmt.Sprintf("/machinetype/%d/tools", id)

		title := fmt.Sprintf("Machine Tool Details - %s - %s",
			machineType.Name, machineTypeTool.Name)
		form := formulate.EditForm{}
		form.New("fa-cubes", title)

		// Layout the fields
		form.Row(3).
			AddNumber(1, "Position", "ID", "1").
			AddInput(2, "Name", "Name")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			machineType.ID = id
			go func() {
				done := false
				rpcClient.Call("MachineRPC.DeleteMachineTypeTool", shared.MachineTypeToolRPCData{
					Channel:   Session.Channel,
					MachineID: id,
					ID:        tool,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machineTypeTool)
			go func() {
				done := false
				rpcClient.Call("MachineRPC.UpdateMachineTypeTool", shared.MachineTypeToolRPCData{
					Channel:         Session.Channel,
					MachineID:       id,
					ID:              tool,
					MachineTypeTool: &machineTypeTool,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machineTypeTool)

	}()

	print("TODO - machineTypeToolEdit", id, tool)
}

func machineTypeParts(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	print("TODO - machineTypeParts", id)
}

func machineTypeMachines(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		data := []shared.Machine{}
		rpcClient.Call("MachineRPC.MachineOfType", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &data)
		machineType := shared.MachineType{}
		rpcClient.Call("MachineRPC.GetMachineType", shared.MachineRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &machineType)

		BackURL := fmt.Sprintf("/machinetype/%d", id)

		form := formulate.ListForm{
			Draggable: true,
		}
		form.New("fa-gears", fmt.Sprintf("Machines of type %s", machineType.Name))
		// form.KeyField = "MachineID"

		// Define the layout
		form.Column("Site", "SiteName")
		form.Column("Name", "Name")
		form.Column("Status", "Status")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.RowEvent(func(key string) {
			print("clicked on key", key)
			Session.Navigate(fmt.Sprintf("/machine/%s", key))
		})

		form.Render("machine-type-machines", "main", data)

	}()
}
