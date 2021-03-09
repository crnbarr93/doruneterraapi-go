package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
)

const (
	baseURL     string = "https://dd.b.pvp.net/latest/"
	maxKnownSet int    = 4
)

type DDCard struct {
	AssociatedCardRefs    []string `json:"associatedCardRefs" bson:"associatedCardRefs"`
	Region                string   `json:"region" bson:"region"`
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
	Set                   string   `json:"card_set" bson:"card_set"`
}

func (c DDCard) getSetInteger() (int, error) {
	setString := string(c.Set[len(c.Set)-1:])
	return strconv.Atoi(setString)
}

func (c *DDCard) toCard() types.Card {
	cardSet, err := c.getSetInteger()
	if err != nil {
		log.Fatalln("Could not convert set string to integer")
	}

	return types.Card{
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

func getSetData(set int) []types.Card {
	setURL := getSetURL(4)

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

	//SEPARATE TO NEW FUNC
	cards := make([]types.Card, len(ddCards))
	for i, card := range ddCards {
		cards[i] = card.toCard()
	}

	return cards
}