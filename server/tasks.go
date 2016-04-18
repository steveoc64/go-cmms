package main

import (
	"fmt"
	"log"
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
	err := DB.SQL(`select * from sched_task where machine_id=$1`, machineID).QueryStructs(tasks)

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
		fmt.Sprintf("id %d", id),
		fmt.Sprintf("%s %s", task.Freq, task.Descr))

	return nil
}

func (t *TaskRPC) UpdateSched(data *shared.SchedTaskUpdateData, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("sched_task").
		SetWhitelist(data.SchedTask,
			"comp_type", "tool_id",
			"component", "descr", "startdate", "oneoffdate", "freq", "days", "week", "count",
			"labour_cost", "material_cost", "duration_days").
		Where("id = $1", data.SchedTask.ID).
		Exec()

	logger(start, "Task.UpdateSched",
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
			data.Channel, data.SchedTask.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %d %s %s",
			data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr))

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
		fmt.Sprintf("Channel %d, Task %d, User %d %s %s",
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

	DB.InsertInto("sched_task").
		Whitelist("machine_id", "comp_type", "tool_id",
			"component", "descr", "startdate", "oneoffdate", "freq", "days", "week", "count",
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

	return nil
}

func (t *TaskRPC) List(channel int, tasks *[]shared.Task) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read the sites that this user has access to
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

	logger(start, "Task.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Tasks", len(*tasks)))

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
	start := time.Now()

	numTasks := 0

	month := runDate.Month()
	year := runDate.Year()

	nextWeek := runDate.AddDate(0, 0, 7)
	tommorow := runDate.AddDate(0, 0, 1)
	priorWeek := runDate.AddDate(0, 0, -7)

	log.Printf("»»» SchedTask Generate run for %s", runDate.Format(rfc3339DateLayout))
	log.Printf(".. Next Week = %s", nextWeek.Format(rfc3339DateLayout))
	log.Printf(".. Prior Week = %s", priorWeek.Format(rfc3339DateLayout))

	// work out which week of the month we are in
	firstOfTheMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	firstday := firstOfTheMonth.Weekday()
	firstWeek := firstOfTheMonth

	log.Printf(".. first of the month falls on a %s", firstday)
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
	fifthWeek := firstWeek.AddDate(0, 0, 28)

	log.Printf(".. 1st Week is %s", firstWeek.Format(rfc3339DateLayout))
	log.Printf(".. 2nd Week is %s", secondWeek.Format(rfc3339DateLayout))
	log.Printf(".. 3rd Week is %s", thirdWeek.Format(rfc3339DateLayout))
	log.Printf(".. 4th Week is %s", fourthWeek.Format(rfc3339DateLayout))

	// Go through each scheduled task in turn
	scheds := []shared.SchedTask{}
	newTask := shared.Task{}

	DB.SQL(`select * from sched_task order by id`).QueryStructs(&scheds)

	for _, st := range scheds {
		// log.Printf("consider the case of task %d with freq %s", st.ID, st.Freq)
		switch st.Freq {
		case "Monthly":
			dueDate := firstWeek
			lastDate := secondWeek
			if st.Week == nil {
				log.Printf("Error - Monthly Task %d has null week", st.ID)
				break
			}
			switch *st.Week {
			case 1:
				dueDate = firstWeek
				lastDate = secondWeek
			case 2:
				dueDate = secondWeek
				lastDate = thirdWeek
			case 3:
				dueDate = thirdWeek
				lastDate = fourthWeek
			case 4:
				dueDate = fourthWeek
				lastDate = fifthWeek
			}

			if dueDate.After(priorWeek) && dueDate.Before(tommorow) {
				log.Printf("Task %d On Week %d Next Due on %s last date %s",
					st.ID, *st.Week,
					dueDate.Format(rfc3339DateLayout),
					lastDate.Format(rfc3339DateLayout))
				// Generate a new Task record
				genTask(st, &newTask, dueDate, lastDate)
				numTasks++

			} else {
				log.Printf("Task %d On Week %d Next Due on %s last date %s (not due yet)",
					st.ID, *st.Week,
					dueDate.Format(rfc3339DateLayout),
					lastDate.Format(rfc3339DateLayout))
			}
		case "Yearly":
			// If the one off date is within the window
			if st.StartDate == nil {
				log.Printf("Error - Task %d is yearly but has a null startdate", st.ID)
			} else {
				if st.StartDate.After(priorWeek) && st.StartDate.Before(nextWeek) {
					log.Printf("Task %d Yearly date %s is in the window %s - %s",
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
		case "One Off":
			// If the one off date is within the window
			if st.OneOffDate == nil {
				log.Printf("Error - Task %d is yearly but has a null startdate", st.ID)
			} else {
				if st.OneOffDate.After(priorWeek) && st.OneOffDate.Before(nextWeek) {
					log.Printf("Task %d OneOff date %s is in the window %s - %s",
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
		case "Every N Days":
			// Get the last one generated
			// If there is none, then create the first one
			// Else, calculate the next date, and check if its in the window
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

	log.Printf("»»» Generating Task from Sched %d", st.ID)

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

	return nil
}
