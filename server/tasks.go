package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/steveoc64/go-cmms/shared"
	// "github.com/jung-kurt/gofpdf"
	// "github.com/signintech/gopdf"
)

const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

type TaskRPC struct{}

// Get all the tasks for the given machine
func (t *TaskRPC) ListMachineSched(machineID int, tasks *[]shared.SchedTask) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(`select * from sched_task where machine_id=$1 order by id`, machineID).QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Task.ListMachineSched",
		fmt.Sprintf("Machine %d", machineID),
		fmt.Sprintf("%d tasks", len(*tasks)))

	return nil
}

// Get all the tasks that contain the given hash
func (t *TaskRPC) ListHashSched(id int, tasks *[]shared.SchedTask) error {
	start := time.Now()

	hashname := ""
	err := DB.SQL(`select name from hashtag where id=$1`, id).QueryScalar(&hashname)

	if err != nil {
		log.Println(err.Error())
	} else {

		// Read the sites that this user has access to
		err = DB.SQL(`select t.*
		from sched_task t
		where lower(descr) like lower($1)
		order by id`, "%#"+hashname+"%").QueryStructs(tasks)

		if err != nil {
			log.Println(err.Error())
		}
	}

	logger(start, "Task.ListHashSched",
		fmt.Sprintf("Hash %d", id),
		fmt.Sprintf("%d tasks", len(*tasks)))

	return nil
}

func (t *TaskRPC) GetSched(id int, task *shared.SchedTask) error {
	start := time.Now()

	err := DB.SQL(`select * from sched_task where id=$1`, id).QueryStruct(task)

	if err != nil {
		log.Println(err.Error())
	} else {
		// Get the parts allowed from the PartClass of the machine
		partClass := 0
		DB.SQL(`select part_class from machine where id=$1`, task.MachineID).
			QueryScalar(&partClass)

		// if partClass != 0 {
		// 	DB.SQL(`select * from part where class=$1 order by part.name`, partClass).
		// 		QueryStructs(&task.PartsAllowed)
		// }

		log.Println("key task", task.ID, "class", partClass)
		// Get the parts used in this sched
		DB.SQL(`select 
			p.id as part_id,p.stock_code as stock_code,p.name as name,p.qty_type as qty_type,
			u.qty as qty,u.notes as notes
			from part p
			left join sched_task_part u on u.part_id=p.id and u.task_id=$1
			where p.class=$2
			order by p.name`, task.ID, partClass).
			QueryStructs(&task.PartsRequired)

		// dereference all the qty and notes fields
		for i, v := range task.PartsRequired {
			if v.QtyPtr != nil {
				task.PartsRequired[i].Qty = *v.QtyPtr
			}
			if v.NotesPtr != nil {
				task.PartsRequired[i].Notes = *v.NotesPtr
			}
		}
	}

	logger(start, "Task.GetSched",
		fmt.Sprintf("Sched %d", id),
		fmt.Sprintf("%s %s", task.Freq, task.Descr))

	return nil
}

func (t *TaskRPC) UpdateSched(data *shared.SchedTaskUpdateData, ok *bool) error {
	start := time.Now()

	if data.SchedTask.Freq == "Every N Days" {
		if data.SchedTask.Days == nil {
			i := 1
			data.SchedTask.Days = &i
		} else {
			if *data.SchedTask.Days < 1 {
				*data.SchedTask.Days = 1
			}
		}
	}
	if data.SchedTask.DurationDays < 1 {
		data.SchedTask.DurationDays = 1
	}

	conn := Connections.Get(data.Channel)

	DB.Update("sched_task").
		SetWhitelist(data.SchedTask,
			"comp_type", "tool_id",
			"component", "descr", "startdate", "oneoffdate",
			"freq", "days", "week", "weekday", "count", "user_id",
			"labour_cost", "material_cost", "duration_days").
		Where("id = $1", data.SchedTask.ID).
		Exec()

	logger(start, "Task.UpdateSched",
		fmt.Sprintf("Channel %d, Sched %d, User %d %s %s",
			data.Channel, data.SchedTask.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %d %s %s",
			data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr))

	*ok = true

	newTasks := 0
	schedTaskScan(time.Now(), &newTasks)
	return nil
}

