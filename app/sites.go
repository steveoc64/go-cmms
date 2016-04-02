package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

type SiteMapData struct {
	Status shared.SiteStatusReport
	Sites  []shared.Site
}

// Display a map showing the status of all sites, with buttons to
// jump to specific sites
func siteMap(context *router.Context) {

	// Get a list of sites
	go func() {
		data := SiteMapData{}
		err := rpcClient.Call("SiteRPC.UserList", channelID, &data.Sites)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			// Get the site statuses
			rpcClient.Call("SiteRPC.StatusReport", channelID, &data.Status)

			w := dom.GetWindow()
			doc := w.Document()
			loadTemplate("sitemap", "main", data)

			// Attach listeners for each button
			for _, v := range data.Sites {
				mbtn := doc.GetElementByID(fmt.Sprintf("%d", v.ID)).(*dom.HTMLInputElement)
				mbtn.AddEventListener("click", false, func(evt dom.Event) {
					id := evt.Target().GetAttribute("id") // site id is a string
					evt.PreventDefault()
					r.Navigate("/sitemachines/" + id)
				})
			}
		}
	}()
}

type SiteMachineData struct {
	MultiSite bool
	Site      shared.Site
	Status    shared.SiteStatusReport
	Machines  []shared.Machine
}

type RaiseIssueData struct {
	MachineID int
	CompID    int
}

// Show all the machines at a Site, for the case where we have > 1 site
// Param:  site
func siteMachines(context *router.Context) {
	idStr := context.Params["site"]
	// print("in the machineDet function", idStr)
	id, _ := strconv.Atoi(idStr)

	// Get a list of machines at this site
	req := shared.MachineReq{
		Channel: channelID,
		SiteID:  id,
	}
	data := SiteMachineData{}

	nonTools := []string{
		"Electrical", "Hydraulic", "Lube", "Printer", "Console", "Uncoiler", "Rollbed",
	}

	go func() {
		data.MultiSite = true
		rpcClient.Call("SiteRPC.Get", id, &data.Site)
		// print("Site =", data.Site)

		err := rpcClient.Call("SiteRPC.MachineList", &req, &data.Machines)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			// Get the site statuses
			rpcClient.Call("SiteRPC.StatusReport", channelID, &data.Status)
			// print("SiteMachine status report", data.Status)

			loadTemplate("sitemachines", "main", data)
			w := dom.GetWindow()
			doc := w.Document()
			austmap := doc.GetElementByID("austmap")
			austmap.AddEventListener("click", false, func(evt dom.Event) {
				r.Navigate("/")
			})

			// Attach a menu opener for each machine
			for _, v := range data.Machines {
				mid := fmt.Sprintf("machine-div-%d", v.ID)
				machinediv := doc.GetElementByID(mid)
				machinediv.AddEventListener("click", false, func(evt dom.Event) {
					machine_id, _ := strconv.Atoi(machinediv.GetAttribute("machine-id"))
					machinemenu := doc.GetElementByID("machine-menu").(*dom.BasicHTMLElement)

					// get the machine and construct a new menu based on the components
					for _, m := range data.Machines {
						if m.ID == machine_id {
							menu := fmt.Sprintf(`<h3 id="machine-comp-title">%s</h3>`, m.Name)
							// add the tool components
							for i, c := range m.Components {
								menu += fmt.Sprintf(`<a href="#" id="machine-comp-%d" machine="%d" comp="%d" class="%s">%d %s</a>`,
									c.ID, m.ID, c.ID, c.GetClass(), i+1, c.Name)
							}
							// add the non-tool components
							for i, c := range nonTools {
								menu += fmt.Sprintf(`<a href="#" id="machine-nontool-%d" machine="%d" class="%s" comp="%s">%s</a>`,
									i, m.ID, m.GetClass(m.Electrical), c, c)
							}
							machinemenu.SetInnerHTML(menu)
							tk := machinemenu.Class()
							if !tk.Contains("cbp-spmenu-open") {
								tk.Add("cbp-spmenu-open")
							}
							// attach a backout option on the title
							a := doc.GetElementByID("machine-comp-title")
							a.AddEventListener("click", false, func(evt dom.Event) {
								evt.PreventDefault()
								tk.Remove("cbp-spmenu-open")
							})
							// attach event listeners to each tool menu item
							for _, c := range m.Components {
								a := doc.GetElementByID(fmt.Sprintf("machine-comp-%d", c.ID))
								a.AddEventListener("click", false, func(evt dom.Event) {
									evt.PreventDefault()

									// Get the details of the component that we clicked on
									t := evt.Target()
									machineID, _ := strconv.Atoi(t.GetAttribute("machine"))
									compID, _ := strconv.Atoi(t.GetAttribute("comp"))
									print("machine id", machineID, "comp id", compID)

									// hide the side menu
									tk.Remove("cbp-spmenu-open")

									// load a raise issue form from a template
									d := RaiseIssueData{
										MachineID: machineID,
										CompID:    compID,
									}
									loadTemplate("raise-comp-issue", "#raise-comp-issue", d)
									doc.QuerySelector("#raise-comp-issue").Class().Add("md-show")

									// Add the machine diagram to the form
									for _, m := range data.Machines {
										print("compare", m.ID, machineID)
										if m.ID == machineID {
											print("got match", m)
											loadTemplate("machine-diag", "#issue-machine-diag", m)
											break
										}
									}

									// Handle button clicks
									doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
										print("cancel the event")
										doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
									})
									doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
										print("TODO - save the event details")
										doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
									})
								})
							}
							// attach event listeners to each non-tool menu item
							for i, _ := range nonTools {
								a := doc.GetElementByID(fmt.Sprintf("machine-nontool-%d", i))
								a.AddEventListener("click", false, func(evt dom.Event) {
									evt.PreventDefault()
									tk.Remove("cbp-spmenu-open")
								})
							}
						} // is matching machine
					}

				})
			} // range
		} // else
	}()
}

// Show all the machines at a Site, for the case where we have only 1 site
func homeSite(context *router.Context) {

	// Get a list of machines at this site
	data := SiteMachineData{}

	go func() {
		data.MultiSite = false
		rpcClient.Call("SiteRPC.GetHome", channelID, &data.Site)
		// print("Site =", data.Site)

		err := rpcClient.Call("SiteRPC.HomeMachineList", channelID, &data.Machines)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			loadTemplate("sitemachines", "main", data)
		}
	}()
}

// Show a list of all sites
func siteList(context *router.Context) {
	loadTemplate("sitelist", "main", nil)
}
