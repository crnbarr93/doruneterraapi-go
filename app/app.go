package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/handler"
)

type App struct {
	Router *echo.Echo
	DB     *db.Database
}

func New(router *echo.Echo, db *db.Database) *App {
	return &App{
		Router: router,
		DB:     db,
	}
}

func (a *App) Run(port string) error {
	err := a.DB.Connect()
	if err != nil {
		return err
	}

	// Middleware
	a.Router.Use(middleware.Logger())
	a.Router.Use(middleware.Recover())
	a.Router.Pre(middleware.RemoveTrailingSlash())

	//Routes
	a.Router.GET("/",handler.Hello)
		

	// cardRoutes := a.Router.Group("/cards")
	// cardRoutes.GET("/all", handler.GetCards)

	// Start server
	return a.Router.Start(port)
}
