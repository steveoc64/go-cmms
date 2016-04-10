package main

import "github.com/go-humble/router"

func adminDashboard(context *router.Context) {
	loadTemplate("admin-dashboard", "main", nil)
}
