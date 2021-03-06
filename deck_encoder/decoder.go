package deck_encoder

import (
	"encoding/base32"
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidCode = errors.New("deckcode invalid")
	ErrOldVersion  = errors.New("The provided code requires a higher version of this library; please update.")
)

func fixDeckcodeLength(dc string) string {
	length := len(dc)
	if length%8 != 0 {
		for i := 0; i < 8-length%8; i++ {
			dc = dc + string(base32.StdPadding)
		}
	}
	return dc
}

//each deckcode starts with 0001 0001 - 4 bits for format(currently only 1) and 4 bits for version(currently only 1)
// format1 version1 is represented by 00010001 = 17
func decodeHeader(bs []byte) ([]byte, error) {
	byteFormatVersion, c := binary.Uvarint(bs)
	version := byteFormatVersion & 0xF
	if c != 1 {
		return bs, ErrInvalidCode
	}
	if int(version) > MAX_KNOWN_VERSION {
		return bs, ErrOldVersion
	}
	bs = bs[c:]
	return bs, nil
}

func decodeByteStream(bs []byte) (Deck, error) {
	var deck Deck
	for len(bs) > 0 {
		for i := 0; i < MAX_CARD_COUNT; i++ {
			var err error
			var cards []CardInDeck
			bs, cards, err = decodeSetFactionCombinations(bs, MAX_CARD_COUNT-i)
			if err != nil {
				return Deck{}, err
			}
			deck.Cards = append(deck.Cards, cards...)
		}
		if len(bs) != 0 {
			return Deck{}, ErrInvalidCode
		}
	}
	return deck, nil
}

func decodeSetFactionCombinations(bs []byte, count int) ([]byte, []CardInDeck, error) {
	var returnCards []CardInDeck
	combinationCount, c := binary.Uvarint(bs)
	bs = bs[c:]
	for j := 0; j < int(combinationCount); j++ {
		var cards []CardInDeck
		var err error
		bs, cards, err = decodeSetFactionCombinationCards(bs, count)
		if err != nil {
			return []byte{}, []CardInDeck{}, err
		}
		returnCards = append(returnCards, cards...)
	}
	return bs, returnCards, nil
}

func decodeSetFactionCombinationCards(bs []byte, count int) ([]byte, []CardInDeck, error) {
	var cards []CardInDeck
	countOfUniqueCards, c := binary.Uvarint(bs)
	bs = bs[c:]
	set, c := binary.Uvarint(bs)
	bs = bs[c:]
	faction, c := binary.Uvarint(bs)
	bs = bs[c:]
	for i := 0; i < int(countOfUniqueCards); i++ {
		cardNumber, c := binary.Uvarint(bs)
		bs = bs[c:]
		card := CardInDeck{
			Card: Card{
				Set:     int(set),
				Faction: int(faction),
				Number:  int(cardNumber),
			},
			Count: count,
		}
		cards = append(cards, card)
	}
	return bs, cards, nil
}
