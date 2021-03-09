package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
)

func GetCards(c echo.Context) error {
	data := db.GetAllCards()

	return c.JSON(http.StatusOK, data)
}
