package shared

import "sync"

type CellType int

const (
	CellEmpty CellType = iota
	CellWallDestructible
	CellWallIndestructible
	CellBomb
	CellExplosion
	CellPowerUp
)

type Cell struct {
	Type CellType

	Bomb *Bomb // only when CellBomb

	PowerUp *PowerUp // only when CellPowerUp
}

type Field struct {
	Width, Height int
	Cells         [][]Cell
	mutex         sync.RWMutex
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

	cell := f.GetCellAtPos(pos.X, pos.Y)
	if cell == nil || cell.Type != CellEmpty {
		return
	}

	cell.Type = CellBomb
	cell.Bomb = &Bomb{
		Owner:              player,
		Position:           pos,
		Strength:           player.Bomb.Strength,
		TicksTillExplosion: player.Bomb.Strength,
		PierceWalls:        player.Bomb.PierceWalls,
	}
}

func GenerateEmptyField(width, height int) *Field {
	field := &Field{
		Width:  width,
		Height: height,
		Cells:  make([][]Cell, height),
	}

	for y := range height {
		field.Cells[y] = make([]Cell, width)
		for x := range width {
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
