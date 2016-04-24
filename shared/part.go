package shared

import (
	"fmt"
	"time"
)

type PartClass struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Descr string `db:"descr"`
}

type PartClassUpdateData struct {
	Channel   int
	PartClass *PartClass
}

type PartListReq struct {
	Channel int
	Class   int
}

type Part struct {
	ID                   int        `db:"id"`
	Class                int        `db:"class"`
	Name                 string     `db:"name"`
	Descr                string     `db:"descr"`
	StockCode            string     `db:"stock_code"`
	ReorderStocklevel    float64    `db:"reorder_stocklevel"`
	ReorderQty           float64    `db:"reorder_qty"`
	LatestPrice          float64    `db:"latest_price"`
	LastPriceDate        *time.Time `db:"last_price_date"`
	LastPriceDateDisplay string     `db:"last_price_date_display"`
	QtyType              string     `db:"qty_type"`
	Picture              string     `db:"picture"`
	Notes                string     `db:"notes"`
}

type PartUpdateData struct {
	Channel int
	Part    *Part
}

func (p *Part) ReorderDetails() string {
	return fmt.Sprintf("%g / %g", p.ReorderStocklevel, p.ReorderQty)
}

func (p *Part) DisplayPrice() string {
	return fmt.Sprintf("$ %8.2f", p.LatestPrice)
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
