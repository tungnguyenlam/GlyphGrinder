package main

import (
	"fmt"
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
	TileStairsDown
)

// ItemType defines the classification of an item.
type ItemType uint8

const (
	ItemPotion ItemType = iota
	ItemWeapon
	ItemAmulet
)

// Item represents an object on the map or in an entity's inventory.
type Item struct {
	ID          int
	Name        string
	Pos         Position
	ItemType    ItemType
	HealAmount  int
	DamageBonus int
}

// GameMap holds the grid.
type GameMap struct {
	Width    int
	Height   int
	Tiles    [][]TileType
	Visible  [][]bool // tiles currently in line-of-sight from the player
	Explored [][]bool // tiles that have been seen at least once
}

// Entity represents any movable actor (Player, Monster).
type Entity struct {
	ID        int
	Name      string
	Pos       Position
	IsPlayer  bool
	Rune      string
	Color     string
	Health    int
	MaxHealth int
	Damage    int
	IsRanged  bool
	Range     int
	Inventory []Item
}

// Action represents an intent or action performed by the player on a turn.
type Action uint8

const (
	ActionNone Action = iota
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionDescend
	ActionPickup
	ActionUseItem
	ActionUseItem1
	ActionUseItem2
	ActionUseItem3
	ActionUseItem4
	ActionUseItem5
	ActionUseItem6
	ActionUseItem7
	ActionUseItem8
	ActionUseItem9
)

// GameState is the flat root state of the game engine.
type GameState struct {
	Seed      int64
	Depth     int
	Map       GameMap
	Player    Entity
	Entities  []Entity
	Items     []Item
	Log       []string
	IsVictory bool
}

