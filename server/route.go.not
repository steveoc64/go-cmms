package main

import (
	"github.com/steveoc64/go-cmms/shared"
	"log"
	"time"
)

type RouteRPC struct{}

func (r *RouteRPC) Get(req *shared.RouteReq, res *shared.RouteResponse) error {

	start := time.Now()

	conn := Connections.Get(req.Channel)
	log.Println("Conn =", conn)

	switch conn.UserRole {
	case "Admin":
		getAdminTemplate(req.Name, res)
	case "Worker":
		getWorkerTemplate(req.Name, res)
	case "Site Manager":
		getSiteMgrTemplate(req.Name, res)
	case "Service Contractor":
		getServiceContractorTemplate(req.Name, res)
	case "Floor":
		getFloorTemplate(req.Name, res)
	}

	log.Printf(`Route.Get ->
    » (%d,%s)
    « (%s) 
    = %s`,
		req.Channel, req.Name,
		res.Template,
		time.Since(start))

	return nil
}

func getAdminTemplate(name string, res *shared.RouteResponse) {

}

func getWorkerTemplate(name string, res *shared.RouteResponse) {

	switch name {
	case "SiteMap":
		res.Template = `
<div>
This is a little template
</div>
		`
	}
}

func getSiteMgrTemplate(name string, res *shared.RouteResponse) {

}

func getServiceContractorTemplate(name string, res *shared.RouteResponse) {

}

func getFloorTemplate(name string, res *shared.RouteResponse) {

}
