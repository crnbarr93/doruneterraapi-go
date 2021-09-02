package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	database := InitializeDatabase()
	database.DropCollection("decks")
	InitModels(database)
}

func TestGetCardNumber(t *testing.T) {
	card := Card{
		ID: "01SI045",
	}

	cardNumber := card.CardNumber()

	assert.Equal(t, 45, cardNumber)
}

func TestCacheCards(t *testing.T) {
	database := InitializeDatabase()
	collection := database.Collection("cards")
	model := NewCardModel(collection)

	err := model.CacheCards()

	assert.Nil(t, err)
	assert.NotEmpty(t, model.Cards)
}

func TestGetAll(t *testing.T) {
	cards := make([]Card, 1)
	cards[0] = Card{ID: "test"}

	model := CardModel{Cards: cards}

	expected := cards[0].ID
	received := model.GetAll()[0].ID

	assert.Equal(t, expected, received)
}

func TestGetCard(t *testing.T) {
	expected := "01IO012"
	received := Cards.GetCard("01IO012").ID

	assert.Equal(t, expected, received)

	nilCard := Cards.GetCard("doesntexist")

	assert.Nil(t, nilCard)
}

func TestUpdateCards(t *testing.T) {
	cardUpdates := make([]Card, 1)
	cardUpdates[0] = Card{
		ID:                 "01FR024",
		AssociatedCardRefs: make([]string, 0),
		Region:             "Freljord",
		RegionRef:          "Freljord",
		Supertype:          "Champion",
		CardSet:            1,
	}

	database := InitializeDatabase()

	model := NewCardModel(database.Collection("cards"))
	model.UpdateCards(cardUpdates)

	updatedCard, err := model.GetCardFromDB(cardUpdates[0].ID)

	assert.Nil(t, err)
	assert.Equal(t, cardUpdates[0].Region, updatedCard.Region)
}
