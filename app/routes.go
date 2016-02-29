package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
)

var appFn map[string]router.Handler

func enableRoutes(Role string) {

	print("enabling routes for role", Role)

	switch Role {
	case "Admin":
		appFn = map[string]router.Handler{
			"dashboard": adminDashboard,
			"sites":     siteList,
			"machines":  machineList,
			"events":    eventList,
			"tools":     toolList,
			"parts":     partsList,
			"vendors":   vendorList,
			"users":     usersList,
			"reports":   adminReports,
		}
	case "Worker":
		appFn = map[string]router.Handler{
			"sitemap":  siteMap,
			"machines": machineList,
		}
	}
}

var r *router.Router

func initRouter() {
	r = router.New()
	r.ShouldInterceptLinks = true
	r.HandleFunc("/", defaultRoute)
	r.Start()
}

func defaultRoute(context *router.Context) {
	print("default route")
}

func loadRoutes(Role string, Routes []shared.UserRoute) {

	print("Loading new routing table")
	enableRoutes(Role)
	if r != nil {
		r.Stop()
	}
	r = router.New()
	r.ShouldInterceptLinks = true

	for _, v := range Routes {
		// print("route:", v.Route, v.Func)
		if f, ok := appFn[v.Func]; ok {
			// print("found a function called", v.Func)
			print("adding route", v.Route, v.Func)
			r.HandleFunc(v.Route, f)
		}
	}
	r.Start()

}
