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
		form.Column("Date", "GetDateSent")
		form.Column("Number", "GetNumber")
		form.Column("Message", "Message")
		form.Column("Reference", "Ref")
		// form.Column("Phone", "Phone")
		form.Column("Status", "GetStatus")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/util")
		})

		form.PrintEvent(func(evt dom.Event) {
			evt.PreventDefault()
			dom.GetWindow().Print()
		})

		form.Render("sms-list", "main", smsTrans)
	}()
}
