package utils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

const (
	baseURL     string = "https://dd.b.pvp.net/latest/"
	maxKnownSet int    = 5
)

//Data Dragon Card - Structure received from Endpoint
type DDCard struct {
	AssociatedCardRefs    []string `json:"associatedCardRefs" bson:"associatedCardRefs"`
	Region                string   `json:"region" bson:"region"`
	Regions               []string `json:"regions" bson:"regions"`
	RegionRef             string   `json:"regionRef" bson:"regionRef"`
	Attack                int      `json:"attack" bson:"attack"`
	Cost                  int      `json:"cost" bson:"cost"`
	Health                int      `json:"health" bson:"health"`
	Description           string   `json:"description" bson:"description"`
	DescriptionRaw        string   `json:"descriptionRaw" bson:"descriptionRaw"`
	LevelUpDescription    string   `json:"levelupDescription" bson:"levelupDescription"`
	LevelUpDescriptionRaw string   `json:"levelupDescriptionRaw" bson:"levelupDescriptionRaw"`
	FlavorText            string   `json:"flavorText" bson:"flavorText"`
	ArtistName            string   `json:"artistName" bson:"artistName"`
	Name                  string   `json:"name" bson:"name"`
	CardCode              string   `json:"cardCode,omitempty" bson:"cardCode,omitempty"`
	Keywords              []string `json:"keywords" bson:"keywords"`
	KeywordRefs           []string `json:"keywordRefs" bson:"keywordRefs"`
	SpellSpeed            string   `json:"spellSpeed" bson:"spellSpeed"`
	SpellSpeedRef         string   `json:"spellSpeedRef" bson:"spellSpeedRef"`
	Rarity                string   `json:"rarity" bson:"rarity"`
	RarityRef             string   `json:"rarityRef" bson:"rarityRef"`
	Subtype               string   `json:"subtype" bson:"subtype"`
	Supertype             string   `json:"supertype" bson:"supertype"`
	Type                  string   `json:"type" bson:"type"`
	Collectible           bool     `json:"collectible" bson:"collectible"`
	Set                   string   `json:"set" bson:"set"`
}

func (c DDCard) getSetInteger() (int, error) {
	setString := string(c.Set[len(c.Set)-1:])
	return strconv.Atoi(setString)
}

//Map to Card Type
func (c *DDCard) toCard() models.Card {
	cardSet, err := c.getSetInteger()
	if err != nil {
		log.Fatalln("Could not convert set string to integer")
	}

	return models.Card{
		ID:                    c.CardCode,
		AssociatedCardRefs:    c.AssociatedCardRefs,
		Region:                c.Region,
		RegionRef:             c.RegionRef,
		Attack:                c.Attack,
		Cost:                  c.Cost,
		Health:                c.Health,
		Description:           c.Description,
		DescriptionRaw:        c.DescriptionRaw,
		LevelUpDescription:    c.LevelUpDescription,
		LevelUpDescriptionRaw: c.LevelUpDescriptionRaw,
		FlavorText:            c.FlavorText,
		ArtistName:            c.ArtistName,
		Name:                  c.Name,
		CardCode:              c.CardCode,
		Keywords:              c.Keywords,
		KeywordRefs:           c.KeywordRefs,
		SpellSpeed:            c.SpellSpeed,
		SpellSpeedRef:         c.SpellSpeedRef,
		Rarity:                c.Rarity,
		RarityRef:             c.RarityRef,
		Subtype:               c.Subtype,
		Supertype:             c.Supertype,
		Type:                  c.Type,
		Collectible:           c.Collectible,
		CardSet:               cardSet,
	}
}

func getSetURL(set int) string {
	setString := strconv.Itoa(set)

	return baseURL + "set" + setString + "/en_us/data/set" + setString + "-en_us.json"
}

func getSetData(set int) []models.Card {
	setURL := getSetURL(set)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(setURL)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Could not retrieve set data")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var ddCards []DDCard
	err = json.Unmarshal(body, &ddCards)
	if err != nil {
		log.Fatalln(err)
	}

	cards := make([]models.Card, len(ddCards))
	for i, card := range ddCards {
		cards[i] = card.toCard()
	}

	return cards
}

func hasCardChanged(model *models.CardModel, card models.Card) bool {
	savedCard := model.GetCard(card.CardCode)
	return savedCard == nil || !savedCard.Compare(card)
}

func getCardsToUpdate(model *models.CardModel, cards []models.Card) []models.Card {
	var updatedCards []models.Card

	for _, card := range cards {
		hasUpdated := hasCardChanged(model, card)
		if hasUpdated {
			updatedCards = append(updatedCards, card)
		}
	}

	return updatedCards
}

func updateSet(model *models.CardModel, set int) {
	setData := getSetData(set)
	setUpdates := getCardsToUpdate(model, setData)
	if len(setUpdates) > 0 {
		model.UpdateCards(setUpdates)
	}
	log.Printf("Updated %v cards for Set %v", len(setUpdates), set)
}

func UpdateAllSets(db *db.Database) {
	db.WaitForConnection()
	collection := db.Collection("cards")
	model := models.NewCardModel(collection)
	if len(model.Cards) == 0 {
		model.CacheCards()
	}

	for i := 1; i <= maxKnownSet; i++ {
		go updateSet(model, i)
	}

}
