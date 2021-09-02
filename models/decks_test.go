package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/deck_encoder"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
)

func init() {
	database := InitializeDatabase()
	database.DropCollection("decks")
	InitModels(database)
}

func TestDeckCardCount(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "2", Quantity: 2}, {CardID: "1", Quantity: 3}},
	}

	received := deck.CardCount()

	assert.Equal(t, 5, received)
}

func TestDeckChampionCount(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received := deck.ChampionCount()

	assert.Equal(t, 3, received)
}

func TestDeckCalculateRegions(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received := deck.CalculateRegions()

	assert.Equal(t, []string{"Freljord", "Noxus"}, received)
}

func TestAllCardsValid(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "1", Quantity: 3}},
	}

	received, err := deck.AllCardsValid()

	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Card with ID 1 does not exist"), err)
}

func TestIsDeckValid(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received, err := deck.IsValid(false, false, false)

	assert.Equal(t, true, received)
	assert.Nil(t, err)

	deck = Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck must include 40 cards to be published"), err)

	deck = Deck{}

	received, err = deck.IsValid(true, false, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck must include at least 1 card"), err)

	deck = Deck{Cards: []CardQuantity{{CardID: "01FR024", Quantity: 9}, {CardID: "01IO012", Quantity: 31}}}
	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck can only contain at most 6 Champion Cards"), err)

	deck = Deck{Cards: []CardQuantity{{CardID: "01IO012", Quantity: 40}}}
	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck can only contain, at most, 3 of any individual card"), err)
}

func TestEncodeDeck(t *testing.T) {
	deck := Deck{
		Cards: []CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	code := deck.Encode()

	assert.Equal(t, "CQBACAIBDAAQCAYMAAAA", code)

	decoded, err := deck_encoder.Decode(code)

	assert.Nil(t, err)
	assert.Equal(t, deck.ToEncodableDeck(), decoded)
}

func saveDeck() (*Deck, error) {
	newDeck := Deck{
		Title:         "Some Test Deck",
		OwnerUsername: "TestUser",
		Owner:         "1",
	}

	return Decks.SaveDeck(newDeck)
}

func TestSaveDeck(t *testing.T) {
	received, err := saveDeck()
	assert.Nil(t, err)
	assert.NotNil(t, received.ID)
}

func TestGetDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := Decks.GetDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received.Title)
}

func TestGetDecksByOwner(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := Decks.GetDecksByOwner(deck.OwnerUsername)

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = Decks.GetDecksByOwner(strings.ToLower(deck.OwnerUsername))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = Decks.GetDecksByOwner("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestGetDecksByOwnerID(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := Decks.GetDecksByOwnerID(deck.Owner)

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = Decks.GetDecksByOwnerID("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestSearchDecks(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := Decks.SearchDecks(strings.ToLower(deck.Title))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = Decks.SearchDecks(strings.ToLower(deck.OwnerUsername))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = Decks.SearchDecks("no deck exists with this search")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestUpdateDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	updatedDeck := deck
	updatedDeck.Title = "New Title"
	received, err := Decks.UpdateDeck(*updatedDeck)

	assert.Nil(t, err)
	assert.Equal(t, updatedDeck.Title, received.Title)
}

func TestDeleteDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	deletedDeck, err := Decks.DeleteDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, false, deletedDeck.Published)
	assert.Equal(t, true, deletedDeck.Deleted)
}

func TestPublishDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	publishedDeck, err := Decks.PublishDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, true, publishedDeck.Published)

	_, err = Decks.DeleteDeck(deck.ID)
	if err != nil {
		panic(err)
	}

	_, err = Decks.PublishDeck(deck.ID)
	assert.NotNil(t, err)
}
