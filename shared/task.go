package shared

import (
	"fmt"
	"time"
)

type SchedTask struct {
	ID            int        `db:"id"`
	MachineID     int        `db:"machine_id"`
	CompType      string     `db:"comp_type"`
	ToolID        int        `db:"tool_id"`
	Component     string     `db:"component"`
	Descr         string     `db:"descr"`
	StartDate     *time.Time `db:"startdate"`
	OneOffDate    *time.Time `db:"oneoffdate"`
	Freq          string     `db:"freq"`
	ParentTask    *int       `db:"parent_task"`
	Days          *int       `db:"days"`
	Count         *int       `db:"count"`
	Week          *int       `db:"week"`
	WeekDay       *int       `db:"weekday"`
	UserID        int        `db:"user_id"`
	DurationDays  int        `db:"duration_days"`
	LabourCost    float64    `db:"labour_cost"`
	MaterialCost  float64    `db:"material_cost"`
	OtherCostDesc *[]string  `db:"other_cost_desc"`
	OtherCost     *[]float64 `db:"other_cost"`
	LastGenerated *time.Time `db:"last_generated"`
	Paused        bool       `db:"paused"`
	PartsRequired []PartReq  `db:"parts_required"`
}

type SchedTaskUpdateData struct {
	Channel   int
	SchedTask *SchedTask
}

// NOTE - the times passed in the Format function are REFERENCE dates for the
// formatter, these are not dates of your choosing. Please read docs for time.Parse()
func (t *SchedTask) ShowFrequency() string {
	switch t.Freq {
	case "Monthly":
		return fmt.Sprintf("Monthly - Week %d", *t.Week)
	case "Yearly":
		return fmt.Sprintf("Yearly - %s", t.StartDate.Format("Mon, Jan 2 2006"))
	case "Every N Days":
		return fmt.Sprintf("Every %d Days", *t.Days)
	case "One Off":
		return fmt.Sprintf("Once at - %s", t.OneOffDate.Format("Mon, Jan 2 2006"))
	case "Job Count":
		return fmt.Sprintf("Job Count > %d", *t.Count)
	}
	return fmt.Sprintf("%s %d", t.Freq, t.Week)
}

func (t *SchedTask) ShowComponent(m Machine) string {
	switch t.CompType {
	case "A":
		return "General Maint."
	case "T":
		// for _, c := range m.Components {
		// 	if c.ID == t.ToolID {
		// 		return c.Name
		// 	}
		// }
		return fmt.Sprintf("Tool %d", t.ToolID)
	case "C":
		return t.Component
	}
	return fmt.Sprintf("%s:%d:%s", t.CompType, t.ToolID, t.Component)
}

func (t *SchedTask) ShowPaused() string {
	if t.Paused {
		return "PAUSED"
	}
	return ""
}

type SchedTaskPart struct {
	TaskID int     `db:"task_id"`
	PartID int     `db:"part_id"`
	Qty    float64 `db:"qty"`
	Notes  string  `db:"notes"`
}

type TaskPart struct {
	TaskID    int     `db:"task_id"`
	PartID    int     `db:"part_id"`
	PartName  string  `db:"part_name"`
	StockCode string  `db:"stock_code"`
	Qty       float64 `db:"qty"`
	QtyUsed   float64 `db:"qty_used"`
	QtyType   string  `db:"qty_type"`
	Notes     string  `db:"notes"`
}

type TaskCheck struct {
	TaskID   int        `db:"task_id"`
	Seq      int        `db:"seq"`
	Descr    string     `db:"descr"`
	Done     bool       `db:"done"`
	DoneDate *time.Time `db:"done_date"`
}

type TaskCheckUpdate struct {
	Channel   int
	TaskCheck *TaskCheck
}

func (t *TaskCheck) ShowDoneDate() string {
	if t.DoneDate == nil {
		return ""
	}
	return t.DoneDate.Format(dateDisplayFormat)
}

type PartReq struct {
	PartID    int      `db:"part_id"`
	StockCode string   `db:"stock_code"`
	Name      string   `db:"name"`
	QtyType   string   `db:"qty_type"`
	QtyPtr    *float64 `db:"qty"`
	NotesPtr  *string  `db:"notes"`
	Qty       float64  `db:"qty_deref"`
	Notes     string   `db:"notes_deref"`
}

type PartReqEdit struct {
	Channel int
	Task    SchedTask
	Part    *PartReq
}

type Task struct {
	ID                int         `db:"id"`
	MachineID         int         `db:"machine_id"`
	MachineName       string      `db:"machine_name"`
	SiteID            int         `db:"site_id"`
	SiteName          string      `db:"site_name"`
	SchedID           int         `db:"sched_id"`
	EventID           int         `db:"event_id"`
	CompType          string      `db:"comp_type"`
	ToolID            int         `db:"tool_id"`
	Component         string      `db:"component"`
	Descr             string      `db:"descr"`
	CreatedDate       time.Time   `db:"created_date"`
	StartDate         *time.Time  `db:"startdate"`
	DisplayStartDate  string      `db:"display_startdate"`
	Log               string      `db:"log"`
	DueDate           *time.Time  `db:"due_date"`
	DisplayDueDate    string      `db:"display_duedate"`
	EscalateDate      *time.Time  `db:"escalate_date"`
	AssignedBy        *int        `db:"assigned_by"`
	AssignedTo        *int        `db:"assigned_to"`
	Username          *string     `db:"username"`
	DisplayUsername   string      `db:"display_username"`
	AssignedDate      *time.Time  `db:"assigned_date"`
	CompletedDate     *time.Time  `db:"completed_date"`
	HasIssue          bool        `db:"has_issue"`
	IssueResolvedDate *time.Time  `db:"issue_resolved_date"`
	LabourEst         float64     `db:"labour_est"`
	LabourHrs         float64     `db:"labour_hrs"`
	MaterialEst       float64     `db:"material_est"`
	LabourCost        float64     `db:"labour_cost"`
	MaterialCost      float64     `db:"material_cost"`
	OtherCostDesc     *[]string   `db:"other_cost_desc"`
	OtherCost         *[]float64  `db:"other_cost"`
	Parts             []TaskPart  `db:"parts"`
	Checks            []TaskCheck `db:"checks"`
	AllDone           bool        `db:"all_done"`
}

type TaskUpdateData struct {
	Channel int
	Task    *Task
}

const (
	dateDisplayFormat = "Mon, Jan 2 2006"
)

func (t *Task) GetID() string {
	return fmt.Sprintf("%06d", t.ID)
}

func (t *Task) GetStartDate() string {
	if t.StartDate == nil {
		return ""
	}
	return t.StartDate.Format(dateDisplayFormat)
}

func (t *Task) GetDueDate() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(dateDisplayFormat)
}

func (t *Task) GetCompletedDate() string {
	if t.CompletedDate == nil {
		return ""
	}
	return t.CompletedDate.Format(dateDisplayFormat)
}

func (t *Task) DurationDays() string {
	d := t.DueDate.Sub(*t.StartDate)
	days := 1 + (d / (time.Hour * 24))
	if days == 1 {
		return "1 Day"
	}
	return fmt.Sprintf("%d Days", days)
}

func (t *Task) DurationHrs() string {
	d := t.DueDate.Sub(*t.StartDate)
	hrs := d / (time.Hour)
	if hrs == 1 {
		return "1 Hour"
	}
	return fmt.Sprintf("%d Hours", hrs)
}

type Hashtag struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Descr string `db:"descr"`
}

func (h *Hashtag) HashName() string {
	return "#" + h.Name
}

type HashtagUpdateData struct {
	Channel int
	Hashtag *Hashtag
}
