package deck_encoder

import (
	"fmt"
)

func (c Card) String() string {
	return fmt.Sprintf("%02d%s%03d", c.Set, c.GetFaction(), c.Number)
}

func (c Card) GetFaction() Faction {
	factions := []Faction{DEMACIA, FRELJORD, IONIA, NOXUS, PILTOVERZAUN, SHADOWISLES, BILGEWATER, SHURIMA, MOUNTTARGON, BANDLECITY}
	if c.Faction >= len(factions) {
		return UNKNOWN
	}
	return factions[c.Faction]
}

func FactionNumberFromName(reg string) int {
	factionNames := map[string]int{"Demacia": 0, "Freljord": 1, "Ionia": 2, "Noxus": 3, "Piltover & Zaun": 4, "Shadow Isles": 5, "Bilgewater": 6, "Shurima": 7, "Mount Targon": 9, "Bandle City": 10}
	if val, ok := factionNames[reg]; ok {
		return val
	}

	return -1
}
