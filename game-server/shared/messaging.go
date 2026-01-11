package shared

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func BuildMessage(msgType string, message string) string {
	return fmt.Sprintf(`{"type": "%s", "message": "%s"}`, msgType, message)
}

func BroadcastMessage(msgType string, message string) {
	PlayersMutex.RLock()
	defer PlayersMutex.RUnlock()

	msg := BuildMessage(msgType, message)
	for _, player := range Players {
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Broadcast error:", err)
		}
	}
}

// TODO: Format the game state better to not get this: {0 <nil> <nil>} or this: {{{} {0 0}} 0 0 {{} 0} {{} 0}}}
func BroadcastGameState(field *Field) {
	fmt.Println("Broadcasting game state...")
	fmt.Println("Game state:", Players, field)
	BroadcastMessage("game_state", fmt.Sprintf(`{"players": %v, "field": %v}`, Players, field))
}

func SendMessageToClientByID(clientID int, msgType string, message string) {
	PlayersMutex.RLock()
	defer PlayersMutex.RUnlock()

	player, ok := Players[clientID]
	if !ok {
		log.Printf("Client %d not found", clientID)
		return
	}

	msg := BuildMessage(msgType, message)
	err := player.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("Send error:", err)
	}
}
