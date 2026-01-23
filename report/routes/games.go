package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type queryOptions struct {
	ID string `query:"id"`

	showFullGames  bool `query:"showFullGames"`
	showEmptyGames bool `query:"showEmptyGames"`
	showJoinable   bool `query:"showJoinable"`
	showRunning    bool `query:"showRunning"`
}

func GetGames(c echo.Context) error {
	var options queryOptions
	if err := c.Bind(&options); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Scheme options provided!")
	}

	// TODO: implement

	return nil
}
