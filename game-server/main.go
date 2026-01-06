package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       int
	Conn     *websocket.Conn
	PlayerID int // Server-assigned player ID
}

type GameState int

const (
    GameStateWaiting GameState = iota  // 0
    GameStatePlaying                    // 1
    GameStatePaused                     // 2
    GameStateFinished                   // 3
)

var currentState GameState = GameStateWaiting

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.RWMutex
	nextClientID = 1
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Create and register client
	clientsMutex.Lock()
	client := &Client{
		ID:       nextClientID,
		Conn:     conn,
		PlayerID: nextClientID, // Use same ID for player
	}
	clients[client.ID] = client
	nextClientID++
	clientsMutex.Unlock()

	log.Printf("Client %d connected (total clients: %d)", client.ID, len(clients))

	// Cleanup on disconnect
	defer func() {
		clientsMutex.Lock()
		delete(clients, client.ID)
		clientsMutex.Unlock()
		log.Printf("Client %d disconnected (remaining clients: %d)", client.ID, len(clients))
	}()

	// Read messages from client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Client %d sent: %s", client.ID, message)

		// Echo back with client ID
		response := fmt.Sprintf("Server received from client %d: %s", client.ID, message)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

var gameWasStarted bool = false
func main() {
	fmt.Println("Game Server started on :8080")

	http.HandleFunc("/ws", handleWebSocket)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if currentState == GameStateFinished {
				return
			}
			update() // This runs 20 times per second
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func update() {

	if gameWasStarted && len(clients) == 0 {
		currentState = GameStateFinished
		return
	}

	switch currentState {
		case GameStateWaiting:
			lobbyState()
		case GameStatePlaying:
			playingState()
	}
}

func lobbyState() {

}

func playingState() {

}