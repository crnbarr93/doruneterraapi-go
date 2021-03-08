package db

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client    *mongo.Client
	database  *mongo.Database
	connected bool
)

func connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Config.Database.Address))
	client = mgoClient
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	connect()
	Ping()

	database = client.Database(config.Config.Database.Database)
	if database == nil {
		log.Fatal("No database found")
	}

	log.Println(color.GreenString("Connected to MongoDB!"))
	connected = true
}

func GetClient() *mongo.Client {
	return client
}

func GetDatabase() *mongo.Database {
	return database
}

func GetIsConnected() bool {
	return connected
}

func Ping() error {
	client := GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())

	return err
}