func (t *TaskRPC) SchedPart(data shared.PartReqEdit, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Kill any existing relationship, and then create a new one
	DB.SQL(`delete from sched_task_part where task_id=$1 and part_id=$2`,
		data.Task.ID, data.Part.PartID).Exec()

	record := shared.SchedTaskPart{
		TaskID: data.Task.ID,
		PartID: data.Part.PartID,
		Qty:    data.Part.Qty,
		Notes:  data.Part.Notes,
	}

	DB.InsertInto("sched_task_part").
		Whitelist("task_id", "part_id", "qty", "notes").
		Record(record).
		Exec()

	logger(start, "Task.SchedPart",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Part %d Qty %f Notes %s",
			data.Part.PartID, data.Part.Qty, data.Part.Notes))

	*ok = true

	return nil
}

func (t *TaskRPC) SchedPlay(id int, ok *bool) error {
	start := time.Now()

	DB.SQL(`update sched_task set paused=false,last_generated=now() where id=$1`, id).Exec()

	logger(start, "Task.SchedPlay",
		fmt.Sprintf("Sched %d", id),
		"Now Running")

	*ok = true

	newTasks := 0
	schedTaskScan(time.Now(), &newTasks)

	return nil
}

func (t *TaskRPC) SchedPause(id int, ok *bool) error {
	start := time.Now()

	DB.SQL(`update sched_task set paused=true where id=$1`, id).Exec()

	logger(start, "Task.SchedPause",
		fmt.Sprintf("Sched %d", id),
		"Now Paused")

	*ok = true

	return nil
}

func (t *TaskRPC) DeleteSched(data *shared.SchedTaskUpdateData, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("sched_task").
		Where("id=$1", data.SchedTask.ID).
		Exec()

	logger(start, "Task.DeleteSched",
		fmt.Sprintf("Channel %d, Sched %d, User %d %s %s",
			data.Channel, data.SchedTask.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %d %s %s",
			data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr))

	*ok = true
	return nil
}

func (t *TaskRPC) InsertSched(data *shared.SchedTaskUpdateData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	if data.SchedTask.Freq == "Every N Days" {
		if data.SchedTask.Days == nil {
			i := 1
			data.SchedTask.Days = &i
		} else {
			if *data.SchedTask.Days < 1 {
				*data.SchedTask.Days = 1
			}
		}
	}
	if data.SchedTask.DurationDays < 1 {
		data.SchedTask.DurationDays = 1
	}

	// Default the schedule to paused, so we can fine tune it before starting
	// the first generation
	data.SchedTask.Paused = true

	DB.InsertInto("sched_task").
		Whitelist("machine_id", "comp_type", "tool_id",
			"component", "descr", "startdate", "oneoffdate",
			"freq", "days", "week", "weekday", "count", "user_id",
			"labour_cost", "material_cost", "duration_days", "paused").
		Record(data.SchedTask).
		Returning("id").
		QueryScalar(id)

	logger(start, "Task.InsertSched",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %d %s %d %s %s",
			*id, data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr))

	newTasks := 0
	schedTaskScan(time.Now(), &newTasks)

	return nil
}

func (t *TaskRPC) Update(data shared.TaskUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	oldTask := shared.Task{}
	DB.SQL(`select * from task where id=$1`, data.Task.ID).QueryStruct(&oldTask)

	if conn.UserRole == "Admin" {

		// Admin can re-assing the task to another user
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

	for _, v := range data.Task.Parts {
		// log.Println("part = ", v)

		DB.Update("task_part").
			SetWhitelist(v, "notes", "qty_used").
			Where("task_id=$1 and part_id=$2", data.Task.ID, v.PartID).
			Exec()
	}

	logger(start, "Task.Update",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%f %f %f %s",
			data.Task.LabourCost, data.Task.MaterialCost, data.Task.LabourHrs, data.Task.Log))

	*done = true
	return nil
}

