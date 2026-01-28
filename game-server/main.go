package main

import (
	"bomberman-game-server/game"
	networking "bomberman-game-server/networking"
	"bomberman-game-server/report"
	"bomberman-game-server/shared"
	"fmt"
	"log"
	"net/http"
	"time"
)

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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// TODO: Handle shutdown in a nice way
var startTime time.Time

func main() {
	startTime = time.Now()
	fmt.Println("Game Server started on :8080")

	http.HandleFunc("/ws", corsMiddleware(networking.HandleWebSocket))
	http.HandleFunc("/health", corsMiddleware(handleHealthCheck))

	// This stops the game when nobody is connected and the game was never started
	idleChecker := time.NewTicker(30 * time.Second)
	shared.GlobalWaitGroup.Add(1)
	go func() {
		defer idleChecker.Stop()
		defer shared.GlobalWaitGroup.Done()

		for range idleChecker.C {

			if game.CurrentGameState == game.GameStateFinished {
				break
			}

			uptime := time.Since(startTime)
			log.Printf("Server uptime: %s", uptime.String())
			if uptime > 2*time.Minute && len(shared.Players) == 0 && !game.GetGameWasStarted() {
				log.Println("No players connected and game hasn't started for 2 minutes. Shutting down server.")
				game.CurrentGameState = game.GameStateFinished
				break
			}
		}

	}()

	report.NewStatusReporter()

	log.Fatal(http.ListenAndServe(":8080", nil))
	shared.GlobalWaitGroup.Wait() // Waiting for all the other threads to shut down when game state it set to finished
	fmt.Println("Game Server stopped!")
}
