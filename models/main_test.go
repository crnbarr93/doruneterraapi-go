package models_test

import (
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

var SavedArchetypes []*models.Archetype
var SavedDecks []*models.Deck

func saveDecks() {
	deckOne := models.Deck{
		Cards:   []models.CardQuantity{{CardID: "01FR024", Quantity: 2}},
		Regions: []string{"Noxus", "Demacia"},
	}
	deckTwo := models.Deck{Regions: []string{"Noxus", "Freljord"}, Cards: []models.CardQuantity{{CardID: "01IO012", Quantity: 3}}}

	savedOne, err := models.Decks.SaveDeck(deckOne)
	if err != nil {
		panic(err)
	}
	savedTwo, err := models.Decks.SaveDeck(deckTwo)
	if err != nil {
		panic(err)
	}

	SavedDecks = append(SavedDecks, savedOne)
	SavedDecks = append(SavedDecks, savedTwo)
}

func saveArchetypes() {
	archetypeOne := models.Archetype{Decks: []string{SavedDecks[0].ID}}
	keyCardOne := models.CardInArchetype{CardID: "01FR024", Quantity: 2, QuantityAppears: []int{0, 2, 0}}
	archetypeTwo := models.Archetype{Decks: []string{SavedDecks[1].ID}, KeyCards: []models.CardInArchetype{keyCardOne}}

	savedOne, err := models.Archetypes.SaveArchetype(archetypeOne)
	if err != nil {
		panic(err)
	}
	savedTwo, err := models.Archetypes.SaveArchetype(archetypeTwo)
	if err != nil {
		panic(err)
	}

	SavedArchetypes = append(SavedArchetypes, savedOne)
	SavedArchetypes = append(SavedArchetypes, savedTwo)
}

func init() {
	database := InitializeDatabase()
	database.DropCollection("decks")
	database.DropCollection("archetypes")
	database.DropCollection("users")
	models.InitModels(database)
	saveDecks()
	saveArchetypes()
}

func InitializeDatabase() *db.Database {
	database := db.New(config.TestConfig.Database)
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.WaitForConnection()
	return database
}
