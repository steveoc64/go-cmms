package shared

type LoginCredentials struct {
	Username   string
	Password   string
	RememberMe bool
	Channel    int
}

type LoginReply struct {
	Result string
	Role   string
	Site   string
	Home   string
	Menu   []string
}
