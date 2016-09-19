package main

import (
	"errors"

	"itrak-cmms/shared"

	"github.com/go-humble/router"
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
			"stops":                 stops,
			"site-list":             siteList,
			"site-edit":             siteEdit,
			"site-add":              siteAdd,
			"site-sched-list":       siteSchedList,
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
			"task-invoices":         taskInvoices,
			"task-invoice":          taskInvoice,
			"task-invoice-add":      taskInvoiceAdd,
			"stoppage-list":         stoppageList,
			"stoppage-edit":         stoppageEdit,
			"stoppage-complete":     stoppageComplete,
			"stoppage-new-task":     stoppageNewTask,
			"stoppage-task-list":    stoppageTaskList,
			// "class-select":          classSelect,
			"class-select":           partsList,
			"class-add":              classAdd,
			"part-list":              partList,
			"part-edit":              partEdit,
			"part-add":               partAdd,
			"user-list":              userList,
			"user-edit":              userEdit,
			"user-add":               userAdd,
			"reports":                adminReports,
			"util":                   adminUtils,
			"hashtags":               hashtagList,
			"hashtag-add":            hashtagAdd,
			"hashtag-edit":           hashtagEdit,
			"hashtag-used":           hashtagUsed,
			"sms-list":               SMSList,
			"machine-types":          machineTypes,
			"machine-type-add":       machineTypeAdd,
			"machine-type-edit":      machineTypeEdit,
			"machine-type-machines":  machineTypeMachines,
			"machine-type-stoppages": machineTypeStoppages,
			"machine-type-tools":     machineTypeTools,
			"machine-type-tool-add":  machineTypeToolAdd,
			"machine-type-tool-edit": machineTypeToolEdit,
			"machine-type-parts":     machineTypeParts,
			"phototest":              phototest,
			"phototest-edit":         phototestEdit,
			"phototest-add":          phototestAdd,
			"testeditor":             testeditor,
			"usersonline":            usersOnline,
		}
	case "Technician":
		Session.AppFn = map[string]router.Handler{
			"sitemap":        siteMap,
			"sitemachines":   siteMachines,
			"task-list":      taskList,
			"task-edit":      taskEdit,
			"task-part-list": taskPartList,
			"stoppages":      stoppageList,
			"parts":          partList,
			"reports":        technicianReports,
			"diary":          technicianDiary,
		}
	case "Floor":
		Session.AppFn = map[string]router.Handler{
			"sitemap":      siteMap,
			"sitemachines": siteMachines,
		}
	}

	w := dom.GetWindow()
	doc := w.Document()

	if el := doc.QuerySelector("#show-image"); el != nil {
		// print("Adding click event for photo view")
		el.AddEventListener("click", false, func(evt dom.Event) {
			el.Class().Remove("md-show")
			// doc.QuerySelector("#show-image").Class().Remove("md-show")
		})
	}

}

func initRouter() {
	// print("initRouter")
	Session.Subscriptions = make(map[string]MessageFunction)
	Session.ID = make(map[string]int)

	Session.Router = router.New()
	Session.Router.ShouldInterceptLinks = true
	Session.Router.HandleFunc("/", defaultRoute)
	Session.Router.Start()

}

func defaultRoute(context *router.Context) {
	print("Nav to Default Route")
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
