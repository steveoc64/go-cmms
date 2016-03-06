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
	sites := []shared.Site{}
	err := rpcClient.Call("SiteRPC.List", channelID, &sites)
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
}

// Show all the machines at a Site
func siteMachines(context *router.Context) {
	idStr := context.Params["site"]
	print("in the machineDet function", idStr)
	id, _ := strconv.Atoi(idStr)

	// Get a list of machines at this site
	req := shared.MachineReq{
		Channel: channelID,
		SiteID:  id,
	}
	machines := []shared.Machine{}

	go func() {
		err := rpcClient.Call("SiteRPC.MachineList", &req, &machines)
		if err != nil {
			print("RPC error", err.Error())
		} else {
			loadTemplate("sitemachines", machines)
		}
	}()

}

// Show a list of all sites
func siteList(context *router.Context) {
	loadTemplate("sitelist", nil)
}
