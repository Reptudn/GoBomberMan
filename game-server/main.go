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

	http.HandleFunc("/ws", networking.HandleWebSocket)
	http.HandleFunc("/health", handleHealthCheck)

	go func() {
		time.Sleep(2 * time.Second)
		serverReady = true
		fmt.Println("Server ready")
	}()

	// This stops the game when nobody is connected and the game was never started
	idleChecker := time.NewTicker(30 * time.Second)
	go func(startTime time.Time) {
		defer idleChecker.Stop()

		for range idleChecker.C {
			uptime := time.Since(startTime)
			log.Printf("Server uptime: %s", uptime.String())
			if uptime > 2*time.Minute && len(shared.Players) == 0 && !game.GetGameWasStarted() {
				log.Println("No players connected and game not started for 2 minutes. Shutting down server.")
				os.Exit(0)
			}
		}

	}(startTime)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
