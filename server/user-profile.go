package main

import (
	"fmt"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type UserProfileRPC struct{}

///////////////////////////////////////////////////////////
// SQL
const UserProfileQuery = `select 
u.id,u.username,u.passwd,u.email,u.role,u.sms,u.name
	from users u
	where id=$1`

///////////////////////////////////////////////////////////
// Code

// Get the user profile for the given channel
func (s *UserProfileRPC) Get(channel int, prof *shared.UserProfile) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(UserProfileQuery, conn.UserID).QueryStruct(prof)

	logger(start, "UserProfile.Get",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s %s %s", prof.Email, prof.SMS, prof.Name))

	return nil
}

// Set the user profile
func (s *UserProfileRPC) Set(req *shared.UserProfileUpdate, done *bool) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	DB.Update("users").
		SetWhitelist(req, "name", "passwd", "email", "sms").
		Where("id = $1", req.ID).
		Exec()

	logger(start, "UserProfile.Set",
		fmt.Sprintf("Channel %d, User %d %s %s",
			req.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s %s %s %s", req.Email, req.SMS, req.Name, req.Passwd))

	*done = true

	return nil
}
