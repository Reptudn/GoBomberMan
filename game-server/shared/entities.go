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

var DefaultPlayerMoveDelay = 5

type Player struct {
	ID   int
	Conn *websocket.Conn

	Color string

	Pos                Pos
	NextPos            Pos
	TicksSinceLastMove int

	Alive bool
	Speed float64

	BombCount        int
	MaxBombCount     int
	Bomb             Bomb
	WantsToPlaceBomb bool

	WriteMutex sync.Mutex
}

func (p *Player) CanMove() bool {
	return (float64)(p.TicksSinceLastMove) >= float64(DefaultPlayerMoveDelay)
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
var MinBombStrenght = 2
var MaxBombStrength = 15
var MinBombCount = 1
var MaxBombCount = 10

func (p *Player) GetBasicBomb() *Bomb {
	return &Bomb{
		Owner:              p,
		Position:           p.Pos,
		Strength:           MinBombStrenght,
		TicksTillExplosion: 20,
		PierceWalls:        false,
	}
}
