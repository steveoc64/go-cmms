package main

import "github.com/steveoc64/go-cmms/shared"

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machines"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/tasks", Func: "tasks"},
			{Route: "/stoppages", Func: "stoppages"},
			{Route: "/parts", Func: "parts"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/users", Func: "users"},
			{Route: "/user/{id}", Func: "user-edit"},
			{Route: "/reports", Func: "reports"},
		}
	case "Worker":
		// If this user has 1 site, they ony get 1 route
		// Otherwise they get a map and 1 route to show the machines at each site
		numSites := 1
		DB.SQL(`select count(*) from user_site where user_id=$1`, uid).QueryScalar(&numSites)
		if numSites == 1 {
			return []shared.UserRoute{
				{Route: "/", Func: "homesite"},
			}
		}
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/tasks", Func: "tasks"},
			{Route: "/stoppages", Func: "stoppages"},
			{Route: "/parts", Func: "parts"},
			{Route: "/reports", Func: "reports"},
		}
	case "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machines"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/tasks", Func: "tasks"},
			{Route: "/stoppages", Func: "stoppages"},
			{Route: "/parts", Func: "parts"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/users", Func: "users"},
			{Route: "/user/{id}", Func: "user-edit"},
			{Route: "/reports", Func: "reports"},
		}
	case "Floor":
		return []shared.UserRoute{
			{Route: "/", Func: "machines"},
		}
	case "Service Contractor":
		return []shared.UserRoute{
			{Route: "/", Func: "dashboard"},
			{Route: "/machines", Func: "machines"},
		}
	}
	return []shared.UserRoute{}

}
