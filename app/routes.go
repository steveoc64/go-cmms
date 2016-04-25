package main

import (
	"errors"

	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func fixLinks() {
	Session.Router.InterceptLinks()
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
	Session.Router.InterceptLinks()
	return nil
}

func enableRoutes(Role string) {

	print("enabling routes for role", Role)

	switch Role {
	case "Admin", "Site Manager":
		Session.AppFn = map[string]router.Handler{
			"sitemap":               siteMap,
			"sitemachines":          siteMachines,
			"site-list":             siteList,
			"site-edit":             siteEdit,
			"site-add":              siteAdd,
			"site-machine-list":     siteMachineList,
			"site-machine-add":      siteMachineAdd,
			"site-user-list":        siteUserList,
			"site-task-list":        siteTaskList,
			"site-reports":          siteReports,
			"machine-edit":          machineEdit,
			"machine-sched-list":    machineSchedList,
			"machine-sched-add":     machineSchedAdd,
			"machine-reports":       machineReports,
			"machine-stoppage-list": machineStoppageList,
			"sched-edit":            schedEdit,
			"sched-task-list":       schedTaskList,
			"task-list":             taskList,
			"task-edit":             taskEdit,
			"task-part-list":        taskPartList,
			"task-complete":         taskComplete,
			"stoppage-list":         stoppageList,
			"stoppage-edit":         stoppageEdit,
			"stoppage-complete":     stoppageComplete,
			"stoppage-new-task":     stoppageNewTask,
			"stoppage-task-list":    stoppageTaskList,
			"class-select":          classSelect,
			"class-add":             classAdd,
			"part-list":             partList,
			"part-edit":             partEdit,
			"part-add":              partAdd,
			"user-list":             userList,
			"user-edit":             userEdit,
			"user-add":              userAdd,
			"reports":               adminReports,
			"util":                  adminUtils,
		}
	case "Technician":
		Session.AppFn = map[string]router.Handler{
			"sitemap":        siteMap,
			"sitemachines":   siteMachines,
			"task-list":      taskList,
			"task-edit":      taskEdit,
			"task-part-list": taskPartList,
			"task-complete":  taskComplete,
			"stoppages":      stoppageList,
			"parts":          partList,
			"reports":        technicianReports,
		}
	case "Floor":
		Session.AppFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
		}
	}
}

func initRouter() {
	print("initRouter")
	Session.Router = router.New()
	Session.Router.ShouldInterceptLinks = true
	Session.Router.HandleFunc("/", defaultRoute)
	Session.Router.Start()
}

func defaultRoute(context *router.Context) {
	print("default route")
}

func loadRoutes(Role string, Routes []shared.UserRoute) {

	// print("Loading new routing table")
	if Session.Router != nil {
		Session.Router.Stop()
	}
	Session.Router = router.New()
	Session.Router.ShouldInterceptLinks = true
	enableRoutes(Role)

	for _, v := range Routes {
		if f, ok := Session.AppFn[v.Func]; ok {
			// print("found a function called", v.Func)
			// print("adding route", v.Route, v.Func)
			Session.Router.HandleFunc(v.Route, f)
		}
	}
	Session.Router.Start()
}
