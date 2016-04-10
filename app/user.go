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
			Session.Router.Navigate("/user/new")
		})
	}()
}

type UserEditData struct {
	User  shared.User
	Title string
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
		data.Title = "User Details - " + data.User.Name
		loadTemplate("user-edit", "main", data)
		doc.QuerySelector("#focusme").(*dom.HTMLInputElement).Focus()

		// Add handlers for this form
		doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
			Session.Router.Navigate("/users")
		})
		doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			Session.Router.Navigate("/users")
		})
		// Allow ESC to close dialog
		doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				evt.PreventDefault()
				Session.Router.Navigate("/users")
			}
		})
		doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
			// go func() {
			evt.PreventDefault()
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

			req := shared.UserUpdate{
				Channel:  Session.Channel,
				ID:       id,
				Username: data.User.Username,
				Name:     data.User.Name,
				Passwd:   data.User.Passwd,
				Email:    data.User.Email,
				SMS:      data.User.SMS,
			}

			d := false
			go func() {
				rpcClient.Call("UserRPC.Save", &req, &d)
			}()
			Session.Router.Navigate("/users")
			// }()

		})

	}()

}

func siteUserList(context *router.Context) {
	print("TODO - siteUserList")
}

func userNew(context *router.Context) {

	w := dom.GetWindow()
	doc := w.Document()

	data := UserEditData{}
	// rpcClient.Call("SiteRPC.List", Session.Channel, &data.Sites)
	data.Title = "Add New User"
	loadTemplate("user-edit", "main", data)
	doc.QuerySelector("#focusme").(*dom.HTMLInputElement).Focus()

	// Add handlers for this form
	doc.QuerySelector("#legend").AddEventListener("click", false, func(evt dom.Event) {
		Session.Router.Navigate("/users")
	})
	doc.QuerySelector(".md-close").AddEventListener("click", false, func(evt dom.Event) {
		evt.PreventDefault()
		Session.Router.Navigate("/users")
	})
	doc.QuerySelector(".grid-form").AddEventListener("keyup", false, func(evt dom.Event) {
		if evt.(*dom.KeyboardEvent).KeyCode == 27 {
			Session.Router.Navigate("/users")
		}
	})
	doc.QuerySelector(".md-save").AddEventListener("click", false, func(evt dom.Event) {
		// go func() {
		evt.PreventDefault()
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

		req := shared.UserUpdate{
			Channel:  Session.Channel,
			ID:       0,
			Username: data.User.Username,
			Name:     data.User.Name,
			Passwd:   data.User.Passwd,
			Email:    data.User.Email,
			SMS:      data.User.SMS,
		}

		d := 0
		go func() {
			rpcClient.Call("UserRPC.Insert", &req, &d)
			print("new record = ", d)
		}()
		Session.Router.Navigate("/users")

	})

}
