package main

import (
	"math/rand"
	"time"
)

// Position represents a 2D coordinate on the grid.
type Position struct {
	X int
	Y int
}

// TileType defines the physical nature of a grid cell.
type TileType uint8

const (
	TileFloor TileType = iota
	TileWall
)

// GameMap holds the grid.
type GameMap struct {
	Width  int
	Height int
	Tiles  [][]TileType
}

// Entity represents any movable actor (Player, Monster).
type Entity struct {
	ID        int
	Pos       Position
	IsPlayer  bool
	Rune      string
	Color     string
	Health    int
	MaxHealth int
	Damage    int
}

// Action represents an intent or action performed by the player on a turn.
type Action uint8

const (
	ActionNone Action = iota
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
)

// GameState is the flat root state of the game engine.
type GameState struct {
	Seed     int64
	Map      GameMap
	Player   Entity
	Entities []Entity
	Log      []string
}

// Step processes a single game turn based on the provided player action.
func (s GameState) Step(act Action) GameState {
	var dx, dy int
	switch act {
	case ActionMoveUp:
		dy = -1
	case ActionMoveDown:
		dy = 1
	case ActionMoveLeft:
		dx = -1
	case ActionMoveRight:
		dx = 1
	case ActionNone:
		return s
	}

	newX := s.Player.Pos.X + dx
	newY := s.Player.Pos.Y + dy

	if newX < 0 || newX >= s.Map.Width || newY < 0 || newY >= s.Map.Height {
		return s
	}
	if s.Map.Tiles[newY][newX] == TileWall {
		return s
	}

	s.Player.Pos.X = newX
	s.Player.Pos.Y = newY
	return s
}

// Rect represents a rectangular region on the grid.
type Rect struct {
	X int
	Y int
	W int
	H int
}

// Center returns the center position of the rectangle.
func (r Rect) Center() Position {
	return Position{
		X: r.X + r.W/2,
		Y: r.Y + r.H/2,
	}
}

// Intersects returns true if r overlaps with or touches another rectangle.
func (r Rect) Intersects(other Rect) bool {
	return r.X-1 <= other.X+other.W && r.X+r.W+1 >= other.X &&
		r.Y-1 <= other.Y+other.H && r.Y+r.H+1 >= other.Y
}

// GenerateMap creates a map with rooms and corridors using the given random source.
func GenerateMap(width, height int, rng *rand.Rand) (GameMap, []Rect) {
	m := GameMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]TileType, height),
	}
	for y := range m.Tiles {
		m.Tiles[y] = make([]TileType, width)
		for x := range m.Tiles[y] {
			m.Tiles[y][x] = TileWall
		}
	}

	minSize := 3
	maxW := width / 3
	if maxW < minSize {
		maxW = minSize
	} else if maxW > 8 {
		maxW = 8
	}
	maxH := height / 3
	if maxH < minSize {
		maxH = minSize
	} else if maxH > 6 {
		maxH = 6
	}

	maxRooms := 10
	var rooms []Rect

	for i := 0; i < maxRooms; i++ {
		w := minSize
		if maxW > minSize {
			w += rng.Intn(maxW - minSize + 1)
		}
		h := minSize
		if maxH > minSize {
			h += rng.Intn(maxH - minSize + 1)
		}

		if width-1-w < 1 || height-1-h < 1 {
			continue
		}
		x := 1 + rng.Intn(width-1-w)
		y := 1 + rng.Intn(height-1-h)

		newRoom := Rect{X: x, Y: y, W: w, H: h}
		overlap := false
		for _, other := range rooms {
			if newRoom.Intersects(other) {
				overlap = true
				break
			}
		}
		if overlap {
			continue
		}

		// Carve room floor
		for ry := newRoom.Y; ry < newRoom.Y+newRoom.H; ry++ {
			for rx := newRoom.X; rx < newRoom.X+newRoom.W; rx++ {
				m.Tiles[ry][rx] = TileFloor
			}
		}

		// Connect to previous room center with corridor
		if len(rooms) > 0 {
			prevCenter := rooms[len(rooms)-1].Center()
			newCenter := newRoom.Center()

			if rng.Intn(2) == 0 {
				carveHCorridor(&m, prevCenter.X, newCenter.X, prevCenter.Y)
				carveVCorridor(&m, prevCenter.Y, newCenter.Y, newCenter.X)
			} else {
				carveVCorridor(&m, prevCenter.Y, newCenter.Y, prevCenter.X)
				carveHCorridor(&m, prevCenter.X, newCenter.X, newCenter.Y)
			}
		}

		rooms = append(rooms, newRoom)
	}

	// Fallback if no rooms generated
	if len(rooms) == 0 {
		w := width - 2
		if w < 1 {
			w = 1
		}
		h := height - 2
		if h < 1 {
			h = 1
		}
		fallbackRoom := Rect{X: 1, Y: 1, W: w, H: h}
		for ry := fallbackRoom.Y; ry < fallbackRoom.Y+fallbackRoom.H; ry++ {
			for rx := fallbackRoom.X; rx < fallbackRoom.X+fallbackRoom.W; rx++ {
				m.Tiles[ry][rx] = TileFloor
			}
		}
		rooms = append(rooms, fallbackRoom)
	}

	return m, rooms
}

