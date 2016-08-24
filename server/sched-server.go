package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	// "github.com/jung-kurt/gofpdf"
	"itrak-cmms/shared"
	// "github.com/jung-kurt/gofpdf"
	// "github.com/signintech/gopdf"
)

// Get all the tasks for the given machine
func (t *TaskRPC) ListMachineSched(data shared.MachineRPCData, tasks *[]shared.SchedTask) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select * from sched_task where machine_id=$1 order by id`, data.ID).QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	// Get the latest thumbnails for this task, if present
	for i, v := range *tasks {

		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='sched' and entity_id=$1) 
			order by type, id desc`, v.ID).
			QueryStructs(&photos)
		(*tasks)[i].Photos = photos
	}

	logger(start, "Task.ListMachineSched",
		fmt.Sprintf("Machine %d", data.ID),
		fmt.Sprintf("%d tasks", len(*tasks)),
		data.Channel, conn.UserID, "machine", 0, false)

	return nil
}

// Get all the sched tasks for the given site
func (t *TaskRPC) ListSiteSched(data shared.TaskRPCData, tasks *[]shared.SchedTask) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select t.*,m.name as machine_name
	  from sched_task t
		left join machine m on m.id=t.machine_id
		where m.site_id=$1
		order by m.name`, data.ID).QueryStructs(tasks)

	if err != nil {
		log.Println(err.Error())
	}

	// Get the latest thumbnails for this task, if present
	for i, v := range *tasks {

		photos := []shared.Photo{}
		DB.SQL(`select id,thumb 
			from photo 
			where (entity='sched' and entity_id=$1) 
			order by type, id desc`, v.ID).
			QueryStructs(&photos)
		(*tasks)[i].Photos = photos
	}

	logger(start, "Task.ListSiteSched",
		fmt.Sprintf("Data %d", data.ID),
		fmt.Sprintf("%d tasks", len(*tasks)),
		data.Channel, conn.UserID, "sched_task", 0, false)

	return nil
}

// Get all the tasks that contain the given hash
func (t *TaskRPC) ListHashSched(data shared.HashtagRPCData, tasks *[]shared.SchedTask) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	hashname := ""
	err := DB.SQL(`select name from hashtag where id=$1`, data.ID).QueryScalar(&hashname)

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

		// Get the latest thumbnails for this task, if present
		for i, v := range *tasks {

			photos := []shared.Photo{}
			DB.SQL(`select id,thumb 
			from photo 
			where (entity='sched' and entity_id=$1) 
			order by type, id desc`, v.ID).
				QueryStructs(&photos)
			(*tasks)[i].Photos = photos
		}

	}

	logger(start, "Task.ListHashSched",
		fmt.Sprintf("Hash %d", data.ID),
		fmt.Sprintf("%d tasks", len(*tasks)),
		data.Channel, conn.UserID, "sched_task", 0, false)

	return nil
}

func (t *TaskRPC) GetSched(data shared.TaskRPCData, task *shared.SchedTask) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	err := DB.SQL(`select * from sched_task where id=$1`, data.ID).QueryStruct(task)

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

		// log.Println("key task", task.ID, "class", partClass)
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

	// Get the last 8 photo previews for this task
	photos := []shared.Photo{}
	DB.SQL(`select id,preview,type,datatype,filename,entity,entity_id,notes
	 from photo
	 where (entity='sched' and entity_id=$1) 
	 order by type, id desc`, data.ID).
		QueryStructs(&photos)

	task.Photos = photos

	logger(start, "Task.GetSched",
		fmt.Sprintf("Sched %d", data.ID),
		fmt.Sprintf("%s %s", task.Freq, task.Descr),
		data.Channel, conn.UserID, "sched_task", data.ID, false)

	return nil
}

func (t *TaskRPC) UpdateSched(data shared.SchedTaskRPCData, ok *bool) error {
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

	if data.SchedTask.Freq == "Every N Months" {
		if data.SchedTask.Months == nil {
			i := 1
			data.SchedTask.Months = &i
		} else {
			if *data.SchedTask.Months < 1 {
				*data.SchedTask.Months = 1
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
			"freq", "days", "months", "week", "weekday", "count", "user_id",
			"labour_cost", "material_cost", "duration_days").
		Where("id = $1", data.SchedTask.ID).
		Exec()

	// fmt.Printf("passed in newphoto %v\n", data.SchedTask.NewPhoto)
	// If there is a new photo to be added to the task, then add it
	if data.SchedTask.NewPhoto.Data != "" {
		photo := shared.Photo{
			Data:     data.SchedTask.NewPhoto.Data,
			Filename: data.SchedTask.NewPhoto.Filename,
			Entity:   "sched",
			EntityID: data.SchedTask.ID,
		}

		// decodePhoto(photo.Data, &photo.Preview, &photo.Thumb)
		decodePhoto(&photo)
		DB.InsertInto("photo").
			Columns("entity", "entity_id", "photo", "thumb", "preview", "type", "datatype", "filename").
			Record(photo).
			Exec()
	}

	logger(start, "Task.UpdateSched",
		fmt.Sprintf("Channel %d, Sched %d, User %d %s %s",
			data.Channel, data.SchedTask.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %s %d %s %s",
			data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr),
		data.Channel, conn.UserID, "sched_task", data.SchedTask.ID, true)

	*ok = true

	newTasks := 0
	schedTaskScan(data.Channel, conn.UserID, time.Now(), &newTasks)
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
			data.Part.PartID, data.Part.Qty, data.Part.Notes),
		data.Channel, conn.UserID, "sched_task_part", 0, true)

	*ok = true

	return nil
}

func (t *TaskRPC) SchedPlay(data shared.SchedTaskRPCData, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`update sched_task set paused=false,last_generated=now() where id=$1`, data.ID).Exec()

	logger(start, "Task.SchedPlay",
		fmt.Sprintf("Sched %d", data.ID),
		"Now Running",
		data.Channel, conn.UserID, "sched_task", data.ID, true)

	*ok = true

	newTasks := 0
	schedTaskScan(data.Channel, conn.UserID, time.Now(), &newTasks)

	return nil
}

func (t *TaskRPC) SchedPause(data shared.SchedTaskRPCData, ok *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	DB.SQL(`update sched_task set paused=true where id=$1`, data.ID).Exec()

	logger(start, "Task.SchedPause",
		fmt.Sprintf("Sched %d", data.ID),
		"Now Paused",
		data.Channel, conn.UserID, "sched_task", data.ID, true)

	*ok = true

	return nil
}

func (t *TaskRPC) DeleteSched(data shared.SchedTaskRPCData, ok *bool) error {
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
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr),
		data.Channel, conn.UserID, "sched_task", data.SchedTask.ID, true)

	*ok = true
	return nil
}

func (t *TaskRPC) InsertSched(data shared.SchedTaskRPCData, id *int) error {
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

	if data.SchedTask.Months == nil {
		println("months is nil")
	} else {
		println("before", *data.SchedTask.Months)

	}
	if data.SchedTask.Freq == "Every N Months" {
		if data.SchedTask.Months == nil {
			i := 1
			data.SchedTask.Months = &i
		} else {
			if *data.SchedTask.Months < 1 {
				*data.SchedTask.Months = 1
			}
		}
	}
	if data.SchedTask.Months != nil {
		println("after", *data.SchedTask.Months)
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
			"freq", "days", "months", "week", "weekday", "count", "user_id",
			"labour_cost", "material_cost", "duration_days", "paused").
		Record(data.SchedTask).
		Returning("id").
		QueryScalar(id)

	logger(start, "Task.InsertSched",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %d %s %d %s %s",
			*id, data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr),
		data.Channel, conn.UserID, "sched_task", *id, true)

	newTasks := 0
	schedTaskScan(data.Channel, conn.UserID, time.Now(), &newTasks)

	return nil
}

// Get a list of tasks for a given stoppage event
func (t *TaskRPC) SchedList(data shared.TaskRPCData, tasks *[]shared.Task) error {
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
		where t.sched_id=$1
		order by t.startdate`, data.ID).
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
			order by id desc limit 8`, v.ID, v.EventID).
			QueryStructs(&photos)
		(*tasks)[k].Photos = photos
	}

	logger(start, "Task.SchedList",
		fmt.Sprintf("Sched %d", data.ID),
		fmt.Sprintf("%d Tasks", len(*tasks)),
		data.Channel, conn.UserID, "task", 0, false)

	return nil
}

func (t *TaskRPC) Generate(runDate time.Time, count *int) error {
	return schedTaskScan(0, 0, runDate, count)
}

var GenerateMutex sync.Mutex

func autoGenerate() {

	log.Printf("... Running task scheduler")
	go func() {
		newTasks := 0

		// Init Hours to be the hour of the day
		hours := time.Now().Hour()

		for {
			schedTaskScan(0, 0, time.Now(), &newTasks)
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

func schedTaskScan(channel int, user_id int, runDate time.Time, count *int) error {

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
		case "Every N Months":
			if st.Months == nil {
				log.Printf("Error - Task %d on every N months has no month specified", st.ID)
			} else {

				// Get the last one generated
				// If there is none, then create the first one
				if st.LastGenerated == nil {
					log.Printf("»»» Task %d Every %d months, first entry - start now",
						st.ID,
						*st.Days)

					// Generate a new Task record
					genTask(st, &newTask, today, today.AddDate(0, 0, st.DurationDays))
					numTasks++
					st.LastGenerated = &today

				}
				// Else, calculate the next date, and check if its in the window
				allDone := false
				nextDate := st.LastGenerated.AddDate(0, *st.Months, 0)

				for !allDone {
					if nextDate.After(priorWeek) && nextDate.Before(nextWeek) {
						log.Printf("»»» Task %d Every %d months, next due at %s is within %s - %s",
							st.ID,
							*st.Months,
							nextDate.Format(rfc3339DateLayout),
							priorWeek.Format(rfc3339DateLayout),
							nextWeek.Format(rfc3339DateLayout))

						// Generate a new Task record
						genTask(st, &newTask, nextDate, nextDate.AddDate(0, 0, st.DurationDays))
						numTasks++

						// keep looping, looking at the next date
						nextDate = nextDate.AddDate(0, *st.Months, 0)
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
		fmt.Sprintf("%d New Tasks Generated", *count),
		channel, user_id, "task", 0, true)

	return nil
}

type machineLookup struct {
	MachineUser int `db:"machine_user"`
	SiteUser    int `db:"site_user"`
}

func genTask(st shared.SchedTask, task *shared.Task, startDate time.Time, dueDate time.Time) error {

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

	if st.UserID != 0 {
		userID = st.UserID
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

	// expand out the hashtags of the SM before we do anything else

	hashes := []shared.Hashtag{}
	DB.SQL(`select * from hashtag order by length(name) desc`).QueryStructs(&hashes)

	// print("hashes", hashes)
	desc := task.Descr
	if strings.Contains(desc, "#") {
		// Keep looping through doing text conversions until there is
		// nothing left to expand
		stillLooking := true
		for stillLooking {
			stillLooking = false
			for _, v := range hashes {
				theHash := "#" + v.Name
				if strings.Contains(desc, theHash) {
					desc = strings.Replace(desc, theHash, v.Descr, -1)
					stillLooking = true
				}
			}
		}
	} // contains hashes

	println("HashExpand", task.Descr, "to", desc)
	task.Descr = desc

	DB.InsertInto("task").
		Whitelist("machine_id", "sched_id", "comp_type", "tool_id", "component",
			"descr", "startdate", "due_date", "escalate_date",
			"assigned_to", "assigned_date", "labour_est", "material_est").
		Record(task).
		Returning("id").
		QueryScalar(&task.ID)

	DB.SQL(`update sched_task set last_generated=$2 where id=$1`, st.ID, startDate).Exec()
	lines := strings.Split(desc, "\n")
	println("lines =", lines)

	chekbox := 1
	for _, line := range lines {

		// Does this line have a checkbox ?
		if x := strings.Index(line, "["); x > -1 {
			// println("x = ", x)
			if x2 := strings.Index(line[x+1:], "]"); x2 > -1 {
				x2 += x + 1
				println("We have a chekbox number", chekbox, "defined between ", x, "to", x2)

				check := shared.TaskCheck{
					TaskID: task.ID,
					Seq:    chekbox,
					Descr:  line[x+1 : x2],
					Done:   false,
				}

				DB.InsertInto("task_check").
					Whitelist("task_id", "seq", "descr", "done").
					Record(check).
					Exec()

				fmt.Printf("Added task %v\n", check)

				chekbox++
				continue
			}
		}
	}

	return nil
}

func genTaskOld(st shared.SchedTask, task *shared.Task, startDate time.Time, dueDate time.Time) error {

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

	if st.UserID != 0 {
		userID = st.UserID
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
		theLine := strings.TrimSpace(l)
		if strings.HasPrefix(theLine, "- ") {
			check := shared.TaskCheck{
				TaskID: task.ID,
				Seq:    seq,
				Descr:  theLine[2:],
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
