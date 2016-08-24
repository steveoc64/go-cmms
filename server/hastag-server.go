package main

import (
	"fmt"
	"time"

	// "github.com/jung-kurt/gofpdf"
	"itrak-cmms/shared"
	// "github.com/jung-kurt/gofpdf"
	// "github.com/signintech/gopdf"
)

func (t *TaskRPC) HashtagList(channel int, hashtags *[]shared.Hashtag) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(`select * from hashtag order by name`).QueryStructs(hashtags)

	logger(start, "Task.HashtagList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Hashtags", len(*hashtags)),
		channel, conn.UserID, "hashtag", 0, false)

	return nil
}

// Get a list of hashtags - longest first, to enable recursive processing of expansions
func (t *TaskRPC) HashtagListByLen(channel int, hashtags *[]shared.Hashtag) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(`select * from hashtag order by length(name) desc`).QueryStructs(hashtags)

	logger(start, "Task.HashtagListByLen",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Hashtags", len(*hashtags)),
		channel, conn.UserID, "hashtag", 0, false)

	return nil
}

func (t *TaskRPC) HashtagGet(data shared.HashtagRPCData, hashtag *shared.Hashtag) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select * from hashtag where id=$1`, data.ID).QueryStruct(hashtag)

	logger(start, "Task.HashtagGet",
		fmt.Sprintf("ID %d", data.ID),
		fmt.Sprintf("Name %s", hashtag.Name),
		data.Channel, conn.UserID, "hashtag", data.ID, false)

	return nil
}

func (t *TaskRPC) HashtagInsert(data shared.HashtagRPCData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.InsertInto("hashtag").
		Whitelist("name", "descr").
		Record(data.Hashtag).
		Returning("id").
		QueryScalar(id)

	logger(start, "Task.HashtagInsert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID %d Name %s",
			data.Hashtag.ID, data.Hashtag.Name),
		data.Channel, conn.UserID, "hashtag", *id, true)

	conn.Broadcast("hashtag", "insert", *id)
	return nil
}

func (t *TaskRPC) HashtagDelete(data shared.HashtagRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`delete from hashtag where id=$1`, data.Hashtag.ID).Exec()

	logger(start, "Task.HashtagDelete",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID %d Name %s",
			data.Hashtag.ID, data.Hashtag.Name),
		data.Channel, conn.UserID, "hashtag", data.Hashtag.ID, true)

	*done = true
	conn.Broadcast("hashtag", "delete", data.Hashtag.ID)

	return nil
}

func (t *TaskRPC) HashtagUpdate(data shared.HashtagRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`update hashtag set name=$2,descr=$3 where id=$1`,
		data.Hashtag.ID, data.Hashtag.Name, data.Hashtag.Descr).Exec()

	logger(start, "Task.HashtagUpdate",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID %d Name %s",
			data.Hashtag.ID, data.Hashtag.Name),
		data.Channel, conn.UserID, "hashtag", data.Hashtag.ID, true)

	*done = true

	conn.Broadcast("hashtag", "update", data.Hashtag.ID)

	return nil
}

// Convert 'ABCDEFG' to, for example, 'A,BCD,EFG'
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}
