package models_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

var savedUser *models.User

func init() {
	saveUser()
}

func saveUser() {
	email := "test@test.com"
	username := "test_user"
	password := "password"

	user, err := models.Users.Register(username, email, password)
	if err != nil {
		panic(err)
	}

	savedUser = user
}

func TestRegister(t *testing.T) {
	email := "test@test.com"
	username := "test_user"
	password := "password"

	assert.Equal(t, savedUser.Email, email)
	assert.Equal(t, savedUser.Username, username)

	failUser, err := models.Users.Register(username, email, password)

	assert.Nil(t, failUser)
	assert.NotNil(t, err)
}

func TestLogin(t *testing.T) {

	email := savedUser.Email
	password := "password"

	loggedInUser, err := models.Users.Login(email, password)

	assert.Nil(t, err)
	assert.Equal(t, loggedInUser.Email, email)

	wrongPassword := "wrong"

	failUser, err := models.Users.Login(email, wrongPassword)

	assert.Nil(t, failUser)
	assert.NotNil(t, err)
}

func TestGetUserById(t *testing.T) {
	correctID := savedUser.UserID()
	correctUser, err := models.Users.GetUserById(correctID)
	assert.Nil(t, err)
	assert.Equal(t, correctUser.ID.Hex(), correctID)
}

func TestGetUserByEmail(t *testing.T) {
	correctEmail := savedUser.Email
	correctUser, err := models.Users.GetUserByEmail(correctEmail)
	assert.Nil(t, err)
	assert.Equal(t, correctUser.Email, correctEmail)

	incompleteEmail := savedUser.Email[:len(savedUser.Email)-2]
	incompleteUser, err := models.Users.GetUserByEmail(incompleteEmail)
	assert.Nil(t, incompleteUser)
	assert.NotNil(t, err)

	caseSensitiveEmail := strings.ToUpper(savedUser.Email)
	caseSensitiveUser, err := models.Users.GetUserByEmail(caseSensitiveEmail)
	assert.Nil(t, err)
	assert.Equal(t, caseSensitiveUser.Email, correctEmail)
}

func TestGetUserByUsername(t *testing.T) {
	correctUsername := savedUser.Username
	correctUser, err := models.Users.GetUserByUsername(correctUsername)
	assert.Nil(t, err)
	assert.Equal(t, correctUser.Username, correctUsername)

	incompleteUsername := savedUser.Username[:len(savedUser.Username)-2]
	incompleteUser, err := models.Users.GetUserByUsername(incompleteUsername)
	assert.Nil(t, incompleteUser)
	assert.NotNil(t, err)

	caseSensitiveUsername := strings.ToUpper(savedUser.Username)
	caseSensitiveUser, err := models.Users.GetUserByUsername(caseSensitiveUsername)
	assert.Nil(t, err)
	assert.Equal(t, caseSensitiveUser.Username, correctUsername)
}

func TestSearch(t *testing.T) {
	email := savedUser.Email[:len(savedUser.Email)-2]
	username := savedUser.Username[:len(savedUser.Username)-1]

	usersByEmail, err := models.Users.SearchUsers("", email)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(usersByEmail))
	assert.True(t, strings.Contains(usersByEmail[0].Email, email))

	usersByUsername, err := models.Users.SearchUsers(username, "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(usersByUsername))
	assert.True(t, strings.Contains(usersByEmail[0].Username, username))
}

func TestUpdateUser(t *testing.T) {
	socials := models.SocialLinks{
		Instagram: "@someuser",
	}
	userToUpdate := savedUser
	userToUpdate.Socials = socials

	received, err := models.Users.UpdateUser(userToUpdate)

	assert.Nil(t, err)
	assert.Equal(t, received.Socials, socials)
}
