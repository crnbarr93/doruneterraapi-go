package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var deckModel *DeckModel

func initializeDeckCollection() *mongo.Collection {
	database := InitializeDatabase()
	collection := database.Collection("decks")
	return collection
}

func init() {
	collection := initializeDeckCollection()
	deckModel = NewDeckModel(collection)
}

func TestGetDeck(t *testing.T) {
	expected := "Something"
	received, err := deckModel.GetDeck("abgiPojD4")

	assert.Nil(t, err)
	assert.Equal(t, expected, received.Title)
}

func TestGetDecksByOwner(t *testing.T) {
	expected := "Something"
	received, err := deckModel.GetDecksByOwner("Fyasco")

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	expected = "Something"
	received, err = deckModel.GetDecksByOwner("fyasco")

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = deckModel.GetDecksByOwner("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestGetDecksByOwnerID(t *testing.T) {
	expected := "Something"
	received, err := deckModel.GetDecksByOwnerID("604a0c8858e9489cf763465f")

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = deckModel.GetDecksByOwnerID("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestSearchDecks(t *testing.T) {
	expected := "Something"
	received, err := deckModel.SearchDecks("something")

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	expected = "Something"
	received, err = deckModel.SearchDecks("fyasco")

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = deckModel.SearchDecks("no deck exists with this search")

	assert.Nil(t, err)
	assert.Empty(t, received)
}
