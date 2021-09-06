package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

func TestJWT(t *testing.T) {
	user := models.User{
		Username: "testuser",
	}

	token, err := CreateToken(user)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	decoded, err := DecodeToken(token)
	assert.Nil(t, err)
	assert.Equal(t, &user, decoded)
}
