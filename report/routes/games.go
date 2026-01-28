package routes

import (
	"bomberman-report-layer/shared"
	"net/http"

	"github.com/labstack/echo/v4"
)

type queryOptions struct {
	ID string `query:"id"`

	showFullGames  bool `query:"showFullGames"`
	showEmptyGames bool `query:"showEmptyGames"`
	showRunning    bool `query:"showRunning"`
}

func GetGames(c echo.Context) error {
	var options queryOptions
	if err := c.Bind(&options); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Scheme options provided!")
	}

	var games []shared.GameServer = make([]shared.GameServer, 0)

	// TODO: only append games based on the given query criteria
	for _, game := range shared.Games {
		var include bool = true

		// check if game is full
		if game.CurrPlayerCount == game.MaxPlayerCount {
			include = false

			if options.showFullGames {
				include = true
			}

		}

		// check if game is empty
		if game.CurrPlayerCount == 0 {
			include = false

			if options.showEmptyGames {
				include = true
			}
		}

		// check if the game is joinable aka not in any other state than the lobby state
		if game.GameState != shared.GameStateLobby {
			include = false

			if options.showRunning {
				include = true
			}
		}

		if include {
			games = append(games, game)
		}
	}

	return c.JSON(http.StatusOK, games)
}
