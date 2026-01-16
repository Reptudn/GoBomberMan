package actions

import "bomberman-game-server/shared"

func handlePlaceBomb(player *shared.Player) {
	if player.BombCount <= 0 {
		return
	}
	player.WantsToPlaceBomb = true
}
