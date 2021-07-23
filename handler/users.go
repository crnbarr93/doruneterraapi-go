package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

type LoginRequest struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func Login(c echo.Context) error {
	u := new(LoginRequest)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	email := u.Email
	password := u.Password

	user, err := models.Users.Login(email, password)
	if err != nil {
		return err
	}

	return c.JSON(200, user)
}

type RegisterRequest struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func Register(c echo.Context) error {
	u := new(RegisterRequest)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	username := u.Username
	email := u.Email
	password := u.Password

	user, err := models.Users.Register(username, email, password)
	if err != nil {
		return err
	}

	return c.JSON(200, user)
}

func SearchUsers(c echo.Context) error {
	email := c.QueryParam("email")
	username := c.QueryParam("username")

	users, err := models.Users.SearchUsers(username, email)
	if err != nil {
		return err
	}

	return c.JSON(200, users)
}

func ValidateEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if len(email) == 0 {
		return echo.ErrBadRequest
	}

	user, _ := models.Users.GetUserByEmail(email)

	return c.JSON(200, user == nil)
}

func ValidateUsername(c echo.Context) error {
	username := c.QueryParam("username")
	if len(username) == 0 {
		return echo.ErrBadRequest
	}

	user, _ := models.Users.GetUserByUsername(username)

	return c.JSON(200, user == nil)
}
