package types

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

type Card struct {
	ID                    string   `json:"_id,omitempty" bson:"_id,omitempty"`
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
	CardSet               int      `json:"card_set" bson:"card_set"`
	CardSubset            int      `json:"card_subset,omitempty" bson:"card_subset,omitempty"`
}

func (c Card) Compare(b Card) bool {
	return cmp.Equal(c, b)
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

type User struct {
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Access      int       `json:"access" bson:"access"`
	Username    string    `json:"username" bson:"username"`
	Email       string    `json:"email" bson:"email"`
	Password    string    `json:"password" bson:"password"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}
