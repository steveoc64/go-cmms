package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	// "github.com/jung-kurt/gofpdf"
	"itrak-cmms/shared"
	// "github.com/jung-kurt/gofpdf"
	// "github.com/signintech/gopdf"
)

const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

type TaskRPC struct{}

// Update a Task
func (t *TaskRPC) Update(data shared.TaskRPCData, updatedTask *shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	canAllocate := false
	DB.SQL(`select can_allocate from users where id=$1`, conn.UserID).QueryScalar(&canAllocate)
	useRole := conn.UserRole
	if canAllocate {
		useRole = "Admin"
	}

	oldTask := shared.Task{}
	DB.SQL(`select * from task where id=$1`, data.Task.ID).QueryStruct(&oldTask)

	if useRole == "Admin" {

		// Admin can re-assign the task to another user
		DB.Update("task").
			SetWhitelist(data.Task,
				"log", "assigned_to",
				"labour_cost", "material_cost", "labour_hrs").
			Where("id = $1", data.Task.ID).
			Exec()

	} else {
		DB.Update("task").
			SetWhitelist(data.Task,
				"log",
				"labour_hrs").
			Where("id = $1", data.Task.ID).
			Exec()
	}

	// If there is a new photo to be added to the task, then add it
	if data.Task.NewPhoto.Data != "" {
		println("Adding new photo", data.Task.NewPhoto.Data[:22])
		photo := shared.Photo{
			Data:     data.Task.NewPhoto.Data,
			Filename: data.Task.NewPhoto.Filename,
			Entity:   "task",
			EntityID: data.Task.ID,
		}

		// decodePhoto(photo.Data, &photo.Preview, &photo.Thumb)
		decodePhoto(&photo)
		DB.InsertInto("photo").
			Columns("entity", "entity_id", "photo", "thumb", "preview", "type", "datatype", "filename").
			Record(photo).
			Exec()
	}

	// If assigned to, then re-calc the labour cost if the hours have changed
	if oldTask.LabourHrs != data.Task.LabourHrs {
		if data.Task.AssignedTo != nil {
			user := shared.User{}
			DB.SQL(`select hourly_rate
				from users
				where id=$1`, *data.Task.AssignedTo).
				QueryStruct(&user)

			data.Task.LabourCost = user.HourlyRate * data.Task.LabourHrs
			DB.SQL(`update task 
				set labour_cost=$2 
				where id=$1`, data.Task.ID, data.Task.LabourCost).
				Exec()
		}
	}

	if false {
		// dont want to do this anymore, and the parts are not so tightly coupled to the task now

		for _, v := range data.Task.Parts {
			// log.Println("part = ", v)

			DB.Update("task_part").
				SetWhitelist(v, "notes", "qty_used").
				Where("task_id=$1 and part_id=$2", data.Task.ID, v.PartID).
				Exec()
		}
	}

	logger(start, "Task.Update",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%f %f %f %s",
			data.Task.LabourCost, data.Task.MaterialCost, data.Task.LabourHrs, data.Task.Log),
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	DB.SQL(`select * from task where id=$1`, data.Task.ID).QueryStruct(updatedTask)

	conn.Broadcast("task", "update", data.Task.ID)
	return nil
}

