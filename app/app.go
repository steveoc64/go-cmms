package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	"github.com/steveoc64/go-cmms/shared"
)

type MessageFunction func(*shared.AsyncMessage)

type GlobalSessionData struct {
	Username      string
	UserRole      string
	UserID        int
	Channel       int
	Router        *router.Router
	AppFn         map[string]router.Handler
	Subscriptions map[string]MessageFunction
}

var Session GlobalSessionData

func (s *GlobalSessionData) Navigate(url string) {
	s.Subscriptions = make(map[string]MessageFunction)
	s.Router.Navigate(url)
}

func (s *GlobalSessionData) Subscribe(msg string, fn MessageFunction) {
	s.Subscriptions[msg] = fn
}

func main() {

	initRouter()
	formulate.Templates(GetTemplate)
	websocketInit()
	initLoginForm()
	showLoginForm()
}
