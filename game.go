package main

import (
	"fmt"
	"math/rand"
	"sort"
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
	Width    int
	Height   int
	Tiles    [][]TileType
	Visible  [][]bool
	Explored [][]bool
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
	GameOver bool
}

// Step resolves a player action and returns the resulting game state.
func (g GameState) Step(action Action) GameState {
	if g.GameOver {
		return g
	}

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

	previousPlayerPos := g.Player.Pos
	g = g.resolvePlayerMove(dx, dy)
	if g.Player.Pos != previousPlayerPos {
		g = g.refreshVisibility()
	}
	return g.resolveMonsterTurns()
}

func (g GameState) resolvePlayerMove(dx, dy int) GameState {
	newX := g.Player.Pos.X + dx
	newY := g.Player.Pos.Y + dy
	if newX < 0 || newX >= g.Map.Width || newY < 0 || newY >= g.Map.Height {
		return g
	}
	if g.Map.Tiles[newY][newX] == TileWall {
		return g
	}
	for i, entity := range g.Entities {
		if entity.Pos == (Position{X: newX, Y: newY}) {
			return g.attackMonster(i)
		}
	}

	g.Player.Pos = Position{X: newX, Y: newY}
	return g
}

func (g GameState) resolveMonsterTurns() GameState {
	g.Entities = append([]Entity(nil), g.Entities...)
	occupied := make(map[Position]struct{}, len(g.Entities))
	turnOrder := make([]int, 0, len(g.Entities))
	for i, entity := range g.Entities {
		if entity.Health <= 0 {
			continue
		}
		occupied[entity.Pos] = struct{}{}
		turnOrder = append(turnOrder, i)
	}
	sort.Slice(turnOrder, func(i, j int) bool {
		return g.Entities[turnOrder[i]].ID < g.Entities[turnOrder[j]].ID
	})

	for _, index := range turnOrder {
		monster := &g.Entities[index]
		if manhattanDistance(monster.Pos, g.Player.Pos) == 1 {
			g.Player.Health = max(0, g.Player.Health-monster.Damage)
			g.Log = appendLog(g.Log, fmt.Sprintf("Monster %d hits you for %d damage.", monster.ID, monster.Damage))
			if g.Player.Health == 0 {
				g.GameOver = true
				g.Log = appendLog(g.Log, "You die.")
				break
			}
			continue
		}

		for _, delta := range approachSteps(monster.Pos, g.Player.Pos) {
			next := Position{X: monster.Pos.X + delta.X, Y: monster.Pos.Y + delta.Y}
			if next.X < 0 || next.X >= g.Map.Width || next.Y < 0 || next.Y >= g.Map.Height {
				continue
			}
			if next == g.Player.Pos || g.Map.Tiles[next.Y][next.X] == TileWall {
				continue
			}
			if _, blocked := occupied[next]; blocked {
				continue
			}

			delete(occupied, monster.Pos)
			monster.Pos = next
			occupied[next] = struct{}{}
			break
		}
	}

	return g
}

func approachSteps(from, to Position) []Position {
	dx := to.X - from.X
	dy := to.Y - from.Y
	horizontal := Position{X: sign(dx)}
	vertical := Position{Y: sign(dy)}

	steps := make([]Position, 0, 2)
	if abs(dx) >= abs(dy) {
		if horizontal.X != 0 {
			steps = append(steps, horizontal)
		}
		if vertical.Y != 0 {
			steps = append(steps, vertical)
		}
		return steps
	}
	if vertical.Y != 0 {
		steps = append(steps, vertical)
	}
	if horizontal.X != 0 {
		steps = append(steps, horizontal)
	}
	return steps
}

