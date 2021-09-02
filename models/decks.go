package models

import (
	"context"
	"fmt"
	"time"

	"github.com/teris-io/shortid"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/deck_encoder"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeckModel struct {
	collection *mongo.Collection
}

type CardQuantity struct {
	CardID   string `json:"cardId" bson:"cardId"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

type DeckBadge struct {
	Color string `json:"color" bson:"color"`
	Text  string `json:"text" bson:"text"`
}

type Deck struct {
	ID             string         `json:"_id,omitempty" bson:"_id,omitempty"`
	Cards          []CardQuantity `json:"cards" bson:"cards"`
	DeckCode       string         `json:"deckCode" bson:"deckCode"`
	Title          string         `json:"title" bson:"title"`
	OwnerUsername  string         `json:"ownerUsername" bson:"ownerUsername"`
	Owner          string         `json:"owner" bson:"owner"`
	DateCreated    time.Time      `json:"dateCreated" bson:"dateCreated"`
	DateUpdated    time.Time      `json:"dateUpdated" bson:"dateUpdated"`
	DatePublished  time.Time      `json:"datePublished,omitempty" bson:"datePublished,omitempty"`
	DateDeleted    time.Time      `json:"dateDeleted,omitempty" bson:"dateDeleted,omitempty"`
	PageViews      int            `json:"pageviews" bson:"pageviews"`
	Guide          string         `json:"guide" bson:"guide"`
	Published      bool           `json:"published" bson:"published"`
	Deleted        bool           `json:"deleted" bson:"deleted"`
	Regions        []string       `json:"regions" bson:"regions"`
	FeaturedPlayer string         `json:"featuredPlayer,omitempty" bson:"featuredPlayer,omitempty"`
	Badge          DeckBadge      `json:"deckBadge,omitempty" bson:"deckBadge,omitempty"`
	Sandbox        bool           `json:"sandbox" bson:"sandbox"`
}

func (d Deck) DeckID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(d.ID)
}

func containsRegion(regions []string, region string) bool {
	for _, reg := range regions {
		if reg == region {
			return true
		}
	}

	return false
}

func (d Deck) CalculateRegions() []string {
	regions := []string{}

	for _, cardQuant := range d.Cards {
		card := Cards.GetCard(cardQuant.CardID)

		if !containsRegion(regions, card.Region) {
			regions = append(regions, card.Region)
		}
	}

	return regions
}

func (d Deck) AllCardsValid() (bool, error) {
	for _, cardQuant := range d.Cards {
		card := Cards.GetCard(cardQuant.CardID)

		if card == nil {
			return false, types.InvalidDeckErrorFromString(fmt.Sprintf("Card with ID %s does not exist", cardQuant.CardID))
		}
	}

	return true, nil
}

func (d Deck) CardCount() int {
	count := 0

	for _, cardQuant := range d.Cards {
		count += cardQuant.Quantity
	}

	return count
}

func (d Deck) ChampionCount() int {
	count := 0

	for _, cardQuant := range d.Cards {
		card := Cards.GetCard(cardQuant.CardID)

		if card.Supertype == "Champion" {
			count += cardQuant.Quantity
		}
	}

	return count
}

func (d Deck) IsValid(strict, publish, sandbox bool) (bool, error) {
	if valid, cardID := d.AllCardsValid(); !valid {
		return false, types.InvalidDeckErrorFromString(fmt.Sprintf("Card with ID %d does not exist", cardID))
	}

	cardCount := d.CardCount()
	if publish && !sandbox && cardCount != 40 {
		return false, types.InvalidDeckErrorFromString("Deck must include 40 cards to be published")
	}

	if strict && cardCount == 0 {
		return false, types.InvalidDeckErrorFromString("Deck must include at least 1 card")
	}

	if d.ChampionCount() > 6 && publish {
		return false, types.InvalidDeckErrorFromString("Deck can only contain at most 6 Champion Cards")
	}

	for _, cardQuant := range d.Cards {
		if cardQuant.Quantity > 3 {
			return false, types.InvalidDeckErrorFromString("Deck can only contain, at most, 3 of any individual card")
		}
	}

	return true, nil
}

func (d Deck) mapCardsForEncoding() []deck_encoder.CardInDeck {
	cardsInDeck := make([]deck_encoder.CardInDeck, 0)
	for _, cardQuant := range d.Cards {
		card := Cards.GetCard(cardQuant.CardID)

		encodingCard := card.ToEncodableCardInDeck(cardQuant.Quantity)
		cardsInDeck = append(cardsInDeck, encodingCard)
	}

	return cardsInDeck
}

func (d Deck) ToEncodableDeck() deck_encoder.Deck {
	return deck_encoder.Deck{
		Cards: d.mapCardsForEncoding(),
	}
}

func (d Deck) Encode() string {
	encodable := d.ToEncodableDeck()

	code := deck_encoder.Encode(encodable)

	return code
}

func InitDeckModel(d *db.Database) *DeckModel {
	collection := d.Collection("decks")
	m := NewDeckModel(collection)
	return m
}

func NewDeckModel(collection *mongo.Collection) *DeckModel {
	return &DeckModel{
		collection: collection,
	}
}

func (m DeckModel) SaveDeck(deck Deck) (*Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	newDeck := deck
	newID, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	newDeck.ID = newID

	_, err = m.collection.InsertOne(ctx, newDeck)
	if err != nil {
		return nil, err
	}

	return &newDeck, nil
}

func (m DeckModel) UpdateDeck(deck Deck) (*Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	deckID := deck.ID

	deck.ID = ""

	var updatedDeck Deck
	after := options.After
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	curr := m.collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: deckID}}, bson.M{"$set": deck}, &options)
	err := curr.Decode(&updatedDeck)
	if err != nil {
		return nil, err
	}

	return &updatedDeck, nil
}

func (m DeckModel) GetDeck(deckID string) (*Deck, error) {
	var deck Deck

	result := m.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: deckID}})
	err := result.Decode(&deck)
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (m DeckModel) GetDecksByOwner(ownerName string) ([]*Deck, error) {
	var data []*Deck
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

func (m DeckModel) GetDecksByOwnerID(ownerID string) ([]*Deck, error) {
	var data []*Deck
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

func (m DeckModel) SearchDecks(search string) ([]*Deck, error) {
	var data []*Deck
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

func (m DeckModel) DeleteDeck(deckID string) (*Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var deletedDeck Deck
	filter := bson.M{"_id": deckID}
	update := bson.M{"$set": bson.M{"published": false, "deleted": true}}
	after := options.After
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	curr := m.collection.FindOneAndUpdate(ctx, filter, update, &options)
	err := curr.Decode(&deletedDeck)

	return &deletedDeck, err
}

func (m DeckModel) PublishDeck(deckID string) (*Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var deletedDeck Deck
	filterParams := make([]bson.M, 2)
	filterParams[0] = bson.M{"_id": deckID}
	filterParams[1] = bson.M{"deleted": false}

	filter := bson.M{"$and": filterParams}
	update := bson.M{"$set": bson.M{"published": true}}
	after := options.After
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	curr := m.collection.FindOneAndUpdate(ctx, filter, update, &options)
	err := curr.Decode(&deletedDeck)

	return &deletedDeck, err
}
