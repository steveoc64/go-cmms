package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
)

var appFn map[string]router.Handler

func enableRoutes(Role string) {

	print("enabling routes for role", Role)

	appFn = map[string]router.Handler{
		"dashboard": dashboard,
		"machines":  machineList,
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

	enableRoutes(Role)

	for _, v := range Routes {
		// print("route:", v.Route, v.Func)
		if f, ok := appFn[v.Func]; ok {
			// print("found a function called", v.Func)
			r.HandleFunc(v.Route, f)
		}
	}

}
