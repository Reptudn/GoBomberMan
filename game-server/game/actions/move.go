package actions

import (
	"bomberman-game-server/game"
	"bomberman-game-server/shared"
	"fmt"
)

func handlePlayerMovement(player *shared.Player, moveData *shared.MoveData) {

	newPos := player.Pos

	switch moveData.Direction {
	case "up":
		newPos.Y -= 1
	case "down":
		newPos.Y += 1
	case "left":
		newPos.X -= 1
	case "right":
		newPos.X += 1
	default:
		fmt.Printf("Invalid direction: %s", moveData.Direction)
		return
	}

	if game.PlayingField.IsWalkable(newPos) {
		player.NextPos = newPos
	}
}
