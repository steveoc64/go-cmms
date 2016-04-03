package main

import (
	"errors"

	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

var appFn map[string]router.Handler

var r *router.Router

func fixLinks() {
	r.InterceptLinks()
}

// Load a template and attach it to the specified element in the doc
func loadTemplate(template string, selector string, data interface{}) error {
	w := dom.GetWindow()
	doc := w.Document()

	t, err := GetTemplate(template)
	if t == nil {
		print("Failed to load template", template)
		return errors.New("Invalid template")
	}
	if err != nil {
		print(err.Error())
		return err
	}

	el := doc.QuerySelector(selector)
	if el == nil {
		print("Could not find selector", selector)
		return errors.New("Invalid selector")
	}
	if err := t.ExecuteEl(el, data); err != nil {
		print(err.Error())
		return err
	}
	r.InterceptLinks()
	return nil
}

func enableRoutes(Role string) {

	print("enabling routes for role", Role)

	switch Role {
	case "Admin":
		appFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
			"dashboard":    adminDashboard,
			"sites":        siteList,
			"machines":     machineList,
			"events":       eventList,
			"workorders":   workOrderList,
			"tools":        toolList,
			"parts":        partsList,
			"vendors":      vendorList,
			"users":        usersList,
			"reports":      adminReports,
			"skills":       skillsList,
		}
	case "SiteManager":
		appFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
			"tasks":        homeSite,
			"users":        usersList,
			"events":       eventList,
			"workorders":   workOrderList,
		}
	case "Worker":
		appFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
			"homesite":     homeSite,
		}
	case "Floor":
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

	// print("Loading new routing table")
	enableRoutes(Role)
	if r != nil {
		r.Stop()
	}
	r = router.New()
	r.ShouldInterceptLinks = true

	for _, v := range Routes {
		if f, ok := appFn[v.Func]; ok {
			// print("found a function called", v.Func)
			// print("adding route", v.Route, v.Func)
			r.HandleFunc(v.Route, f)
		}
	}
	r.Start()
}
