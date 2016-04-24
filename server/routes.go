package main

import "github.com/steveoc64/go-cmms/shared"

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin":
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
			{Route: "/sched/parts/{id}", Func: "sched-part-list"},
			{Route: "/sched/task/{id}", Func: "sched-task-list"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/task/complete/{id}", Func: "task-complete"},
			{Route: "/stoppages", Func: "stoppage-list"},
			{Route: "/stoppage/{id}", Func: "stoppage-edit"},
			{Route: "/stoppage/complete/{id}", Func: "stoppage-complete"},
			{Route: "/stoppage/newtask/{id}", Func: "stoppage-new-task"},
			{Route: "/stoppage/tasks/{id}", Func: "stoppage-task-list"},
			{Route: "/class/add", Func: "class-add"},
			{Route: "/class/select", Func: "class-select"},
			{Route: "/parts/{id}", Func: "part-list"},
			{Route: "/part/add/{id}", Func: "part-add"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/users", Func: "user-list"},
			{Route: "/user/{id}", Func: "user-edit"},
			{Route: "/user/add", Func: "user-add"},
			{Route: "/reports", Func: "reports"},
			{Route: "/util", Func: "util"},
		}
	case "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "site-list"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machine-list"},
			{Route: "/site/users/{id}", Func: "site-user-list"},
			{Route: "/site/tasks/{id}", Func: "site-task-list"},
			{Route: "/site/reports/{id}", Func: "site-reports"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/machine/sched/{machine}", Func: "machine-sched-list"},
			{Route: "/machine/reports/{machine}", Func: "machine-reports"},
			{Route: "/machine/stoppages/{machine}", Func: "machine-stoppage-list"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/stoppages", Func: "stoppage-list"},
			{Route: "/stoppage/{id}", Func: "stoppage-edit"},
			{Route: "/stoppage/tasks/{id}", Func: "stoppage-task-list"},
			{Route: "/class/select", Func: "class-select"},
			{Route: "/parts/{id}", Func: "part-list"},
			{Route: "/part/add/{id}", Func: "part-add"},
			{Route: "/part/{id}", Func: "part-edit"},
			{Route: "/reports", Func: "reports"},
		}
	case "Technician":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/task/complete/{id}", Func: "task-complete"},
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
