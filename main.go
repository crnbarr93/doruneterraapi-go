package main

import (
	"github.com/labstack/echo"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/app"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/utils"
)

func main() {
	database := db.New(config.Config.Database)
	go utils.UpdateAllSets(database)
	e := echo.New()

	app := app.New(e, database)

	err := app.Run(":1323")
	if err != nil {
		panic(err)
	}
}
