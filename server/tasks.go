package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/steveoc64/go-cmms/shared"
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

func (t *TaskRPC) GetSched(id int, task *shared.SchedTask) error {
	start := time.Now()

	err := DB.SQL(`select * from sched_task where id=$1`, id).QueryStruct(task)

	if err != nil {
		log.Println(err.Error())
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

func (t *TaskRPC) SchedPlay(id int, ok *bool) error {
	start := time.Now()

	DB.SQL(`update sched_task set paused=false,last_generated=now() where id=$1`, id).Exec()

	logger(start, "Task.SchedPlay",
		fmt.Sprintf("Sched %d", id),
		"Now Running")

	*ok = true

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

	DB.InsertInto("sched_task").
		Whitelist("machine_id", "comp_type", "tool_id",
			"component", "descr", "startdate", "oneoffdate",
			"freq", "days", "week", "weekday", "count", "user_id",
			"labour_cost", "material_cost", "duration_days").
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

	DB.Update("task").
		SetWhitelist(data.Task,
			"log",
			"labour_cost", "material_cost").
		Where("id = $1", data.Task.ID).
		Exec()

	logger(start, "Task.Update",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.Task.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%f %f %s",
			data.Task.LabourCost, data.Task.MaterialCost, data.Task.Log))

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
	case "Worker":
		// Limit the tasks to only our own tasks
		err := DB.SQL(`select 
		t.*,m.name as machine_name,s.name as site_name,u.username as username
		from task t 
			left join machine m on m.id=t.machine_id
			left join site s on s.id=m.site_id
			left join users u on u.id=t.assigned_to
		where t.assigned_to=$1
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
		where m.site_id in $1
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

	logger(start, "Task.Get",
		fmt.Sprintf("ID %d", id),
		task.Descr)

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

func (t *TaskRPC) Generate(runDate time.Time, count *int) error {
	return schedTaskScan(runDate, count)
}

var GenerateMutex sync.Mutex

func autoGenerate() {

	log.Printf("... Running task scheduler")
	go func() {
		newTasks := 0
		for {
			schedTaskScan(time.Now(), &newTasks)
			time.Sleep(1 * time.Hour)
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

	return nil
}
