package game

import (
	"bomberman-game-server/game/generation"
	"bomberman-game-server/shared"
	"fmt"
	"strconv"
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
var endGameReason string

func initializeGame(width, height int) {
	fmt.Println("Initializing game...")
	// PlayingField = shared.GenerateEmptyField(width, height)
	PlayingField = generation.GenerateClassic(width, height, 0.3)

	fmt.Println("Setting spawn positions...")
	for _, player := range shared.Players {
		player.Pos = PlayingField.GetRandomSpawnPos()
		player.NextPos = player.Pos
		player.Bomb = *player.GetBasicBomb()
		player.Alive = true
		player.Speed = 1.0
		player.BombCount = 2
		player.MaxBombCount = 2
		player.TicksSinceLastMove = 0
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
	initializeGame(25, 25)

	for _, player := range shared.Players {
		shared.SendMessageToClientByID(player.ID, "id", strconv.Itoa(player.ID))
	}

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
				endGame()
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

	if isGameOver() {
		fmt.Println("Game over!")
		CurrentGameState = GameStateFinished
		endGameReason = "Game Over!"
	}

	tickAllPlayers()
	tickAllBombs()
	tickAllExplosions()

	shared.BroadcastGameState(PlayingField)

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

			PlayingField.ExplodeBomb(bomb, true)
		}
	}
}

func tickAllExplosions() {
	for y := 0; y < PlayingField.Height; y++ {
		for x := 0; x < PlayingField.Width; x++ {
			cell := &PlayingField.Cells[y][x]
			if cell.Type == shared.CellExplosion || cell.Type == shared.CellExplosionPierce {
				cell.TicksTillExplosionOver--
				if cell.TicksTillExplosionOver <= 0 {
					cell.Type = shared.CellEmpty
				}
			}
		}
	}
}

func tickAllPlayers() {
	shared.PlayersMutex.RLock()
	defer shared.PlayersMutex.RUnlock()

	for id, player := range shared.Players {

		if !player.Alive {
			continue
		}

		// MOVEMENT
		player.TicksSinceLastMove++
		if player.CanMove() && !player.NextPos.Equal(player.Pos) {
			player.Pos = player.NextPos
			player.TicksSinceLastMove = 0
		}

		// Player death handling and collect powerups
		cell := PlayingField.GetCellAtPos(player.Pos.X, player.Pos.Y)
		if cell != nil {

			if shared.IsPowerUp(cell.Type) {
				shared.HandlePowerUpForPlayer(cell.Type, player)
				cell.Type = shared.CellEmpty
			} else if cell.Type == shared.CellExplosion || cell.Type == shared.CellExplosionPierce {
				player.Alive = false
			}
		}

		// BOMB PLACEMENT
		if player.WantsToPlaceBomb {
			if player.BombCount <= 0 {
				player.BombCount = 0
			}
			if player.BombCount > 0 {
				PlayingField.PlaceBomb(player)
			}
			player.BombCount--
			player.WantsToPlaceBomb = false
		}

		shared.Players[id] = player
		fmt.Printf("%d has %d max bombs\n", id, player.MaxBombCount)
	}
}

/*
func endGame() {

	endMessage := "Game Over!"

	if len(shared.Players) <= 1 {
		endMessage = "No players left!"
	}

	alivePlayers := 0
	for _, player := range shared.Players {
		if player.Alive {
			alivePlayers++
		}
	}

	if alivePlayers <= 1 {
		endMessage = "WINNER! Only one player left!"
	}

	if alivePlayers == 0 {
		endMessage = "TIE! All players have been eliminated!"
	}

	shared.BroadcastMessage("game_over", endMessage, false)

	shared.PlayersMutex.Lock()
	defer shared.PlayersMutex.Unlock()

	for id, player := range shared.Players {
		if player == nil {
			delete(shared.Players, id)
			continue
		}
		// Serialize writes/close to avoid concurrent write/close panics
		player.WriteMutex.Lock()
		if player.Conn != nil {
			_ = player.Conn.Close()
			player.Conn = nil
		}
		player.WriteMutex.Unlock()

		delete(shared.Players, id)
	}
}
*/

func endGame() {

	endMessage := "Game Over!"

	if len(shared.Players) <= 1 {
		endMessage = "No players left!"
	}

	alivePlayers := 0
	for _, player := range shared.Players {
		if player.Alive {
			alivePlayers++
		}
	}

	if alivePlayers <= 1 {
		endMessage = "WINNER! Only one player left!"
	}

	if alivePlayers == 0 {
		endMessage = "TIE! All players have been eliminated!"
	}

	shared.BroadcastMessage("game_over", endMessage, false)

	shared.PlayersMutex.Lock()
	for id, player := range shared.Players {
		if player == nil {
			delete(shared.Players, id)
			continue
		}
		player.WriteMutex.Lock()
		if player.Conn != nil {
			_ = player.Conn.Close()
			player.Conn = nil
		}
		player.WriteMutex.Unlock()

		delete(shared.Players, id)
	}
	shared.PlayersMutex.Unlock()
}
