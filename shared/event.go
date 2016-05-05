package shared

import (
	"fmt"
	"time"
)

type RaiseIssue struct {
	Channel   int
	MachineID int
	CompID    int
	IsTool    bool
	Machine   *Machine
	Component *Component
	NonTool   string
	Descr     string
}

type Event struct {
	ID           int        `db:"id"`
	SiteID       int        `db:"site_id"`
	Type         string     `db:"type"`
	MachineID    int        `db:"machine_id"`
	MachineName  string     `db:"machine_name"`
	SiteName     string     `db:"site_name"`
	ToolID       int        `db:"tool_id"`
	ToolType     string     `db:"tool_type"`
	Priority     int        `db:"priority"`
	Status       string     `db:"status"`
	StartDate    time.Time  `db:"startdate"`
	DisplayDate  string     `db:"display_date"`
	CreatedBy    int        `db:"created_by"`
	AllocatedBy  int        `db:"allocated_by"`
	Username     string     `db:"username"`
	AllocatedTo  int        `db:"allocated_to"`
	Completed    *time.Time `db:"completed"`
	LabourCost   float64    `db:"labour_cost"`
	MaterialCost float64    `db:"material_cost"`
	OtherCost    float64    `db:"other_cost"`
	Notes        string     `db:"notes"`
	ParentEvent  int        `db:"parent_event"`
	AssignedTo   []string   `db:"assigned_to"`
}

type AssignEvent struct {
	Channel     int
	Event       *Event
	AssignTo    int
	StartDate   *time.Time
	DueDate     *time.Time
	LabourEst   float64
	MaterialEst float64
	SiteName    string
	MachineName string
	ToolType    string
	DisplayDate string
	Username    string
	Notes       string
}

type EventUpdateData struct {
	Channel int
	Event   *Event
}

const (
	datetimeDisplayFormat = "Mon, Jan 2 2006 15:04:05"
)

func (e *Event) GetStartDate() string {
	return e.StartDate.Format(datetimeDisplayFormat)
}

func (e *Event) GetStatus() string {
	switch e.Status {
	case "":
		return "Pending"
	case "Assigned":
		status := "Assigned To: "
		for i, j := range e.AssignedTo {
			if i > 0 {
				status += ", "
			}
			status += j
		}
		return status
	default:
		return e.Status
	}
}

func (e *Event) GetCompleted() string {
	if e.Completed != nil {
		return fmt.Sprintf("%s", e.Completed.Format(datetimeDisplayFormat))
	} else {
		return ""
	}
}
