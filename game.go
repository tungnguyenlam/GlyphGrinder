package main

import "math/rand"

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

// Action describes one turn requested by the player.
type Action uint8

const (
	ActionNone Action = iota
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
)

// GameMap holds the grid.
type GameMap struct {
	Width  int
	Height int
	Tiles  [][]TileType
}

type room struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r room) center() Position {
	return Position{X: r.X + r.Width/2, Y: r.Y + r.Height/2}
}

func (r room) overlaps(other room) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
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

// GameState is the flat root state of the game engine.
type GameState struct {
	Map      GameMap
	Player   Entity
	Entities []Entity
	Log      []string
}

// Step resolves a player action and returns the resulting game state.
func (g GameState) Step(action Action) GameState {
	dx, dy := 0, 0
	switch action {
	case ActionMoveUp:
		dy = -1
	case ActionMoveDown:
		dy = 1
	case ActionMoveLeft:
		dx = -1
	case ActionMoveRight:
		dx = 1
	default:
		return g
	}

	newX := g.Player.Pos.X + dx
	newY := g.Player.Pos.Y + dy
	if newX < 0 || newX >= g.Map.Width || newY < 0 || newY >= g.Map.Height {
		return g
	}
	if g.Map.Tiles[newY][newX] == TileWall {
		return g
	}

	g.Player.Pos = Position{X: newX, Y: newY}
	return g
}

// NewGame generates a dungeon using rng and places the player in its first room.
func NewGame(width, height int, rng *rand.Rand) GameState {
	gameMap, rooms := generateDungeon(width, height, rng)

	return GameState{
		Map: gameMap,
		Player: Entity{
			IsPlayer:  true,
			Pos:       rooms[0].center(),
			Rune:      "@",
			Color:     "#00FF00",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
	}
}

func generateDungeon(width, height int, rng *rand.Rand) (GameMap, []room) {
	if width < 5 || height < 5 {
		panic("dungeon dimensions must be at least 5x5")
	}
	if rng == nil {
		panic("dungeon generation requires an RNG")
	}

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

	const (
		roomAttempts = 80
		minRoomSize  = 3
	)
	maxRoomWidth := min(7, width-2)
	maxRoomHeight := min(5, height-2)
	rooms := make([]room, 0, 8)

	for range roomAttempts {
		candidate := room{
			Width:  minRoomSize + rng.Intn(maxRoomWidth-minRoomSize+1),
			Height: minRoomSize + rng.Intn(maxRoomHeight-minRoomSize+1),
		}
		candidate.X = 1 + rng.Intn(width-candidate.Width-1)
		candidate.Y = 1 + rng.Intn(height-candidate.Height-1)

		overlaps := false
		for _, existing := range rooms {
			if candidate.overlaps(existing) {
				overlaps = true
				break
			}
		}
		if overlaps {
			continue
		}

		carveRoom(&m, candidate)
		if len(rooms) > 0 {
			carveCorridor(&m, rooms[len(rooms)-1].center(), candidate.center(), rng)
		}
		rooms = append(rooms, candidate)
	}

	return m, rooms
}

func carveRoom(m *GameMap, r room) {
	for y := r.Y; y < r.Y+r.Height; y++ {
		for x := r.X; x < r.X+r.Width; x++ {
			m.Tiles[y][x] = TileFloor
		}
	}
}

func carveCorridor(m *GameMap, from, to Position, rng *rand.Rand) {
	if rng.Intn(2) == 0 {
		carveHorizontal(m, from.X, to.X, from.Y)
		carveVertical(m, from.Y, to.Y, to.X)
		return
	}

	carveVertical(m, from.Y, to.Y, from.X)
	carveHorizontal(m, from.X, to.X, to.Y)
}

func carveHorizontal(m *GameMap, fromX, toX, y int) {
	if fromX > toX {
		fromX, toX = toX, fromX
	}
	for x := fromX; x <= toX; x++ {
		m.Tiles[y][x] = TileFloor
	}
}

func carveVertical(m *GameMap, fromY, toY, x int) {
	if fromY > toY {
		fromY, toY = toY, fromY
	}
	for y := fromY; y <= toY; y++ {
		m.Tiles[y][x] = TileFloor
	}
}
