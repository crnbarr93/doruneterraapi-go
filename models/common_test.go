package models

import (
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
)

func InitializeDatabase() *db.Database {
	database := db.New(config.TestConfig.Database)
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.WaitForConnection()
	return database
}
