package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/form"
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

type SiteMapData struct {
	Status shared.SiteStatusReport
	Sites  []shared.Site
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
	IsTool    bool
	Machine   *shared.Machine
	Component *shared.Component
	NonTool   string
}

// Display a map showing the status of all sites, with buttons to
// jump to specific sites
func siteMap(context *router.Context) {

	// Get a list of sites
	go func() {
		data := SiteMapData{}
		err := rpcClient.Call("SiteRPC.UserList", Session.Channel, &data.Sites)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			// Get the site statuses
			rpcClient.Call("SiteRPC.StatusReport", Session.Channel, &data.Status)

			w := dom.GetWindow()
			doc := w.Document()
			loadTemplate("sitemap", "main", data)

			// Attach listeners for each button
			for _, v := range data.Sites {
				mbtn := doc.GetElementByID(fmt.Sprintf("%d", v.ID)).(*dom.HTMLDivElement)
				mbtn.AddEventListener("click", false, func(evt dom.Event) {
					id := evt.CurrentTarget().GetAttribute("id") // site id is a string
					evt.PreventDefault()
					Session.Router.Navigate("/sitemachines/" + id)
				})
			}

			// Add an Action Grid depending on which role the user is logged in as
			print("user role =", Session.UserRole)
			switch Session.UserRole {
			case "Admin", "Site Manager":
				loadTemplate("admin-actions", "#action-grid", nil)
			case "Worker":
				loadTemplate("worker-actions", "#action-grid", nil)
			}

			// Add a click handler to navigate to the page
			for _, ai := range doc.QuerySelectorAll(".action__item") {
				url := ai.(*dom.HTMLDivElement).GetAttribute("url")
				if url != "" {
					ai.AddEventListener("click", false, func(evt dom.Event) {
						url := evt.CurrentTarget().GetAttribute("url")
						Session.Router.Navigate(url)
					})
				}
			}
		}
	}()
}

// Show all the machines at a Site, for the case where we have > 1 site
// Param:  site
func siteMachines(context *router.Context) {
	idStr := context.Params["site"]
	id, _ := strconv.Atoi(idStr)

	// Get a list of machines at this site
	req := shared.MachineReq{
		Channel: Session.Channel,
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
			rpcClient.Call("SiteRPC.StatusReport", Session.Channel, &data.Status)
			// print("SiteMachine status report", data.Status)

			loadTemplate("sitemachines", "main", data)
			w := dom.GetWindow()
			doc := w.Document()
			austmap := doc.GetElementByID("austmap")
			austmap.AddEventListener("click", false, func(evt dom.Event) {
				Session.Router.Navigate("/")
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
									i, m.ID, m.GetClass(m.GetStatus(c)), c, c)
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
									d := RaiseIssueData{}
									d.MachineID, _ = strconv.Atoi(t.GetAttribute("machine"))
									d.CompID, _ = strconv.Atoi(t.GetAttribute("comp"))
									d.IsTool = true
									for _, m1 := range data.Machines {
										if m1.ID == d.MachineID {
											d.Machine = &m1

											// now get the component
											for _, c1 := range m1.Components {
												if c1.ID == d.CompID {
													d.Component = &c1
													break
												}
											}
											break
										}
									}

									// hide the side menu
									tk.Remove("cbp-spmenu-open")

									// load a raise issue form from a template
									loadTemplate("raise-comp-issue", "#raise-comp-issue", d)
									doc.QuerySelector("#raise-comp-issue").Class().Add("md-show")

									// Add the machine diagram to the form
									if d.Machine != nil {
										loadTemplate("machine-diag", "#issue-machine-diag", d)
									}

									// Handle button clicks
									doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
										print("TODO - cancel the event, cleanup any temp attachments")
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

									// Get the details of the component that we clicked on
									t := evt.Target()
									d := RaiseIssueData{}
									d.MachineID, _ = strconv.Atoi(t.GetAttribute("machine"))
									d.NonTool = t.GetAttribute("comp")
									d.IsTool = false
									for _, m1 := range data.Machines {
										if m1.ID == d.MachineID {
											d.Machine = &m1
											break
										}
									}

									// hide the side menu
									tk.Remove("cbp-spmenu-open")

									// load a raise issue form from a template
									loadTemplate("raise-comp-issue", "#raise-comp-issue", d)
									doc.QuerySelector("#raise-comp-issue").Class().Add("md-show")

									// Add the machine diagram to the form
									if d.Machine != nil {
										loadTemplate("machine-diag", "#issue-machine-diag", d)
									}

									// Handle button clicks
									doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
										print("TODO - cancel the event, cleanup any temp attachments")
										doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
									})
									doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
										print("TODO - save the event details")
										doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
									})
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
		rpcClient.Call("SiteRPC.GetHome", Session.Channel, &data.Site)
		// print("Site =", data.Site)

		err := rpcClient.Call("SiteRPC.HomeMachineList", Session.Channel, &data.Machines)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			loadTemplate("sitemachines", "main", data)
			dom.GetWindow().ScrollTo(1, 1)
		}
	}()
}

