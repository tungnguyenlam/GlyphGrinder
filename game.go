package main

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

// NewGame initializes a blank game state.
func NewGame(width, height int) GameState {
	m := GameMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]TileType, height),
	}
	for y := range m.Tiles {
		m.Tiles[y] = make([]TileType, width)
		for x := range m.Tiles[y] {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				m.Tiles[y][x] = TileWall
			} else {
				m.Tiles[y][x] = TileFloor // Walkable floor
			}
		}
	}

	return GameState{
		Map: m,
		Player: Entity{
			IsPlayer:  true,
			Pos:       Position{X: width / 2, Y: height / 2},
			Rune:      "@",
			Color:     "#00FF00",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
	}
}
