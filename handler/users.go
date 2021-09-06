package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/utils"
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

	jwtCookie := utils.CreateJWTCookie(*user)
	if jwtCookie != nil {
		c.SetCookie(jwtCookie)
	}

	return c.JSON(200, user)
}

func Logout(c echo.Context) error {
	emptyCookie := &http.Cookie{
		Name:    "authtoken",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	c.SetCookie(emptyCookie)

	return c.JSON(200, true)
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

func Auth(c echo.Context) error {
	cookie := utils.GetJWTCookie(c.Cookies())

	if cookie == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid auth token")
	}

	user, err := utils.DecodeToken(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not decode auth token")
	}

	return c.JSON(200, user)
}
