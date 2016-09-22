package shared

import "time"

type Report struct {
	DateFrom time.Time `db:"date_from"`
	DateTo   time.Time `db:"date_to"`
}

type ReportRPCData struct {
	Channel int
	Report  *Report
}
