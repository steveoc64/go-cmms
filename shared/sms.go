package shared

import (
	"fmt"
	"time"
)

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
	Local      bool      `db:"local"`
}

func (s *SMSTrans) GetNumber() string {
	return fmt.Sprintf("%s %s", s.NumberTo, s.NumberUsed)
}

// const (
// 	dateDisplayFormat = "Mon, Jan 2 2006"
// )

func (s *SMSTrans) GetDateSent() string {
	return s.DateSent.Format(dateDisplayFormat)
}

func (s *SMSTrans) GetStatus() string {
	return fmt.Sprintf("%s %s", s.Status, s.Error)
}
