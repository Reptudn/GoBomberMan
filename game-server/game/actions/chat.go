package actions

import (
	"bomberman-game-server/shared"
)

func handleChat(action *shared.ChatData) error {
	shared.BroadcastMessage("chat", action.Message, false)
	return nil
}
