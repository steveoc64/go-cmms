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
	ID         int        `db:"id"`
	SiteId     int        `db:"site_id"`
	Name       string     `db:"name"`
	Descr      string     `db:"descr"`
	Make       string     `db:"make"`
	Model      string     `db:"model"`
	Serialnum  string     `db:"serialnum"`
	IsRunning  bool       `db:"is_running"`
	Status     string     `db:"status"`
	Stopped    *time.Time `db:"stopped_at"`
	Started    *time.Time `db:"started_at"`
	Alert      *time.Time `db:"alert_at"`
	Picture    string     `db:"picture"`
	SiteName   *string    `db:"site_name"`
	Span       *string    `db:"span"`
	Notes      string     `db:"notes"`
	Electrical string     `db:"electrical"`
	Hydraulic  string     `db:"hydraulic"`
	Printer    string     `db:"printer"`
	Console    string     `db:"console"`
	Rollbed    string     `db:"rollbed"`
	Uncoiler   string     `db:"uncoiler"`
	Lube       string     `db:"lube"`
	Components []Component
}

func (c *Machine) GetClass(status string) string {
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

func (m *Machine) SVGWidth1() string {
	i := 250 + (len(m.Components) * 50)
	return fmt.Sprintf("%d", i)
}

func (m *Machine) SVGWidth2() string {
	i := 170 + (len(m.Components) * 50)
	return fmt.Sprintf("%d", i)
}

func (m *Machine) SVGX() string {
	i := 250 + (len(m.Components) * 50) - 26
	return fmt.Sprintf("%d", i)
}

type Component struct {
	MachineID   int    `db:"machine_id"`
	Position    int    `db:"position"`
	ZIndex      int    `db:"zindex"`
	ID          int    `db:"id"`
	SiteId      int    `db:"site_id"`
	Name        string `db:"name"`
	Descr       string `db:"descr"`
	Make        string `db:"make"`
	Model       string `db:"model"`
	Qty         int    `db:"qty"`
	StockCode   string `db:"stock_code"`
	Serialnum   string `db:"serialnum"`
	Picture     string `db:"picture"`
	Notes       string `db:"notes"`
	SiteName    string `db:"site_name"`
	MachineName string `db:"machine_name"`
	Status      string `db:"status"`
	IsRunning   bool   `db:"is_running"`
}

func (c *Component) SVGX(index int) string {
	return fmt.Sprintf("%d", 250+(index*50))
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

type Part struct {
	ID                int     `db:"id"`
	Name              string  `db:"name"`
	Descr             string  `db:"descr"`
	StockCode         string  `db:"stock_code"`
	ReorderStocklevel float64 `db:"reorder_stocklevel"`
	ReorderQty        float64 `db:"reorder_qty"`
	LatestPrice       float64 `db:"latest_price"`
	QtyType           string  `db:"qty_type"`
	Picture           string  `db:"picture"`
	Notes             string  `db:"notes"`
}

type PartComponents struct {
	ComponentID int     `db:"component_id"`
	PartID      int     `db:"part_id"`
	Qty         int     `db:"qty"`
	StockCode   *string `db:"stock_code"` // component stock code and name
	Name        *string `db:"name"`
	MachineName string  `db:"machine_name"`
	SiteName    string  `db:"site_name"`
	MachineID   int     `db:"machine_id"`
	SiteID      int     `db:"site_id"`
}

type PartVendors struct {
	VendorId    int     `db:"vendor_id"`
	Name        string  `db:"name"`
	Descr       string  `db:"descr"`
	Address     string  `db:"address"`
	VendorCode  string  `db:"vendor_code"`
	LatestPrice float64 `db:"latest_price"`
}
