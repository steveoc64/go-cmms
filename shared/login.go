package shared

type LoginCredentials struct {
	Username   string
	Password   string
	RememberMe bool
	Channel    int
}

type LoginReply struct {
	Result string
	Token  string
	Menu   []string
}
