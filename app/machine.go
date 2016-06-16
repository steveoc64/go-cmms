package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// "github.com/steveoc64/go-cmms/shared"
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
		classes := []shared.PartClass{}

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
		rpcClient.Call("PartRPC.ClassList", Session.Channel, &classes)

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
			AddSelect(1, "Machine Type (for classification of Parts)", "PartClass",
				classes, "ID", "Name",
				1, machine.PartClass)

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
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machine)

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
			Add(1, "Descrpition", "text", "Descr", "")

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

		print("got machine types", data)
		BackURL := "/"

		form := formulate.ListForm{}
		form.New("fa-cubes", "Machine Types")

		// Define the layout
		form.Column("Name", "Name")

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
			AddCheck(1, "RollBed", "Rollbed")

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

		print("got machine type", machineType)

		BackURL := "/machinetypes"
		title := fmt.Sprintf("Machine Type Details - %s", machineType.Name)
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
			AddCheck(1, "RollBed", "Rollbed")

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
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machineType)

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

		form := formulate.ListForm{}
		form.New("fa-cubes", fmt.Sprintf("Tools - %s", machineType.Name))
		// form.KeyField = "MachineID"

		// Define the layout
		form.Column("Name", "Name")

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

	}()
}

func machineTypeToolAdd(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	print("TODO - machineTypeToolAdd", id)
}

func machineTypeToolEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	print("TODO - machineTypeParts", id)
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
	print("TODO - machineTypeMachines", id)
}
