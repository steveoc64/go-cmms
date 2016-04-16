package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
)

type GlobalSessionData struct {
	Username string
	UserRole string
	UserID   int
	Channel  int
	Router   *router.Router
	AppFn    map[string]router.Handler
}

var Session GlobalSessionData

func main() {

	initRouter()
	formulate.Templates(GetTemplate)
	websocketInit()
	initLoginForm()
	showLoginForm()
}
