package models_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/deck_encoder"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
)

func TestDeckCardCount(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "2", Quantity: 2}, {CardID: "1", Quantity: 3}},
	}

	received := deck.CardCount()

	assert.Equal(t, 5, received)
}

func TestDeckChampionCount(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received := deck.ChampionCount()

	assert.Equal(t, 3, received)
}

func TestDeckCalculateRegions(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received := deck.CalculateRegions()

	assert.Equal(t, []string{"Freljord", "Noxus"}, received)
}

func TestAllCardsValid(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "1", Quantity: 3}},
	}

	received, err := deck.AllCardsValid()

	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Card with ID 1 does not exist"), err)
}

func TestIsDeckValid(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received, err := deck.IsValid(false, false, false)

	assert.Equal(t, true, received)
	assert.Nil(t, err)

	deck = models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck must include 40 cards to be published"), err)

	deck = models.Deck{}

	received, err = deck.IsValid(true, false, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck must include at least 1 card"), err)

	deck = models.Deck{Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 9}, {CardID: "01IO012", Quantity: 31}}}
	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck can only contain at most 6 Champion Cards"), err)

	deck = models.Deck{Cards: []models.CardQuantity{{CardID: "01IO012", Quantity: 40}}}
	received, err = deck.IsValid(false, true, false)
	assert.Equal(t, false, received)
	assert.Equal(t, types.InvalidDeckErrorFromString("Deck can only contain, at most, 3 of any individual card"), err)
}

func TestEncodeDeck(t *testing.T) {
	deck := models.Deck{
		Cards: []models.CardQuantity{{CardID: "01FR024", Quantity: 3}, {CardID: "01IO012", Quantity: 3}},
	}

	code := deck.Encode()

	assert.Equal(t, "CQBACAIBDAAQCAYMAAAA", code)

	decoded, err := deck_encoder.Decode(code)

	assert.Nil(t, err)
	assert.Equal(t, deck.ToEncodableDeck(), decoded)
}

func saveDeck() (*models.Deck, error) {
	newDeck := models.Deck{
		Title:         "Some Test Deck",
		OwnerUsername: "TestUser",
		Owner:         "1",
	}

	return models.Decks.SaveDeck(newDeck)
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
	received, err := models.Decks.GetDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received.Title)
}

func TestGetDecksByOwner(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := models.Decks.GetDecksByOwner(deck.OwnerUsername)

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = models.Decks.GetDecksByOwner(strings.ToLower(deck.OwnerUsername))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = models.Decks.GetDecksByOwner("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestGetDecksByOwnerID(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := models.Decks.GetDecksByOwnerID(deck.Owner)

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = models.Decks.GetDecksByOwnerID("userdoesntexist")

	assert.Nil(t, err)
	assert.Empty(t, received)
}

func TestSearchDecks(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	expected := deck.Title
	received, err := models.Decks.SearchDecks(strings.ToLower(deck.Title))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = models.Decks.SearchDecks(strings.ToLower(deck.OwnerUsername))

	assert.Nil(t, err)
	assert.Greater(t, len(received), 0)
	assert.Equal(t, expected, received[0].Title)

	received, err = models.Decks.SearchDecks("no deck exists with this search")

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
	received, err := models.Decks.UpdateDeck(*updatedDeck)

	assert.Nil(t, err)
	assert.Equal(t, updatedDeck.Title, received.Title)
}

func TestDeleteDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	deletedDeck, err := models.Decks.DeleteDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, false, deletedDeck.Published)
	assert.Equal(t, true, deletedDeck.Deleted)
}

func TestPublishDeck(t *testing.T) {
	deck, err := saveDeck()
	if err != nil {
		panic(err)
	}

	publishedDeck, err := models.Decks.PublishDeck(deck.ID)

	assert.Nil(t, err)
	assert.Equal(t, true, publishedDeck.Published)

	_, err = models.Decks.DeleteDeck(deck.ID)
	if err != nil {
		panic(err)
	}

	_, err = models.Decks.PublishDeck(deck.ID)
	assert.NotNil(t, err)
}

func TestSearchPopularDecks(t *testing.T) {
	deckOne := models.Deck{
		PageViews:     1000,
		DatePublished: time.Now().Add(time.Hour * -2),
		Regions:       []string{"Demacia"},
		Published:     true,
		Cards:         []models.CardQuantity{{CardID: "test", Quantity: 1}},
	}
	deckTwo := models.Deck{
		PageViews:     1000,
		DatePublished: time.Now().Add(time.Hour * -1),
		Regions:       []string{"Noxus"},
		Published:     true,
	}
	deckThree := models.Deck{
		PageViews:     5,
		DatePublished: time.Now().Add(time.Hour * -3),
		Regions:       []string{"Demacia", "Noxus"},
		Published:     true,
	}
	decks := []models.Deck{deckOne, deckTwo, deckThree}
	var savedDecks []models.Deck

	for _, deck := range decks {
		saved, err := models.Decks.SaveDeck(deck)
		if err != nil {
			panic(err)
		}

		savedDecks = append(savedDecks, *saved)
	}

	baseQuery := models.SearchPopularDecksQuery{}
	resp, err := models.Decks.GetPopularDecks(baseQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(savedDecks), len(resp))
	assert.Equal(t, savedDecks[0].ID, resp[1].ID)
	assert.Equal(t, savedDecks[1].ID, resp[0].ID)

	limitQuery := models.SearchPopularDecksQuery{Limit: 2}
	resp, err = models.Decks.GetPopularDecks(limitQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, limitQuery.Limit, len(resp))
	assert.Equal(t, savedDecks[0].ID, resp[1].ID)
	assert.Equal(t, savedDecks[1].ID, resp[0].ID)

	paginatedQuery := models.SearchPopularDecksQuery{Limit: 1, Page: 2}
	resp, err = models.Decks.GetPopularDecks(paginatedQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, paginatedQuery.Limit, len(resp))
	assert.Equal(t, savedDecks[2].ID, resp[0].ID)

	searchCardQuery := models.SearchPopularDecksQuery{Cards: []string{"test"}}
	resp, err = models.Decks.GetPopularDecks(searchCardQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(resp))
	assert.Equal(t, savedDecks[0].ID, resp[0].ID)

	searchRegionQuery := models.SearchPopularDecksQuery{Regions: []string{"Noxus"}}
	resp, err = models.Decks.GetPopularDecks(searchRegionQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 2, len(resp))
	assert.Equal(t, savedDecks[1].ID, resp[0].ID)
	assert.Equal(t, savedDecks[2].ID, resp[1].ID)

	searchMultiRegionQuery := models.SearchPopularDecksQuery{Regions: []string{"Noxus", "Demacia"}}
	resp, err = models.Decks.GetPopularDecks(searchMultiRegionQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(resp))
	assert.Equal(t, savedDecks[2].ID, resp[0].ID)

	sortedQuery := models.SearchPopularDecksQuery{Sorting: "pageViews", SortAsc: -1}
	resp, err = models.Decks.GetPopularDecks(sortedQuery)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(savedDecks), len(resp))
	assert.Equal(t, savedDecks[2].ID, resp[0].ID)
	assert.True(t, resp[0].PageViews <= resp[1].PageViews)
	assert.True(t, resp[1].PageViews <= resp[2].PageViews)
}
