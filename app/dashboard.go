package main

import "github.com/go-humble/router"

// "github.com/steveoc64/go-cmms/shared"
// "honnef.co/go/js/dom"

func adminDashboard(context *router.Context) {
	print("in the admin dashboard function", context.Path)

	// w := dom.GetWindow()
	// doc := w.Document()

	loadTemplate("admin_dashboard", "main", nil)
}