// Add an attachment to a Task
func (t *TaskRPC) AddAttach(data shared.TaskRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// If there is a new photo to be added to the task, then add it
	if data.Task.NewPhoto.Data != "" {
		println("Adding new photo", data.Task.NewPhoto.Data[:22])
		photo := shared.Photo{
			Data:     data.Task.NewPhoto.Data,
			Filename: data.Task.NewPhoto.Filename,
			Entity:   "task",
			EntityID: data.Task.ID,
		}

		// decodePhoto(photo.Data, &photo.Preview, &photo.Thumb)
		decodePhoto(&photo)
		DB.InsertInto("photo").
			Columns("entity", "entity_id", "photo", "thumb", "preview", "type", "datatype", "filename").
			Record(photo).
			Exec()
	} else {
		println("attachment is empty !!!")
	}

	logger(start, "Task.AddAttach",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Task.NewPhoto.Filename,
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	conn.Broadcast("task", "update", data.Task.ID)
	*done = true
	return nil
}

// Update a Task - just the hours
func (t *TaskRPC) UpdateHours(data shared.TaskRPCData, updatedTask *shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	DB.Update("task").
		SetWhitelist(data.Task, "assigned_to", "labour_hrs").
		Where("id = $1", data.Task.ID).
		Exec()

	if data.Task.AssignedTo != nil {
		user := shared.User{}
		DB.SQL(`select hourly_rate
				from users
				where id=$1`, *data.Task.AssignedTo).
			QueryStruct(&user)

		data.Task.LabourCost = user.HourlyRate * data.Task.LabourHrs
		DB.SQL(`update task 
				set labour_cost=$2 
				where id=$1`, data.Task.ID, data.Task.LabourCost).
			Exec()
	}

	logger(start, "Task.UpdateHours",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%f %f %f %s",
			data.Task.LabourCost, data.Task.LabourHrs),
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	DB.SQL(`select * from task where id=$1`, data.Task.ID).QueryStruct(updatedTask)

	conn.Broadcast("task", "update", data.Task.ID)
	return nil
}

// Update a Task - just the notes
func (t *TaskRPC) UpdateNotes(data shared.TaskRPCData, updatedTask *shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	DB.SQL(`update task 
				set log=$2 
				where id=$1`, data.Task.ID, data.Task.Log).
		Exec()

	logger(start, "Task.UpdateNotes",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Task.Log,
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	DB.SQL(`select * from task where id=$1`, data.Task.ID).QueryStruct(updatedTask)

	conn.Broadcast("task", "update", data.Task.ID)
	return nil
}

func (t *TaskRPC) Delete(data shared.TaskRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("task").
		Where("id=$1", data.Task.ID).
		Exec()

	logger(start, "Task.Delete",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		"",
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	*done = true
	return nil
}

func (t *TaskRPC) List(channel int, tasks *[]shared.Task) error {
	start := time.Now()

	conn := Connections.Get(channel)

	switch conn.UserRole {
	case "Technician":
		// Limit the tasks to only our own tasks
		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.assigned_to=$1 and completed_date is null
		order by t.startdate desc, id desc`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Site Manager":
		// Limit the tasks to just the sites that we are in control of
		sites := []int{}

		DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&sites)

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username,x.highlight as site_highlight
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
			left join user_site x on x.user_id=$2 and x.site_id=m.site_id
		where m.site_id in $1 and completed_date is null
		order by t.startdate desc, id desc`, sites, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Admin":

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username,x.highlight as site_highlight
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
			left join user_site x on x.user_id=$1 and x.site_id=m.site_id
		where completed_date is null
		order by t.startdate desc, id desc`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// get a set of hashtags for easy processing
	hashes := []shared.Hashtag{}
	DB.SQL(`select * from hashtag order by length(name) desc`).QueryStructs(&hashes)

	for i, v := range *tasks {

		// convert the description into a listable description

		// First, expand out the hashtags
		l := v.Descr
		if false && strings.Contains(l, "#") {
			stillLooking := true
			for stillLooking {
				stillLooking = false
				for _, h := range hashes {
					theHash := "#" + h.Name
					if strings.Contains(l, theHash) {
						l = strings.Replace(l, theHash, h.Descr, -1)
						stillLooking = true
					}
				}
			}
		}

		// trim the description
		if len(l) > 80 {
			(*tasks)[i].Descr = fmt.Sprintf("%s ...", l[:80])
		}

		// Get the latest thumbnails for this task, if present
		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='task' and entity_id=$1) 
			or (entity='event' and entity_id=$2) 
			or (entity='sched' and entity_id=$3) 
			order by type,id desc`, v.ID, v.EventID, v.SchedID).
			QueryStructs(&photos)
		(*tasks)[i].Photos = photos
	}

	// // Read the sites that this user has access to
	// err := DB.SQL(`select
	// 	t.*,m.name as machine_name,s.name as site_name,u.username as username
	// 	from task t
	// 		left join machine m on m.id=t.machine_id
	// 		left join site s on s.id=m.site_id
	// 		left join users u on u.id=t.assigned_to
	// 	order by t.startdate`).
	// 	QueryStructs(tasks)

	logger(start, "Task.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Tasks", len(*tasks)),
		channel, conn.UserID, "task", 0, false)

	return nil
}

func (t *TaskRPC) ListCompleted(channel int, tasks *[]shared.Task) error {
	start := time.Now()

	conn := Connections.Get(channel)

	switch conn.UserRole {
	case "Technician":
		// Limit the tasks to only our own tasks
		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.assigned_to=$1 and t.completed_date is not null
			and t.completed_date is not null
			and t.startdate > NOW() - INTERVAL '1 month'
		order by t.startdate desc, id desc`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Site Manager":
		// Limit the tasks to just the sites that we are in control of
		sites := []int{}

		DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&sites)

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username,x.highlight as site_highlight
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
			left join user_site x on x.user_id=$2 and x.site_id=m.site_id
		where m.site_id in $1 and t.completed_date is not null
			and t.completed_date is not null
			and t.startdate > NOW() - INTERVAL '1 month'
		order by t.startdate desc, id desc`, sites, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Admin":

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username,x.highlight as site_highlight 
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
			left join user_site x on x.user_id=$1 and x.site_id=m.site_id
		where t.completed_date is not null
		  and t.startdate > NOW() - INTERVAL '1 month'
		order by t.startdate desc, id desc`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// trim the descr fields
	for k, v := range *tasks {
		if len(v.Descr) > 80 {
			(*tasks)[k].Descr = fmt.Sprintf("%s ...", v.Descr[:80])
		}

		// Get the latest thumbnails for this task, if present
		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='task' and entity_id=$1) 
			or (entity='event' and entity_id=$2) 
			or (entity='sched' and entity_id=$3) 
			order by type,id desc`, v.ID, v.EventID, v.SchedID).
			QueryStructs(&photos)
		(*tasks)[k].Photos = photos
	}

	// // Read the sites that this user has access to
	// err := DB.SQL(`select
	// 	t.*,m.name as machine_name,s.name as site_name,u.username as username
	// 	from task t
	// 		left join machine m on m.id=t.machine_id
	// 		left join site s on s.id=m.site_id
	// 		left join users u on u.id=t.assigned_to
	// 	order by t.startdate`).
	// 	QueryStructs(tasks)

	logger(start, "Task.ListCompleted",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Tasks", len(*tasks)),
		channel, conn.UserID, "task", 0, false)

	return nil
}

