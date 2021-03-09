package types

import "github.com/google/go-cmp/cmp"

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
