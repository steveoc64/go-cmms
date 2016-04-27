package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type EventRPC struct{}

// Get all the tasks for the given machine
func (t *EventRPC) Raise(issue shared.RaiseIssue, id *int) error {
	start := time.Now()

	conn := Connections.Get(issue.Channel)
	*id = 0

	ToolName := issue.NonTool
	if issue.IsTool {
		ToolName = issue.Component.Name
	}

	// Create 1 event record - which includes details of both tool and machine
	evt := &shared.Event{
		SiteID:    issue.Machine.SiteID,
		Type:      "Alert",
		MachineID: issue.Machine.ID,
		ToolID:    issue.CompID,
		ToolType:  ToolName,
		CreatedBy: conn.UserID,
		Notes:     issue.Descr,
		Priority:  1,
	}

	DB.InsertInto("event").
		Whitelist("site_id", "type", "machine_id", "tool_id", "tool_type", "created_by", "notes", "priority").
		Record(evt).
		Returning("id").
		QueryScalar(id)

	DB.SQL(`update machine 
			set alert_at=localtimestamp, status=$2 
			where id=$1`,
		issue.Machine.ID,
		`Needs Attention`).
		Exec()

	// if its a tool, then update the tool record, otherwise update the non-tool field on the machine record
	if evt.ToolID == 0 {
		// is a non-tool.
		fieldName := ""
		switch evt.ToolType {
		case "Electrical":
			fieldName = "electrical"
		case "Hydraulic":
			fieldName = "hydraulic"
		case "Lube":
			fieldName = "lube"
		case "Printer":
			fieldName = "printer"
		case "Console":
			fieldName = "console"
		case "Uncoiler":
			fieldName = "uncoiler"
		case "Rollbed":
			fieldName = "rollbed"
		}
		if fieldName != "" {
			DB.SQL(fmt.Sprintf("update machine set %s='Needs Attention' where id=$1", fieldName), evt.MachineID).Exec()
		}
	} else {
		// is a tool
		DB.SQL(`update component
			set status='Needs Attention'
			where id=$1`, evt.ToolID).
			Exec()
	}

	// Generate an SMS
	// err = SendSMS("0417824950",
	// 	fmt.Sprintf("%s on Machine %s %s: %s", "Alert", issue.Machine.Name, ToolName, issue.Descr),
	// 	fmt.Sprintf("%d", id))

	// Generate an Email with same details as SMS

	// Patch in any attached documents
	// _, err = DB.SQL(`update doc
	// 	set ref_id=$1, name=$3, type='toolevent'
	// 	where type='temptoolevent' and ref_id=$2
	// 	`, evt.ID, evt.ToolId, evt.Notes).Exec()

	logger(start, "Event.Raise",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			issue.Channel, issue.Machine.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Event %d Tool %d:%s Desc %s", *id, evt.ToolID, ToolName, evt.Notes))

	// send SMS to Peter only at this stage

	siteName := ""
	DB.SQL(`select name from site where id=$1`, issue.Machine.SiteID).QueryScalar(&siteName)

	if Config.SMSOn {
		SendSMS("0439086713",
			fmt.Sprintf("Alert at Site %s on Machine %s on %s: %s",
				siteName,
				issue.Machine.Name,
				ToolName,
				issue.Descr),
			fmt.Sprintf("%d", evt.ID))
	} else {
		willSend := fmt.Sprintf("Alert at Site %s on Machine %s on %s: %s",
			siteName,
			issue.Machine.Name,
			ToolName,
			issue.Descr)

		log.Println("will send SMS", willSend)

	}

	return nil
}

func (e *EventRPC) List(channel int, events *[]shared.Event) error {
	start := time.Now()

	conn := Connections.Get(channel)

	switch conn.UserRole {
	case "Site Manager":
		// Limit the tasks to just the sites that we are in control of
		sites := []int{}

		DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&sites)

		err := DB.SQL(`select 
		e.*,m.name as machine_name,s.name as site_name,u.username as username
		from event e
			left join machine m on m.id=e.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=e.created_by
		where m.site_id in $1 and completed is null
		order by e.startdate`, sites).
			QueryStructs(events)

		if err != nil {
			log.Println(err.Error())
		}
	case "Admin":
		err := DB.SQL(`select 
		e.*,m.name as machine_name,s.name as site_name,u.username as username
		from event e
			left join machine m on m.id=e.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=e.created_by		
		order by e.completed desc,e.startdate`).
			QueryStructs(events)

		if err != nil {
			log.Println(err.Error())
		}
	}

	logger(start, "Event.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Events", len(*events)))

	return nil
}

