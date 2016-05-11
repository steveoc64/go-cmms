package shared

import "time"

type SMSTrans struct {
	ID         int       `db:"id"`
	NumberTo   string    `db:"number_to"`
	NumberUsed string    `db:"number_used"`
	UserID     int       `db:"user_id"`
	Message    string    `db:"message"`
	DateSent   time.Time `db:"date_sent"`
	Ref        string    `db:"ref"`
	Status     string    `db:"status"`
	Error      string    `db:"error"`
}
