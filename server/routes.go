package main

import "github.com/steveoc64/go-cmms/shared"

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin", "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "site-list"},
			{Route: "/site/add", Func: "site-add"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machine-list"},
			{Route: "/site/machine/add/{id}", Func: "site-machine-add"},
			{Route: "/site/users/{id}", Func: "site-user-list"},
			{Route: "/site/tasks/{id}", Func: "site-task-list"},
			{Route: "/site/reports/{id}", Func: "site-reports"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/machine/sched/{machine}", Func: "machine-sched-list"},
			{Route: "/machine/sched/add/{machine}", Func: "machine-sched-add"},
			{Route: "/machine/reports/{machine}", Func: "machine-reports"},
			{Route: "/machine/stoppages/{machine}", Func: "machine-stoppage-list"},
			{Route: "/sched/{id}", Func: "sched-edit"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/stoppages", Func: "stoppage-list"},
			{Route: "/stoppage/{id}", Func: "stoppage-edit"},
			{Route: "/parts", Func: "part-list"},
			{Route: "/part/add", Func: "part-add"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/users", Func: "user-list"},
			{Route: "/user/{id}", Func: "user-edit"},
			{Route: "/user/add", Func: "user-add"},
			{Route: "/reports", Func: "reports"},
		}
	case "Technician":
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
