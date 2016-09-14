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
	X              int     `db:"x"`
	Y              int     `db:"y"`
	AlertsTo       int     `db:"alerts_to"`
	TasksTo        int     `db:"tasks_to"`
	Manager        int     `db:"manager"`
	Highlight      bool    `db:"highlight"`
}

func (s *Site) GetKey() int {
	return s.ID
}

func (s *Site) GetDisplay() string {
	return s.Name

}

type SiteRPCData struct {
	Channel int
	ID      int
	Site    *Site
}

type SiteStatusReport struct {
	Edinburgh string
	Minto     string
	Tomago    string
	Chinderah string
	USA       string
}

func ButtonColor(status string) string {
	switch status {
	case "Running", "":
		return "GreenBtn"
	case "Needs Attention":
		return "YellowBtn"
	case "Stopped":
		return "RedBtn"
	}
	return ""
}

func (s SiteStatusReport) EButton() string {
	return ButtonColor(s.Edinburgh)
}

func (s SiteStatusReport) MButton() string {
	return ButtonColor(s.Minto)
}
func (s SiteStatusReport) TButton() string {
	return ButtonColor(s.Tomago)
}
func (s SiteStatusReport) CButton() string {
	return ButtonColor(s.Chinderah)
}
func (s SiteStatusReport) UButton() string {
	return ButtonColor(s.USA)
}
