package routes

import (
	"bomberman-report-layer/shared"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterGame(c echo.Context) error {

	fmt.Println("Trying to register a game.")

	var newGame shared.GameServer
	if err := c.Bind(&newGame); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	if err := shared.AddGame(newGame); err != nil {
		return c.JSON(http.StatusFound, err.Error())
	}

	fmt.Println("Registered Game!")
	return c.JSON(http.StatusCreated, fmt.Sprintf("The Game with the UUID of %s has been registered successfully!", newGame.UUID))
}
