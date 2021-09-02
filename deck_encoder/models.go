package deck_encoder

type Deck struct {
	Cards []CardInDeck
}

type CardInDeck struct {
	Card  Card
	Count int
}

type Card struct {
	Set     int
	Faction int
	Number  int
}

type Faction string
type FactionName string

const (
	DEMACIA      Faction = "DE"
	FRELJORD     Faction = "FR"
	IONIA        Faction = "IO"
	NOXUS        Faction = "NX"
	PILTOVERZAUN Faction = "PZ"
	SHADOWISLES  Faction = "SI"
	BILGEWATER   Faction = "BW"
	SHURIMA      Faction = "SH"
	MOUNTTARGON  Faction = "MT"
	BANDLECITY   Faction = "BC"
	UNKNOWN      Faction = "XX"
)

const MAX_CARD_COUNT = 3

const MAX_KNOWN_VERSION = 4
