package main

import (
	"fmt"
	"strconv"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
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
	Session.Subscribe("sitestatus", _siteMap)
	Session.Subscribe("usersites", _siteMap)
	go _siteMap("show", 1)
}

func _siteMap(action string, id int) {

	if action == "usersite" {
		if id != Session.UserID {
			return
		}
	}

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
				Session.Navigate("/sitemachines/" + id)
			})
		}

		// Add an Action Grid depending on which role the user is logged in as
		// print("user role =", Session.UserRole)
		switch Session.UserRole {
		case "Admin":
			loadTemplate("admin-actions", "#action-grid", nil)
		case "Site Manager":
			if Session.CanAllocate {
				loadTemplate("site-supermanager-actions", "#action-grid", nil)

			} else {
				loadTemplate("site-manager-actions", "#action-grid", nil)

			}
		case "Technician":
			loadTemplate("technician-actions", "#action-grid", nil)
		}

		// Add a click handler to navigate to the page
		for _, ai := range doc.QuerySelectorAll(".action__item") {
			url := ai.(*dom.HTMLDivElement).GetAttribute("url")
			if url != "" {
				ai.AddEventListener("click", false, func(evt dom.Event) {
					url := evt.CurrentTarget().GetAttribute("url")
					Session.Navigate(url)
				})
			}
		}
	}
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
			Session.Navigate("/")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/site/add")
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/site/" + key)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
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
			Session.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&site)
			go func() {
				newID := 0
				rpcClient.Call("SiteRPC.Insert", shared.SiteRPCData{
					Channel: Session.Channel,
					Site:    &site,
				}, &newID)
				print("added site", newID)
				Session.Navigate(BackURL)
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

		rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &site)
		rpcClient.Call("SiteRPC.List", Session.Channel, &sites)
		rpcClient.Call("UserRPC.GetManagers", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      0,
		}, &allManagers)
		rpcClient.Call("UserRPC.GetManagers", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &managers)
		rpcClient.Call("UserRPC.GetTechnicians", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &technicians)

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
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			site.ID = id
			go func() {
				done := false
				rpcClient.Call("SiteRPC.Delete", shared.SiteRPCData{
					Channel: Session.Channel,
					Site:    &site,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&site)
			go func() {
				done := false
				rpcClient.Call("SiteRPC.Update", shared.SiteRPCData{
					Channel: Session.Channel,
					Site:    &site,
				}, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &site)

		// And attach actions
		form.ActionGrid("site-actions", "#action-grid", site.ID, func(url string) {
			Session.Navigate(url)
		})

	}()

}

func siteReports(context *router.Context) {
	print("TODO - siteReports")
}

// Show all the machines at a Site, for the case where we have > 1 site
// Param:  site

func siteMachines(context *router.Context) {
	id, err := strconv.Atoi(context.Params["site"])
	if err != nil {
		print(err.Error())
		return
	}
	Session.ID["site"] = id

	Session.Subscribe("sitestatus", _siteMachines)
	go _siteMachines("edit", id)
}

func _siteMachines(action string, id int) {

	BackURL := "/"

	switch action {
	case "edit":
		print("manually edit")
	case "delete":
		if id != Session.ID["site"] {
			return
		}
		print("current record has been deleted")
		Session.Navigate(BackURL)
		return
	default:
		// run it anyway, because there is a chance that the
		// map lights are going to change, even if the event
		// is not for this site !!!
		print("refresh due to site status change")
		id = Session.ID["site"]

		// if id != Session.ID["site"] {
		// 	return
		// }
	}

	// Get a list of machines at this site
	data := SiteMachineData{}

	nonTools := []string{
		"Electrical", "Hydraulic", "Pnuematic", "Lube", "Printer", "Console", "Uncoiler", "Rollbed", "Conveyor", "Encoder", "StripGuide",
	}

	RefreshURL := fmt.Sprintf("/sitemachines/%d", id)

	data.MultiSite = true
	rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
		Channel: Session.Channel,
		ID:      id,
	}, &data.Site)
	// print("Site =", data.Site)

	err := rpcClient.Call("SiteRPC.MachineList", shared.SiteRPCData{
		Channel: Session.Channel,
		ID:      id,
	}, &data.Machines)

	print("got machines", data.Machines)
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
			Session.Navigate("/")
		})

		switch Session.UserRole {
		case "Admin", "Technician", "Site Manager":
			// Attach a menu opener for each machine
			for _, v := range data.Machines {
				mid := fmt.Sprintf("machine-div-%d", v.ID)
				machinediv := doc.GetElementByID(mid)
				machinediv.AddEventListener("click", false, func(evt dom.Event) {
					machine_id, _ := strconv.Atoi(machinediv.GetAttribute("machine-id"))
					machinemenu := doc.GetElementByID("machine-menu").(*dom.BasicHTMLElement)
					print("click", evt)

					// evt.Target points to the SVG sub element that took the click
					t1 := evt.Target()
					print("t1", t1)

					// now we have to walk UP the tree until we find a containing parent
					// element that has a declared tool type
					foundOne := false
					hitEnd := false
					tooltype := ""

					for !foundOne && !hitEnd {
						tooltype = t1.GetAttribute("tooltype")
						print("tooltype", tooltype)
						if tooltype != "" {
							print("clickd on ", t1.TagName(), " with tooltype =", tooltype)
							foundOne = true
						} else {
							t1 = t1.ParentElement()
							print("stepping up to parent", t1.TagName())
							switch t1.TagName() {
							case "div", "body", "DIV", "BODY", "HTML":
								hitEnd = true
							}
						}
					}

					d := shared.RaiseIssue{}
					d.Channel = Session.Channel
					d.MachineID = machine_id
					for _, m1 := range data.Machines {
						if m1.ID == d.MachineID {
							d.Machine = &m1
						}
					}

					doitNow := false

					if foundOne {
						switch tooltype {
						case "Tool":
							toolid := t1.GetAttribute("toolid")
							print("tool id =", toolid)
							d.CompID, _ = strconv.Atoi(toolid)
							d.IsTool = true
							doitNow = true
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
						case "Component":
							component := t1.GetAttribute("component")
							print("component =", component)
							d.NonTool = component
							doitNow = true
						default:
							print("clicked on something other than a tool or component", tooltype)
						}
					}

					if doitNow {
						// load a raise issue form from a template
						loadTemplate("raise-comp-issue", "#raise-comp-issue", d)
						doc.QuerySelector("#raise-comp-issue").Class().Add("md-show")

						// Add the machine diagram to the form
						if d.Machine != nil {
							loadTemplate("machine-diag", "#issue-machine-diag", d)
						}

						// Handle button clicks
						doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
							evt.PreventDefault()
							doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
						})
						doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
							evt.PreventDefault()
							d.Channel = Session.Channel
							d.Descr = doc.QuerySelector("#evtdesc").(*dom.HTMLTextAreaElement).Value
							d.Photo = doc.QuerySelector("[name=PhotoPreview]").(*dom.HTMLImageElement).GetAttribute("src")
							go func() {
								newID := 0
								rpcClient.Call("EventRPC.Raise", d, &newID)
								print("Raised new event", newID)
								Session.Navigate(RefreshURL)
							}()
							doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
						})

						// add a handler on the photo field
						setPhotoOnlyField("Photo")

					} else {

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
											evt.PreventDefault()
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
												Session.Navigate(RefreshURL)
											}()
											doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
										})
									})
								}
								// attach event listeners to each non-tool menu item
								for i := range nonTools {
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
											evt.PreventDefault()
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
												Session.Navigate(RefreshURL)
											}()
											doc.QuerySelector("#raise-comp-issue").Class().Remove("md-show")
										})
									})
								}
							} // is matching machine
						} // for loop to construct the side menu
					} // if we clicked outside of a known tool or component

				})
			} // range
		} // switch user role
	} // else
}

func siteSchedList(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		site := shared.Site{}
		rpcClient.Call("SiteRPC.Get", shared.SiteRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &site)
		// print("site", site)
		tasks := []shared.SchedTask{}
		rpcClient.Call("TaskRPC.ListSiteSched", shared.TaskRPCData{
			Channel: Session.Channel,
			ID:      id,
		}, &tasks)
		// print("tasks", tasks)

		BackURL := fmt.Sprintf("/site/%d", id)

		form := formulate.ListForm{}
		form.New("fa-wrench", "All Sched Maints for - "+site.Name)

		// Define the layout
		form.Column("Machine", "MachineName")
		form.Column("Tool / Component", "Component")
		form.Column("Frequency", "ShowFrequency")
		form.Column("Description", "Descr")
		form.Column("$ Labour", "LabourCost")
		form.Column("$ Materials", "MaterialCost")
		form.Column("Duration", "DurationDays")
		form.Column("Job Status", "ShowPaused")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			print("TODO - popup a machine selection dialog, and then nav to /machine/sched/add/machine_id")
			// Session.Navigate(fmt.Sprintf("/machine/sched/add/%d", id))
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/sched/" + key)
		})

		form.Render("site-sched-list", "main", tasks)
	}()
}