func carveHCorridor(m *GameMap, x1, x2, y int) {
	minX, maxX := x1, x2
	if x1 > x2 {
		minX, maxX = x2, x1
	}
	for x := minX; x <= maxX; x++ {
		if x > 0 && x < m.Width-1 && y > 0 && y < m.Height-1 {
			m.Tiles[y][x] = TileFloor
		}
	}
}

func carveVCorridor(m *GameMap, y1, y2, x int) {
	minY, maxY := y1, y2
	if y1 > y2 {
		minY, maxY = y2, y1
	}
	for y := minY; y <= maxY; y++ {
		if x > 0 && x < m.Width-1 && y > 0 && y < m.Height-1 {
			m.Tiles[y][x] = TileFloor
		}
	}
}

// NewGoblin creates a goblin monster entity.
func NewGoblin(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Pos:       pos,
		IsPlayer:  false,
		Rune:      "g",
		Color:     "#00FF66",
		Health:    10,
		MaxHealth: 10,
		Damage:    3,
	}
}

// NewOrc creates an orc monster entity.
func NewOrc(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Pos:       pos,
		IsPlayer:  false,
		Rune:      "o",
		Color:     "#FF5555",
		Health:    20,
		MaxHealth: 20,
		Damage:    6,
	}
}

// NewGameWithSeed initializes game state using a deterministic seed for map generation.
func NewGameWithSeed(width, height int, seed int64) GameState {
	rng := rand.New(rand.NewSource(seed))
	gameMap, rooms := GenerateMap(width, height, rng)
	spawnPos := rooms[0].Center()

	nextID := 1
	var entities []Entity

	for _, r := range rooms {
		numMonsters := rng.Intn(2) + 1
		for m := 0; m < numMonsters; m++ {
			rx := r.X + rng.Intn(r.W)
			ry := r.Y + rng.Intn(r.H)
			pos := Position{X: rx, Y: ry}

			if pos == spawnPos {
				continue
			}
			occupied := false
			for _, e := range entities {
				if e.Pos == pos {
					occupied = true
					break
				}
			}
			if occupied {
				continue
			}

			var monster Entity
			if rng.Intn(10) < 7 {
				monster = NewGoblin(nextID, pos)
			} else {
				monster = NewOrc(nextID, pos)
			}
			nextID++
			entities = append(entities, monster)
		}
	}

	return GameState{
		Seed: seed,
		Map:  gameMap,
		Player: Entity{
			ID:        0,
			IsPlayer:  true,
			Pos:       spawnPos,
			Rune:      "@",
			Color:     "#00FF00",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
		Entities: entities,
	}
}

// NewGame initializes a game state with a random seed.
func NewGame(width, height int) GameState {
	return NewGameWithSeed(width, height, time.Now().UnixNano())
}
