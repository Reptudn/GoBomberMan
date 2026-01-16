package shared

import (
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var PlayersMutex sync.RWMutex
var Players = make(map[int]*Player)

func playersAsJSON() string {
	var builder strings.Builder
	builder.WriteString("[")
	index := 0
	for _, player := range Players {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(player.ToJSON())
		index++
	}
	builder.WriteString("]")
	return builder.String()
}

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

	WriteMutex sync.Mutex
}

func (p *Player) ToJSON() string {
	return `{"id":` + strconv.Itoa(p.ID) +
		`,"color":"` + p.Color +
		`","pos":` + p.Pos.ToJSON() +
		`,"alive":` + strconv.FormatBool(p.Alive) + `}`
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
		Strength:           2,
		TicksTillExplosion: 20,
		PierceWalls:        false,
	}
}

type PowerUp struct {
	ID     int
	Type   string
	Effect func(*Player)
}
