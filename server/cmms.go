package main

import (
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/steveoc64/godev/config"
	"github.com/steveoc64/godev/db"
	"github.com/steveoc64/godev/echocors"
	"github.com/steveoc64/godev/smt"
	"log"
)

var e *echo.Echo

func main() {

	Config := config.LoadConfig()
	cpus := smt.Init()
	fmt.Printf("Go-CMMS running on %d CPU cores\n", cpus)

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

	db.Init(Config.DataSourceName)

	// Add the all important Websocket handler
	Connections = new(ConnectionsList)
	registerRPC()
	e.WebSocket("/ws", webSocket)
	//go pinger()

	// Start the web server
	if Config.Debug {
		log.Printf("... Starting Web Server on port %d", Config.WebPort)
	}
	e.Run(fmt.Sprintf(":%d", Config.WebPort))
}
