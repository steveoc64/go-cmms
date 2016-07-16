package main

import (
	"itrak-cmms/shared"

	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"honnef.co/go/js/dom"
)

func eventList(context *router.Context) {
	print("TODO - eventList")
}

func workOrderList(context *router.Context) {
	print("TODO - workOrderList")
}

func SMSList(context *router.Context) {
	Session.Subscribe("sms", _SMSList)
	go _SMSList("list", 0)
}

func _SMSList(action string, id int) {

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
}
