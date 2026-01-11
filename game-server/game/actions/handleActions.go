package actions

import (
	"bomberman-game-server/game"
	"bomberman-game-server/shared"
	"fmt"
)

func HandlePlayerAction(player *shared.Player, action *shared.Action) (string, error) {

	fmt.Println("Handling player action of type:", action.Type)

	switch action.Type {
	case "move":
		{
			fmt.Println("Handling move action")
			moveData, err := action.GetMoveData()
			if err != nil {
				return "", fmt.Errorf("Failed to get move data: %w", err)
			}
			handlePlayerMovement(player, moveData)
			break
		}
	case "place_bomb":
		{
			fmt.Println("Handling place bomb action")
			handlePlaceBomb(player)
			break
		}
	case "chat":
		{
			fmt.Println("Handling chat action")
			chatData, err := action.GetChatData()
			if err != nil {
				return "", fmt.Errorf("Failed to get chat data: %w", err)
			}
			handleChat(chatData)
			break
		}
	case "start_game":
		{
			fmt.Println("Handling start game action")
			err := game.StartGame()
			if err != nil {
				return "", fmt.Errorf("Failed to start game: %w", err)
			}
			break
		}
	default:
		return "", fmt.Errorf("Invalid action: %s", action.Type)
	}

	return "OK", nil
}