func (t *TaskRPC) Delete(data shared.TaskUpdateData, done *bool) error {
	log.Printf("TODO TaskRPC.Delete")
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
		order by t.startdate`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Site Manager":
		// Limit the tasks to just the sites that we are in control of
		sites := []int{}

		DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&sites)

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where m.site_id in $1 and completed_date is null
		order by t.startdate`, sites).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Admin":
		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where completed_date is null
		order by t.startdate`).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
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
		fmt.Sprintf("%d Tasks", len(*tasks)))

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
		where t.assigned_to=$1 and completed_date is not null
		order by t.startdate`, conn.UserID).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Site Manager":
		// Limit the tasks to just the sites that we are in control of
		sites := []int{}

		DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&sites)

		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where m.site_id in $1 and completed_date is not null
		order by t.startdate`, sites).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
	case "Admin":
		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where completed_date is not null
		order by t.startdate`).
			QueryStructs(tasks)
		if err != nil {
			log.Println(err.Error())
		}
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
		fmt.Sprintf("%d Tasks", len(*tasks)))

	return nil
}

func (t *TaskRPC) Get(id int, task *shared.Task) error {
	start := time.Now()

	err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.id=$1`, id).
		QueryStruct(task)

	if err != nil {
		log.Println(err.Error())
	}

	// Now get all the parts for this task
	DB.SQL(`select 
		t.*,p.name as part_name,p.stock_code as stock_code,p.qty_type as qty_type
		from task_part t
		left join part p on p.id=t.part_id
		where t.task_id=$1`, id).QueryStructs(&task.Parts)

	// Now get all the checks for this task
	DB.SQL(`select * from task_check where task_id=$1`, id).QueryStructs(&task.Checks)

	logger(start, "Task.Get",
		fmt.Sprintf("ID %d", id),
		task.Descr)

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
			data.TaskCheck.TaskID, data.TaskCheck.Seq))

	*done = true
	return nil
}

func (t *TaskRPC) SiteList(id int, tasks *[]shared.Task) error {
	start := time.Now()

	// conn := Connections.Get(channel)

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
		order by t.startdate`, id).
		QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	// logger(start, "Task.SiteList",
	// 	fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
	// 		channel, id, conn.UserID, conn.Username, conn.UserRole),
	// 	fmt.Sprintf("%d Tasks", len(*tasks)))
	logger(start, "Task.SiteList",
		fmt.Sprintf("Site %d", id),
		fmt.Sprintf("%d Tasks", len(*tasks)))

	return nil
}