func (t *TaskRPC) Get(data shared.TaskRPCData, task *shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.id=$1`, data.ID).
		QueryStruct(task)

	if err != nil {
		log.Println(err.Error())
	}

	// Now get all the parts for this task
	DB.SQL(`select 
		t.*,p.name as part_name,p.stock_code as stock_code,p.qty_type as qty_type
		from task_part t
		left join part p on p.id=t.part_id
		where t.task_id=$1`, data.ID).QueryStructs(&task.Parts)

	// Now get all the checks for this task
	DB.SQL(`select * from task_check where task_id=$1 order by task_id,seq`, data.ID).QueryStructs(&task.Checks)

	// Get the photo previews for this task
	photos := []shared.Photo{}

	DB.SQL(`select id,preview,type,datatype,filename,entity,entity_id,notes
	 from photo
	 where (entity='task' and entity_id=$1) 
	 or (entity='event' and entity_id=$2) 
	 or (entity='sched' and entity_id=$3) 
	 order by type, id desc`, data.ID, task.EventID, task.SchedID).
		QueryStructs(&photos)
	task.Photos = photos

	// Now, if the user requesting this read is the person assigned to, then
	// stamp the task as having been read
	if !task.IsRead && task.AssignedTo != nil && conn.UserID == *task.AssignedTo {
		println("Marking task as read")
		DB.SQL(`update task set is_read=true, read_date=now() where id=$1`, data.ID).Exec()
		conn.Broadcast("task", "update", data.ID)
	}

	logger(start, "Task.Get",
		fmt.Sprintf("ID %d", data.ID),
		task.Descr,
		data.Channel, conn.UserID, "task", data.ID, false)

	// log.Printf("task %v\n", *task)

	return nil
}

func (t *TaskRPC) Check(data shared.TaskCheckUpdate, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`update task_check 
		set done=true,done_date=now()
		where task_id=$1 and seq=$2`,
		data.TaskCheck.TaskID,
		data.TaskCheck.Seq).Exec()

	logger(start, "Task.Check",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Task %d Seq %d Checked",
			data.TaskCheck.TaskID, data.TaskCheck.Seq),
		data.Channel, conn.UserID, "task_check", data.TaskCheck.TaskID, true)

	*done = true
	return nil
}

// Get a list of tasks at a given Site
func (t *TaskRPC) SiteList(data shared.TaskRPCData, tasks *[]shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select 
		t.*,
		m.name as machine_name,
		s.name as site_name,s.id as site_id,
		u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where s.id=$1
		order by t.startdate`, data.ID).
		QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	// trim the descr fields
	for k, v := range *tasks {
		if len(v.Descr) > 80 {
			(*tasks)[k].Descr = fmt.Sprintf("%s ...", v.Descr[:80])
		}

		// Get the latest thumbnails for this task, if present
		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='task' and entity_id=$1) 
			or (entity='event' and entity_id=$2) 
			or (entity='sched' and entity_id=$3)
			order by type,id desc`, v.ID, v.EventID, v.SchedID).
			QueryStructs(&photos)
		(*tasks)[k].Photos = photos
	}

	// logger(start, "Task.SiteList",
	// 	fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
	// 		channel, id, conn.UserID, conn.Username, conn.UserRole),
	// 	fmt.Sprintf("%d Tasks", len(*tasks)))
	logger(start, "Task.SiteList",
		fmt.Sprintf("Site %d", data.ID),
		fmt.Sprintf("%d Tasks", len(*tasks)),
		data.Channel, conn.UserID, "task", 0, false)

	return nil
}

// Get a list of tasks for a given stoppage event
func (t *TaskRPC) StoppageList(data shared.TaskRPCData, tasks *[]shared.Task) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select 
		t.*,
		m.name as machine_name,
		s.name as site_name,s.id as site_id,
		u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.event_id=$1
		order by t.startdate desc,id desc`, data.ID).
		QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	for k, v := range *tasks {

		// trim the descr field
		if len(v.Descr) > 80 {
			(*tasks)[k].Descr = v.Descr[:80] + "..."
		}

		// Get the latest thumbnails for this task, if present
		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='task' and entity_id=$1) 
			or (entity='event' and entity_id=$2) 
			or (entity='sched' and entity_id=$3) 
			order by type,id desc`, v.ID, v.EventID, v.SchedID).
			QueryStructs(&photos)
		(*tasks)[k].Photos = photos
	}

	logger(start, "Task.StoppageList",
		fmt.Sprintf("Stoppage Event %d", data.ID),
		fmt.Sprintf("%d Tasks", len(*tasks)),
		data.Channel, conn.UserID, "task", 0, false)

	return nil
}

func (t *TaskRPC) Complete(data shared.TaskRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Mark the task as complete
	DB.SQL(`update task 
		set completed_date=now()
		where id=$1`, data.Task.ID).Exec()

	// Decrement the stock values for any parts used
	for _, v := range data.Task.Parts {
		if v.QtyUsed != 0 {
			thePart := shared.Part{}
			DB.SQL(`select * from part where id=$1`, v.PartID).QueryStruct(&thePart)
			thePart.CurrentStock -= v.QtyUsed
			DB.SQL(`update part set current_stock=$2 where id=$1`, thePart.ID, thePart.CurrentStock).Exec()
			DB.SQL(`insert into part_stock (part_id, stock_level, descr) values ($1, $2, $3)`,
				thePart.ID,
				thePart.CurrentStock,
				fmt.Sprintf("Used %.02f on task %06d : %s", v.QtyUsed, data.Task.ID, v.Notes)).
				Exec()
		}

	}

	conn.Broadcast("task", "update", data.Task.ID)

	// If the task has a parent event, then clear the event IF there are
	// no incomplete tasks left against that event.

	if data.Task.EventID != 0 {

		// are there any incomplete tasks still attached to this event ?
		numTasks := 0
		DB.SQL(`select count(*) 
			from task 
			where event_id=$1 and completed_date is null`, data.Task.EventID).
			QueryScalar(&numTasks)

		if numTasks == 0 {
			// Mark the event as complete
			DB.SQL(`update event 
				set completed=now(), status='Complete'
				where id=$1`, data.Task.EventID).Exec()

			conn.Broadcast("event", "update", data.Task.EventID)

			event := &shared.Event{}
			id := data.Task.EventID
			DB.SQL(`select
				e.*,m.name as machine_name,s.name as site_name,u.username as username
				from event e
					left join machine m on m.id=e.machine_id
					left join site s on s.id=m.site_id
					left join users u on u.id=e.created_by
				where e.id=$1`, id).QueryStruct(event)

			// fetch all assignments
			DB.SQL(`select u.username
				from task t
				left join users u on u.id=t.assigned_to
				where t.event_id=$1`, id).
				QueryStructs(&event.AssignedTo)

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
				case "Pnuematic":
					fieldName = "pnuematic"
				case "Conveyor":
					fieldName = "conveyor"
				case "Encoder":
					fieldName = "encoder"
				case "StripGuide":
					fieldName = "strip_guide"

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
				machine.Pnuematic != "Running" ||
				machine.Conveyor != "Running" ||
				machine.Encoder != "Running" ||
				machine.StripGuide != "Running" ||
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
					conn.Broadcast("machine", "update", event.MachineID)
				}
			}

		} // after clearing this task, there are no more tasks attached to the stoppage
	} // Task is linked to a stoppage

	// 2 SMS's to generate :
	// - 1 to the person that allocated the task to the tech
	// - 1 to the person that raised the original alert
	// Note that Scheduled Tasks will generate neither

	machine := shared.Machine{}
	DB.SQL(`select * from machine where id=$1`, data.Task.MachineID).QueryStruct(&machine)

	phoneNumber := ""

	if data.Task.AssignedBy != nil {
		smsMsg := fmt.Sprintf("Task %06d Completed:\n %s - %s",
			data.Task.ID,
			machine.Name,
			data.Task.Component)

		DB.SQL(`select sms from users where id=$1`, data.Task.AssignedBy).QueryScalar(&phoneNumber)

		if Config.SMSOn {

			if phoneNumber != "" {
				SendSMS(phoneNumber, smsMsg, fmt.Sprintf("%d", data.Task.ID), *data.Task.AssignedBy)
			} else {
				log.Println("No Phone Number for SMS:", smsMsg)
			}
		} else {
			log.Println("Will send SMS:", smsMsg, "to", phoneNumber)
		}
	}

	if data.Task.EventID != 0 {
		smsMsg := fmt.Sprintf("Task %06d Completed:\n %s - %s",
			data.Task.ID,
			machine.Name,
			data.Task.Component)

		event := shared.Event{}
		DB.SQL(`select * from event where id=$1`, data.Task.EventID).QueryStruct(&event)

		phoneNumber2 := ""
		DB.SQL(`select sms from users where id=$1`, event.CreatedBy).QueryScalar(&phoneNumber2)

		if Config.SMSOn {

			if phoneNumber2 != "" {

				if phoneNumber == phoneNumber2 {
					log.Println("Stoppage raiser and Task Assigner are the same person .. dont need 2 SMS to the same person")
				} else {
					SendSMS(phoneNumber2, smsMsg, fmt.Sprintf("%d", data.Task.ID), event.CreatedBy)
				}
			} else {
				log.Println("No Phone Number for SMS:", smsMsg)
			}
		} else {
			log.Println("Will send SMS:", smsMsg, "to", phoneNumber2)
		}

	}

	logger(start, "Task.Complete",
		fmt.Sprintf("Channel %d, Task %d User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		"Task marked as complete",
		data.Channel, conn.UserID, "task", data.Task.ID, true)

	*done = true
	return nil
}

func (t *TaskRPC) Retransmit(data shared.TaskRPCData, result *string) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	println("Task = ", data.Task)

	machine := shared.Machine{}
	DB.SQL(`select * from machine where id=$1`, data.Task.MachineID).QueryStruct(&machine)

	phoneNumber := ""
	DB.SQL(`select sms from users where id=$1`, data.Task.AssignedTo).QueryScalar(&phoneNumber)
	println("Phone =", phoneNumber)

	if phoneNumber != "" {
		useMobile := true
		DB.SQL(`select use_mobile from users where id=$1`, data.Task.AssignedTo).QueryScalar(&useMobile)

		notes := data.Task.Descr
		if len(notes) > 40 {
			notes = notes[:40] + "..."
		}
		smsMsg := fmt.Sprintf("Task %06d:\n %s - %s : %s",
			data.Task.ID,
			notes,
			machine.Name,
			data.Task.CompType)

		if Config.SMSOn {

			if useMobile {
				SendSMS(phoneNumber, smsMsg, fmt.Sprintf("%d", data.Task.ID), *data.Task.AssignedTo)
				*result = "Sent: " + smsMsg + " to " + phoneNumber
			} else {
				*result = "User Has Requested no SMS, otherwise we would send: " + smsMsg + " to " + phoneNumber
			}
		} else {
			log.Println("Will send SMS:", smsMsg, "to", phoneNumber)
			*result = "SMS is turned off, but will send: " + smsMsg + " to " + phoneNumber
		}

	} else {
		*result = "User has no phone number registered for SMS"
	}

	logger(start, "Task.Retransmit",
		fmt.Sprintf("Channel %d, Task %d User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		"Retrans message",
		data.Channel, conn.UserID, "task", data.Task.ID, false)

	return nil
}

func (t *TaskRPC) AddParts(data shared.TaskRPCPartData, partsUsedUpdate *shared.PartsUsedUpdate) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// get the existing task_part, then remove it
	oldTaskPart := shared.TaskPart{}
	DB.SQL(`select * from task_part where task_id=$1 and part_id=$2`, data.ID, data.Part).QueryStruct(&oldTaskPart)
	DB.SQL(`delete from task_part where task_id=$1 and part_id=$2`, data.ID, data.Part).Exec()

	// insert a new task_part if the qut is not zero
	if data.Qty != 0.0 {
		DB.SQL(`insert into task_part
		(task_id,part_id,qty_used,qty)
		values ($1,$2,$3,0)`,
			data.ID, data.Part, data.Qty).Exec()
	}

	// Calculate the stock difference
	delta := data.Qty - oldTaskPart.QtyUsed

	// println("OldQty", oldTaskPart.QtyUsed, "NewQty", data.Qty, "delta", delta)

	// Update the stock on hand value on the part
	oldStockOnHand := 0.0
	DB.SQL(`select current_stock from part where id=$1`, data.Part).QueryScalar(&oldStockOnHand)
	newStockOnHand := oldStockOnHand - delta

	DB.SQL(`update part set current_stock=$2 where id=$1`, data.Part, newStockOnHand).Exec()

	// Insert a stock audit record against the part
	DB.SQL(`insert into part_stock
		(part_id,stock_level,descr)
		values ($1,$2,$3)`,
		data.Part,
		newStockOnHand,
		fmt.Sprintf("Used %.1f on Task %06d", delta, data.ID)).Exec()

	// Get the new total material cost for this whole task

	totalMaterialCost := 0.0
	DB.SQL(`select 
		sum(t.qty_used * p.latest_price) as totalm 
		from task_part t 
		left join part p on p.id=t.part_id 
		where t.task_id=$1`, data.ID).QueryScalar(&totalMaterialCost)
	DB.SQL(`update task set material_cost=$2 where id=$1`, data.ID, totalMaterialCost).Exec()

	logger(start, "Task.AddParts",
		fmt.Sprintf("Channel %d, Task %d Part %d Qty %.2f User %d %s %s",
			data.Channel, data.ID, data.Part, data.Qty, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Delta %.2f", delta),
		data.Channel, conn.UserID, "task_part", data.ID, true)

	partsUsedUpdate.NewStockOnHand = newStockOnHand
	partsUsedUpdate.TotalMaterialCost = totalMaterialCost

	return nil
}

func (t *TaskRPC) GetQtyUsed(data shared.TaskRPCPartData, qty *float64) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// get the existing task_part qty
	DB.SQL(`select qty_used from task_part where task_id=$1 and part_id=$2`, data.ID, data.Part).QueryScalar(qty)

	logger(start, "Task.GetQtyUsed",
		fmt.Sprintf("Channel %d, Task %d Part %d User %d %s %s",
			data.Channel, data.ID, data.Part, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Qty %.2f", *qty),
		data.Channel, conn.UserID, "task_part", data.ID, false)

	return nil
}

func (t *TaskRPC) GetParts(data shared.TaskRPCData, parts *[]shared.TaskPart) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select t.*,p.name as part_name,p.stock_code
		from task_part t
		left join part p on p.id=t.part_id
		where t.task_id=$1
		order by p.name`, data.ID).
		QueryStructs(parts)

	logger(start, "Task.GetParts",
		fmt.Sprintf("Channel %d, Task %d User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Got %d parts", len(*parts)),
		data.Channel, conn.UserID, "task_part", data.ID, false)

	return nil

}
