package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/steveoc64/godev/config"
	"github.com/steveoc64/godev/db"
	"github.com/steveoc64/godev/echocors"
	"github.com/steveoc64/godev/smt"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
)

var e *echo.Echo
var DB *runner.DB

func main() {

	Config := config.LoadConfig()
	cpus := smt.Init()
	fmt.Printf("Go-CMMS running on %d CPU cores\n", cpus)

	// Make sure the SMS stuff is all working before we go too far
	// smsbal, smserr := sms.GetBalance()
	// if smserr != nil {
	// 	log.Fatal("Cannot retrieve SMS account info", smserr.Error())
	// }
	// log.Println("... Remaining SMS Balance =", smsbal)

	// Start up the basic web server
	e = echo.New()
	e.Index("public/index.html")
	e.ServeDir("/", "public/")

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	if Config.Debug {
		e.SetDebug(true)
	}
	echocors.Init(e, Config.Debug)

	// Connect to the database
	DB = db.Init(Config.DataSourceName)

	// Add the all important Websocket handler
	Connections = new(ConnectionsList)
	registerRPC()
	e.WebSocket("/ws", webSocket)

	// Start the web server
	if Config.Debug {
		log.Printf("... Starting Web Server on port %d", Config.WebPort)
	}
	e.Run(fmt.Sprintf(":%d", Config.WebPort))
}
