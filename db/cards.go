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
	cards      []types.Card
)

func init() {
	go setup()
}

func setup() {
	<-GetIsConnected()
	assignCollection()
	cacheCards()
}

func assignCollection() {
	collection = database.Collection("cards")
}

func cacheCards() {
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

	cards = data
}

func GetAllCards() []types.Card {
	return cards
}
