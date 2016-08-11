package shared

import (
	"fmt"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

type MachineReq struct {
	Channel int
	SiteID  int
}

type Machine struct {
	ID              int         `db:"id"`
	SiteID          int         `db:"site_id"`
	Name            string      `db:"name"`
	Descr           string      `db:"descr"`
	Make            string      `db:"make"`
	Model           string      `db:"model"`
	Serialnum       string      `db:"serialnum"`
	IsRunning       bool        `db:"is_running"`
	Status          string      `db:"status"`
	Stopped         *time.Time  `db:"stopped_at"`
	Started         *time.Time  `db:"started_at"`
	Alert           *time.Time  `db:"alert_at"`
	SiteName        *string     `db:"site_name"`
	Span            *string     `db:"span"`
	Notes           string      `db:"notes"`
	Electrical      string      `db:"electrical"`
	Hydraulic       string      `db:"hydraulic"`
	Pnuematic       string      `db:"pnuematic"`
	Conveyor        string      `db:"conveyor"`
	Printer         string      `db:"printer"`
	Console         string      `db:"console"`
	Rollbed         string      `db:"rollbed"`
	Uncoiler        string      `db:"uncoiler"`
	Lube            string      `db:"lube"`
	Encoder         string      `db:"encoder"`
	StripGuide      string      `db:"strip_guide"`
	AlertsTo        int         `db:"alerts_to"`
	TasksTo         int         `db:"tasks_to"`
	Components      []Component `db:"components"`
	PartClass       int         `db:"part_class"`
	MachineType     int         `db:"machine_type"`
	MachineTypeData MachineType `db:"machine_type_data"`
}

type MachineRPCData struct {
	Channel int
	ID      int
	Machine *Machine
}

func (m *Machine) GetClass(status string) string {
	switch status {
	case "Needs Attention":
		return "needs_attention"
	case "Maintenance Pending":
		return "pending"
	case "Stopped":
		return "stopped"
	default:
		return "running"
	}
}

func (m *Machine) GetStatus(nontool string) string {
	switch nontool {
	case "Electrical":
		return m.Electrical
	case "Hydraulic":
		return m.Hydraulic
	case "Printer":
		return m.Printer
	case "Console":
		return m.Console
	case "Rollbed":
		return m.Rollbed
	case "Uncoiler":
		return m.Uncoiler
	case "Lube":
		return m.Lube
	case "Encoder":
		return m.Encoder
	case "StripGuide":
		return m.StripGuide
	default:
		return "Running"
	}
}

func (m *Machine) SVGWidth1() string {
	i := 300 + (len(m.Components) * 50)
	if i < 400 {
		i = 400
	}
	return fmt.Sprintf("%d", i)
}

func (m *Machine) SVGWidth2() string {
	i := 220 + (len(m.Components) * 50)
	if i < 320 {
		i = 320
	}
	return fmt.Sprintf("%d", i)
}

func (m *Machine) SVGX() string {
	i := 300 + (len(m.Components) * 50)
	if i < 400 {
		i = 400
	}
	return fmt.Sprintf("%d", i-26)
}

func (m *Machine) ConveyorWheel(wheel int) string {
	i := 200 + (len(m.Components) * 50)
	if i < 300 {
		i = 300
	}
	switch wheel {
	case 1:
		return "90"
	case 2:
		return fmt.Sprintf("%d", 90+(i/3))
	case 3:
		return fmt.Sprintf("%d", 90+(2*i/3))
	case 4:
		return fmt.Sprintf("%d", 90+i)
	}
	return "90"
}

func (m *Machine) SVGStatus() string {
	switch m.Status {
	case "Needs Attention":
		return "url(#YellowBtn)"
	case "Maintenance Pending":
		return "pending"
	case "Stopped":
		return "url(#RedBtn)"
	default:
		return "url(#GreenBtn)"
	}
}

func (m *Machine) NonToolBg(status string) string {
	switch status {
	case "Needs Attention":
		return "url(#YellowBtn)"
	case "Maintenance Pending":
		return "pending"
	case "Stopped":
		return "url(#RedBtn)"
	default:
		return "url(#bgrad)"
	}
}

type Component struct {
	MachineID int    `db:"machine_id"`
	MTT       int    `db:"mtt_id"`
	Position  int    `db:"position"`
	ZIndex    int    `db:"zindex"`
	ID        int    `db:"id"`
	SiteId    int    `db:"site_id"`
	Name      string `db:"name"`
	Descr     string `db:"descr"`
	Make      string `db:"make"`
	Model     string `db:"model"`
	Qty       int    `db:"qty"`
	StockCode string `db:"stock_code"`
	Serialnum string `db:"serialnum"`
	// Picture     string `db:"picture"`
	Notes       string `db:"notes"`
	SiteName    string `db:"site_name"`
	MachineName string `db:"machine_name"`
	Status      string `db:"status"`
	IsRunning   bool   `db:"is_running"`
}

func (c *Component) SVGX(index int) string {
	return fmt.Sprintf("%d", 300+(index*50))
}

func (c *Component) SVGName(index int) string {
	return fmt.Sprintf("%d", index+1)
}

func (c *Component) SVGFill() string {
	// print("getting fill for status", c.Status)
	switch c.Status {
	case "Needs Attention":
		return "#fff176"
	case "Maintenance Pending":
		return "#9e9d24"
	case "Stopped":
		return "#ff7043"
	default:
		return "white"
	}
}

func (c *Component) SVGFill2(id int) string {
	// print("getting fill for status", c.Status, id)
	if c.ID == id {
		return "url(#BlueBtn)"
	}

	switch c.Status {
	case "Needs Attention":
		return "#fff176"
	case "Maintenance Pending":
		return "#9e9d24"
	case "Stopped":
		return "#ff7043"
	default:
		return "white"
	}
}

func (c *Component) GetClass() string {
	switch c.Status {
	case "Needs Attention":
		return "needs_attention"
	case "Maintenance Pending":
		return "pending"
	case "Stopped":
		return "stopped"
	default:
		return "running"
	}
}

type MachineType struct {
	ID             int               `db:"id"`
	Name           string            `db:"name"`
	Photo          string            `db:"photo"`
	PhotoPreview   string            `db:"photo_preview"`
	PhotoThumbnail string            `db:"photo_thumbnail"`
	Electrical     bool              `db:"electrical"`
	Hydraulic      bool              `db:"hydraulic"`
	Pnuematic      bool              `db:"pnuematic"`
	Lube           bool              `db:"lube"`
	Printer        bool              `db:"printer"`
	Console        bool              `db:"console"`
	Uncoiler       bool              `db:"uncoiler"`
	Rollbed        bool              `db:"rollbed"`
	Conveyor       bool              `db:"conveyor"`
	Encoder        bool              `db:"encoder"`
	StripGuide     bool              `db:"strip_guide"`
	NumTools       int               `db:"num_tools"`
	Tools          []MachineTypeTool `db:"tools"`
	SelectedTool   int               `db:"selected_tool"`
}

func (m *MachineType) SVGWidth1() string {
	i := 300 + (m.NumTools * 50)
	if i < 400 {
		i = 400
	}
	return fmt.Sprintf("%d", i)
}

func (m *MachineType) SVGWidth2() string {
	i := 220 + (m.NumTools * 50)
	if i < 320 {
		i = 320
	}
	return fmt.Sprintf("%d", i)
}

func (m *MachineType) SVGX() string {
	i := 300 + (m.NumTools * 50)
	if i < 400 {
		i = 400
	}
	return fmt.Sprintf("%d", i-26)
}

func (m *MachineType) SVGStatus() string {
	return "url(#GreenBtn)"
}

func (m *MachineType) NonToolBg() string {
	return "url(#bgrad)"
}

type MachineTypeRPCData struct {
	Channel     int
	ID          int
	MachineType *MachineType
}

type MachineTypeTool struct {
	MachineID   int          `db:"machine_id"`
	MachineType *MachineType `db:"machine_type"`
	ID          int          `db:"id"`
	Position    int          `db:"position"`
	Name        string       `db:"name"`
}

func (c *MachineTypeTool) SVGX(index int) string {
	return fmt.Sprintf("%d", 300+(index*50))
}

func (c *MachineTypeTool) SVGName(index int) string {
	return fmt.Sprintf("%d", index+1)
}

func (c *MachineTypeTool) SVGFill() string {
	return "white"
}

func (c *MachineTypeTool) SVGFill2(id int) string {
	if c.MachineType != nil {
		if c.ID == c.MachineType.SelectedTool {
			return "cyan"
		}
	}
	return "white"
}

func (c *MachineTypeTool) GetClass() string {
	return "running"
}

type MachineTypeToolRPCData struct {
	Channel         int
	MachineID       int
	ID              int
	MachineTypeTool *MachineTypeTool
}

func (m *MachineTypeTool) GetName() string {
	return fmt.Sprintf("%d) %s", m.Position, m.Name)
}
