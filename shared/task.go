package shared

import "time"

type SchedTask struct {
	ID            int        `db:"id"`
	MachineID     int        `db:"machine_id"`
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
