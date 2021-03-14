package models

import (
	"context"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeckModel struct {
	collection *mongo.Collection
}

func NewDeckModel(collection *mongo.Collection) *DeckModel {
	return &DeckModel{
		collection: collection,
	}
}

func (m DeckModel) GetDeck(deckID string) (*types.Deck, error) {
	var deck types.Deck
	result := m.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: deckID}})
	err := result.Decode(&deck)
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (m DeckModel) GetDecksByOwner(ownerName string) ([]*types.Deck, error) {
	var data []*types.Deck
	regex := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: ownerName, Options: "i"}}}
	filter := bson.D{{Key: "ownerUsername", Value: regex}}
	cur, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (m DeckModel) GetDecksByOwnerID(ownerID string) ([]*types.Deck, error) {
	var data []*types.Deck
	filter := bson.D{{Key: "owner", Value: ownerID}}
	cur, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (m DeckModel) SearchDecks(search string) ([]*types.Deck, error) {
	var data []*types.Deck
	regex := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: search, Options: "i"}}}
	searchParams := make([]bson.D, 2)
	searchParams[0] = bson.D{{Key: "ownerUsername", Value: regex}}
	searchParams[1] = bson.D{{Key: "title", Value: regex}}

	filter := bson.D{{Key: "$or", Value: searchParams}}
	cur, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	return data, nil
}
