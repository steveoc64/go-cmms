package shared

type UserProfile struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Passwd   string `db:"passwd"`
	Email    string `db:"email"`
	Role     string `db:"role"`
	SMS      string `db:"sms"`
	Sites    []Site
}

type UserProfileUpdate struct {
	Channel int    `db:"channel"`
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Passwd  string `db:"passwd"`
	Email   string `db:"email"`
	SMS     string `db:"sms"`
}
