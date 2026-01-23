package routes

import (
	"bomberman-report-layer/shared"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateGame(c echo.Context) error {

	var updatedGame shared.GameServer
	if err := c.Bind(&updatedGame); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := shared.UpdateGame(updatedGame.UUID, updatedGame); err != nil {
		return c.JSON(http.StatusFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Upated game with uuid %s", updatedGame.UUID))
}