// Step processes a single game turn based on the provided player action.
func (s GameState) Step(act Action) GameState {
	if act == ActionNone || s.IsVictory || s.Player.Health <= 0 {
		return s
	}

	if act == ActionDescend {
		if s.Map.Tiles[s.Player.Pos.Y][s.Player.Pos.X] == TileStairsDown {
			nextDepth := s.Depth + 1
			if nextDepth <= 0 {
				nextDepth = 2
			}
			nextSeed := s.Seed + int64(s.Depth)*10007
			nextState := NewGameWithSeedAndDepth(s.Map.Width, s.Map.Height, nextSeed, nextDepth)
			nextState.Player.Health = s.Player.Health
			nextState.Player.MaxHealth = s.Player.MaxHealth
			nextState.Player.Damage = s.Player.Damage
			if s.Player.Name != "" {
				nextState.Player.Name = s.Player.Name
			}
			if len(s.Player.Inventory) > 0 {
				nextState.Player.Inventory = append([]Item(nil), s.Player.Inventory...)
			}
			nextState.Log = append([]string(nil), s.Log...)
			nextState.Log = append(nextState.Log, fmt.Sprintf("You descend to dungeon level %d.", nextDepth))
			return nextState
		} else {
			s.Log = append([]string(nil), s.Log...)
			s.Log = append(s.Log, "There are no stairs here.")
			return s
		}
	}

	if act == ActionPickup {
		s.Log = append([]string(nil), s.Log...)
		s.Entities = append([]Entity(nil), s.Entities...)
		s.Items = append([]Item(nil), s.Items...)
		s.Player.Inventory = append([]Item(nil), s.Player.Inventory...)

		itemIdx := -1
		for i, it := range s.Items {
			if it.Pos == s.Player.Pos {
				itemIdx = i
				break
			}
		}

		if itemIdx == -1 {
			s.Log = append(s.Log, "There is nothing here to pick up.")
			return s
		}

		// Deep-copy Explored so value semantics are preserved.
		oldExplored := s.Map.Explored
		s.Map.Explored = makeBoolGrid(s.Map.Width, s.Map.Height)
		for y := range oldExplored {
			copy(s.Map.Explored[y], oldExplored[y])
		}

		picked := s.Items[itemIdx]
		newItems := make([]Item, 0, len(s.Items)-1)
		for i, it := range s.Items {
			if i != itemIdx {
				newItems = append(newItems, it)
			}
		}
		s.Items = newItems

		picked.Pos = Position{X: -1, Y: -1}
		s.Player.Inventory = append(s.Player.Inventory, picked)
		if picked.ItemType == ItemWeapon {
			s.Player.Damage += picked.DamageBonus
		} else if picked.ItemType == ItemAmulet {
			s.IsVictory = true
			s.Log = append(s.Log, "*** YOU HAVE FOUND THE AMULET OF YENDOR! VICTORY IS YOURS! ***")
		}
		s.Log = append(s.Log, fmt.Sprintf("You pick up the %s.", picked.Name))

		s.runMonsterTurns()
		s.Map.ComputeFOV(s.Player.Pos, FOVRadius)
		return s
	}

	if act == ActionUseItem || (act >= ActionUseItem1 && act <= ActionUseItem9) {
		s.Log = append([]string(nil), s.Log...)
		s.Entities = append([]Entity(nil), s.Entities...)
		s.Items = append([]Item(nil), s.Items...)
		s.Player.Inventory = append([]Item(nil), s.Player.Inventory...)

		if len(s.Player.Inventory) == 0 {
			s.Log = append(s.Log, "You have no items to use.")
			return s
		}

		targetIdx := -1
		if act == ActionUseItem {
			for i, it := range s.Player.Inventory {
				if it.ItemType == ItemPotion {
					targetIdx = i
					break
				}
			}
			if targetIdx == -1 {
				targetIdx = 0
			}
		} else {
			targetIdx = int(act - ActionUseItem1)
		}

		if targetIdx < 0 || targetIdx >= len(s.Player.Inventory) {
			s.Log = append(s.Log, "You have no item in that slot.")
			return s
		}

		// Deep-copy Explored grid
		oldExplored := s.Map.Explored
		s.Map.Explored = makeBoolGrid(s.Map.Width, s.Map.Height)
		for y := range oldExplored {
			copy(s.Map.Explored[y], oldExplored[y])
		}

		item := s.Player.Inventory[targetIdx]
		if item.ItemType == ItemPotion {
			heal := item.HealAmount
			if s.Player.Health+heal > s.Player.MaxHealth {
				heal = s.Player.MaxHealth - s.Player.Health
			}
			s.Player.Health += heal

			newInv := make([]Item, 0, len(s.Player.Inventory)-1)
			for i, it := range s.Player.Inventory {
				if i != targetIdx {
					newInv = append(newInv, it)
				}
			}
			s.Player.Inventory = newInv
			s.Log = append(s.Log, fmt.Sprintf("You drink the %s and recover %d HP.", item.Name, heal))
		} else if item.ItemType == ItemWeapon {
			s.Log = append(s.Log, fmt.Sprintf("You are wielding the %s (+%d damage).", item.Name, item.DamageBonus))
		}

		s.runMonsterTurns()
		s.Map.ComputeFOV(s.Player.Pos, FOVRadius)
		return s
	}

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
	}

	newX := s.Player.Pos.X + dx
	newY := s.Player.Pos.Y + dy

	if newX < 0 || newX >= s.Map.Width || newY < 0 || newY >= s.Map.Height {
		return s
	}
	if s.Map.Tiles[newY][newX] == TileWall {
		return s
	}

	s.Entities = append([]Entity(nil), s.Entities...)
	s.Log = append([]string(nil), s.Log...)
	s.Items = append([]Item(nil), s.Items...)
	if len(s.Player.Inventory) > 0 {
		s.Player.Inventory = append([]Item(nil), s.Player.Inventory...)
	}

	// Deep-copy Explored so value semantics are preserved.
	oldExplored := s.Map.Explored
	s.Map.Explored = makeBoolGrid(s.Map.Width, s.Map.Height)
	for y := range oldExplored {
		copy(s.Map.Explored[y], oldExplored[y])
	}

	targetPos := Position{X: newX, Y: newY}
	targetIdx := -1
	for i, e := range s.Entities {
		if e.Pos == targetPos {
			targetIdx = i
			break
		}
	}

	if targetIdx != -1 {
		target := s.Entities[targetIdx]
		target.Health -= s.Player.Damage

		playerName := s.Player.Name
		if playerName == "" {
			playerName = "Player"
		}
		monsterName := target.Name
		if monsterName == "" {
			monsterName = "Monster"
		}

		s.Log = append(s.Log, fmt.Sprintf("%s hits %s for %d damage.", playerName, monsterName, s.Player.Damage))

		newEntities := make([]Entity, 0, len(s.Entities))
		for i, e := range s.Entities {
			if i == targetIdx {
				if target.Health > 0 {
					newEntities = append(newEntities, target)
				}
			} else {
				newEntities = append(newEntities, e)
			}
		}

		if target.Health <= 0 {
			s.Log = append(s.Log, fmt.Sprintf("%s dies.", monsterName))
		}

		s.Entities = newEntities
	} else {
		s.Player.Pos = targetPos
		if s.Map.Tiles[newY][newX] == TileStairsDown {
			s.Log = append(s.Log, "You see a staircase leading down here. (Press > or Enter to descend)")
		} else {
			for _, it := range s.Items {
				if it.Pos == targetPos {
					s.Log = append(s.Log, fmt.Sprintf("You see a %s here. (Press g or , to pick up)", it.Name))
					break
				}
			}
		}
	}

	s.runMonsterTurns()

	// Recalculate FOV from the player's new position.
	s.Map.ComputeFOV(s.Player.Pos, FOVRadius)

	return s
}

