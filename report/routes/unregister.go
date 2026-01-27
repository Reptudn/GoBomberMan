package routes

import (
	"bomberman-report-layer/shared"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type unregisterGameBody struct {
	UUID string `json:"uuid"`
}

func UnregisterGame(c echo.Context) error {
	var request unregisterGameBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format.")
	}

	if err := shared.RemoveGame(request.UUID); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("The Game with the uuid of %s has been unregistered successfully!", request.UUID))
}
