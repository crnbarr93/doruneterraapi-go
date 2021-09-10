package models

import (
	"context"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/teris-io/shortid"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardInArchetype struct {
	CardID          string `json:"card" bson:"card"`
	Quantity        int    `json:"quantity" bson:"quantity"`
	QuantityAppears []int  `json:"quantityAppears" bson:"quantityAppears"`
}

type KeywordsInArchetype struct {
	Keyword  string  `json:"keyword" bson:"keyword"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Pct      float32 `json:"pct" bson:"pct"`
}

type Archetype struct {
	ID             string                `json:"_id,omitempty" bson:"_id,omitempty"`
	Decks          []string              `json:"decks" bson:"decks"`
	KeyCards       []CardInArchetype     `json:"keyCards,omitempty" bson:"keyCards"`
	Regions        []string              `json:"regions,omitempty" bson:"regions"`
	Title          string                `json:"title" bson:"title"`
	SanitizedTitle string                `json:"sanitizedTitle" bson:"sanitizedTitle"`
	Description    string                `json:"description" bson:"description"`
	Meta           string                `json:"meta" bson:"meta"`
	Keywords       []KeywordsInArchetype `json:"keywords,omitempty" bson:"keywords"`
	Background     string                `json:"background" bson:"background"`
	Deleted        bool                  `json:"deleted" bson:"deleted"`
	Hidden         bool                  `json:"hidden" bson:"hidden"`
}

type PopulatedArchetype struct {
	Archetype
	Decks []Deck `json:"decks" bson:"decks"`
}

type ArchetypesModel struct {
	collection *mongo.Collection
}

func NewArchetypesModel(collection *mongo.Collection) *ArchetypesModel {
	return &ArchetypesModel{
		collection: collection,
	}
}

func InitArchetypesModel(d *db.Database) *ArchetypesModel {
	collection := d.Collection("archetypes")
	return NewArchetypesModel(collection)
}

func (a Archetype) PopulateDecks() (*PopulatedArchetype, error) {
	decks, err := Decks.GetDecks(a.Decks)
	if err != nil {
		return nil, err
	}

	return &PopulatedArchetype{
		a,
		decks,
	}, nil
}

func (a *Archetype) CalculateDetails() error {
	popArch, err := a.PopulateDecks()
	if err != nil {
		return err
	}

	a.KeyCards = popArch.CalculateKeyCards()
	a.Regions = popArch.CalculateRegions()
	a.SanitizedTitle = a.SanitizeTitle()
	a.Keywords = popArch.CalculateKeywords()

	return nil
}

func addCardQuantity(quantities []CardInArchetype, quantity CardQuantity) []CardInArchetype {
	for i := 0; i < len(quantities); i++ {
		quant := quantities[i]
		if quant.CardID == quantity.CardID {
			quant.Quantity += quantity.Quantity
			for j := 0; j < len(quant.QuantityAppears); j++ {
				if j+1 == quantity.Quantity {
					quant.QuantityAppears[j] += quantity.Quantity
				}
			}

			quantities[i] = quant

			return quantities
		}
	}

	newQuantity := CardInArchetype{
		CardID:          quantity.CardID,
		Quantity:        quantity.Quantity,
		QuantityAppears: []int{0, 0, 0},
	}

	newQuantity.QuantityAppears[quantity.Quantity-1] += quantity.Quantity

	quantities = append(quantities, newQuantity)

	return quantities
}

func (a *PopulatedArchetype) CalculateKeyCards() []CardInArchetype {
	var cardQuantities []CardInArchetype

	for _, deck := range a.Decks {
		cards := deck.Cards

		for _, quant := range cards {
			cardQuantities = addCardQuantity(cardQuantities, quant)
		}
	}

	return cardQuantities
}

func (a *PopulatedArchetype) CalculateRegions() []string {
	var regions []string

	for _, deck := range a.Decks {
		deckRegions := deck.Regions

		for _, reg := range deckRegions {
			if !containsRegion(regions, reg) {
				regions = append(regions, reg)
			}
		}
	}
	sort.Strings(regions)
	return regions
}

func addKeyword(keywords []KeywordsInArchetype, keyword string, quantity int) []KeywordsInArchetype {
	for i := 0; i < len(keywords); i++ {
		quant := keywords[i]
		if quant.Keyword == keyword {
			quant.Quantity += quantity
			keywords[i] = quant
			return keywords
		}
	}

	newQuantity := KeywordsInArchetype{
		Keyword:  keyword,
		Quantity: quantity,
	}

	keywords = append(keywords, newQuantity)

	return keywords
}

func (a *PopulatedArchetype) CalculateKeywords() []KeywordsInArchetype {
	var keywords []KeywordsInArchetype
	var total int

	for _, deck := range a.Decks {
		cards := deck.Cards

		for _, quant := range cards {
			card := Cards.GetCard(quant.CardID)
			for _, keyword := range card.Keywords {
				keywords = addKeyword(keywords, keyword, quant.Quantity)
				total += quant.Quantity
			}
		}
	}

	for i := 0; i < len(keywords); i++ {
		keywordQuant := keywords[i]
		keywordQuant.Pct = float32(keywordQuant.Quantity) / float32(total)
		keywords[i] = keywordQuant
	}

	return keywords
}

func (a *Archetype) SanitizeTitle() string {
	title := a.Title

	title = strings.TrimSpace(title)

	regexSpace := regexp.MustCompile(` `)
	title = regexSpace.ReplaceAllString(title, "-")

	regexNonAlphanumeric := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	title = regexNonAlphanumeric.ReplaceAllString(title, "")

	title = strings.ToLower(title)
	title = url.QueryEscape(title)

	regex := regexp.MustCompile(`[!'()]`)
	title = regex.ReplaceAllString(title, "")

	regexTwo := regexp.MustCompile(`\*`)
	title = regexTwo.ReplaceAllString(title, "%2A")

	return title
}

func PopulateArchetypes(archetypes []Archetype) ([]PopulatedArchetype, error) {
	populatedArchetypes := make([]PopulatedArchetype, 0)
	for _, archetype := range archetypes {
		populatedArchetype, err := archetype.PopulateDecks()
		if err != nil {
			return nil, err
		}
		populatedArchetypes = append(populatedArchetypes, *populatedArchetype)
	}

	return populatedArchetypes, nil
}

func (m *ArchetypesModel) queryPopulatedArchetypes(query bson.M) ([]PopulatedArchetype, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var archetypes []Archetype

	cur, err := m.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &archetypes); err != nil {
		return nil, err
	}

	populatedArchetypes, err := PopulateArchetypes(archetypes)
	if err != nil {
		return nil, err
	}

	return populatedArchetypes, nil
}

func (m *ArchetypesModel) SaveArchetype(archetype Archetype) (*Archetype, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	newArchetype := archetype
	newID, err := shortid.Generate()
	if err != nil {
		return nil, err
	}

	newArchetype.ID = newID

	_, err = m.collection.InsertOne(ctx, newArchetype)
	if err != nil {
		return nil, err
	}

	return &newArchetype, nil
}

func (m *ArchetypesModel) GetArchetypes() ([]PopulatedArchetype, error) {
	return m.queryPopulatedArchetypes(bson.M{"deleted": false})
}

func (m *ArchetypesModel) GetArchetypesRaw() ([]*Archetype, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var archetypes []*Archetype

	cur, err := m.collection.Find(ctx, bson.M{"deleted": false})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &archetypes); err != nil {
		return nil, err
	}

	return archetypes, nil
}

func (m *ArchetypesModel) GetDeckArchetypes(deckID string) ([]PopulatedArchetype, error) {
	return m.queryPopulatedArchetypes(bson.M{"deleted": false, "decks": deckID})
}

func (m *ArchetypesModel) GetCardArchetypes(cardID string) ([]PopulatedArchetype, error) {
	return m.queryPopulatedArchetypes(bson.M{"deleted": false, "keyCards.card": cardID})
}
