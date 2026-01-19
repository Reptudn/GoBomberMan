package shared

import (
	"math/rand"
)

func randomEffect(player *Player) {
	p := CreateRandomPowerUp()
	if f := getEffectFuncByType(p); f != nil {
		f(player)
	}
}

func removeSpeedEffect(player *Player) {
	player.Speed--
}

func addSpeedEffect(player *Player) {
	player.Speed++
}

func addBombEffect(player *Player) {
	player.BombCount++
}

func removeBombEffect(player *Player) {
	if player.BombCount-1 >= MinBombCount {
		player.BombCount--
	}
}

func normalBombEffect(player *Player) {
	player.Bomb.PierceWalls = true
}

func piercingBombEffect(player *Player) {
	player.Bomb.PierceWalls = true
}

func addStrenghtEffect(player *Player) {
	player.Bomb.Strength++
}

func removeStrenghtEffect(player *Player) {
	if player.Bomb.Strength-1 >= MinBombStrenght {
		player.Bomb.Strength--
	}
}

func IsPowerUp(cellType CellType) bool {
	return cellType >= CellPowerUpRandom
}

func getEffectFuncByType(powerupType CellType) func(*Player) {
	if !IsPowerUp(powerupType) {
		return nil
	}

	switch powerupType {
	case CellPowerUpRandom:
		return randomEffect
	case CellPowerUpAddBomb:
		return addBombEffect
	case CellPowerUpRemoveBomb:
		return removeBombEffect
	case CellPowerUpNormalBomb:
		return normalBombEffect
	case CellPowerUpPiercingBomb:
		return piercingBombEffect
	case CellPowerUpAddStrenght:
		return addStrenghtEffect
	case CellPowerUpRemoveStrenght:
		return removeBombEffect
	case CellPowerUpAddSpeed:
		return addSpeedEffect
	case CellPowerUpRemoveSpeed:
		return removeSpeedEffect
	default:
		return nil
	}
}

func CreateRandomPowerUp() CellType {
	powerups := []CellType{
		CellPowerUpAddBomb,
		CellPowerUpRemoveBomb,
		CellPowerUpNormalBomb,
		CellPowerUpPiercingBomb,
		CellPowerUpAddStrenght,
		CellPowerUpRemoveStrenght,
		CellPowerUpAddSpeed,
		CellPowerUpRemoveSpeed,
	}
	return powerups[rand.Intn(len(powerups))]
}

func HandlePowerUpForPlayer(powerupType CellType, p *Player) {
	if !IsPowerUp(powerupType) {
		return
	}
	if effect := getEffectFuncByType(powerupType); effect != nil {
		effect(p)
	}
}
