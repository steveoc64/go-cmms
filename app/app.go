package main

import (
	"github.com/go-humble/router"
	"github.com/steveoc64/formulate"
	// "honnef.co/go/js/dom"
	"itrak-cmms/shared"
)

type MessageFunction func(string, int)

type GlobalSessionData struct {
	Username      string
	UserRole      string
	UserID        int
	CanAllocate   bool
	Channel       int
	Router        *router.Router
	AppFn         map[string]router.Handler
	Subscriptions map[string]MessageFunction
	ID            map[string]int
}

var Session GlobalSessionData

func (s *GlobalSessionData) Navigate(url string) {
	// On navigate, clear out any subscriptions on events
	s.Subscriptions = make(map[string]MessageFunction)
	s.Router.Navigate(url)
	go rpcClient.Call("LoginRPC.Nav", shared.Nav{
		Channel: s.Channel,
		Route:   url,
	}, &url)
}

func (s *GlobalSessionData) Subscribe(msg string, fn MessageFunction) {
	s.Subscriptions[msg] = fn
}

func (s *GlobalSessionData) Reload(context *router.Context) {
	s.Router.Navigate(context.Path)
}

func main() {

	initRouter()
	formulate.Templates(GetTemplate)
	websocketInit()
	initLoginForm()
	showLoginForm()
}