func manhattanDistance(a, b Position) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func sign(value int) int {
	switch {
	case value < 0:
		return -1
	case value > 0:
		return 1
	default:
		return 0
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (g GameState) attackMonster(index int) GameState {
	entities := append([]Entity(nil), g.Entities...)
	entities[index].Health -= g.Player.Damage
	if entities[index].Health > 0 {
		g.Entities = entities
		g.Log = appendLog(g.Log, fmt.Sprintf("You hit monster %d for %d damage.", entities[index].ID, g.Player.Damage))
		return g
	}

	defeatedID := entities[index].ID
	g.Entities = append(entities[:index], entities[index+1:]...)
	g.Log = appendLog(g.Log, fmt.Sprintf("You kill monster %d.", defeatedID))
	return g
}

func appendLog(log []string, message string) []string {
	next := make([]string, len(log), len(log)+1)
	copy(next, log)
	return append(next, message)
}

const visibilityRadius = 6

func (g GameState) refreshVisibility() GameState {
	visible := makeBoolGrid(g.Map.Width, g.Map.Height)
	explored := cloneBoolGrid(g.Map.Explored, g.Map.Width, g.Map.Height)
	radiusSquared := visibilityRadius * visibilityRadius

	for y := max(0, g.Player.Pos.Y-visibilityRadius); y <= min(g.Map.Height-1, g.Player.Pos.Y+visibilityRadius); y++ {
		for x := max(0, g.Player.Pos.X-visibilityRadius); x <= min(g.Map.Width-1, g.Player.Pos.X+visibilityRadius); x++ {
			dx := x - g.Player.Pos.X
			dy := y - g.Player.Pos.Y
			if dx*dx+dy*dy > radiusSquared {
				continue
			}
			if !hasLineOfSight(g.Map, g.Player.Pos, Position{X: x, Y: y}) {
				continue
			}
			visible[y][x] = true
			explored[y][x] = true
		}
	}

	g.Map.Visible = visible
	g.Map.Explored = explored
	return g
}

func makeBoolGrid(width, height int) [][]bool {
	grid := make([][]bool, height)
	for y := range grid {
		grid[y] = make([]bool, width)
	}
	return grid
}

func cloneBoolGrid(source [][]bool, width, height int) [][]bool {
	clone := makeBoolGrid(width, height)
	for y := 0; y < min(height, len(source)); y++ {
		copy(clone[y], source[y])
	}
	return clone
}

func hasLineOfSight(m GameMap, from, to Position) bool {
	x, y := from.X, from.Y
	dx := abs(to.X - from.X)
	dy := abs(to.Y - from.Y)
	sx := sign(to.X - from.X)
	sy := sign(to.Y - from.Y)
	err := dx - dy

	for {
		if x == to.X && y == to.Y {
			return true
		}
		if (Position{X: x, Y: y}) != from && m.Tiles[y][x] == TileWall {
			return false
		}

		twiceError := 2 * err
		if twiceError > -dy {
			err -= dy
			x += sx
		}
		if twiceError < dx {
			err += dx
			y += sy
		}
	}
}

// NewGame generates a dungeon using rng and places the player in its first room.
func NewGame(width, height int, rng *rand.Rand) GameState {
	gameMap, rooms := generateDungeon(width, height, rng)
	playerPos := rooms[0].center()

	g := GameState{
		Map:      gameMap,
		Entities: placeMonsters(gameMap, playerPos, rng),
		Player: Entity{
			IsPlayer:  true,
			Pos:       playerPos,
			Rune:      "@",
			Color:     "#00FF00",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
	}
	return g.refreshVisibility()
}

func placeMonsters(m GameMap, playerPos Position, rng *rand.Rand) []Entity {
	openFloors := make([]Position, 0, m.Width*m.Height)
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			pos := Position{X: x, Y: y}
			if m.Tiles[y][x] == TileFloor && pos != playerPos {
				openFloors = append(openFloors, pos)
			}
		}
	}
	rng.Shuffle(len(openFloors), func(i, j int) {
		openFloors[i], openFloors[j] = openFloors[j], openFloors[i]
	})

	monsterCount := min(3, len(openFloors))
	monsters := make([]Entity, monsterCount)
	for i := range monsters {
		monsters[i] = Entity{
			ID:        i + 1,
			Pos:       openFloors[i],
			Rune:      "g",
			Color:     "#FF5F5F",
			Health:    20,
			MaxHealth: 20,
			Damage:    5,
		}
	}
	return monsters
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
