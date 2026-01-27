package routes

import "github.com/labstack/echo/v4"

func Ping(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"message": "Pong from Main Backend",
	})
}
