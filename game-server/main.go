package main

import (
	"bomberman-game-server/game"
	networking "bomberman-game-server/networking"
	"bomberman-game-server/shared"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var serverReady bool = false

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
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

// TODO: Handle shutdown in a nice way
var startTime time.Time

func main() {
	startTime = time.Now()
	fmt.Println("Game Server started on :8080")

	http.HandleFunc("/ws", corsMiddleware(networking.HandleWebSocket))
	http.HandleFunc("/health", corsMiddleware(handleHealthCheck))

	go func() {
		time.Sleep(2 * time.Second)
		serverReady = true
		fmt.Println("Server ready")
	}()

	// This stops the game when nobody is connected and the game was never started
	idleChecker := time.NewTicker(30 * time.Second)
	go func() {
		defer idleChecker.Stop()

		for range idleChecker.C {
			uptime := time.Since(startTime)
			log.Printf("Server uptime: %s", uptime.String())
			if uptime > 2*time.Minute && len(shared.Players) == 0 && !game.GetGameWasStarted() {
				log.Println("No players connected and game not started for 2 minutes. Shutting down server.")
				os.Exit(0)
			}
		}

	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
