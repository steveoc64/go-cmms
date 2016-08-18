package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"itrak-cmms/shared"
)

type LoginRPC struct{}

type dbLoginResponse struct {
	ID          int            `db:"id"`
	Username    string         `db:"username"`
	Name        string         `db:"name"`
	Role        string         `db:"role"`
	SiteID      int            `db:"site_id"`
	SiteName    sql.NullString `db:"sitename"`
	CanAllocate bool           `db:"can_allocate"`
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

		// log.Println(`select u.id,u.username,u.name,u.role,u.site_id,s.name as sitename
		// 	from users u
		// 	left join site s on (s.id = u.site_id)
		// 	where lower(u.username) = lower('`, lc.Username, `') and lower(passwd) = lower('`, lc.Password, `')`)

		err := DB.
			Select("u.id,u.username,u.name,u.role,u.site_id,s.name as sitename,u.can_allocate as can_allocate").
			From(`users u
			left join site s on (s.id = u.site_id)`).
			Where("lower(u.username) = lower($1) and lower(passwd) = lower($2)",
				lc.Username, lc.Password).
			QueryStruct(res)
		// log.Println(res)

		if err != nil {
			log.Println("Login Failed:", err.Error())
			lr.Result = "Failed"
			lr.Token = ""
			// lr.Menu = []shared.UserMenu{}
			lr.Routes = []shared.UserRoute{}
			lr.Role = ""
			lr.Site = ""
		} else {
			// log.Println("Login OK")
			lr.Result = "OK"
			lr.Token = fmt.Sprintf("%d", lc.Channel)

			//lr.Menu = []string{"RPC Dashboard", "Events", "Sites", "Machines", "Tools", "Parts", "Vendors", "Users", "Skills", "Reports"}
			// lr.Menu = getMenu(res.Role)
			lr.Routes = getRoutes(res.ID, res.Role)
			lr.Role = res.Role
			lr.ID = res.ID
			lr.CanAllocate = res.CanAllocate
			if res.SiteName.Valid {
				lr.Site = res.SiteName.String
			}
			conn.Login(lc.Username, res.ID, res.Role)
			Connections.Show("connections after new login")
		}
	}

	logger(start, "Login.Login",
		fmt.Sprintf("%s,%s,%t,%d", lc.Username, lc.Password, lc.RememberMe, lc.Channel),
		fmt.Sprintf("%s,%s,%s", lr.Result, lr.Role, lr.Site),
		lc.Channel, lr.ID, "users", lr.ID, false)

	return nil
}
