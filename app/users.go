package main

import (
	"github.com/go-humble/form"
	"github.com/go-humble/router"
	"github.com/steveoc64/go-cmms/shared"
	"honnef.co/go/js/dom"
)

// Display modal dialog of user profile, with the ability to update info
func userProfile() {

	go func() {
		w := dom.GetWindow()
		doc := w.Document()
		// print("edit user profile")

		// TODO - get the user profile data from the backend
		data := shared.UserProfile{}
		err := rpcClient.Call("UserProfileRPC.Get", channelID, &data)
		if err != nil {
			print("RPC error", err.Error())
			return
		}

		loadTemplate("user-profile", "#user-profile", data)
		el := doc.QuerySelector("#user-profile")
		el.Class().Add("md-show")

		// Setup the close button
		closeBtn := doc.QuerySelector(".md-up-close")
		if closeBtn != nil {
			closeBtn.AddEventListener("click", false, func(evt dom.Event) {
				evt.PreventDefault()
				el.Class().Remove("md-show")
			})
		}

		// Setup the save button
		saveBtn := doc.QuerySelector(".md-up-save")
		if saveBtn != nil {
			saveBtn.AddEventListener("click", false, func(evt dom.Event) {

				go func() {
					evt.PreventDefault()

					formEl := doc.QuerySelector("#user-profile-form")
					f, err := form.Parse(formEl)
					if err != nil {
						print("Form Parse Error", err.Error())
					}
					data.Name, err = f.GetString("Name")
					if err != nil {
						print("Name", err.Error())
					}
					data.Email, err = f.GetString("Email")
					if err != nil {
						print("email", err.Error())
					}
					data.SMS, err = f.GetString("SMS")
					if err != nil {
						print("sms", err.Error())
					}
					p1, _ := f.GetString("p1")
					p2, _ := f.GetString("p2")

					// print("updated data =", data, p1, p2)
					if p1 != p2 {
						w.Alert("Passwords do not match")
					} else {
						if p1 != "" {
							data.Passwd = p1
						}
						d := false
						req := shared.UserProfileUpdate{
							Channel: channelID,
							ID:      data.ID,
							Name:    data.Name,
							Passwd:  data.Passwd,
							Email:   data.Email,
							SMS:     data.SMS,
						}
						// print("passing update req", req)
						rpcClient.Call("UserProfileRPC.Set", &req, &d)
						el.Class().Remove("md-show")
					}

				}()
			})
		}
	}()
}

// Display a list of users
func usersList(context *router.Context) {
	print("in the usersList function")
}
