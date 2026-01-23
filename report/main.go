package report

import (
	"bomberman-report-layer/routes"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	isProd := os.Getenv("ENV") == "production"

	if isProd {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	} else {
		e.Debug = true
		e.Use(middleware.RequestLogger())
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"Content-Length", "Accept-Encoding", "X-CSRF-Token", "Cache-Control", "X-Requested-With",
		},
		AllowCredentials: true,
	}))

	e.GET("/games", routes.GetGames)
	e.POST("/register", routes.RegisterGame)
	e.POST("/update", routes.UpdateGame)
	e.DELETE("/unregister", routes.UnregisterGame)

	e.Logger.Fatal(e.Start(":8080"))
}