// Show a list of all sites
func siteList(context *router.Context) {

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := []shared.Site{}
		rpcClient.Call("SiteRPC.List", Session.Channel, &data)
		loadTemplate("site-list", "main", data)

		// Add a handler for clicking on a row
		doc.GetElementByID("site-list").AddEventListener("click", false, func(evt dom.Event) {
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			Session.Router.Navigate("/site/" + key)
		})

		// Add a handler for clicking on the add butto
		doc.QuerySelector(".data-add-btn").AddEventListener("click", false, func(evt dom.Event) {
			print("add new site")
		})
	}()
}

type SiteEditData struct {
	Site  shared.Site
	Sites []shared.Site
}

// Show an edit form for the given site
func siteEdit(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}
	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := SiteEditData{}
		rpcClient.Call("SiteRPC.Get", id, &data.Site)
		rpcClient.Call("SiteRPC.List", Session.Channel, &data.Sites)
		loadTemplate("site-edit", "main", data)

		// Add handlers for this form
		doc.QuerySelector("legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate("/sites")
		})
		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			print("cancel edit site")
			Session.Router.Navigate("/sites")
		})
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			// Parse the form element and get a form.Form object in return.
			f, err := form.Parse(doc.QuerySelector(".grid-form"))
			if err != nil {
				print("form parse error", err.Error())
				return
			}
			if err := f.Bind(&data.Site); err != nil {
				print("form bind error", err.Error())
				return
			}
			// manually get the textarea
			data.Site.Notes = doc.GetElementByID("notes").(*dom.HTMLTextAreaElement).Value

			// manually get the selected options for now
			parentSite := doc.GetElementByID("parentSite").(*dom.HTMLSelectElement).SelectedIndex
			data.Site.ParentSite = 0
			if parentSite > 0 {
				data.Site.ParentSite = data.Sites[parentSite-1].ID
			}
			data.Site.StockSite = 0
			stockSite := doc.GetElementByID("stockSite").(*dom.HTMLSelectElement).SelectedIndex
			if stockSite > 0 {
				data.Site.StockSite = data.Sites[stockSite-1].ID
			}
			updateData := &shared.SiteUpdateData{
				Channel: Session.Channel,
				Site:    &data.Site,
			}
			go func() {
				retval := 0
				rpcClient.Call("SiteRPC.Save", updateData, &retval)
				Session.Router.Navigate("/sites")
			}()
		})

		// Add an Action Grid
		loadTemplate("site-actions", "#action-grid", id)
		for _, ai := range doc.QuerySelectorAll(".action__item") {
			url := ai.(*dom.HTMLDivElement).GetAttribute("url")
			if url != "" {
				ai.AddEventListener("click", false, func(evt dom.Event) {
					url := evt.CurrentTarget().GetAttribute("url")
					Session.Router.Navigate(url)
				})
			}
		}

	}()

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

	print("show machine list for site", id)

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		req := shared.MachineReq{
			Channel: Session.Channel,
			SiteID:  id,
		}
		data := SiteMachineData{}

		rpcClient.Call("SiteRPC.Get", id, &data.Site)
		rpcClient.Call("SiteRPC.MachineList", &req, &data.Machines)
		loadTemplate("site-machine-list", "main", data)

		// Add a back handler on the header
		doc.GetElementByID("header").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate(fmt.Sprintf("/site/%d", id))
		})

		// Add a handler for clicking on a row
		t := doc.GetElementByID("machine-list")
		t.AddEventListener("click", false, func(evt dom.Event) {
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			Session.Router.Navigate("/machine/" + key)
		})

	}()
}
