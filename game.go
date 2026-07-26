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

// GameState is the flat root state of the game engine.
type GameState struct {
	Map      GameMap
	Player   Entity
	Entities []Entity
	Log      []string
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