func (t *TaskRPC) StoppageList(id int, tasks *[]shared.Task) error {
	start := time.Now()

	// conn := Connections.Get(channel)

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
		order by t.startdate`, id).
		QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	// logger(start, "Task.SiteList",
	// 	fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
	// 		channel, id, conn.UserID, conn.Username, conn.UserRole),
	// 	fmt.Sprintf("%d Tasks", len(*tasks)))
	logger(start, "Task.StoppageList",
		fmt.Sprintf("Stoppage Event %d", id),
		fmt.Sprintf("%d Tasks", len(*tasks)))

	return nil
}

func (t *TaskRPC) Generate(runDate time.Time, count *int) error {
	return schedTaskScan(runDate, count)
}

var GenerateMutex sync.Mutex

func autoGenerate() {

	log.Printf("... Running task scheduler")
	go func() {
		newTasks := 0

		// Init Hours to be the hour of the day
		hours := time.Now().Hour()

		for {
			schedTaskScan(time.Now(), &newTasks)
			time.Sleep(1 * time.Hour)

			hours++
			if hours >= 24 {
				hours = time.Now().Hour()
				log.Println("24 Hours - db backup")
				out, err := exec.Command("../scripts/cmms-backup.sh").Output()
				if err != nil {
					log.Println(err)
				} else {
					log.Println(string(out))
				}
			}
		}
	}()
}

func schedTaskScan(runDate time.Time, count *int) error {

	GenerateMutex.Lock()
	defer GenerateMutex.Unlock()

	start := time.Now()

	numTasks := 0

	month := runDate.Month()
	year := runDate.Year()

	today := time.Now()
	nextWeek := runDate.AddDate(0, 0, 7)
	tommorow := runDate.AddDate(0, 0, 1)
	priorWeek := runDate.AddDate(0, 0, -7)

	log.Printf("»»» SchedTask Generate run for %s", runDate.Format(rfc3339DateLayout))

	// work out which week of the month we are in
	firstOfTheMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	firstday := firstOfTheMonth.Weekday()
	firstWeek := firstOfTheMonth

	dd := (int)(firstday)
	switch dd {
	case 0:
		firstWeek = firstWeek.AddDate(0, 0, 1)
	case 1:
		// already set
	default:
		firstWeek = firstWeek.AddDate(0, 0, 8-dd)
	}

	secondWeek := firstWeek.AddDate(0, 0, 7)
	thirdWeek := firstWeek.AddDate(0, 0, 14)
	fourthWeek := firstWeek.AddDate(0, 0, 21)
	// fifthWeek := firstWeek.AddDate(0, 0, 28)

	log.Printf(".. first of the month falls on a %s", firstday)
	log.Printf(".. 1st Week is %s", firstWeek.Format(rfc3339DateLayout))
	log.Printf(".. 2nd Week is %s", secondWeek.Format(rfc3339DateLayout))
	log.Printf(".. 3rd Week is %s", thirdWeek.Format(rfc3339DateLayout))
	log.Printf(".. 4th Week is %s", fourthWeek.Format(rfc3339DateLayout))
	log.Printf(".. Next Week = %s", nextWeek.Format(rfc3339DateLayout))
	log.Printf(".. Prior Week = %s", priorWeek.Format(rfc3339DateLayout))
	log.Printf(".. Tomorrow = %s", tommorow.Format(rfc3339DateLayout))

	// Go through each scheduled task in turn
	scheds := []shared.SchedTask{}
	newTask := shared.Task{}

	DB.SQL(`select * from sched_task where paused=false order by id`).QueryStructs(&scheds)

	for _, st := range scheds {
		doit := true

		// if st.LastGenerated == nil {
		// 	log.Printf("--------- Processing Sched Task %d with freq %s never generated yet", st.ID, st.Freq)
		// } else {
		// 	log.Printf("--------- Processing Sched Task %d with freq %s, last gen on %s", st.ID, st.Freq, st.LastGenerated.Format(rfc3339DateLayout))
		// }

		switch st.Freq {
		case "Monthly":
			dueDate := firstWeek
			lastDate := secondWeek
			if st.Week == nil {
				log.Printf("Error - Monthly Task %d has null week", st.ID)
				break
			}
			if st.WeekDay == nil {
				log.Printf("Error - Monthly Task %d has null weekday", st.ID)
			}
			if *st.WeekDay < 1 {
				*st.WeekDay = 1
			}
			if *st.WeekDay > 5 {
				*st.WeekDay = 5
			}
			switch *st.Week {
			case 1:
				dueDate = firstWeek
			case 2:
				dueDate = secondWeek
			case 3:
				dueDate = thirdWeek
			case 4:
				dueDate = fourthWeek
			}

			lastDate = dueDate.AddDate(0, 0, 4) // the friday of the week
			realDueDate := dueDate.AddDate(0, 0, *st.WeekDay-1)

			// now that we knov the duedate, check that we havent already generated this task
			if st.LastGenerated != nil {
				if st.LastGenerated.Format(rfc3339DateLayout) == realDueDate.Format(rfc3339DateLayout) {
					// log.Printf("Task %d has already been generated for %s", st.ID, dueDate.Format(rfc3339DateLayout))
					doit = false
				}
			}

			if doit {
				if dueDate.After(priorWeek) && dueDate.Before(tommorow) {

					// Excellent - the Week that the task belongs to falls inside the window
					// so its now safe to increment the actual start date by day of the week

					log.Printf("»»» Task %d On Week %d Next Due on %s last date %s",
						st.ID, *st.Week,
						realDueDate.Format(rfc3339DateLayout),
						lastDate.Format(rfc3339DateLayout))

					// Generate a new Task record
					genTask(st, &newTask, realDueDate, lastDate)
					numTasks++
				}
			}
		case "Yearly":
			// If the one off date is within the window
			if st.StartDate == nil {
				log.Printf("Error - Task %d is yearly but has a null startdate", st.ID)
			} else {
				if st.LastGenerated != nil {

					if st.LastGenerated.Format(rfc3339DateLayout) == st.StartDate.Format(rfc3339DateLayout) {
						// log.Printf("Task %d has already been generated for %s", st.ID, st.StartDate.Format(rfc3339DateLayout))
						doit = false
					}
				}
				if doit {
					if st.StartDate.After(priorWeek) && st.StartDate.Before(nextWeek) {
						log.Printf("»»» Task %d Yearly date %s is in the window %s - %s",
							st.ID,
							st.StartDate.Format(rfc3339DateLayout),
							priorWeek.Format(rfc3339DateLayout),
							nextWeek.Format(rfc3339DateLayout))

						// Generate a new Task record
						dueDate := *st.StartDate
						genTask(st, &newTask, dueDate, dueDate.AddDate(0, 0, st.DurationDays))
						numTasks++
					}
				}
			}
		case "One Off":
			// If the one off date is within the window
			if st.OneOffDate == nil {
				log.Printf("Error - Task %d is yearly but has a null startdate", st.ID)
			} else {
				if st.LastGenerated != nil {
					if st.LastGenerated.Format(rfc3339DateLayout) == st.OneOffDate.Format(rfc3339DateLayout) {
						// log.Printf("Task %d has already been generated for %s", st.ID, st.OneOffDate.Format(rfc3339DateLayout))
						doit = false
					}
				}
				if doit {
					if st.OneOffDate.After(priorWeek) && st.OneOffDate.Before(nextWeek) {
						log.Printf("»»» Task %d OneOff date %s is in the window %s - %s",
							st.ID,
							st.OneOffDate.Format(rfc3339DateLayout),
							priorWeek.Format(rfc3339DateLayout),
							nextWeek.Format(rfc3339DateLayout))

						// Generate a new Task record
						dueDate := *st.OneOffDate
						genTask(st, &newTask, dueDate, dueDate.AddDate(0, 0, st.DurationDays))
						numTasks++
					}
				}
			}
		case "Every N Days":
			if st.Days == nil {
				log.Printf("Error - Task %d on every N days has no days specified", st.ID)
			} else {

				// Get the last one generated
				// If there is none, then create the first one
				if st.LastGenerated == nil {
					log.Printf("»»» Task %d Every %d days, first entry - start now",
						st.ID,
						*st.Days)

					// Generate a new Task record
					genTask(st, &newTask, today, today.AddDate(0, 0, st.DurationDays))
					numTasks++
					st.LastGenerated = &today

				}
				// Else, calculate the next date, and check if its in the window
				allDone := false
				nextDate := st.LastGenerated.AddDate(0, 0, *st.Days)

				for !allDone {
					if nextDate.After(priorWeek) && nextDate.Before(nextWeek) {
						log.Printf("»»» Task %d Every %d days, next due at %s is within %s - %s",
							st.ID,
							*st.Days,
							nextDate.Format(rfc3339DateLayout),
							priorWeek.Format(rfc3339DateLayout),
							nextWeek.Format(rfc3339DateLayout))

						// Generate a new Task record
						genTask(st, &newTask, nextDate, nextDate.AddDate(0, 0, st.DurationDays))
						numTasks++

						// keep looping, looking at the next date
						nextDate = nextDate.AddDate(0, 0, *st.Days)
					} else {
						allDone = true
					}
				}
			}
		case "Job Count":
			// Get the last generated job count
			// If there is none, then create the first one
			// Else, calculate the next count, and check if it has been exceeded
		}
	}

	*count = numTasks
	logger(start, "Task.Generate",
		fmt.Sprintf("As of date %s", runDate.Format(rfc3339DateLayout)),
		fmt.Sprintf("%d New Tasks Generated", *count))

	return nil
}

type machineLookup struct {
	MachineUser int `db:"machine_user"`
	SiteUser    int `db:"site_user"`
}

func genTask(st shared.SchedTask, task *shared.Task, startDate time.Time, dueDate time.Time) error {

	// log.Printf("»»» Generating Task from Sched %d Freq %s for date %s", st.ID, st.Freq, startDate.Format(rfc3339DateLayout))

	// Calculate the receiving user for this task
	userIDs := machineLookup{}
	DB.SQL(`select 
		m.tasks_to as machine_user,s.tasks_to as site_user
		from machine m 
		left join site s on s.id=m.site_id
		where m.id=$1`, st.MachineID).
		QueryStruct(&userIDs)

	userID := userIDs.SiteUser
	if userIDs.MachineUser != 0 {
		userID = userIDs.MachineUser
	}

	escDate := startDate.AddDate(0, 1, 0)

	task.MachineID = st.MachineID
	task.SchedID = st.ID
	task.CompType = st.CompType
	task.ToolID = st.ToolID
	task.Component = st.Component
	task.Descr = st.Descr
	task.StartDate = &startDate
	task.DueDate = &dueDate
	task.EscalateDate = &escDate
	task.AssignedTo = &userID
	task.AssignedDate = &startDate
	task.LabourEst = st.LabourCost
	task.MaterialEst = st.MaterialCost

	DB.InsertInto("task").
		Whitelist("machine_id", "sched_id", "comp_type", "tool_id", "component",
			"descr", "startdate", "due_date", "escalate_date",
			"assigned_to", "assigned_date", "labour_est", "material_est").
		Record(task).
		Returning("id").
		QueryScalar(&task.ID)

	DB.SQL(`update sched_task set last_generated=$2 where id=$1`, st.ID, startDate).Exec()

	// Now copy across the parts usage from the sched
	schedParts := []shared.SchedTaskPart{}
	DB.SQL(`select * from sched_task_part where task_id=$1`, st.ID).QueryStructs(&schedParts)

	for _, s := range schedParts {

		taskPart := shared.TaskPart{
			TaskID: task.ID,
			PartID: s.PartID,
			Qty:    s.Qty,
			Notes:  s.Notes,
		}

		DB.InsertInto("task_part").
			Whitelist("task_id", "part_id", "qty", "notes").
			Record(taskPart).
			Exec()
	}

	// Expand out using the hashtags
	hasHashtag := false
	oldDescr := st.Descr

	if strings.Contains(oldDescr, "#") {
		hashes := []shared.Hashtag{}
		// Apply the longest hashtag first

		DB.SQL(`select * from hashtag order by length(name) desc`).QueryStructs(&hashes)

		// Keep looping through doing text conversions until there is
		// nothing left to expand
		stillLooking := true
		for stillLooking {
			stillLooking = false
			for _, v := range hashes {
				theHash := "#" + v.Name
				if strings.Contains(oldDescr, theHash) {
					oldDescr = strings.Replace(oldDescr, theHash, v.Descr, -1)
					hasHashtag = true
					stillLooking = true
				}
			}
		}
	}

	// Now generate the task check items based on the description field of the schedtask
	lines := strings.Split(oldDescr, "\n")
	seq := 1
	descr := ""

	for _, l := range lines {
		if strings.HasPrefix(l, "- ") {
			check := shared.TaskCheck{
				TaskID: task.ID,
				Seq:    seq,
				Descr:  l[2:],
				Done:   false,
			}

			DB.InsertInto("task_check").
				Whitelist("task_id", "seq", "descr", "done").
				Record(check).
				Exec()
			seq++
		} else {
			descr += l
			descr += "\n"
		}
	}
	log.Println("Modded desc from", st.Descr, "to", descr)
	if hasHashtag || seq > 1 {
		DB.SQL(`update task set descr=$1 where id=$2`, descr, task.ID).Exec()
	}

	return nil
}

func (t *TaskRPC) HashtagList(channel int, hashtags *[]shared.Hashtag) error {
	start := time.Now()

	conn := Connections.Get(channel)

	DB.SQL(`select * from hashtag order by name`).QueryStructs(hashtags)

	logger(start, "Task.HashtagList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Hashtags", len(*hashtags)))

	return nil
}

func (t *TaskRPC) HashtagGet(id int, hashtag *shared.Hashtag) error {
	start := time.Now()

	DB.SQL(`select * from hashtag where id=$1`, id).QueryStruct(hashtag)

	logger(start, "Task.HashtagGet",
		fmt.Sprintf("ID %d", id),
		fmt.Sprintf("Name %s", hashtag.Name))

	return nil
}

func (t *TaskRPC) HashtagInsert(data shared.HashtagUpdateData, id *int) error {
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
			data.Hashtag.ID, data.Hashtag.Name))

	conn.Broadcast("hashtag", "insert", *id)
	return nil
}

func (t *TaskRPC) HashtagDelete(data shared.HashtagUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`delete from hashtag where id=$1`, data.Hashtag.ID).Exec()

	logger(start, "Task.HashtagDelete",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID %d Name %s",
			data.Hashtag.ID, data.Hashtag.Name))

	*done = true
	conn.Broadcast("hashtag", "delete", data.Hashtag.ID)

	return nil
}

func (t *TaskRPC) HashtagUpdate(data shared.HashtagUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`update hashtag set name=$2,descr=$3 where id=$1`,
		data.Hashtag.ID, data.Hashtag.Name, data.Hashtag.Descr).Exec()

	logger(start, "Task.HashtagUpdate",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID %d Name %s",
			data.Hashtag.ID, data.Hashtag.Name))

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

func (t *TaskRPC) Diary(channel int, done *bool) error {

	os.Mkdir("public/pdf", 0700)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Technicians Weekly Task List")
	// fileStr := example.Filename("basic")
	err := pdf.OutputFileAndClose("public/pdf/diary.pdf")
	if err != nil {
		log.Println("PDF Err", err.Error())
	}

	pdf = gofpdf.New("P", "mm", "A4", "")
	type countryType struct {
		nameStr, capitalStr, areaStr, popStr string
	}
	countryList := make([]countryType, 0, 8)
	header := []string{"Country", "Capital", "Area (sq km)", "Pop. (thousands)"}
	loadData := func(fileStr string) {
		fl, err := os.Open(fileStr)
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var c countryType
			for scanner.Scan() {
				// Austria;Vienna;83859;8075
				lineStr := scanner.Text()
				list := strings.Split(lineStr, ";")
				if len(list) == 4 {
					c.nameStr = list[0]
					c.capitalStr = list[1]
					c.areaStr = list[2]
					c.popStr = list[3]
					countryList = append(countryList, c)
				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(countryList) == 0 {
				err = fmt.Errorf("error loading data from %s", fileStr)
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}
	// Simple table
	basicTable := func() {
		for _, str := range header {
			pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		for _, c := range countryList {
			pdf.CellFormat(40, 6, c.nameStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.capitalStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.areaStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.popStr, "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}
	// Better table
	improvedTable := func() {
		// Column widths
		w := []float64{40.0, 35.0, 40.0, 45.0}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		// 	Header
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		// Data
		for _, c := range countryList {
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.Ln(-1)
		}
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	// Colored table
	fancyTable := func() {
		// Colors, line width and bold font
		pdf.SetFillColor(0, 64, 96)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(0, 64, 96)
		pdf.SetLineWidth(.3)
		pdf.SetFont("", "B", 0)
		// 	Header
		w := []float64{40, 35, 40, 45}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		// Color and font restoration
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)
		// 	Data
		fill := false
		for _, c := range countryList {
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.Ln(-1)
			fill = !fill
		}
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	loadData("countries.txt")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	basicTable()
	pdf.AddPage()
	improvedTable()
	pdf.AddPage()
	fancyTable()
	err = pdf.OutputFileAndClose("public/pdf/countries.pdf")
	return nil
}

func (t *TaskRPC) Complete(data shared.TaskUpdateData, done *bool) error {
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

		} // after clearing this task, there are no more tasks attached to the stoppage
	} // Task is linked to a stoppage

	logger(start, "Task.Complete",
		fmt.Sprintf("Channel %d, Task %d User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		"Task marked as complete")

	*done = true
	return nil
}
