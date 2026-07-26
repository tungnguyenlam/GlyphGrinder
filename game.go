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
