package models

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/deck_encoder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Card struct {
	ID                    string   `json:"_id,omitempty" bson:"_id,omitempty"`
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
	CardSet               int      `json:"card_set" bson:"card_set"`
	CardSubset            int      `json:"card_subset,omitempty" bson:"card_subset,omitempty"`
}

func (c Card) Compare(b Card) bool {
	return cmp.Equal(c, b)
}

func (c Card) CardNumber() int {
	substr := c.ID[len(c.ID)-3:]
	cardNum, err := strconv.Atoi(substr)
	if err != nil {
		return -1
	}
	return cardNum
}

func (c Card) ToEncodableCard() deck_encoder.Card {
	return deck_encoder.Card{
		Faction: deck_encoder.FactionNumberFromName(c.Region),
		Set:     c.CardSet,
		Number:  c.CardNumber(),
	}
}

func (c Card) ToEncodableCardInDeck(quantity int) deck_encoder.CardInDeck {
	return deck_encoder.CardInDeck{
		Card:  c.ToEncodableCard(),
		Count: quantity,
	}
}

type CardModel struct {
	collection *mongo.Collection
	Cards      []Card
}

func InitCardModel(d *db.Database) *CardModel {
	collection := d.Collection("cards")
	m := NewCardModel(collection)
	m.CacheCards()
	return m
}

func NewCardModel(collection *mongo.Collection) *CardModel {
	return &CardModel{
		collection: collection,
		Cards:      make([]Card, 0),
	}
}

func (m *CardModel) CacheCards() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var data []Card

	cur, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())

	if err := cur.All(ctx, &data); err != nil {
		return err
	}

	m.Cards = data

	return nil
}

func (m *CardModel) GetAll() []Card {
	return m.Cards
}

func (m *CardModel) GetCard(cardCode string) *Card {
	cards := m.GetAll()
	for _, card := range cards {
		if card.CardCode == cardCode {
			return &card
		}
	}

	cardInDB, err := m.GetCardFromDB(cardCode)
	if err != nil {
		return nil
	}

	return cardInDB
}

func (m *CardModel) GetCardFromDB(cardCode string) (*Card, error) {
	var card Card
	result := m.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: cardCode}})
	err := result.Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (m *CardModel) UpdateCards(cards []Card) {
	var operations []mongo.WriteModel

	for _, card := range cards {
		filter := bson.M{"_id": card.ID}
		update := bson.M{"$set": card}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)
		operations = append(operations, operation)
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	bulkContext := context.Background()
	_, err := m.collection.BulkWrite(bulkContext, operations, &bulkOption)
	if err != nil {
		log.Fatal(err)
	}
	go m.CacheCards()
}
