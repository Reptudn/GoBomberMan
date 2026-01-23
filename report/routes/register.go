package routes

import (
	"bomberman-report-layer/shared"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterGame(c echo.Context) error {

	var newGame shared.GameServer
	if err := c.Bind(&newGame); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	if err := shared.AddGame(newGame); err != nil {
		return c.JSON(http.StatusFound, fmt.Sprintf("The Game with the uuid of %s has alreabdy been registered", newGame.UUID))
	}

	return c.JSON(http.StatusCreated, fmt.Sprintf("The Game with the UUID of %s has been registered successfully!", newGame.UUID))
}
