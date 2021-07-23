package app

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/handler"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
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
	a.DB.WaitForConnection()
	models.InitModels(a.DB)

	if err != nil {
		return err
	}

	// Middleware
	a.Router.Validator = &handler.Validator{Validator: validator.New()}
	a.Router.Use(middleware.Logger())
	a.Router.Use(middleware.Recover())
	a.Router.Pre(middleware.RemoveTrailingSlash())
	a.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	a.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

	//Routes
	cardRoutes := a.Router.Group("/cards")
	cardRoutes.GET("/all", handler.GetCards)
	cardRoutes.GET("/:id", handler.GetCard)

	userRoutes := a.Router.Group("/users")
	userRoutes.POST("/login", handler.Login)
	userRoutes.POST("/", handler.Register)
	userRoutes.GET("/search", handler.SearchUsers)
	userRoutes.GET("/validate/email", handler.ValidateEmail)
	userRoutes.GET("/validate/username", handler.ValidateUsername)

	// Start server
	return a.Router.Start(port)
}
