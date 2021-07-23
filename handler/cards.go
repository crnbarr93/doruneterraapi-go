package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

func GetCards(c echo.Context) error {
	data := models.Cards.GetAll()

	return c.JSON(http.StatusOK, data)
}

func GetCard(c echo.Context) error {
	id := c.Param("id")
	data := models.Cards.GetCard(id)

	if data == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, data)
}
