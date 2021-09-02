package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

func TestGetSetURL(t *testing.T) {
	expected := "https://dd.b.pvp.net/latest/set4/en_us/data/set4-en_us.json"
	received := getSetURL(4)

	assert.Equal(t, expected, received)
}

func TestGetSetInteger(t *testing.T) {
	card := DDCard{Set: "set1"}

	expected := 1
	received, err := card.getSetInteger()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestGetSetData(t *testing.T) {
	data := getSetData(1)

	expected := 1
	received := data[0].CardSet

	assert.Equal(t, expected, received)
}

func TestHasCardChanged(t *testing.T) {
	cards := make([]models.Card, 1)
	cards[0] = models.Card{ID: "test", CardCode: "test"}

	model := models.CardModel{Cards: cards}

	shouldNotHaveChanged := hasCardChanged(&model, cards[0])

	assert.False(t, shouldNotHaveChanged)

	changedCard := models.Card{ID: "test2", CardCode: "test"}

	shouldHaveChanged := hasCardChanged(&model, changedCard)
	assert.True(t, shouldHaveChanged)
}

func TestUpdateSetData(t *testing.T) {
	savedCards := make([]models.Card, 1)
	cardUpdates := make([]models.Card, 1)
	savedCards[0] = models.Card{ID: "test", CardCode: "test"}
	cardUpdates[0] = models.Card{ID: "test2", CardCode: "test"}

	model := models.CardModel{Cards: savedCards}

	updatedCards := getCardsToUpdate(&model, cardUpdates)

	assert.Equal(t, cardUpdates[0], updatedCards[0])
}
