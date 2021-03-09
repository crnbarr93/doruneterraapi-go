package models

import (
	"context"
	"time"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardModel struct {
	collection *mongo.Collection
	cards      []types.Card
}

func New(collection *mongo.Collection) *CardModel {
	return &CardModel{
		collection: collection,
		cards:      make([]types.Card, 0),
	}
}

func (m *CardModel) cacheCards() error {
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

	m.cards = data

	return nil
}

func (m *CardModel) GetAll() []types.Card {
	return m.cards
}
