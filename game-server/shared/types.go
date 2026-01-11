package shared

type Pos struct {
	X, Y int
}

// In this program this is an empty position or a none position
func (p Pos) Empty() Pos {
	return Pos{X: -1, Y: -1}
}

func (p Pos) IsEmptyPos() bool {
	return p.X == -1 && p.Y == -1
}

func (p Pos) Add(other Pos) Pos {
	return Pos{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Pos) Sub(other Pos) Pos {
	return Pos{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p Pos) Equal(other Pos) bool {
	return p.X == other.X && p.Y == other.Y
}
