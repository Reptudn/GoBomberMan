package main

import (
	networking "bomberman-game-server/networking"
	"fmt"
	"log"
	"net/http"
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
