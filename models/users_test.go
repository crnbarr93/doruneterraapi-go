package models_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/config"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

var users *models.UserModel

func InitializeDatabase() *db.Database {
	database := db.New(config.TestConfig.Database)

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.WaitForConnection()

	err = database.DropCollection("users")
	if err != nil {
		panic(err)
	}
	return database
}

func initUserModel(d *db.Database) {
	users = models.InitUserModel(d)
}

func TestRegister(t *testing.T) {
	db := InitializeDatabase()
	initUserModel(db)

	email := "test@test.com"
	username := "test_user"
	password := "password"

	user, err := users.Register(username, email, password)

	assert.Nil(t, err)
	assert.Equal(t, user.Email, email)
	assert.Equal(t, user.Username, username)

	failUser, err := users.Register(username, email, password)

	assert.Nil(t, failUser)
	assert.NotNil(t, err)
}

func TestLogin(t *testing.T) {
	TestRegister(t)

	email := "test@test.com"
	password := "password"

	user, err := users.Login(email, password)

	assert.Nil(t, err)
	assert.Equal(t, user.Email, email)

	wrongPassword := "wrong"

	failUser, err := users.Login(email, wrongPassword)

	assert.Nil(t, failUser)
	assert.NotNil(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	TestRegister(t)

	correctEmail := "test@test.com"
	correctUser, err := users.GetUserByEmail(correctEmail)
	assert.Nil(t, err)
	assert.Equal(t, correctUser.Email, correctEmail)

	incompleteEmail := "test@test."
	incompleteUser, err := users.GetUserByEmail(incompleteEmail)
	assert.Nil(t, incompleteUser)
	assert.NotNil(t, err)

	caseSensitiveEmail := "Test@test.com"
	caseSensitiveUser, err := users.GetUserByEmail(caseSensitiveEmail)
	assert.Nil(t, err)
	assert.Equal(t, caseSensitiveUser.Email, correctEmail)
}

func TestGetUserByUsername(t *testing.T) {
	TestRegister(t)

	correctUsername := "test_user"
	correctUser, err := users.GetUserByUsername(correctUsername)
	assert.Nil(t, err)
	assert.Equal(t, correctUser.Username, correctUsername)

	incompleteUsername := "test_"
	incompleteUser, err := users.GetUserByUsername(incompleteUsername)
	assert.Nil(t, incompleteUser)
	assert.NotNil(t, err)

	caseSensitiveUsername := "TesT_UseR"
	caseSensitiveUser, err := users.GetUserByUsername(caseSensitiveUsername)
	assert.Nil(t, err)
	assert.Equal(t, caseSensitiveUser.Username, correctUsername)
}

func TestSearch(t *testing.T) {
	TestRegister(t)

	email := "test@tes"
	username := "test_us"

	usersByEmail, err := users.SearchUsers("", email)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(usersByEmail))
	assert.True(t, strings.Contains(usersByEmail[0].Email, email))

	usersByUsername, err := users.SearchUsers(username, "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(usersByUsername))
	assert.True(t, strings.Contains(usersByEmail[0].Username, username))
}