func (s *GameState) runMonsterTurns() {
	playerName := s.Player.Name
	if playerName == "" {
		playerName = "Player"
	}

	for i := 0; i < len(s.Entities); i++ {
		m := s.Entities[i]
		mName := m.Name
		if mName == "" {
			mName = "Monster"
		}

		pDiffX := s.Player.Pos.X - m.Pos.X
		pDiffY := s.Player.Pos.Y - m.Pos.Y
		absX := abs(pDiffX)
		absY := abs(pDiffY)
		distMax := absX
		if absY > distMax {
			distMax = absY
		}

		if absX+absY == 1 {
			prevHealth := s.Player.Health
			s.Player.Health -= m.Damage
			s.Log = append(s.Log, fmt.Sprintf("%s hits %s for %d damage.", mName, playerName, m.Damage))
			if prevHealth > 0 && s.Player.Health <= 0 {
				s.Log = append(s.Log, fmt.Sprintf("%s dies.", playerName))
			}
		} else if m.IsRanged && distMax <= m.Range && s.Map.Visible != nil && s.Map.Visible[m.Pos.Y][m.Pos.X] {
			prevHealth := s.Player.Health
			s.Player.Health -= m.Damage
			s.Log = append(s.Log, fmt.Sprintf("%s shoots %s for %d damage.", mName, playerName, m.Damage))
			if prevHealth > 0 && s.Player.Health <= 0 {
				s.Log = append(s.Log, fmt.Sprintf("%s dies.", playerName))
			}
		} else {
			type stepDir struct{ dx, dy int }
			var primary, secondary stepDir

			if absX >= absY {
				stepX := 0
				if pDiffX > 0 {
					stepX = 1
				} else if pDiffX < 0 {
					stepX = -1
				}
				stepY := 0
				if pDiffY > 0 {
					stepY = 1
				} else if pDiffY < 0 {
					stepY = -1
				}
				primary = stepDir{dx: stepX, dy: 0}
				secondary = stepDir{dx: 0, dy: stepY}
			} else {
				stepX := 0
				if pDiffX > 0 {
					stepX = 1
				} else if pDiffX < 0 {
					stepX = -1
				}
				stepY := 0
				if pDiffY > 0 {
					stepY = 1
				} else if pDiffY < 0 {
					stepY = -1
				}
				primary = stepDir{dx: 0, dy: stepY}
				secondary = stepDir{dx: stepX, dy: 0}
			}

			canOccupy := func(pos Position) bool {
				if pos.X < 0 || pos.X >= s.Map.Width || pos.Y < 0 || pos.Y >= s.Map.Height {
					return false
				}
				if s.Map.Tiles[pos.Y][pos.X] != TileFloor {
					return false
				}
				if pos == s.Player.Pos {
					return false
				}
				for j, other := range s.Entities {
					if j != i && other.Pos == pos {
						return false
					}
				}
				return true
			}

			candPrimary := Position{X: m.Pos.X + primary.dx, Y: m.Pos.Y + primary.dy}
			if primary.dx != 0 || primary.dy != 0 {
				if canOccupy(candPrimary) {
					s.Entities[i].Pos = candPrimary
					continue
				}
			}

			if secondary.dx != 0 || secondary.dy != 0 {
				candSecondary := Position{X: m.Pos.X + secondary.dx, Y: m.Pos.Y + secondary.dy}
				if canOccupy(candSecondary) {
					s.Entities[i].Pos = candSecondary
				}
			}
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// FOVRadius is the default line-of-sight radius around the player.
const FOVRadius = 6

// makeBoolGrid allocates a height×width bool grid initialised to false.
func makeBoolGrid(width, height int) [][]bool {
	grid := make([][]bool, height)
	for y := range grid {
		grid[y] = make([]bool, width)
	}
	return grid
}

// ComputeFOV fills m.Visible using recursive shadowcasting from origin out to
// radius. It also marks every newly visible tile as Explored.
func (m *GameMap) ComputeFOV(origin Position, radius int) {
	m.Visible = makeBoolGrid(m.Width, m.Height)
	if m.Explored == nil {
		m.Explored = makeBoolGrid(m.Width, m.Height)
	}

	// The origin is always visible.
	m.Visible[origin.Y][origin.X] = true
	m.Explored[origin.Y][origin.X] = true

	// Multiplier tables for the eight octants.
	// Each row is {xx, xy, yx, yy} so that:
	//   mapX = origin.X + col*xx + row*xy
	//   mapY = origin.Y + col*yx + row*yy
	octants := [8][4]int{
		{1, 0, 0, 1},
		{0, 1, 1, 0},
		{0, -1, 1, 0},
		{-1, 0, 0, 1},
		{-1, 0, 0, -1},
		{0, -1, -1, 0},
		{0, 1, -1, 0},
		{1, 0, 0, -1},
	}

	for _, mult := range octants {
		m.castLight(origin, radius, 1, 1.0, 0.0, mult)
	}
}

// castLight is the recursive shadowcasting worker for one octant defined by
// the multiplier set. startSlope and endSlope bound the visible arc (slopes
// are in the range [0,1]).
func (m *GameMap) castLight(origin Position, radius, row int, startSlope, endSlope float64, mult [4]int) {
	if startSlope < endSlope {
		return
	}

	nextStartSlope := startSlope
	for j := row; j <= radius; j++ {
		blocked := false
		for col := j; col >= 0; col-- {
			// Map (col, j) through the octant multipliers.
			mapX := origin.X + col*mult[0] + j*mult[1]
			mapY := origin.Y + col*mult[2] + j*mult[3]

			if mapX < 0 || mapX >= m.Width || mapY < 0 || mapY >= m.Height {
				continue
			}

			// Slopes for this cell's left and right edges.
			leftSlope := (float64(col) + 0.5) / (float64(j) - 0.5)
			rightSlope := (float64(col) - 0.5) / (float64(j) + 0.5)

			if startSlope < rightSlope {
				continue
			}
			if endSlope > leftSlope {
				break
			}

			// Mark visible if within circular radius.
			if col*col+j*j <= radius*radius {
				m.Visible[mapY][mapX] = true
				m.Explored[mapY][mapX] = true
			}

			if blocked {
				if m.Tiles[mapY][mapX] == TileWall {
					nextStartSlope = rightSlope
				} else {
					blocked = false
					startSlope = nextStartSlope
				}
			} else if m.Tiles[mapY][mapX] == TileWall {
				blocked = true
				m.castLight(origin, radius, j+1, startSlope, leftSlope, mult)
				nextStartSlope = rightSlope
			}
		}
		if blocked {
			break
		}
	}
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
	} else if maxW > 10 {
		maxW = 10
	}
	maxH := height / 3
	if maxH < minSize {
		maxH = minSize
	} else if maxH > 8 {
		maxH = 8
	}

	// Scale max room attempts proportionally to map area.
	// Baseline: 10 rooms for a 200-tile map (20×10).
	maxRooms := (width * height) / 20
	if maxRooms < 10 {
		maxRooms = 10
	}
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

	// Place down stairs in room furthest from spawn room (rooms[0])
	maxDist := -1
	stairsRoomIdx := 0
	spawnPos := rooms[0].Center()
	for i, r := range rooms {
		c := r.Center()
		dist := abs(c.X-spawnPos.X) + abs(c.Y-spawnPos.Y)
		if dist > maxDist {
			maxDist = dist
			stairsRoomIdx = i
		}
	}
	stairsPos := rooms[stairsRoomIdx].Center()
	m.Tiles[stairsPos.Y][stairsPos.X] = TileStairsDown

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
		Name:      "Goblin",
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
		Name:      "Orc",
		Pos:       pos,
		IsPlayer:  false,
		Rune:      "o",
		Color:     "#FF5555",
		Health:    20,
		MaxHealth: 20,
		Damage:    6,
	}
}

// NewTroll creates a troll monster entity.
func NewTroll(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Name:      "Troll",
		Pos:       pos,
		IsPlayer:  false,
		Rune:      "T",
		Color:     "#D32F2F",
		Health:    40,
		MaxHealth: 40,
		Damage:    10,
	}
}

