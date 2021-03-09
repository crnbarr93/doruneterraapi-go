package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect(t *testing.T) {
	database := New(config.TestConfig.Database)

	assert.NoError(t, database.Connect())
}

func TestWaitForConnection(t *testing.T) {
	database := New(config.TestConfig.Database)
	database.connectionStatus = Connected

	assert.NoError(t, database.WaitForConnection())

	database.connectionStatus = Disconnected
	assert.Error(t, database.WaitForConnection(), "Connection timed out")
}

func TestCollection(t *testing.T) {
	database := New(config.TestConfig.Database)
	database.Connect()
	database.WaitForConnection()

	expected := "cards"
	received := database.Collection(expected)

	assert.Equal(t, expected, received.Name())
}

func TestDropCollection(t *testing.T) {
	database := New(config.TestConfig.Database)
	database.Connect()
	database.WaitForConnection()

	collectionName := "test"
	database.DB.CreateCollection(context.Background(), collectionName)

	database.testing = false
	assert.Error(t, database.DropCollection(collectionName), "Cannot drop collections in a non-test environment")

	database.testing = true
	assert.NoError(t, database.DropCollection(collectionName))

	collections, err := database.DB.ListCollectionNames(context.Background(), &bson.M{"name": collectionName})
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 0, len(collections))
}
