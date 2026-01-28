package networking

import (
	"bomberman-game-server/game"
	"bomberman-game-server/game/actions"
	"bomberman-game-server/report"
	"bomberman-game-server/shared"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	Conn *websocket.Conn
}

var (
	nextClientID = 1
	serverReady  = false
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	player := handleClientConnect(conn)
	defer handleClientDisconnect(player)

	// Read messages from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Client %d sent: %s", player.ID, message)

		handleClientMessage(player, message)
	}
}

func handleClientConnect(conn *websocket.Conn) *shared.Player {

	/*
		if game.CurrentGameState != game.GameStateWaiting {
			fmt.Println("Game is not in waiting state")
			if err := conn.Close(); err != nil {
				log.Println("Error closing connection:", err)
			}
			return nil
			}*/

	shared.PlayersMutex.Lock()
	player := &shared.Player{
		ID:   nextClientID,
		Conn: conn,
	}
	shared.Players[nextClientID] = player
	nextClientID++
	shared.PlayersMutex.Unlock()

	report.TriggerManualStatusUpdate()

	log.Printf("Client %d connected (total clients: %d)", player.ID, len(shared.Players))
	return player
}

func handleClientDisconnect(player *shared.Player) {
	{
		shared.PlayersMutex.Lock()
		defer shared.PlayersMutex.Unlock()
		delete(shared.Players, player.ID)
	}

	log.Printf("Client %d disconnected (remaining clients: %d)", player.ID, len(shared.Players))

	if len(shared.Players) == 0 {
		game.CurrentGameState = game.GameStateFinished
	}
	report.TriggerManualStatusUpdate()
}

func handleClientMessage(player *shared.Player, message []byte) {
	log.Printf("Client %d sent: %s", player.ID, message)

	action, errParse := shared.ParseAction(message)
	if errParse != nil {
		sendMessageToClient(player, "error", "Invalid Message: Parsing failed!")
		return
	}

	msg, errAction := actions.HandlePlayerAction(player, action)
	if errAction != nil {
		sendMessageToClient(player, "error", errAction.Error())
		return
	}

	fmt.Println("Sending successful action message to client")
	sendMessageToClient(player, "success", msg)
}

func sendMessageToClient(player *shared.Player, msgType string, message string) {
	player.WriteMutex.Lock()
	defer player.WriteMutex.Unlock()
	msg := shared.BuildMessage(msgType, message)
	err := player.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("Send error:", err)
	}
}
