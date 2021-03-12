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
	assert.NotEmpty(t, model.Cards)
}

func TestGetAll(t *testing.T) {
	cards := make([]types.Card, 1)
	cards[0] = types.Card{ID: "test"}

	model := CardModel{Cards: cards}

	expected := cards[0].ID
	received := model.GetAll()[0].ID

	assert.Equal(t, expected, received)
}

func TestGetCard(t *testing.T) {
	cards := make([]types.Card, 1)
	cards[0] = types.Card{ID: "test", CardCode: "test"}

	model := CardModel{Cards: cards}

	expected := cards[0].CardCode
	received := model.GetCard(cards[0].CardCode).CardCode

	assert.Equal(t, expected, received)

	nilCard := model.GetCard("doesntexist")

	assert.Nil(t, nilCard)
}

func TestUpdateCards(t *testing.T) {
	cardUpdates := make([]types.Card, 1)
	cardUpdates[0] = types.Card{
		ID:                 "01IO012",
		AssociatedCardRefs: make([]string, 0),
		Region:             "Noxus",
		RegionRef:          "Noxus",
	}

	database := db.New(config.TestConfig.Database)
	database.Connect()
	database.WaitForConnection()

	model := New(database.Collection("cards"))
	model.UpdateCards(cardUpdates)

	updatedCard, err := model.GetCardFromDB(cardUpdates[0].ID)

	assert.Nil(t, err)
	assert.Equal(t, cardUpdates[0].Region, updatedCard.Region)
}
