package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type TaskRPC struct{}

// Get all the parts
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

func (t *TaskRPC) Save(req *shared.SchedTaskEditData, id *int) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	// do some data conversions
	switch req.Task.Freq {
	case "One Off":
		req.Task.StartDate = req.Task.OneOffDate
	case "Job Count":
		req.Task.Days = req.Task.Count
	}

	DB.InsertInto("sched_task").
		Whitelist("machine_id", "comp_type", "tool_id",
			"component", "descr", "startdate", "freq", "days", "week",
			"labour_cost", "material_cost").
		Record(req.Task).
		Returning("id").
		QueryScalar(id)

	logger(start, "Task.Insert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			req.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %d %s %d %s %s",
			*id, req.Task.MachineID, req.Task.CompType, req.Task.ToolID, req.Task.Component, req.Task.Descr))

	return nil
}
