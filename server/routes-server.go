package main

import "itrak-cmms/shared"

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/stops/{site}", Func: "stops"},
			{Route: "/sites", Func: "site-list"},
			{Route: "/site/add", Func: "site-add"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machine-list"},
			{Route: "/site/machine/add/{id}", Func: "site-machine-add"},
			{Route: "/site/users/{id}", Func: "site-user-list"},
			{Route: "/site/tasks/{id}", Func: "site-task-list"},
			{Route: "/site/reports/{id}", Func: "site-reports"},
			{Route: "/site/sched/{id}", Func: "site-sched-list"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/machine/sched/{machine}", Func: "machine-sched-list"},
			{Route: "/machine/sched/add/{machine}", Func: "machine-sched-add"},
			{Route: "/machine/reports/{machine}", Func: "machine-reports"},
			{Route: "/machine/stoppages/{machine}", Func: "machine-stoppage-list"},
			{Route: "/sched/{id}", Func: "sched-edit"},
			{Route: "/sched/task/{id}", Func: "sched-task-list"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/task/complete/{id}", Func: "task-complete"},
			{Route: "/task/invoices/{id}", Func: "task-invoices"},
			{Route: "/task/invoice/{id}", Func: "task-invoice"},
			{Route: "/task/invoice/add/{id}", Func: "task-invoice-add"},
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
			{Route: "/sms", Func: "sms-list"},
			{Route: "/hashtags", Func: "hashtags"},
			{Route: "/hashtag/add", Func: "hashtag-add"},
			{Route: "/hashtag/{id}", Func: "hashtag-edit"},
			{Route: "/hash/used/{id}", Func: "hashtag-used"},
			{Route: "/machinetypes", Func: "machine-types"},
			{Route: "/machinetype/add", Func: "machine-type-add"},
			{Route: "/machinetype/{id}", Func: "machine-type-edit"},
			{Route: "/machinetype/{id}/tools", Func: "machine-type-tools"},
			{Route: "/machinetype/{id}/tool/add", Func: "machine-type-tool-add"},
			{Route: "/machinetype/{id}/tool/{tool}", Func: "machine-type-tool-edit"},
			{Route: "/machinetype/{id}/parts", Func: "machine-type-parts"},
			{Route: "/machinetype/{id}/machines", Func: "machine-type-machines"},
			{Route: "/machinetype/{id}/stoppages", Func: "machine-type-stoppages"},
			{Route: "/phototest", Func: "phototest"},
			{Route: "/phototest/{id}", Func: "phototest-edit"},
			{Route: "/phototest/add", Func: "phototest-add"},
			{Route: "/testeditor", Func: "testeditor"},
			{Route: "/usersonline", Func: "usersonline"},
		}
	case "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/stops/{site}", Func: "stops"},
			{Route: "/sites", Func: "site-list"},
			{Route: "/site/{id}", Func: "site-edit"},
			{Route: "/site/machine/{id}", Func: "site-machine-list"},
			{Route: "/site/users/{id}", Func: "site-user-list"},
			{Route: "/site/tasks/{id}", Func: "site-task-list"},
			{Route: "/site/reports/{id}", Func: "site-reports"},
			{Route: "/site/sched/{id}", Func: "site-sched-list"},
			{Route: "/machine/{id}", Func: "machine-edit"},
			{Route: "/machine/sched/{machine}", Func: "machine-sched-list"},
			{Route: "/machine/reports/{machine}", Func: "machine-reports"},
			{Route: "/machine/stoppages/{machine}", Func: "machine-stoppage-list"},
			{Route: "/sched/{id}", Func: "sched-edit"},
			{Route: "/sched/task/{id}", Func: "sched-task-list"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/stoppages", Func: "stoppage-list"},
			{Route: "/stoppage/{id}", Func: "stoppage-edit"},
			{Route: "/stoppage/tasks/{id}", Func: "stoppage-task-list"},
			{Route: "/stoppage/complete/{id}", Func: "stoppage-complete"},
			{Route: "/stoppage/newtask/{id}", Func: "stoppage-new-task"},
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
			{Route: "/stops/{site}", Func: "stops"},
			{Route: "/sites", Func: "sites"},
			{Route: "/tasks", Func: "task-list"},
			{Route: "/task/{id}", Func: "task-edit"},
			{Route: "/task/parts/{id}", Func: "task-part-list"},
			{Route: "/stoppages", Func: "stoppages"},
			{Route: "/stoppage/{id}", Func: "stoppage-edit"},
			{Route: "/stoppage/tasks/{id}", Func: "stoppage-task-list"},
			{Route: "/stoppage/complete/{id}", Func: "stoppage-complete"},
			{Route: "/stoppage/newtask/{id}", Func: "stoppage-new-task"},
			{Route: "/parts", Func: "parts"},
			{Route: "/reports", Func: "reports"},
			{Route: "/diary", Func: "diary"},
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
