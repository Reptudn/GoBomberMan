package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouteNotFound(e echo.Context) error {
	return e.JSON(http.StatusNotFound, "Invalid route!")
}
