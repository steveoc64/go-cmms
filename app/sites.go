package main

import (
	"fmt"
	"strconv"

	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// Display a map showing the status of all sites, with buttons to
// jump to specific sites
func siteMap(context *router.Context) {

	// Get a list of sites
	go func() {
		sites := []shared.Site{}
		err := rpcClient.Call("SiteRPC.UserList", channelID, &sites)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			// print("got an array of sites !!")

			w := dom.GetWindow()
			doc := w.Document()
			loadTemplate("sitemap", sites)

			// Attach listeners for each button
			for _, v := range sites {
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
			loadTemplate("sitemachines", data)
			w := dom.GetWindow()
			doc := w.Document()
			austmap := doc.GetElementByID("austmap")
			austmap.AddEventListener("click", false, func(evt dom.Event) {
				r.Navigate("/")
			})
		}
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
