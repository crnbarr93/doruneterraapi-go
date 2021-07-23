package db

import (
	"context"
	"errors"
	"log"
	"time"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func New(config config.DatabaseConfig) *Database {
	return &Database{
		config:           config,
		connectionStatus: Disconnected,
		testing:          config.Testing,
	}
}

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.connectionStatus = Connecting
	mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(d.config.Address))
	if err != nil {
		return err
	}
	d.Client = mgoClient

	connErr := d.checkDatabaseConnection()
	if connErr != nil {
		return connErr
	}

	database := mgoClient.Database(d.config.Database)
	if database == nil {
		return errors.New("Could not find configured database")
	}
	d.DB = database

	log.Println("Connected to MongoDB!")
	d.connectionStatus = Connected

	return nil
}

func (d Database) checkDatabaseConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := d.Client.Ping(ctx, readpref.Primary())

	return err
}

func (d Database) WaitForConnection() error {
loop:
	for timeout := time.After(time.Second * 10); ; {
		select {
		case <-timeout:
			return errors.New("Connection timed out")
		default:
		}
		if d.connectionStatus == Connected {
			break loop
		}
	}

	return nil
}

func (d Database) Collection(collection string) *mongo.Collection {
	dbCollection := d.DB.Collection(collection)

	return dbCollection
}

func (d Database) CreateSession() (mongo.Session, error) {
	return d.Client.StartSession()
}

func (d Database) DropCollection(collection string) error {
	println(d.config.Address, d.config.Database)
	if !d.testing {
		return errors.New("Cannot drop collections in a non-test environment")
	}

	err := d.Collection(collection).Drop(context.Background())
	return err
}
