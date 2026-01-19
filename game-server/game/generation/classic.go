package generation

import (
	"bomberman-game-server/shared"
	"math/rand"
)

func GenerateClassic(width, height int, freePerc float64) *shared.Field {
	field := &shared.Field{
		Width:  width,
		Height: height,
		Cells:  make([][]shared.Cell, height),
	}

	// minFreeSpots := int(float64(width*height) * freePerc)

	for y := range field.Cells {
		field.Cells[y] = make([]shared.Cell, width)
		for x := range field.Cells[y] {

			if x%2 == 1 && y%2 == 1 {
				field.Cells[y][x] = shared.Cell{Type: shared.CellWallIndestructible}
				continue
			}

			rand := rand.Intn(100)
			if rand < 15 {
				field.Cells[y][x] = shared.Cell{Type: shared.CellWallDestructible}
			} else if rand < 20 {
				powerup := shared.CreateRandomPowerUp()
				field.Cells[y][x] = shared.Cell{Type: powerup}
			} else {
				field.Cells[y][x] = shared.Cell{Type: shared.CellEmpty}
			}
		}
	}

	return field
}
