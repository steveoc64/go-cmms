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
	Freq          string     `db:"freq"`
	ParentTask    int        `db:"parent_task"`
	Days          int        `db:"days"`
	Week          int        `db:"week"`
	DurationDays  int        `db:"duration_days"`
	LabourCost    float64    `db:"labour_cost"`
	MaterialCost  float64    `db:"material_cost"`
	OtherCostDesc []string   `db:"other_cost_desc"`
	OtherCost     []float64  `db:"other_cost"`
}

type SchedTaskEdit struct {
	ID           int       `db:"id"`
	MachineID    int       `db:"machine_id"`
	CompType     string    `db:"comp_type"`
	ToolID       int       `db:"tool_id"`
	Component    string    `db:"component"`
	Descr        string    `db:"descr"`
	StartDate    time.Time `db:"startdate"`
	OneOffDate   time.Time `db:"oneoffdate"`
	Freq         string    `db:"freq"`
	ParentTask   int       `db:"parent_task"`
	Days         int       `db:"days"`
	Count        int       `db:"count"`
	Week         int       `db:"week"`
	LabourCost   float64   `db:"labour_cost"`
	MaterialCost float64   `db:"material_cost"`
}

type SchedTaskEditData struct {
	Channel int
	Task    SchedTaskEdit
}

func (t *SchedTask) ShowFrequency() string {
	return fmt.Sprintf("%s %d", t.Freq, t.Week)
}
