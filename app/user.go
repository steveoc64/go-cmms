package main

import (
	"fmt"
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
	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	go func() {

		site := shared.Site{}
		rpcClient.Call("SiteRPC.Get", id, &site)

		BackURL := fmt.Sprintf("/site/%d", id)
		form := formulate.EditForm{}
		form.New("fa-user", "User Access for Site - "+site.Name)

		// Layout the fields

		form.Row(1).
			Add(1, "Users", "div", "Users", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &site)

		// And load in a the sites array
		// Need a data set that is a slice of :
		// - Site ID
		// - Site Name
		// - Bool, whether the user has access to this site
		// On toggling a site, send an RPC to the backend to toggle the state
		siteUsers := []shared.SiteUser{}
		req := shared.UserSiteRequest{
			Channel: Session.Channel,
			User:    nil,
			Site:    &site,
		}
		rpcClient.Call("UserRPC.GetSiteUsers", req, &siteUsers)
		loadTemplate("site-users-array", "[name=Users]", siteUsers)

		// add a click handler for the sites array
		w := dom.GetWindow()
		doc := w.Document()

		if el := doc.QuerySelector("[name=Users]"); el != nil {

			el.AddEventListener("click", false, func(evt dom.Event) {
				clickedOn := evt.Target()
				switch clickedOn.TagName() {
				case "INPUT":
					ie := clickedOn.(*dom.HTMLInputElement)
					key, _ := strconv.Atoi(ie.GetAttribute("key"))
					data := shared.UserSiteSetRequest{
						Channel: Session.Channel,
						SiteID:  site.ID,
						UserID:  key,
						Role:    "",
						IsSet:   ie.Checked,
					}

					go func() {
						done := false
						rpcClient.Call("UserRPC.SetSite", data, &done)
					}()
				}

			})
		}

	}()

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

		if Session.UserRole == "Admin" {
			form.Column("Hourly Rate", "HourlyRate")
		}

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/")
		})

		form.NewRowEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate("/user/add")
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		form.RowEvent(func(key string) {
			Session.Navigate("/user/" + key)
		})

		form.Render("user-list", "main", users)

	}()

}

type Roles struct {
	ID   int
	Name string
}

// Edit an existing user
func userEdit(context *router.Context) {

	id, err := strconv.Atoi(context.Params["id"])
	if err != nil {
		print(err.Error())
		return
	}

	roles := []Roles{
		{1, "Admin"},
		{2, "Site Manager"},
		{3, "Technician"},
		{4, "Service Contractor"},
		{5, "Floor"},
		{6, "Public"},
	}

	go func() {
		user := shared.User{}
		rpcClient.Call("UserRPC.Get", id, &user)

		BackURL := "/users"
		form := formulate.EditForm{}
		form.New("fa-user", "User Details - "+user.Name)

		currentRole := 0
		for _, r := range roles {
			if r.Name == user.Role {
				currentRole = r.ID
				break
			}
		}
		// Layout the fields

		form.Row(2).
			Add(1, "Username", "text", "Username", `id="focusme"`).
			Add(1, "Password", "text", "Passwd", "")

		form.Row(3).
			Add(1, "Name", "text", "Name", "").
			Add(1, "Email", "text", "Email", "").
			Add(1, "Mobile", "text", "SMS", "")

		if Session.UserRole == "Admin" {
			form.Row(3).
				AddSelect(1, "Role", "Role", roles, "ID", "Name", 1, currentRole).
				AddNumber(1, "Hourly Rate", "HourlyRate", "1")
		} else {
			form.Row(3).
				AddSelect(1, "Role", "Role", roles, "ID", "Name", 1, currentRole)
		}
		form.Row(1).
			Add(1, "Sites to Access", "div", "Sites", "")

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
		})

		form.DeleteEvent(func(evt dom.Event) {
			evt.PreventDefault()
			go func() {
				data := shared.UserUpdateData{
					Channel: Session.Channel,
					User:    &user,
				}
				done := false
				rpcClient.Call("UserRPC.Delete", data, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.SaveEvent(func(evt dom.Event) {
			evt.PreventDefault()
			form.Bind(&user)

			// now convert the selected role ID into a string
			roleID, _ := strconv.Atoi(user.Role)
			for _, r := range roles {
				if r.ID == roleID {
					user.Role = r.Name
					break
				}
			}

			data := shared.UserUpdateData{
				Channel: Session.Channel,
				User:    &user,
			}
			go func() {
				done := false
				rpcClient.Call("UserRPC.Update", data, &done)
				Session.Navigate(BackURL)
			}()
		})

		form.PrintEvent(func(evt dom.Event) {
			dom.GetWindow().Print()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &user)

		// And load in a the sites array
		// Need a data set that is a slice of :
		// - Site ID
		// - Site Name
		// - Bool, whether the user has access to this site
		// On toggling a site, send an RPC to the backend to toggle the state
		userSites := []shared.UserSite{}
		req := shared.UserSiteRequest{
			Channel: Session.Channel,
			User:    &user,
		}
		rpcClient.Call("UserRPC.GetSites", req, &userSites)
		loadTemplate("user-sites-array", "[name=Sites]", userSites)

		// add a click handler for the sites array
		w := dom.GetWindow()
		doc := w.Document()

		if el := doc.QuerySelector("[name=Sites]"); el != nil {

			el.AddEventListener("click", false, func(evt dom.Event) {
				clickedOn := evt.Target()
				switch clickedOn.TagName() {
				case "INPUT":
					ie := clickedOn.(*dom.HTMLInputElement)
					key, _ := strconv.Atoi(ie.GetAttribute("key"))
					data := shared.UserSiteSetRequest{
						Channel: Session.Channel,
						UserID:  user.ID,
						SiteID:  key,
						Role:    user.Role,
						IsSet:   ie.Checked,
					}

					go func() {
						done := false
						rpcClient.Call("UserRPC.SetSite", data, &done)
					}()
				}

			})

		}

	}()

}

// Add form for a new user
func userAdd(context *router.Context) {

	roles := []Roles{
		{1, "Admin"},
		{2, "Site Manager"},
		{3, "Technician"},
		{4, "Service Contractor"},
		{5, "Floor"},
		{6, "Public"},
	}

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

		form.Row(3).
			Add(1, "Role", "select", "Role", "")
		form.SetSelectOptions("Role", roles, "ID", "Name", 1, 0)

		// Add event handlers
		form.CancelEvent(func(evt dom.Event) {
			evt.PreventDefault()
			Session.Navigate(BackURL)
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
				Session.Navigate(BackURL)
			}()
		})

		// All done, so render the form
		form.Render("edit-form", "main", &user)

	}()

}
