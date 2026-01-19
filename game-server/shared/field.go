package shared

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type CellType int

const (
	CellEmpty CellType = iota
	CellWallDestructible
	CellWallIndestructible
	CellBomb
	CellExplosion
	CellPowerUp
	CellExplosionPierce
)

type Cell struct {
	Type CellType

	Bomb *Bomb // only when CellBomb

	PowerUp *PowerUp // only when CellPowerUp

	TicksTillExplosionOver int
}

func (c *Cell) ToJSON() string {
	return `{"type":` + strconv.Itoa(int(c.Type)) + `}`
}

type Field struct {
	Width, Height int
	Cells         [][]Cell
	mutex         sync.RWMutex
}

func (f *Field) cellsToJSON() string {
	var cellsJSON strings.Builder
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			cellsJSON.WriteString(f.Cells[y][x].ToJSON())
			if x < f.Width-1 || y < f.Height-1 {
				cellsJSON.WriteString(",")
			}
		}
	}
	return cellsJSON.String()
}

func (f *Field) ToJSON() string {
	return `{"width":` + strconv.Itoa(f.Width) +
		`,"height":` + strconv.Itoa(f.Height) +
		`,"cells":[` + f.cellsToJSON() + `]}`
}

func (f *Field) isValidPos(x, y int, width, height int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
}

func (f *Field) GetCellAtPos(x, y int) *Cell {

	f.mutex.RLock()
	defer f.mutex.RUnlock()

	if !f.isValidPos(x, y, f.Width, f.Height) {
		return nil
	}

	return &f.Cells[y][x]
}

func (f *Field) SetCellAtPos(x, y int, cellType CellType) *Cell {

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.isValidPos(x, y, f.Width, f.Height) {
		return nil
	}

	f.Cells[y][x].Type = cellType

	return &f.Cells[y][x]
}

func (f *Field) IsWalkable(pos Pos) bool {
	cell := f.GetCellAtPos(pos.X, pos.Y)
	if cell == nil {
		return false
	}
	return cell.Type == CellEmpty || cell.Type == CellPowerUp || cell.Type == CellBomb
}

func (f *Field) PlaceBomb(player *Player) {
	pos := player.Pos
	if !f.isValidPos(pos.X, pos.Y, f.Width, f.Height) {
		return
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	cell := &f.Cells[pos.Y][pos.X]
	if cell.Type != CellEmpty {
		return
	}

	cell.Type = CellBomb
	cell.Bomb = &Bomb{
		Owner:              player,
		Position:           pos,
		Strength:           player.Bomb.Strength,
		TicksTillExplosion: player.Bomb.TicksTillExplosion,
		PierceWalls:        player.Bomb.PierceWalls,
	}
	fmt.Printf("Bomb placed by %d at position (%d, %d)\n", player.ID, pos.X, pos.Y)
}

func GenerateEmptyField(width, height int) *Field {
	field := &Field{
		Width:  width,
		Height: height,
		Cells:  make([][]Cell, height),
	}

	for y := range field.Cells {
		field.Cells[y] = make([]Cell, width)
		for x := range field.Cells[y] {
			field.Cells[y][x] = Cell{Type: CellEmpty}
		}
	}

	return field
}

func (f *Field) GetAllBombs() []*Bomb {
	var bombs []*Bomb
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			cell := f.GetCellAtPos(x, y)
			if cell != nil && cell.Type == CellBomb {
				bombs = append(bombs, cell.Bomb)
			}
		}
	}
	return bombs
}

func (f *Field) RemoveBomb(bomb *Bomb) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	cell := &f.Cells[bomb.Position.Y][bomb.Position.X]
	if cell.Type != CellBomb || cell.Bomb != bomb {
		return
	}

	cell.Type = CellEmpty
	cell.Bomb = nil
	fmt.Printf("Bomb removed by %d at position (%d, %d)\n", bomb.Owner.ID, bomb.Position.X, bomb.Position.Y)
}

// TODO: check surroundings that player is not softlocked right away
func (f *Field) GetRandomSpawnPos() Pos {
	for {
		x := rand.Intn(f.Width)
		y := rand.Intn(f.Height)
		cell := f.GetCellAtPos(x, y)
		if cell != nil && cell.Type == CellEmpty {
			// TODO: check surroundings
			return Pos{X: x, Y: y}
		}
	}
}

var defaultTicksTillExplosionOver = 5

func (f *Field) ExplodeBomb(bomb *Bomb, lockMutex bool) {

	if bomb == nil {
		return
	}

	if lockMutex {
		f.mutex.Lock()
		defer f.mutex.Unlock()
	}

	cell := &f.Cells[bomb.Position.Y][bomb.Position.X]
	if cell.Type != CellBomb || cell.Bomb != bomb {
		return
	}

	cell.Type = CellExplosion
	cell.Bomb = nil
	cell.TicksTillExplosionOver = defaultTicksTillExplosionOver
	var expandExplosion func(x int, y int, dir string, powerLeft int)
	expandExplosion = func(x int, y int, dir string, powerLeft int) {

		// if no power left
		if powerLeft <= 0 {
			return
		}

		// if out of playing field bounds
		if !f.isValidPos(x, y, f.Width, f.Height) {
			return
		}

		cell := &f.Cells[y][x]

		// if cell is a indestructible wall
		if cell.Type == CellWallIndestructible {
			return
		}

		// if cell is a desctructible wall
		if cell.Type == CellWallDestructible {
			rand := rand.Intn(100)
			if rand < 75 {
				if bomb.PierceWalls {
					cell.Type = CellExplosionPierce
				} else {
					cell.Type = CellExplosion
				}
			} else {
				cell.Type = CellPowerUp
				cell.PowerUp = &PowerUp{ID: 10, Type: "Test", Effect: func(p *Player) {
					p.Bomb.PierceWalls = !p.Bomb.PierceWalls
				}}
			}
			cell.TicksTillExplosionOver = defaultTicksTillExplosionOver
			if !bomb.PierceWalls {
				return
			}
		}

		// if explosion meets another bomb
		if cell.Type == CellBomb && cell.Bomb != nil {
			otherBomb := cell.Bomb
			if bomb.PierceWalls {
				cell.Type = CellExplosionPierce
			} else {
				cell.Type = CellExplosion
			}
			cell.Bomb = nil
			cell.TicksTillExplosionOver = defaultTicksTillExplosionOver
			f.ExplodeBomb(otherBomb, false)
			return
		}

		if bomb.PierceWalls {
			cell.Type = CellExplosionPierce
		} else {
			cell.Type = CellExplosion
		}
		cell.TicksTillExplosionOver = defaultTicksTillExplosionOver

		switch dir {
		case "up":
			expandExplosion(x, y-1, "up", powerLeft-1)
		case "down":
			expandExplosion(x, y+1, "down", powerLeft-1)
		case "left":
			expandExplosion(x-1, y, "left", powerLeft-1)
		case "right":
			expandExplosion(x+1, y, "right", powerLeft-1)
		default:
			expandExplosion(x, y-1, "up", powerLeft-1)
			expandExplosion(x, y+1, "down", powerLeft-1)
			expandExplosion(x-1, y, "left", powerLeft-1)
			expandExplosion(x+1, y, "right", powerLeft-1)
		}
	}

	expandExplosion(bomb.Position.X, bomb.Position.Y, "", bomb.Strength)

	if bomb.Owner != nil {
		bomb.Owner.BombCount++
	}
	fmt.Printf("Bomb exploded by %d at position (%d, %d)\n", bomb.Owner.ID, bomb.Position.X, bomb.Position.Y)
}
