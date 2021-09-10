package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

func TestGetArchetypes(t *testing.T) {
	recv, err := models.Archetypes.GetArchetypes()

	assert.Nil(t, err)
	assert.Equal(t, len(recv), len(SavedArchetypes))
	assert.True(t, len(recv[0].Decks) > 0)
	assert.Equal(t, SavedDecks[0].ID, recv[0].Decks[0].ID)
}

func TestGetArchetypesRaw(t *testing.T) {
	recv, err := models.Archetypes.GetArchetypesRaw()

	assert.Nil(t, err)
	assert.Equal(t, len(recv), len(SavedArchetypes))
	assert.True(t, len(recv[0].Decks) > 0)
	assert.Equal(t, SavedDecks[0].ID, recv[0].Decks[0])
}

func TestGetDeckArchetypes(t *testing.T) {
	recv, err := models.Archetypes.GetDeckArchetypes(SavedDecks[0].ID)

	assert.Nil(t, err)
	assert.NotEmpty(t, recv)
	assert.Equal(t, recv[0].Decks[0].ID, SavedDecks[0].ID)

	recv, err = models.Archetypes.GetDeckArchetypes("not a deck")
	assert.Nil(t, err)
	assert.Empty(t, recv)
}

func TestGetCardArchetypes(t *testing.T) {
	recv, err := models.Archetypes.GetCardArchetypes("01FR024")

	assert.Nil(t, err)
	assert.NotEmpty(t, recv)
	assert.Equal(t, recv[0].Decks[0].ID, SavedDecks[1].ID)

	recv, err = models.Archetypes.GetCardArchetypes("not a card")
	assert.Nil(t, err)
	assert.Empty(t, recv)
}

func TestCalculateKeyCards(t *testing.T) {
	archetype := models.Archetype{
		Decks: []string{SavedDecks[0].ID},
	}

	popArch, err := archetype.PopulateDecks()
	if err != nil {
		panic(err)
	}

	keyCards := popArch.CalculateKeyCards()
	expected := []models.CardInArchetype{
		{CardID: SavedDecks[0].Cards[0].CardID, Quantity: SavedDecks[0].Cards[0].Quantity, QuantityAppears: []int{0, SavedDecks[0].Cards[0].Quantity, 0}},
	}

	assert.NotEmpty(t, keyCards)
	assert.Equal(t, expected, keyCards)
}

func TestCalculateArchetypeRegions(t *testing.T) {
	archetype := models.Archetype{
		Decks: []string{SavedDecks[0].ID, SavedDecks[1].ID},
	}

	popArch, err := archetype.PopulateDecks()
	if err != nil {
		panic(err)
	}

	regions := popArch.CalculateRegions()
	expected := []string{"Demacia", "Freljord", "Noxus"}

	assert.NotEmpty(t, regions)
	assert.Equal(t, expected, regions)
}

func TestSanitizeTitle(t *testing.T) {
	archetype := models.Archetype{
		Title: " New Archetype!&*^@%!() ",
	}

	title := archetype.SanitizeTitle()

	assert.Equal(t, title, "new-archetype")
}

func TestCalculateKeywords(t *testing.T) {
	archetype := models.Archetype{
		Decks: []string{SavedDecks[0].ID, SavedDecks[1].ID},
	}

	popArch, err := archetype.PopulateDecks()
	if err != nil {
		panic(err)
	}

	keywords := popArch.CalculateKeywords()

	assert.NotEmpty(t, keywords)
	for _, keyword := range keywords {
		if keyword.Keyword == "TestKeyword" {
			assert.Equal(t, float32(0.625), keyword.Pct)
			assert.Equal(t, 5, keyword.Quantity)
		} else if keyword.Keyword == "TestKeywordTwo" {
			assert.Equal(t, float32(0.375), keyword.Pct)
			assert.Equal(t, 3, keyword.Quantity)
		}
	}
}

func TestCalculateDetails(t *testing.T) {
	archetype := models.Archetype{
		Decks: []string{SavedDecks[0].ID, SavedDecks[1].ID},
		Title: "New Archetype!",
	}

	err := archetype.CalculateDetails()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, archetype.SanitizedTitle, "new-archetype")
	assert.Equal(t, archetype.Regions, []string{"Demacia", "Freljord", "Noxus"})
	assert.NotEmpty(t, archetype.KeyCards)
	assert.NotEmpty(t, archetype.Keywords)
}

func TestSaveArchetype(t *testing.T) {
	archetypeToSave := models.Archetype{}

	saved, err := models.Archetypes.SaveArchetype(archetypeToSave)

	assert.Nil(t, err)
	assert.NotNil(t, saved.ID)
}
