package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type LoginRPC struct{}

type dbLoginResponse struct {
	ID       int            `db:"id"`
	Username string         `db:"username"`
	Name     string         `db:"name"`
	Role     string         `db:"role"`
	Site_ID  int            `db:"site_id"`
	SiteName sql.NullString `db:"sitename"`
}

func (l *LoginRPC) Login(lc *shared.LoginCredentials, lr *shared.LoginReply) error {
	start := time.Now()

	// do some authentication here

	// send a login reply

	// Get the connection we are on
	// log.Println("channel is", lc.Channel)
	conn := Connections.Get(lc.Channel)
	// log.Println("got conn", conn)
	if conn != nil {
		// validate that username and passwd is correct
		res := &dbLoginResponse{}
		usename := strings.ToLower(lc.Username)
		err := DB.
			Select("u.id,u.username,u.name,u.role,u.site_id,s.name as sitename").
			From(`users u
			left join site s on (s.id = u.site_id)`).
			Where("u.username = $1 and passwd = $2", usename, lc.Password).
			QueryStruct(res)

		if err != nil {
			log.Println("Login Failed:", err.Error())
			lr.Result = "Failed"
			lr.Token = ""
			lr.Menu = []shared.UserMenu{}
			lr.Routes = []shared.UserRoute{}
			lr.Role = ""
			lr.Site = ""
		} else {
			// log.Println("Login OK")
			lr.Result = "OK"
			lr.Token = fmt.Sprintf("%d", lc.Channel)

			//lr.Menu = []string{"RPC Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}
			lr.Menu = getMenu(res.Role)
			lr.Routes = getRoutes(res.ID, res.Role)
			lr.Role = res.Role
			if res.SiteName.Valid {
				lr.Site = res.SiteName.String
			}
			conn.Login(lc.Username, res.ID, res.Role)
			Connections.Show("connections after new login")
		}
	}

	logger(start, "Login.Login",
		fmt.Sprintf("%s,%s,%t,%d", lc.Username, lc.Password, lc.RememberMe, lc.Channel),
		fmt.Sprintf("%s,%s,%s", lr.Result, lr.Role, lr.Site))

	return nil
}

func getMenu(role string) []shared.UserMenu {

	switch role {
	case "Admin":
		return []shared.UserMenu{
			{Title: "Dashboard", Icon: "dashboard", URL: "/"},
			{Title: "Events", Icon: "report_problem", URL: "/events"},
			{Title: "Sites", Icon: "store", URL: "/sites"},
			{Title: "Machines", Icon: "view_array", URL: "/machines"},
			{Title: "Tools", Icon: "view_column", URL: "/tools"},
			{Title: "Parts", Icon: "view_module", URL: "/parts"},
			{Title: "Vendors", Icon: "contact_phone", URL: "/vendors"},
			{Title: "Users", Icon: "supervisor_account", URL: "/users"},
			{Title: "Skills", Icon: "verified_user", URL: "/skills"},
			{Title: "Reports", Icon: "comment", URL: "/reports"},
		}
	case "Worker":
		return []shared.UserMenu{}
	case "Site Manager":
		return []shared.UserMenu{
			{Title: "Dashboard", Icon: "", URL: "/"},
			{Title: "WorkOrders", Icon: "", URL: "/workorders"},
			{Title: "Sites", Icon: "", URL: "/sites"},
			{Title: "Machines", Icon: "", URL: "/machines"},
			{Title: "Users", Icon: "", URL: "/users"},
			{Title: "Reports", Icon: "", URL: "/reports"},
		}
	case "Floor":
		return []shared.UserMenu{}
	case "Service Contractor":
		return []shared.UserMenu{
			{Title: "Dashboard", Icon: "", URL: "/"},
			{Title: "WorkOrders", Icon: "", URL: "/workorders"},
			{Title: "Machines", Icon: "", URL: "/machines"},
			{Title: "Sites", Icon: "", URL: "/sites"},
			{Title: "Reports", Icon: "", URL: "/reports"},
		}
	}
	return []shared.UserMenu{}
}

func getRoutes(uid int, role string) []shared.UserRoute {

	switch role {
	case "Admin":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/machines", Func: "machines"},
			{Route: "/sites", Func: "sites"},
			{Route: "/events", Func: "events"},
			{Route: "/workorders", Func: "workorders"},
			{Route: "/tools", Func: "tools"},
			{Route: "/parts", Func: "parts"},
			{Route: "/vendors", Func: "vendors"},
			{Route: "/users", Func: "users"},
			{Route: "/skills", Func: "skills"},
			{Route: "/reports", Func: "reports"},
		}
	case "Worker":
		// If this user has 1 site, they ony get 1 route
		// Otherwise they get a map and 1 route to show the machines at each site
		numSites := 1
		DB.SQL(`select count(*) from user_site where user_id=$1`, uid).QueryScalar(&numSites)
		if numSites == 1 {
			return []shared.UserRoute{
				{Route: "/", Func: "homesite"},
			}
		}
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
		}
	case "Site Manager":
		return []shared.UserRoute{
			{Route: "/", Func: "sitemap"},
			{Route: "/sitemachines/{site}", Func: "sitemachines"},
			{Route: "/tasks", Func: "tasks"},
			{Route: "/users", Func: "users"},
			{Route: "/events", Func: "events"},
			{Route: "/workorders", Func: "workorders"},
		}
	case "Floor":
		return []shared.UserRoute{
			{Route: "/", Func: "machines"},
		}
	case "Service Contractor":
		return []shared.UserRoute{
			{Route: "/", Func: "dashboard"},
			{Route: "/machines", Func: "machines"},
		}
	}
	return []shared.UserRoute{}

}
