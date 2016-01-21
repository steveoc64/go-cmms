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
	fmt.Printf("Yo Ho Ho, here we Go on %d CPU cores\n", cpus)

	e = echo.New()
	e.Index("./dist/index.html")
	e.ServeDir("/", "./dist")

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	if Config.Debug {
		e.SetDebug(true)
	}
	echocors.Init(e, Config.Debug)

	db.Init(Config.DataSourceName)

	// Start the web server
	if Config.Debug {
		log.Printf("... Starting Web Server on port %d", Config.WebPort)
	}
	e.Run(fmt.Sprintf(":%d", Config.WebPort))
}
