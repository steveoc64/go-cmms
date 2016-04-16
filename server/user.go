package main

import (
	"fmt"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type UserRPC struct{}

///////////////////////////////////////////////////////////
// SQL
const UserGetQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name
	from users u
	where id=$1`

const UserListQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name
	from users u
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
		fmt.Sprintf("%d Users", len(*profs)))

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
		fmt.Sprintf("%s %s %s", prof.Email, prof.SMS, prof.Name))

	return nil
}

// Get the user for the given id
func (u *UserRPC) Get(id int, prof *shared.User) error {
	start := time.Now()

	DB.SQL(UserGetQuery, id).QueryStruct(prof)

	logger(start, "User.Get",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("%s %s %s", prof.Email, prof.SMS, prof.Name))

	return nil
}

// Set the user profile from the popdown list at the top
func (u *UserRPC) Set(req *shared.UserUpdate, done *bool) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	DB.Update("users").
		SetWhitelist(req, "name", "passwd", "email", "sms").
		Where("id = $1", req.ID).
		Exec()

	logger(start, "User.Set",
		fmt.Sprintf("Channel %d, User %d %s %s",
			req.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s %s %s %s", req.Email, req.SMS, req.Name, req.Passwd))

	// *done = true

	return nil
}

// Full update of user record, including username
func (u *UserRPC) Update(data *shared.UserUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("users").
		SetWhitelist(data, "username", "name", "passwd", "email", "sms").
		Where("id = $1", data.User.ID).
		Exec()

	logger(start, "User.Save",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %s %s %s",
			data.User.ID, data.User.Username, data.User.Email, data.User.SMS, data.User.Name, data.User.Passwd))

	*done = true

	return nil
}

// Add a new user record
func (u *UserRPC) Insert(data *shared.UserUpdateData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.InsertInto("users").
		Whitelist("username", "name", "passwd", "email", "sms").
		Record(data.User).
		Returning("id").
		QueryScalar(id)

	logger(start, "User.Insert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %s %s %s %s",
			*id, data.User.Username, data.User.Email, data.User.SMS, data.User.Name, data.User.Passwd))

	return nil
}