func (e *EventRPC) Get(id int, event *shared.Event) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(`select
		e.*,m.name as machine_name,s.name as site_name,u.username as username
		from event e
			left join machine m on m.id=e.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=e.created_by
		where e.id=$1`, id).QueryStruct(event)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Event.Get",
		fmt.Sprintf("ID %d", id),
		event.Notes)

	return nil
}

func (e *EventRPC) Update(data shared.EventUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("event").
		SetWhitelist(data.Event, "notes").
		Where("id = $1", data.Event.ID).
		Exec()

	logger(start, "Event.Update",
		fmt.Sprintf("Channel %d, Event %d User %d %s %s",
			data.Channel, data.Event.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s", data.Event.Notes))

	return nil
}

func (e *EventRPC) Complete(data shared.EventUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	if conn.UserRole != "Admin" {
		return nil
	}

	// Read the sites that this user has access to
	event := shared.Event{}
	DB.SQL(`select
		e.*,m.name as machine_name,s.name as site_name,u.username as username
		from event e
			left join machine m on m.id=e.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=e.created_by
		where e.id=$1`, data.Event.ID).QueryStruct(&event)

	// Mark the event as complete
	DB.SQL(`update event 
		set completed=now(), status='Complete'
		where id=$1`, data.Event.ID).Exec()

	// Reset the affected component - this code is the reverse of
	// the code in the RaiseEvent function above
	if event.ToolID == 0 {
		// Reset the status of the basic component on this machine
		fieldName := ""
		switch event.ToolType {
		case "Electrical":
			fieldName = "electrical"
		case "Hydraulic":
			fieldName = "hydraulic"
		case "Lube":
			fieldName = "lube"
		case "Printer":
			fieldName = "printer"
		case "Console":
			fieldName = "console"
		case "Uncoiler":
			fieldName = "uncoiler"
		case "Rollbed":
			fieldName = "rollbed"
		}
		if fieldName != "" {
			DB.SQL(fmt.Sprintf("update machine set %s='Running' where id=$1", fieldName), event.MachineID).Exec()
		}
	} else {
		// is a tool
		DB.SQL(`update component
			set status='Running'
			where id=$1`, event.ToolID).
			Exec()
	}

	// Reset the whole machine if clear
	machineIsClear := true
	machine := shared.Machine{}
	DB.SQL(`select * from machine where id=$1`, event.MachineID).QueryStruct(&machine)

	if machine.Electrical != "Running" ||
		machine.Hydraulic != "Running" ||
		machine.Printer != "Running" ||
		machine.Console != "Running" ||
		machine.Rollbed != "Running" ||
		machine.Uncoiler != "Running" ||
		machine.Lube != "Running" {
		machineIsClear = false
	}

	if machineIsClear {
		badComps := 0
		DB.SQL(`select count(*) 
			from component 
			where status != 'Running' and machine_id=$1`, machine.ID).
			QueryScalar(&badComps)

		if badComps == 0 {
			DB.SQL("update machine set status='Running' where id=$1", event.MachineID).Exec()
		}
	}

	logger(start, "Event.Complete",
		fmt.Sprintf("Channel %d, Event %d User %d %s %s",
			data.Channel, data.Event.ID, conn.UserID, conn.Username, conn.UserRole),
		"Manually Completed by Admin")

	return nil
}
