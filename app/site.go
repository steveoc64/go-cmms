package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
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

// Display a map showing the status of all sites, with buttons to
// jump to specific sites
func siteMap(context *router.Context) {

	// print("in siteMap", Session)

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
			// print("user role =", Session.UserRole)
			switch Session.UserRole {
			case "Admin":
				loadTemplate("admin-actions", "#action-grid", nil)
			case "Site Manager":
				loadTemplate("site-manager-actions", "#action-grid", nil)
			case "Technician":
				loadTemplate("technician-actions", "#action-grid", nil)
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
		sites := []shared.Site{}
		rpcClient.Call("SiteRPC.List", Session.Channel, &sites)

		form := formulate.ListForm{}
		form.New("fa-industry", "Site List - All Sites")

		// Define the layout
		form.Column("Name", "Name")
		form.Column("Parent Site", "ParentSiteName")
		form.Column("Address", "Address")
		form.Column("Phone", "Phone")
		form.Column("Fax", "Fax")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/site/add")
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/site/" + key)
		})

		form.Render("site-list", "main", sites)
		// form.Render("site-list", "main", data)

	}()
}

type SiteEditData struct {
	Site  shared.Site
	Title string
	Sites []shared.Site
}

// Show an add form for the given site
func siteAdd(context *router.Context) {

	go func() {
		site := shared.Site{}
		sites := []shared.Site{}
		users := []shared.User{}

		rpcClient.Call("SiteRPC.List", Session.Channel, &sites)
		rpcClient.Call("UserRPC.List", Session.Channel, &users)

		BackURL := "/sites"
		title := "Add New Site"
		form := formulate.EditForm{}
		form.New("fa-industry", title)

		// Layout the fields
		form.Row(1).
			AddInput(1, "Name", "Name")

		form.Row(2).
			AddSelect(1, "Parent Site", "ParentSite", sites, "ID", "Name", 0, site.ParentSite).
			AddSelect(1, "Stock Site", "StockSite", sites, "ID", "Name", 0, site.StockSite)

		form.Row(1).
			AddInput(1, "Address", "Address")

		form.Row(2).
			AddInput(1, "Phone", "Phone").
			AddInput(1, "Fax", "Fax")

		form.Row(2).
			AddSelect(1, "Stoppage Alerts To", "AlertsTo", users, "ID", "Name", 0, site.AlertsTo).
			AddSelect(1, "Scheduled Tasks To", "TasksTo", users, "ID", "Name", 0, site.TasksTo)

		form.Row(1).
			AddTextarea(1, "Notes", "Notes")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&site)
			data := shared.SiteUpdateData{
				Channel: Session.Channel,
				Site:    &site,
			}
			go func() {
				newID := 0
				rpcClient.Call("SiteRPC.Insert", data, &newID)
				print("added site", newID)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &site)

	}()

}

// Show an edit form for the given site
func siteEdit(context *router.Context) {
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		site := shared.Site{}
		sites := []shared.Site{}
		allManagers := []shared.User{}
		managers := []shared.User{}
		technicians := []shared.User{}

		rpcClient.Call("SiteRPC.Get", id, &site)
		rpcClient.Call("SiteRPC.List", Session.Channel, &sites)
		rpcClient.Call("UserRPC.GetManagers", 0, &allManagers)
		rpcClient.Call("UserRPC.GetManagers", id, &managers)
		rpcClient.Call("UserRPC.GetTechnicians", id, &technicians)

		BackURL := "/sites"
		title := fmt.Sprintf("Site Details - %s", site.Name)
		form := formulate.EditForm{}
		form.New("fa-industry", title)

		// Layout the fields
		form.Row(2).
			AddInput(1, "Name", "Name").
			AddSelect(1, "Site Manager", "Manager", allManagers, "ID", "Name", 0, site.Manager)
		form.Row(2).
			AddSelect(1, "Parent Site", "ParentSite", sites, "ID", "Name", 0, site.ParentSite).
			AddSelect(1, "Stock Site", "StockSite", sites, "ID", "Name", 0, site.StockSite)

		form.Row(1).
			AddInput(1, "Address", "Address")

		form.Row(2).
			AddInput(1, "Phone", "Phone").
			AddInput(1, "Fax", "Fax")

		form.Row(2).
			AddSelect(1, "Stoppage Alerts To", "AlertsTo", managers, "ID", "Name", 0, site.AlertsTo).
			AddSelect(1, "Send Scheduled Tasks To", "TasksTo", technicians, "ID", "Name", 0, site.TasksTo)

		form.Row(1).
			AddTextarea(1, "Notes", "Notes")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			site.ID = id
			go func() {
				data := shared.SiteUpdateData{
					Channel: Session.Channel,
					Site:    &site,
				}
				done := false
				rpcClient.Call("SiteRPC.Delete", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&site)
			data := shared.SiteUpdateData{
				Channel: Session.Channel,
				Site:    &site,
			}
			go func() {
				done := false
				rpcClient.Call("SiteRPC.Update", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &site)

		// And attach actions
		form.ActionGrid("site-actions", "#action-grid", site.ID, func(url string) {
			Session.Router.Navigate(url)
		})

	}()

}

func siteReports(context *router.Context) {
	print("TODO - siteReports")
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

	RefreshURL := fmt.Sprintf("/sitemachines/%d", id)

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

			switch Session.UserRole {
			case "Admin", "Technician":
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
										d := shared.RaiseIssue{}
										d.Channel = Session.Channel
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
											evt.PreventDefault()
											d.Channel = Session.Channel
											d.Descr = doc.QuerySelector("#evtdesc").(*dom.HTMLTextAreaElement).Value
											go func() {
												newID := 0
												rpcClient.Call("EventRPC.Raise", d, &newID)
												print("Raised new event", newID)
												Session.Router.Navigate(RefreshURL)
											}()
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
										d := shared.RaiseIssue{}
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
											evt.PreventDefault()
											d.Channel = Session.Channel
											d.Descr = doc.QuerySelector("#evtdesc").(*dom.HTMLTextAreaElement).Value

											go func() {
												newID := 0
												rpcClient.Call("EventRPC.Raise", d, &newID)
												print("Raised new event", newID)
												Session.Router.Navigate(RefreshURL)
											}()
											doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
										})
									})
								}
							} // is matching machine
						}

					})
				} // range
			} // switch user role
		} // else
	}()
}
