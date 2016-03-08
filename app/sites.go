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
			loadTemplate("sitemap", data)

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
			print("SiteMachine status report", data.Status)

			loadTemplate("sitemachines", data)
			w := dom.GetWindow()
			doc := w.Document()
			austmap := doc.GetElementByID("austmap")
			austmap.AddEventListener("click", false, func(evt dom.Event) {
				r.Navigate("/")
			})

			// Attach a menu opener for each machine
			for _, v := range data.Machines {
				machinediv := doc.GetElementByID(fmt.Sprintf("machine-div-%d", v.ID))
				machinediv.AddEventListener("click", false, func(evt dom.Event) {
					machinediv := evt.Target().(*dom.BasicHTMLElement)
					print("machinediv =", machinediv)
					mid1 := machinediv.GetAttribute("id")
					mid2 := machinediv.GetAttribute("machine-id")
					print("clicked on machine id", mid1, mid2) // this is a string at this point
					machinemenu := doc.GetElementByID("machine-menu").(*dom.BasicHTMLElement)
					machinemenu.Class().Toggle("cbp-spmenu-open")
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
			loadTemplate("sitemachines", data)
		}
	}()
}

// Show a list of all sites
func siteList(context *router.Context) {
	loadTemplate("sitelist", nil)
}
