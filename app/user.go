package main

import (
	"strconv"

	"github.com/go-humble/form"
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
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
		doc.QuerySelector("#nameField").(*dom.HTMLInputElement).Focus()

		// Setup the close button
		closeBtn := doc.QuerySelector(".md-up-close")
		if closeBtn != nil {
			closeBtn.AddEventListener("click", false, func(evt dom.Event) {
				evt.PreventDefault()
				el.Class().Remove("md-show")
			})
		}

		// Allow ESC to close dialog
		doc.QuerySelector("#user-profile-form").AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				evt.PreventDefault()
				el.Class().Remove("md-show")
			}
		})

		// Setup the save button
		saveBtn := doc.QuerySelector(".md-up-save")
		if saveBtn != nil {
			saveBtn.AddEventListener("click", false, func(evt dom.Event) {

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
					go func() {
						rpcClient.Call("UserRPC.Set", &req, &d)
						el.Class().Remove("md-show")
					}()
				}

			})
		}
	}()
}

func siteUserList(context *router.Context) {
	print("TODO - siteUserList")
}

// Display a list of users
func userList(context *router.Context) {

	go func() {
		users := []shared.User{}
		rpcClient.Call("UserRPC.List", Session.Channel, &users)

		form := formulate.ListForm{}
		form.New("fa-user", "Users List - All Users")

		// Define the layout
		form.Column("Username", "Username")
		form.Column("Name", "Name")
		form.Column("Email", "Email")
		form.Column("Mobile", "SMS")
		form.Column("Role", "Role")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/user/add")
		})

		form.RowEvent(func(key string) {
			Session.Router.Navigate("/user/" + key)
		})

		form.Render("user-list", "main", users)

	}()

}

// Edit an existing user
func userEdit(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {
		user := shared.User{}
		rpcClient.Call("UserRPC.Get", id, &user)

		BackURL := "/users"
		form := formulate.EditForm{}
		form.New("fa-user", "User Details - "+user.Name)

		// Layout the fields

		form.Row(2).
			Add(1, "Username", "text", "Username", `id="focusme"`).
			Add(1, "Password", "text", "Passwd", `id="focusme"`)

		form.Row(3).
			Add(1, "Name", "text", "Name", "").
			Add(1, "Email", "text", "Email", "").
			Add(1, "Mobile", "text", "SMS", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&user)
			data := shared.UserUpdateData{
				Channel: Session.Channel,
				User:    &user,
			}
			go func() {
				done = false
				rpcClient.Call("UserRPC.Update", data, &done)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &user)

	}()

}

// Add form for a new user
func userAdd(context *router.Context) {

	go func() {
		user := shared.User{}

		BackURL := "/users"
		form := formulate.EditForm{}
		form.New("fa-user", "Add New User")

		// Layout the fields

		form.Row(2).
			Add(1, "Username", "text", "Username", `id="focusme"`).
			Add(1, "Password", "text", "Passwd", `id="focusme"`)

		form.Row(3).
			Add(1, "Name", "text", "Name", "").
			Add(1, "Email", "text", "Email", "").
			Add(1, "Mobile", "text", "SMS", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate(BackURL)
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&user)
			data := shared.UserUpdateData{
				Channel: Session.Channel,
				User:    &user,
			}
			go func() {
				newID := 0
				rpcClient.Call("UserRPC.Insert", data, &newID)
				print("added user", newID)
				Session.Router.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &user)

	}()

}
