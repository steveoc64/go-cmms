package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
)

type GlobalSessionData struct {
	Username  string
	UserRole  string
	UserID    int
	Channel   int
	Router    *router.Router
	AppFn     map[string]router.Handler
	Subscribe string
	SFn       func(*shared.AsyncMessage)
}

var Session GlobalSessionData

func main() {

	initRouter()
	formulate.Templates(GetTemplate)
	websocketInit()
	initLoginForm()
	showLoginForm()
}
