package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type UserRPC struct{}

///////////////////////////////////////////////////////////
// SQL
const UserGetQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	where id=$1`

const UserListQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	order by u.username`

const TechniciansListQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	left join user_site x on x.user_id=u.id and x.site_id=$1
	where u.role in ('Worker','Technician') and x.site_id=$1
	order by u.username`

const ManagersListQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	left join user_site x on x.user_id=u.id and x.site_id=$1
	where u.role='Site Manager' and x.site_id=$1
	order by u.username`

const ManagersAllQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	where u.role='Site Manager'
	order by u.username`

const TechniciansAllQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	where u.role='Technician'
	order by u.username`

const AdminsListQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name,u.hourly_rate,u.use_mobile
	from users u
	where u.role='Admin'
	order by u.username`

///////////////////////////////////////////////////////////
// Code

// Get all users
func (u *UserRPC) List(channel int, profs *[]shared.User) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(UserListQuery, conn.UserID).QueryStructs(profs)

	logger(start, "User.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Users", len(*profs)),
		channel, conn.UserID, "users", 0, false)

	return nil
}

// Get the user for the given channel
func (u *UserRPC) Me(channel int, prof *shared.User) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(UserGetQuery, conn.UserID).QueryStruct(prof)

	logger(start, "User.Me",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s %s %s", prof.Email, prof.SMS, prof.Name),
		channel, conn.UserID, "users", 0, false)

	return nil
}

// Get the user for the given id
func (u *UserRPC) Get(data shared.UserRPCData, prof *shared.User) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(UserGetQuery, data.ID).QueryStruct(prof)

	logger(start, "User.Get",
		fmt.Sprintf("%d", data.ID),
		fmt.Sprintf("%s %s %s", prof.Email, prof.SMS, prof.Name),
		data.Channel, conn.UserID, "users", data.ID, false)

	return nil
}

// Set the user profile from the popdown list at the top
func (u *UserRPC) Set(req shared.UserUpdate, done *bool) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	DB.Update("users").
		SetWhitelist(req, "name", "passwd", "email", "sms").
		Where("id = $1", req.ID).
		Exec()

	logger(start, "User.Set",
		fmt.Sprintf("Channel %d, User %d %s %s",
			req.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s %s %s %s", req.Email, req.SMS, req.Name, req.Passwd),
		req.Channel, conn.UserID, "users", req.ID, true)

	// *done = true

	return nil
}

// Full update of user record, including username
func (u *UserRPC) Update(data shared.UserRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("users").
		SetWhitelist(data.User, "username", "name", "passwd", "email", "sms",
			"role", "hourly_rate", "use_mobile").
		Where("id = $1", data.User.ID).
		Exec()

	logger(start, "User.Update",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Role %s Uname %s Eml %s SMS %s Name %s PW %s",
			data.User.ID, data.User.Role, data.User.Username, data.User.Email,
			data.User.SMS, data.User.Name, data.User.Passwd),
		data.Channel, conn.UserID, "users", data.User.ID, true)

	*done = true

	return nil
}

// Add a new user record
func (u *UserRPC) Insert(data shared.UserRPCData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.InsertInto("users").
		Whitelist("username", "name", "passwd", "email", "sms", "hourly_rate").
		Record(data.User).
		Returning("id").
		QueryScalar(id)

	logger(start, "User.Insert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %s %s %s %s",
			*id, data.User.Username, data.User.Email, data.User.SMS, data.User.Name, data.User.Passwd),
		data.Channel, conn.UserID, "users", *id, true)

	return nil
}

// Delete a user
func (u *UserRPC) Delete(data shared.UserRPCData, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	*ok = false
	id := data.User.ID
	DB.DeleteFrom("users").
		Where("id=$1", id).
		Exec()

	logger(start, "User.Delete",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %s %s %s %s",
			id, data.User.Username, data.User.Email, data.User.SMS, data.User.Name, data.User.Passwd),
		data.Channel, conn.UserID, "users", id, true)

	return nil
}

// Get an array of Sites for this user
func (u *UserRPC) GetSites(data shared.UserSiteRequest, userSites *[]shared.UserSite) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select 
		s.id as site_id,s.name as site_name,count(u.*)
		from site s
		left join user_site u
			on u.site_id=s.id
			and u.user_id=$1
		group by s.id
		order by s.name`, data.User.ID).QueryStructs(userSites)

	logger(start, "User.GetSites",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("User %d - %d Sites",
			data.User.ID, len(*userSites)),
		data.Channel, conn.UserID, "users", 0, false)

	return nil
}

// Get an array of Users for this site
func (u *UserRPC) GetSiteUsers(data shared.UserSiteRequest, siteUsers *[]shared.SiteUser) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select 
		u.id as user_id,u.username as username,count(s.*)
		from users u
		left join user_site s
			on s.user_id=u.id
			and s.site_id=$1
		group by u.id
		order by u.username`, data.ID).QueryStructs(siteUsers)

	logger(start, "User.GetSiteUsers",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Site %d - %d Users",
			data.ID, len(*siteUsers)),
		data.Channel, conn.UserID, "users", 0, false)

	return nil
}

// Set the user site relationship
func (u *UserRPC) SetSite(data shared.UserSiteSetRequest, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// delete any existing relationship
	DB.DeleteFrom("user_site").
		Where("user_id=$1 and site_id=$2", data.UserID, data.SiteID).
		Exec()

	if data.IsSet {
		// if the role is undefined, then read it from the user
		if data.Role == "" {
			DB.SQL(`select role from users where id=$1`, data.UserID).QueryScalar(&data.Role)
			log.Println("fetched user role", data.Role)
		}

		DB.SQL(`insert into 
			user_site (user_id,site_id,role)
			values    ($1, $2, $3)`, data.UserID, data.SiteID, data.Role).
			Exec()
	}

	logger(start, "User.SetSite",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("User %d Site %d Role %s %t",
			data.UserID, data.SiteID, data.Role, data.IsSet),
		data.Channel, conn.UserID, "user_site", data.UserID, true)

	conn.Broadcast("usersites", "usersite", data.UserID)
	*done = true
	return nil
}

// Get a list of technicians by Site
func (u *UserRPC) GetTechnicians(data shared.SiteRPCData, users *[]shared.User) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// log.Println(TechniciansListQuery)
	site_id := data.ID
	if site_id == 0 {
		DB.SQL(TechniciansAllQuery).QueryStructs(users)

	} else {
		DB.SQL(TechniciansListQuery, site_id).QueryStructs(users)

	}

	logger(start, "User.GetTechnicians",
		fmt.Sprintf("Site %d", site_id),
		fmt.Sprintf("%d Techs", len(*users)),
		data.Channel, conn.UserID, "users", 0, false)

	return nil
}

// Get a list of technicians by Site
func (u *UserRPC) GetManagers(data shared.SiteRPCData, users *[]shared.User) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// log.Println(ManagersListQuery)
	site_id := data.ID
	if site_id == 0 {
		DB.SQL(ManagersAllQuery).QueryStructs(users)
	} else {
		DB.SQL(ManagersListQuery, site_id).QueryStructs(users)
	}

	// log.Println(AdminsListQuery)
	DB.SQL(AdminsListQuery, site_id).QueryStructs(users)

	logger(start, "User.GetManagers",
		fmt.Sprintf("Site %d", site_id),
		fmt.Sprintf("%d Mananges", len(*users)),
		data.Channel, conn.UserID, "users", 0, false)

	return nil
}
