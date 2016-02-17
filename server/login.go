package main

import (
	"database/sql"
	"github.com/steveoc64/go-cmms/shared"
	"log"
	"strings"
	"time"
)

type LoginRPC struct{}

type loginResponse struct {
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
		res := &loginResponse{}
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
			lr.Menu = []string{}
			lr.Role = ""
			lr.Site = ""
		} else {
			log.Println("Login OK")
			lr.Result = "OK"

			//lr.Menu = []string{"RPC Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}
			lr.Menu = getMenu(res.Role)
			lr.Home = getHome(res.Role)
			lr.Role = res.Role
			if res.SiteName.Valid {
				lr.Site = res.SiteName.String
			}
			conn.Login(lc.Username, res.ID, res.Role, lr.Home)
			Connections.Show("connections after new login")
		}
	}

	log.Printf(`Login.Login ->
    » (%s,%s,%t,%d)
    « (%s,%s,%s) 
    = %s`,
		lc.Username, lc.Password, lc.RememberMe, lc.Channel,
		lr.Result, lr.Role, lr.Site,
		time.Since(start))

	return nil
}

func getMenu(role string) []string {

	switch role {
	case "Admin":
		return []string{"Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}
	case "Worker":
		return []string{}
	case "Site Manager":
		return []string{"Dashboard", "WorkOrders", "Machines", "Sites", "Users", "Reports"}
	case "Floor":
		return []string{}
	case "Service Contractor":
		return []string{"Dashboard", "WorkOrders", "Machines", "Sites", "Reports"}
	}
	return []string{}
}

func getHome(role string) string {

	switch role {
	case "Admin":
		return "Dashboard"
	case "Worker":
		return "SiteMap"
	case "Site Manager":
		return "Dashboard"
	case "Floor":
		return "SiteMap"
	case "Service Contractor":
		return "DashBoard"
	}
	return "Index"
}
