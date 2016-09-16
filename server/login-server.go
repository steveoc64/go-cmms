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

func (l *LoginRPC) Nav(data shared.Nav, r *string) error {
	conn := Connections.Get(data.Channel)
	conn.Routes = append(conn.Routes, data.Route)
	conn.Route = data.Route
	*r = conn.Route
	println("\n----------------------------------")
	log.Printf("%s:%s -> %s\n", conn.Username, conn.UserRole, conn.Route)
	conn.BroadcastAdmin("nav", data.Route, data.Channel)
	return nil
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
			conn.Broadcast("login", "insert", lr.ID)
		}
	}

	logger(start, "Login.Login",
		fmt.Sprintf("%s,%s,%t,%d", lc.Username, lc.Password, lc.RememberMe, lc.Channel),
		fmt.Sprintf("%s,%s,%s", lr.Result, lr.Role, lr.Site),
		lc.Channel, lr.ID, "users", lr.ID, false)

	return nil
}

func (l *LoginRPC) UsersOnline(channel int, u *[]shared.UserOnline) error {
	start := time.Now()
	conn := Connections.Get(channel)
	if conn.UserRole == "Admin" {

		for _, k := range Connections.Keys() {
			v := Connections.Get(k)
			println("k,v = ", k, v)
			req := v.Socket.Request()
			theIP := ""
			if theIP = req.Header.Get("X-Real-Ip"); theIP == "" {
				theIP = req.RemoteAddr
			}
			user := shared.UserOnline{
				ID:          v.UserID,
				Username:    v.Username,
				Browser:     fmt.Sprintf("%s", req.Header["User-Agent"]),
				IP:          theIP,
				Name:        "lookup",
				Email:       "lookup",
				Role:        v.UserRole,
				SMS:         "lookup",
				IsTech:      false,
				CanAllocate: false,
				Route:       v.Route,
				Routes:      v.Routes,
				Duration:    time.Since(v.Time).String(),
				Channel:     v.ID,
			}

			if user.ID == 0 {
				user.Duration = ""
			}
			*u = append(*u, user)

		}

	}

	logger(start, "Login.UsersOnline",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Users Online", len(*u)),
		channel, conn.UserID, "users", 0, false)

	return nil
}
