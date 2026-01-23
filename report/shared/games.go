package shared

import (
	"fmt"
	"sync"
	"time"
)

type GameState int

const (
	Starting = iota
	Lobby
	Running
	Stopping
)

type GameServer struct {
	UUID            string    `json:"uuid"`
	GameState       GameState `json:"gameState"`
	CurrPlayerCount int       `json:"currPlayerCount"`
	MaxPlayerCount  int       `json:"maxPlayerCount"`
	LastUpdateTime  time.Time `json:"-"`
}

var Games map[string]GameServer = make(map[string]GameServer)

var gamesMutex sync.Mutex

func AddGame(newGame GameServer) error {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	_, exists := Games[newGame.UUID]
	if exists {
		return fmt.Errorf("A Game with the UUID of %s has already been registered", newGame.UUID)
	}

	newGame.LastUpdateTime = time.Now()
	Games[newGame.UUID] = newGame
	return nil
}

func RemoveGame(uuid string) error {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	_, exists := Games[uuid]
	if !exists {
		return fmt.Errorf("A Game with the UUID of %s hasn't been registered!", uuid)
	}

	delete(Games, uuid)
	return nil
}

func UpdateGame(uuid string, updatedData GameServer) error {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	if _, exists := Games[uuid]; !exists {
		return fmt.Errorf("Game not registered!")
	}

	updatedData.LastUpdateTime = time.Now()
	Games[uuid] = updatedData
	return nil
}
