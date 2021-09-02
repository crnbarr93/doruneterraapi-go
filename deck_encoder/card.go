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
	factions := []Faction{DEMACIA, FRELJORD, IONIA, NOXUS, PILTOVERZAUN, SHADOWISLES, BILGEWATER, SHURIMA, MOUNTTARGON, BANDLECITY}
	factionNames := []string{"Demacia", "Freljord", "Ionia", "Noxus", "Piltover & Zaun", "Shadow Isles", "Bilgewater", "Shurima", "Mount Targon", "Bandle City"}
	for i, name := range factionNames {
		if name == reg {
			return i
		}
	}

	return len(factions) + 1
}
