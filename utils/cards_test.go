package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetURL(t *testing.T) {
	expected := "https://dd.b.pvp.net/latest/set4/en_us/data/set4-en_us.json"
	received := getSetURL(4)

	assert.Equal(t, expected, received)
}

// func TestGetSetData(t *testing.T) {
// 	data := getSetData(1)

// 	expected := "01IO012"
// }

func TestGetSetInteger(t *testing.T) {
	card := DDCard{Set: "set1"}

	expected := 1
	received, err := card.getSetInteger()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
