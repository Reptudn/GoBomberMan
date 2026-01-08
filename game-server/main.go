package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	GameStateWaiting  GameState = iota // 0
	GameStatePlaying                   // 1
	GameStatePaused                    // 2
	GameStateFinished                  // 3
)

var currentState GameState = GameStateWaiting

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.RWMutex
	nextClientID = 1
	serverReady  = false
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	_ = r
	if serverReady {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Server not ready")
	}
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

		if len(clients) == 0 {
			os.Exit(0)
		}
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
var startTime time.Time

// TODO: Handle shutdown in a nice way

func main() {
	startTime = time.Now()
	fmt.Println("Game Server started on :8080")

	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if currentState == GameStateFinished {
				os.Exit(0)
				return
			}
			update() // This runs 20 times per second
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		serverReady = true
		fmt.Println("Server ready")
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func update() {

	// if the game wasnt started and has no players connected it ends the game
	// if time.Now().After(startTime.Add(time.Second*30)) && !gameWasStarted && len(clients) == 0 {
	// 	fmt.Println("Game ended due to inactivity")
	// 	currentState = GameStateFinished
	// 	return
	// }

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
