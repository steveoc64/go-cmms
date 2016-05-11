package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

func eventList(context *router.Context) {
	print("TODO - eventList")
}

func workOrderList(context *router.Context) {
	print("TODO - workOrderList")
}

func SMSList(context *router.Context) {

	go func() {
		smsTrans := []shared.SMSTrans{}
		rpcClient.Call("SMSRPC.List", Session.Channel, &smsTrans)

		form := formulate.ListForm{}
		form.New("fa-phone", "SMS Traffic Log")

		// Define the layout
		form.Column("Date", "GetDate")
		form.Column("Number", "GetNumber")
		form.Column("Message", "Message")
		form.Column("Reference", "Ref")
		form.Column("Phone", "Phone")
		form.Column("Status", "GetStatus")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/util")
		})

		// form.NewRowEvent(func(evt dom.Event) {
		// 	evt.PreventDefault()
		// 	Session.Router.Navigate("/site/add")
		// })

		// form.RowEvent(func(key string) {
		// 	Session.Router.Navigate("/site/" + key)
		// })

		form.Render("sms-list", "main", smsTrans)
		// form.Render("site-list", "main", data)

	}()
}
