package models

import (
	"context"
	"log"
	"time"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardModel struct {
	collection *mongo.Collection
	Cards      []types.Card
}

func NewCardModel(collection *mongo.Collection) *CardModel {
	return &CardModel{
		collection: collection,
		Cards:      make([]types.Card, 0),
	}
}

func (m *CardModel) CacheCards() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var data []types.Card

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

func (m *CardModel) GetAll() []types.Card {
	return m.Cards
}

func (m *CardModel) GetCard(cardCode string) *types.Card {
	cards := m.GetAll()
	for _, card := range cards {
		if card.CardCode == cardCode {
			return &card
		}
	}

	return nil
}

func (m *CardModel) GetCardFromDB(cardCode string) (*types.Card, error) {
	var card types.Card
	result := m.collection.FindOne(context.Background(), bson.D{{"_id", cardCode}})
	err := result.Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (m *CardModel) UpdateCards(cards []types.Card) {
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
