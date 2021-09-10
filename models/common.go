package models

import "gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"

var Cards *CardModel
var Decks *DeckModel
var Users *UserModel
var Archetypes *ArchetypesModel

func InitModels(d *db.Database) {
	Cards = InitCardModel(d)
	Decks = InitDeckModel(d)
	Users = InitUserModel(d)
	Archetypes = InitArchetypesModel(d)
}
