package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

var appFn map[string]router.Handler

var r *router.Router

func fixLinks() {
	r.InterceptLinks()
}

func loadTemplate(template string, data interface{}) {
	w := dom.GetWindow()
	doc := w.Document()

	t, err := GetTemplate(template)
	if err != nil {
		print(err.Error())
	}

	el := doc.QuerySelector("main")
	if err := t.ExecuteEl(el, data); err != nil {
		print(err.Error())
	}
	r.InterceptLinks()
}

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
			"skills":    skillsList,
		}
	case "Worker":
		appFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
			"homesite":     homeSite,
		}
	}
}

func initRouter() {
	r = router.New()
	r.ShouldInterceptLinks = true
	r.HandleFunc("/", defaultRoute)
	r.Start()
}

func defaultRoute(context *router.Context) {
	// print("default route")
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
		if f, ok := appFn[v.Func]; ok {
			// print("found a function called", v.Func)
			print("adding route", v.Route, v.Func)
			r.HandleFunc(v.Route, f)
		}
	}
	r.Start()
}
