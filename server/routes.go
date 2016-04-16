package main

import "github.com/steveoc64/go-cmms/shared"

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin", "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/site/add", Func: "site-add"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machines"},
			{Route: "/site/users/{id}", Func: "site-users"},
			{Route: "/site/tasks/{id}", Func: "site-tasks"},
			{Route: "/site/reports/{id}", Func: "site-reports"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/site/machine/add/{id}", Func: "site-machine-add"},
			{Route: "/machine/schedules/{machine}", Func: "machine-sched-list"},
			{Route: "/machine/reports/{machine}", Func: "machine-reports"},
			{Route: "/machine/stoppages/{machine}", Func: "machine-stoppage-list"},
			{Route: "/sched/{id}", Func: "sched-edit"},
			{Route: "/tasks", Func: "tasks"},
			{Route: "/stoppages", Func: "stoppages"},
			{Route: "/parts", Func: "parts"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/users", Func: "users"},
			{Route: "/user/{id}", Func: "user-edit"},
			{Route: "/user/new", Func: "user-new"},
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
