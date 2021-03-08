package db

import (
	"context"
	"log"
	"time"

	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection *mongo.Collection
)

func assignCollection() {
	if database == nil {
		return
	}

	collection = database.Collection("cards")
}

func GetAllCards() []types.Card {
	if collection == nil {
		assignCollection()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var data []types.Card
	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	if err := cur.All(ctx, &data); err != nil {
		log.Fatal(err)
	}

	return data
}
