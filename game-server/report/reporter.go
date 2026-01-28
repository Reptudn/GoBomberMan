package report

import (
	"bomberman-game-server/game"
	"bomberman-game-server/shared"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const reportUrl = "http://report-service:8081"

type GameState struct {
	Uuid               string         `json:"uuid"`
	GameState          game.GameState `json:"gameState"`
	CurrentPlayerCount int            `json:"currPlayerCount"`
	MaxPlayerCount     int            `json:"maxPlayerCount"`
}

var httpClient *http.Client = nil
var updateChannel chan int

func NewStatusReporter() {

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	updateChannel = make(chan int)

	register()
	shared.GlobalWaitGroup.Add(2)
	go manualUpdatePushLoop()
	go reporterLoop()
}

func TriggerManualStatusUpdate() {
	if httpClient == nil {
		panic("No Status Reporter was created yet!")
	}
	updateChannel <- 1
}

func sendStatusUpdate() {
	var registerData GameState = GameState{
		Uuid:               game.UUID,
		GameState:          game.CurrentGameState,
		CurrentPlayerCount: len(shared.Players),
		MaxPlayerCount:     shared.MaxPlayers,
	}

	var registerBodyData, err = json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/status", reportUrl), bytes.NewReader(registerBodyData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func register() {

	var registerData GameState = GameState{
		Uuid:               game.UUID,
		GameState:          game.CurrentGameState,
		CurrentPlayerCount: len(shared.Players),
		MaxPlayerCount:     shared.MaxPlayers,
	}

	var registerBodyData, err = json.Marshal(registerData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register", reportUrl), bytes.NewReader(registerBodyData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func unregister() {
	defer shared.GlobalWaitGroup.Done()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/unregister", reportUrl), nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func manualUpdatePushLoop() {
	defer shared.GlobalWaitGroup.Done()

	for range updateChannel {
		// Trigger manual status update
		if game.CurrentGameState == game.GameStateFinished {
			break
		}
		sendStatusUpdate()
	}
}

func reporterLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		if game.CurrentGameState == game.GameStateFinished {
			ticker.Stop()
			break
		}
		sendStatusUpdate()
	}

	unregister()
}
