package main

import (
	"strconv"

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
		data := shared.User{}
		err := rpcClient.Call("UserRPC.Me", Session.Channel, &data)
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
						req := shared.UserUpdate{
							Channel: Session.Channel,
							ID:      data.ID,
							Name:    data.Name,
							Passwd:  data.Passwd,
							Email:   data.Email,
							SMS:     data.SMS,
						}
						// print("passing update req", req)
						rpcClient.Call("UserRPC.Set", &req, &d)
						el.Class().Remove("md-show")
					}

				}()
			})
		}
	}()
}

// Display a list of users
func usersList(context *router.Context) {

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := []shared.User{}
		rpcClient.Call("UserRPC.List", Session.Channel, &data)
		loadTemplate("user-list", "main", data)

		// Add a handler for clicking on a row
		doc.GetElementByID("user-list").AddEventListener("click", false, func(evt dom.Event) {
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			Session.Router.Navigate("/user/" + key)
		})

		// Add a handler for clicking on the add button
		doc.QuerySelector(".data-add-btn").AddEventListener("click", false, func(evt dom.Event) {
			print("add new user")
		})
	}()
}

type UserEditData struct {
	User  shared.User
	Sites []shared.Site
}

func userEdit(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		w := dom.GetWindow()
		doc := w.Document()

		data := UserEditData{}
		rpcClient.Call("UserRPC.Get", id, &data.User)
		rpcClient.Call("SiteRPC.List", Session.Channel, &data.Sites)
		loadTemplate("user-edit", "main", data)

		// Add handlers for this form
		doc.QuerySelector("legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate("/users")
		})
		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/users")
		})
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				// Parse the form element and get a form.Form object in return.
				f, err := form.Parse(doc.QuerySelector(".grid-form"))
				if err != nil {
					print("form parse error", err.Error())
					return
				}
				if err := f.Bind(&data.User); err != nil {
					print("form bind error", err.Error())
					return
				}
				updateData := shared.UserUpdateData{
					Channel: Session.Channel,
					User:    &data.User,
				}
				print("calling user.save = ", updateData)
				// retval := 0
				// rpcClient.Call("UserRPC.Save", &updateData, &retval)
				Session.Router.Navigate("/users")
			}()

		})

	}()

}

func siteUserList(context *router.Context) {
	print("TODO - Site User List")
}
