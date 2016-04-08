package main

import "github.com/go-humble/router"

func adminDashboard(context *router.Context) {
	print("in the admin dashboard function", context.Path)

	loadTemplate("admin-dashboard", "main", nil)
}
