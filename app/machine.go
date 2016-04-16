package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	ff "github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// "github.com/steveoc64/go-cmms/shared"
// "honnef.co/go/js/dom"

func machineList(context *router.Context) {
	print("TODO machineList")
}

func machineEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		machine := shared.Machine{}
		rpcClient.Call("MachineRPC.Get", id, &machine)

		BackURL := fmt.Sprintf("/site/machine/%d", machine.SiteId)
		title := fmt.Sprintf("Machine Details - %s - %s", machine.Name, *machine.SiteName)
		form := ff.EditForm{}
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
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			machine.ID = id
			go func() {
				data := shared.MachineUpdateData{
					Channel: Session.Channel,
					Machine: &machine,
				}
				done := false
				rpcClient.Call("MachineRPC.Delete", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machine)
			go func() {
				data := shared.MachineUpdateData{
					Channel: Session.Channel,
					Machine: &machine,
				}
				done := false
				rpcClient.Call("MachineRPC.Update", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &machine)

		// And attach actions
		form.ActionGrid("machine-actions", "#action-grid", machine.ID, func(url string) {
			print("calling action", url)
			Session.Router.Navigate(url)
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
		rpcClient.Call("SiteRPC.Get", id, &site)
		BackURL := fmt.Sprintf("/site/machine/%d", site.ID)
		title := fmt.Sprintf("Add Machine for Site - %s", site.Name)
		form := ff.EditForm{}
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
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("Delete this record")
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&machine)
			machine.SiteId = site.ID
			machine.Status = "Running"
			go func() {
				data := shared.MachineUpdateData{
					Channel: Session.Channel,
					Machine: &machine,
				}
				newID := 0
				rpcClient.Call("MachineRPC.Insert", data, &newID)
				print("added machine ID", newID)
				Session.Router.Navigate(BackURL)
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

// Show a list of all machines for the given site
func siteMachineList(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		req := shared.MachineReq{
			Channel: Session.Channel,
			SiteID:  id,
		}
		data := SiteMachineData{}

		rpcClient.Call("SiteRPC.Get", id, &data.Site)
		rpcClient.Call("SiteRPC.MachineList", &req, &data.Machines)

		form := ff.ListForm{}
		form.New("fa-cogs", "Machine List for - "+data.Site.Name, &data.Machines)

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Description", "Descr")
		form.Column("Status", "Status")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(fmt.Sprintf("/site/%d", id))
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(fmt.Sprintf("/site/machine/add/%d", id))
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/machine/" + key)
		})

		form.Render("site-machine-list", "main", data)

	}()
}
