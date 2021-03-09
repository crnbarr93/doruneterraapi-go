package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
)

var m CardModel

func initializeDatabase() *db.Database {
	database := db.New(config.TestConfig.Database)
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.WaitForConnection()
	return database
}

func TestCacheCards(t *testing.T) {
	database := initializeDatabase()
	collection := database.Collection("cards")
	model := New(collection)

	err := model.cacheCards()

	assert.Nil(t, err)
	assert.NotEmpty(t, model.cards)
}

func TestGetAll(t *testing.T) {
	cards := make([]types.Card, 1)
	cards[0] = types.Card{ID: "test"}

	model := CardModel{cards: cards}

	expected := cards[0].ID
	received := model.GetAll()[0].ID

	assert.Equal(t, expected, received)
}
