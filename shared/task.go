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
	Days          int        `db:"days"`
	Count         int        `db:"count"`
	Week          int        `db:"week"`
	DurationDays  int        `db:"duration_days"`
	LabourCost    float64    `db:"labour_cost"`
	MaterialCost  float64    `db:"material_cost"`
	OtherCostDesc *[]string  `db:"other_cost_desc"`
	OtherCost     *[]float64 `db:"other_cost"`
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
		return fmt.Sprintf("Monthly - Week %d", t.Week)
	case "Yearly":
		return fmt.Sprintf("Yearly - %s", t.StartDate.Format("Mon, Jan 2 2006"))
	case "Every N Days":
		return fmt.Sprintf("Every %d Days", t.Days)
	case "One Off":
		return fmt.Sprintf("Once at - %s", t.StartDate.Format("Mon, Jan 2 2006"))
	case "Job Count":
		return fmt.Sprintf("Job Count > %d", t.Days)
	}
	return fmt.Sprintf("%s %d", t.Freq, t.Week)
}

func (t *SchedTask) ShowComponent(m Machine) string {
	switch t.CompType {
	case "A":
		return "General Maint."
	case "T":
		for _, c := range m.Components {
			if c.ID == t.ToolID {
				return c.Name
			}
		}
		return fmt.Sprintf("Tool %d", t.ToolID)
	case "C":
		return t.Component
	}
	return fmt.Sprintf("%s:%d:%s", t.CompType, t.ToolID, t.Component)
}

// type SchedTaskEdit struct {
// 	ID           int       `db:"id"`
// 	MachineID    int       `db:"machine_id"`
// 	CompType     string    `db:"comp_type"`
// 	ToolID       int       `db:"tool_id"`
// 	Component    string    `db:"component"`
// 	Descr        string    `db:"descr"`
// 	StartDate    time.Time `db:"startdate"`
// 	OneOffDate   time.Time `db:"oneoffdate"`
// 	Freq         string    `db:"freq"`
// 	ParentTask   int       `db:"parent_task"`
// 	Days         int       `db:"days"`
// 	Count        int       `db:"count"`
// 	Week         int       `db:"week"`
// 	LabourCost   float64   `db:"labour_cost"`
// 	MaterialCost float64   `db:"material_cost"`
// }

// type SchedTaskEditData struct {
// 	Channel int
// 	Task    SchedTaskEdit
// }
