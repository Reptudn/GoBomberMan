package game

import (
	"bomberman-game-server/shared"
	"fmt"
	"os"
	"time"
)

type GameState int

const (
	GameStateWaiting  GameState = iota // 0
	GameStatePlaying                   // 1
	GameStateFinished                  // 2
)

const TicksPerSecond = 20

var CurrentGameState GameState = GameStateWaiting
var gameWasStarted bool = false

func GetGameWasStarted() bool {
	return gameWasStarted
}

var PlayingField *shared.Field

func initializeGame(width, height int) {
	fmt.Println("Initializing game...")
	PlayingField = shared.GenerateEmptyField(width, height)

	// TODO: Set spawn Positions
	fmt.Println("Setting spawn positions...")
	for _, player := range shared.Players {
		player.Pos = shared.Pos{X: width / 2, Y: height / 2}
		player.NextPos = player.Pos
		player.Bomb = *player.GetBasicBomb()
		player.Alive = true
		player.Speed = 1.0
		player.BombCount = 1
		player.MaxBombCount = 2
	}
	fmt.Println("Game initialized.")
}

func StartGame() error {
	if gameWasStarted {
		fmt.Println("Game already started.")
		return fmt.Errorf("Game already started")
	}

	if len(shared.Players) <= 1 {
		fmt.Println("Not enough players to start the game.")
		return fmt.Errorf("Not enough players to start the game (Min 2)")
	}

	gameWasStarted = true
	initializeGame(10, 10)

	fmt.Println("Starting game...")
	CurrentGameState = GameStatePlaying

	shared.BroadcastMessage("game_start", "Game is starting...", false)

	fmt.Println("Starting game loop...")
	runGameLoop()
	return nil
}

func runGameLoop() {
	if !gameWasStarted {
		fmt.Println("Game wasnt initialized yet.")
		return
	}
	ticker := time.NewTicker(time.Second / TicksPerSecond)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			if CurrentGameState == GameStateFinished {
				ticker.Stop()
				fmt.Println("Exiting game loop...")
				endGame("Game is over!")
				os.Exit(0)
				return
			}
			if CurrentGameState == GameStatePlaying {
				gameLoop()
			}
		}
	}()
}

func gameLoop() {

	fmt.Println("Tick")

	tickAllPlayers()
	tickAllBombs()

	shared.BroadcastGameState(PlayingField)

	if isGameOver() {
		fmt.Println("Game over!")
		CurrentGameState = GameStateFinished
	}
}

func isGameOver() bool {

	shared.PlayersMutex.RLock()
	defer shared.PlayersMutex.RUnlock()

	if len(shared.Players) <= 1 {
		return true
	}

	playersAlive := 0
	for _, player := range shared.Players {
		if player.Alive {
			playersAlive++
		}
	}

	return playersAlive <= 1
}

func tickAllBombs() {
	bombs := PlayingField.GetAllBombs()
	for _, bomb := range bombs {
		bomb.TicksTillExplosion--

		if bomb.TicksTillExplosion <= 0 {
			bomb.Explode()
		}
	}
}

func tickAllPlayers() {
	shared.PlayersMutex.RLock()
	defer shared.PlayersMutex.RUnlock()

	for id, player := range shared.Players {

		// BOMB PLACEMENT
		if player.WantsToPlaceBomb {
			player.WantsToPlaceBomb = false
			player.BombCount--
			PlayingField.PlaceBomb(player)
		}

		// MOVEMENT
		if !player.NextPos.Equal(player.Pos) {
			player.Pos = player.NextPos
		}

		shared.Players[id] = player
	}
}

func endGame(endMessage string) {

	shared.BroadcastMessage("game_over", endMessage, false)

	shared.PlayersMutex.RLock()
	defer shared.PlayersMutex.RUnlock()

	for _, player := range shared.Players {
		if player.Conn == nil {
			continue
		}
		player.Conn.Close()
	}
}
