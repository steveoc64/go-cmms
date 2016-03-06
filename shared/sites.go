package shared

type Site struct {
	ID             int     `db:"id"`
	Name           string  `db:"name"`
	Address        string  `db:"address"`
	Phone          string  `db:"phone"`
	Fax            string  `db:"fax"`
	Image          string  `db:"image"`
	ParentSite     int     `db:"parent_site"`
	ParentSiteName *string `db:"parent_site_name"`
	StockSite      int     `db:"stock_site"`
	StockSiteName  *string `db:"stock_site_name"`
	Notes          string  `db:"notes"`
}
