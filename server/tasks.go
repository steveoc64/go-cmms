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

func (t *TaskRPC) Insert(data *shared.SchedTaskUpdateData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.InsertInto("sched_task").
		Whitelist("machine_id", "comp_type", "tool_id",
			"component", "descr", "startdate", "freq", "days", "week",
			"labour_cost", "material_cost").
		Record(data.SchedTask).
		Returning("id").
		QueryScalar(id)

	logger(start, "Task.Insert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %d %s %d %s %s",
			*id, data.SchedTask.MachineID, data.SchedTask.CompType,
			data.SchedTask.ToolID, data.SchedTask.Component, data.SchedTask.Descr))

	return nil
}
