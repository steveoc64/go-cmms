package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
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
