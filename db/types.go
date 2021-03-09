package db

import (
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConnectionStatus int

const (
	Disconnected ConnectionStatus = 0
	Connecting   ConnectionStatus = 1
	Connected    ConnectionStatus = 2
)

type Database struct {
	Client           *mongo.Client
	DB               *mongo.Database
	connectionStatus ConnectionStatus
	config           config.DatabaseConfig
	testing          bool
}
