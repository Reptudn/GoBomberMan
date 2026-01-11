package shared

import (
	"sync"

	"github.com/gorilla/websocket"
)

var PlayersMutex sync.RWMutex
var Players = make(map[int]*Player)

type Player struct {
	ID   int
	Conn *websocket.Conn

	Color string

	Pos     Pos
	NextPos Pos

	Alive bool
	Speed float64

	BombCount        int
	MaxBombCount     int
	Bomb             Bomb
	WantsToPlaceBomb bool
}

type Bomb struct {
	Owner              *Player
	Position           Pos
	Strength           int
	TicksTillExplosion int
	PierceWalls        bool
}

// this is the bomb given to the player on start
func (p *Player) GetBasicBomb() *Bomb {
	return &Bomb{
		Owner:              p,
		Position:           p.Pos,
		Strength:           1,
		TicksTillExplosion: 10,
		PierceWalls:        false,
	}
}

func (b *Bomb) Explode() {
	if b.Owner != nil && b.Owner.BombCount < b.Owner.MaxBombCount {
		b.Owner.BombCount++
	}
}

type PowerUp struct {
	ID     int
	Type   string
	Effect func(*Player)
}