// NewArcher creates an archer monster entity.
func NewArcher(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Name:      "Archer",
		Pos:       pos,
		IsPlayer:  false,
		Rune:      "A",
		Color:     "#FF9800",
		Health:    15,
		MaxHealth: 15,
		Damage:    4,
		IsRanged:  true,
		Range:     5,
	}
}

// NewHealthPotion creates a health potion item.
func NewHealthPotion(id int, pos Position) Item {
	return Item{
		ID:         id,
		Name:       "Health Potion",
		Pos:        pos,
		ItemType:   ItemPotion,
		HealAmount: 25,
	}
}

// NewIronDagger creates an iron dagger weapon item.
func NewIronDagger(id int, pos Position) Item {
	return Item{
		ID:          id,
		Name:        "Iron Dagger",
		Pos:         pos,
		ItemType:    ItemWeapon,
		DamageBonus: 3,
	}
}

// NewAmuletOfYendor creates the Amulet of Yendor goal item.
func NewAmuletOfYendor(id int, pos Position) Item {
	return Item{
		ID:       id,
		Name:     "Amulet of Yendor",
		Pos:      pos,
		ItemType: ItemAmulet,
	}
}

// NewGameWithSeedAndDepth initializes game state with a deterministic seed and floor depth.
func NewGameWithSeedAndDepth(width, height int, seed int64, depth int) GameState {
	if depth < 1 {
		depth = 1
	}
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

			if pos == spawnPos || gameMap.Tiles[pos.Y][pos.X] == TileStairsDown {
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
			roll := rng.Intn(100)
			if depth == 1 {
				if roll < 70 {
					monster = NewGoblin(nextID, pos)
				} else {
					monster = NewOrc(nextID, pos)
				}
			} else if depth == 2 {
				if roll < 40 {
					monster = NewGoblin(nextID, pos)
				} else if roll < 80 {
					monster = NewOrc(nextID, pos)
				} else {
					monster = NewTroll(nextID, pos)
				}
			} else { // depth >= 3
				if roll < 25 {
					monster = NewGoblin(nextID, pos)
				} else if roll < 60 {
					monster = NewOrc(nextID, pos)
				} else if roll < 80 {
					monster = NewTroll(nextID, pos)
				} else {
					monster = NewArcher(nextID, pos)
				}
			}
			nextID++
			entities = append(entities, monster)
		}
	}

	nextItemID := 1
	var items []Item
	var logMsg []string

	if depth >= 5 {
		// On Depth 5, remove down stairs and spawn Amulet of Yendor
		for y := 0; y < gameMap.Height; y++ {
			for x := 0; x < gameMap.Width; x++ {
				if gameMap.Tiles[y][x] == TileStairsDown {
					gameMap.Tiles[y][x] = TileFloor
					amulet := NewAmuletOfYendor(nextItemID, Position{X: x, Y: y})
					nextItemID++
					items = append(items, amulet)
					break
				}
			}
		}
		logMsg = append(logMsg, "You feel a mystic power... The Amulet of Yendor is on this floor!")
	}

	for _, r := range rooms {
		if rng.Intn(10) < 5 {
			rx := r.X + rng.Intn(r.W)
			ry := r.Y + rng.Intn(r.H)
			pos := Position{X: rx, Y: ry}

			if pos != spawnPos && gameMap.Tiles[pos.Y][pos.X] != TileStairsDown {
				occupied := false
				for _, e := range entities {
					if e.Pos == pos {
						occupied = true
						break
					}
				}
				for _, it := range items {
					if it.Pos == pos {
						occupied = true
						break
					}
				}
				if !occupied {
					var item Item
					if rng.Intn(10) < 7 {
						item = NewHealthPotion(nextItemID, pos)
					} else {
						item = NewIronDagger(nextItemID, pos)
					}
					nextItemID++
					items = append(items, item)
				}
			}
		}
	}

	gameMap.Explored = makeBoolGrid(width, height)
	gameMap.ComputeFOV(spawnPos, FOVRadius)

	return GameState{
		Seed:  seed,
		Depth: depth,
		Map:   gameMap,
		Log:   logMsg,
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       spawnPos,
			Rune:      "@",
			Color:     "#00FF00",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
		Entities: entities,
		Items:    items,
	}
}

// NewGameWithSeed initializes game state using a deterministic seed for map generation.
func NewGameWithSeed(width, height int, seed int64) GameState {
	return NewGameWithSeedAndDepth(width, height, seed, 1)
}

// NewGame initializes a game state with a random seed.
func NewGame(width, height int) GameState {
	return NewGameWithSeedAndDepth(width, height, time.Now().UnixNano(), 1)
}
