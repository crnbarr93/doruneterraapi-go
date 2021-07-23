package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
