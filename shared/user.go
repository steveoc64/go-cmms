package shared

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Passwd   string `db:"passwd"`
	Email    string `db:"email"`
	Role     string `db:"role"`
	SMS      string `db:"sms"`
	Sites    []Site `db:"site"`
}

type UserUpdateData struct {
	Channel int
	User    *User
}

type UserUpdate struct {
	Channel  int    `db:"channel"`
	ID       int    `db:"id"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Passwd   string `db:"passwd"`
	Email    string `db:"email"`
	SMS      string `db:"sms"`
}
