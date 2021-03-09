package main

import (
	"github.com/labstack/echo"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/app"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
)

func main() {
	database := db.New(config.Config.Database)
	e := echo.New()

	app := app.New(e, database)

	err := app.Run(":1323")
	if err != nil {
		panic(err)
	}
}
