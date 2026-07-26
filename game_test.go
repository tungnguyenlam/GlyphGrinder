package main

import (
	"strings"
	"testing"
)

func TestNewGameBuildsWalledRoom(t *testing.T) {
	g := NewGame(20, 10)

	if got, want := len(g.Map.Tiles), 10; got != want {
		t.Fatalf("rows = %d, want %d", got, want)
	}
	for y, row := range g.Map.Tiles {
		if got, want := len(row), 20; got != want {
			t.Fatalf("row %d width = %d, want %d", y, got, want)
		}
	}

	for x := 0; x < g.Map.Width; x++ {
		if g.Map.Tiles[0][x] != TileWall || g.Map.Tiles[g.Map.Height-1][x] != TileWall {
			t.Errorf("column %d: expected wall on top and bottom border", x)
		}
	}
	for y := 0; y < g.Map.Height; y++ {
		if g.Map.Tiles[y][0] != TileWall || g.Map.Tiles[y][g.Map.Width-1] != TileWall {
			t.Errorf("row %d: expected wall on left and right border", y)
		}
	}

	// Player position in generated dungeon must be floor
	if g.Map.Tiles[g.Player.Pos.Y][g.Player.Pos.X] != TileFloor {
		t.Errorf("player spawn tile at %+v should be floor", g.Player.Pos)
	}
}

func TestNewGamePlacesPlayerOnFloor(t *testing.T) {
	g := NewGame(20, 10)

	p := g.Player
	if !p.IsPlayer {
		t.Error("player entity should have IsPlayer set")
	}
	if g.Map.Tiles[p.Pos.Y][p.Pos.X] != TileFloor {
		t.Errorf("player position %+v must be a floor tile", p.Pos)
	}
	if p.Health != p.MaxHealth {
		t.Errorf("player starts at %d/%d health, want full", p.Health, p.MaxHealth)
	}
}

func TestMonsterHelpers(t *testing.T) {
	gob := NewGoblin(1, Position{X: 2, Y: 3})
	if gob.ID != 1 || gob.Pos.X != 2 || gob.Pos.Y != 3 {
		t.Errorf("unexpected goblin pos/id: %+v", gob)
	}
	if gob.IsPlayer {
		t.Error("goblin should not be player")
	}
	if gob.Rune != "g" || gob.Health != 10 || gob.Damage != 3 {
		t.Errorf("unexpected goblin stats: %+v", gob)
	}

	orc := NewOrc(2, Position{X: 4, Y: 5})
	if orc.ID != 2 || orc.Pos.X != 4 || orc.Pos.Y != 5 {
		t.Errorf("unexpected orc pos/id: %+v", orc)
	}
	if orc.IsPlayer {
		t.Error("orc should not be player")
	}
	if orc.Rune != "o" || orc.Health != 20 || orc.Damage != 6 {
		t.Errorf("unexpected orc stats: %+v", orc)
	}
}

func TestNewGamePopulatesEntities(t *testing.T) {
	g := NewGameWithSeed(20, 10, 12345)

	if len(g.Entities) == 0 {
		t.Fatal("expected g.Entities to be populated with monsters")
	}

	seenIDs := make(map[int]bool)
	seenPos := make(map[Position]bool)

	seenIDs[g.Player.ID] = true
	seenPos[g.Player.Pos] = true

	for i, e := range g.Entities {
		if e.IsPlayer {
			t.Errorf("entity %d in Entities slice has IsPlayer=true", i)
		}
		if e.Health <= 0 || e.MaxHealth <= 0 || e.Damage <= 0 {
			t.Errorf("entity %d has invalid stats: %+v", i, e)
		}
		if e.Rune != "g" && e.Rune != "o" {
			t.Errorf("entity %d has unexpected rune %q", i, e.Rune)
		}
		if g.Map.Tiles[e.Pos.Y][e.Pos.X] != TileFloor {
			t.Errorf("entity %d at %+v is not on a floor tile", i, e.Pos)
		}
		if seenIDs[e.ID] {
			t.Errorf("duplicate entity ID %d", e.ID)
		}
		seenIDs[e.ID] = true
		if seenPos[e.Pos] {
			t.Errorf("duplicate entity position %+v", e.Pos)
		}
		seenPos[e.Pos] = true
	}
}

func TestBumpToAttackDamage(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 2, Y: 2},
			Damage:    10,
			Health:    100,
			MaxHealth: 100,
		},
		Entities: []Entity{
			NewOrc(1, Position{X: 2, Y: 1}), // Orc has 20 Health
		},
	}

	newState := state.Step(ActionMoveUp)

	// Player should not move
	if newState.Player.Pos != (Position{X: 2, Y: 2}) {
		t.Errorf("player moved to %+v, expected to stay at (2,2)", newState.Player.Pos)
	}

	// Orc should still exist but with 10 health
	if len(newState.Entities) != 1 {
		t.Fatalf("expected 1 entity remaining, got %d", len(newState.Entities))
	}
	if got, want := newState.Entities[0].Health, 10; got != want {
		t.Errorf("orc health = %d, want %d", got, want)
	}

	// Orc counter-attacks on its turn
	if got, want := newState.Player.Health, 94; got != want {
		t.Errorf("player health = %d, want %d", got, want)
	}

	// Log should record player hit and Orc counter-hit
	if len(newState.Log) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(newState.Log))
	}
	if got, want := newState.Log[0], "Player hits Orc for 10 damage."; got != want {
		t.Errorf("log[0] = %q, want %q", got, want)
	}
	if got, want := newState.Log[1], "Orc hits Player for 6 damage."; got != want {
		t.Errorf("log[1] = %q, want %q", got, want)
	}
}

func TestBumpToAttackKillsMonster(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 2, Y: 2},
			Damage:    10,
			Health:    100,
			MaxHealth: 100,
		},
		Entities: []Entity{
			NewGoblin(1, Position{X: 2, Y: 1}), // Goblin has 10 Health
		},
	}

	newState := state.Step(ActionMoveUp)

	// Player should not move
	if newState.Player.Pos != (Position{X: 2, Y: 2}) {
		t.Errorf("player moved to %+v, expected to stay at (2,2)", newState.Player.Pos)
	}

	// Goblin should be removed
	if len(newState.Entities) != 0 {
		t.Errorf("expected 0 entities remaining, got %d", len(newState.Entities))
	}

	// Log should record hit and kill
	if len(newState.Log) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(newState.Log))
	}
	if got, want := newState.Log[0], "Player hits Goblin for 10 damage."; got != want {
		t.Errorf("log[0] = %q, want %q", got, want)
	}
	if got, want := newState.Log[1], "Goblin dies."; got != want {
		t.Errorf("log[1] = %q, want %q", got, want)
	}
}

func TestMonsterMovementTowardsPlayer(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 1, Y: 1},
			Damage:    10,
			Health:    100,
			MaxHealth: 100,
		},
		Entities: []Entity{
			NewGoblin(1, Position{X: 1, Y: 3}),
		},
	}

	// Player moves right to (2, 1)
	newState := state.Step(ActionMoveRight)

	if newState.Player.Pos != (Position{X: 2, Y: 1}) {
		t.Errorf("player position = %+v, want (2, 1)", newState.Player.Pos)
	}

	// Goblin at (1, 3) targeting Player at (2, 1): diffX=1, diffY=-2.
	// Primary axis is Y (absY 2 > absX 1), stepY = -1.
	// Goblin should step to (1, 2).
	if len(newState.Entities) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(newState.Entities))
	}
	if got, want := newState.Entities[0].Pos, (Position{X: 1, Y: 2}); got != want {
		t.Errorf("goblin pos = %+v, want %+v", got, want)
	}
}

func TestMonsterAttackPlayerAndDamageLog(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 2, Y: 2},
			Damage:    10,
			Health:    100,
			MaxHealth: 100,
		},
		Entities: []Entity{
			NewOrc(1, Position{X: 2, Y: 3}), // Orc has 20 Health, 6 Damage
		},
	}

	// Player bump attacks Orc at (2, 3)
	newState := state.Step(ActionMoveDown)

	// Player stays at (2, 2)
	if newState.Player.Pos != (Position{X: 2, Y: 2}) {
		t.Errorf("player pos = %+v, want (2, 2)", newState.Player.Pos)
	}

	// Player dealt 10 damage to Orc (Orc HP 20 -> 10)
	if len(newState.Entities) != 1 || newState.Entities[0].Health != 10 {
		t.Errorf("orc health = %d, want 10", newState.Entities[0].Health)
	}

	// Orc at (2, 3) is adjacent to Player at (2, 2), attacks Player for 6 damage
	if got, want := newState.Player.Health, 94; got != want {
		t.Errorf("player health = %d, want %d", got, want)
	}

	// Log should contain player hit and monster hit
	if len(newState.Log) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(newState.Log))
	}
	if got, want := newState.Log[0], "Player hits Orc for 10 damage."; got != want {
		t.Errorf("log[0] = %q, want %q", got, want)
	}
	if got, want := newState.Log[1], "Orc hits Player for 6 damage."; got != want {
		t.Errorf("log[1] = %q, want %q", got, want)
	}
}

func TestMonsterKillsPlayerLog(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 2, Y: 2},
			Damage:    10,
			Health:    5, // Low health
			MaxHealth: 100,
		},
		Entities: []Entity{
			NewOrc(1, Position{X: 2, Y: 3}), // Orc has 6 Damage
		},
	}

	// Player bump attacks Orc
	newState := state.Step(ActionMoveDown)

	if got := newState.Player.Health; got > 0 {
		t.Fatalf("expected player health <= 0, got %d", got)
	}

	if len(newState.Log) != 3 {
		t.Fatalf("expected 3 log entries, got %d: %v", len(newState.Log), newState.Log)
	}
	if got, want := newState.Log[0], "Player hits Orc for 10 damage."; got != want {
		t.Errorf("log[0] = %q, want %q", got, want)
	}
	if got, want := newState.Log[1], "Orc hits Player for 6 damage."; got != want {
		t.Errorf("log[1] = %q, want %q", got, want)
	}
	if got, want := newState.Log[2], "Player dies."; got != want {
		t.Errorf("log[2] = %q, want %q", got, want)
	}
}

func TestComputeFOVOpenRoom(t *testing.T) {
	// In a small 5×5 room with no interior walls, every tile should be
	// visible from the center.
	gm := GameMap{
		Width:  5,
		Height: 5,
		Tiles: [][]TileType{
			{TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall},
		},
	}
	gm.Explored = makeBoolGrid(5, 5)
	gm.ComputeFOV(Position{X: 2, Y: 2}, FOVRadius)

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if !gm.Visible[y][x] {
				t.Errorf("tile (%d,%d) not visible from center of open room", x, y)
			}
			if !gm.Explored[y][x] {
				t.Errorf("tile (%d,%d) not explored after ComputeFOV", x, y)
			}
		}
	}
}

func TestComputeFOVWallBlocksSight(t *testing.T) {
	// A wall at (2,2) should block sight to (2,3) from origin (2,1).
	//
	//  #####
	//  #.@.#   origin at (2,1)
	//  #.#.#   wall at (2,2)
	//  #...#   (2,3) should be hidden
	//  #####
	gm := GameMap{
		Width:  5,
		Height: 5,
		Tiles: [][]TileType{
			{TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileFloor, TileWall, TileFloor, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall},
		},
	}
	gm.Explored = makeBoolGrid(5, 5)
	gm.ComputeFOV(Position{X: 2, Y: 1}, FOVRadius)

	// The wall itself should be visible (you can see a wall).
	if !gm.Visible[2][2] {
		t.Error("wall at (2,2) should be visible")
	}
	// The tile directly behind the wall should not be visible.
	if gm.Visible[3][2] {
		t.Error("tile at (2,3) behind wall should not be visible")
	}
	// Adjacent floor tiles that aren't blocked should be visible.
	if !gm.Visible[1][1] {
		t.Error("floor at (1,1) should be visible")
	}
}

func TestFOVExploredPersistsAcrossSteps(t *testing.T) {
	// After a player moves, previously explored tiles stay explored
	// even when they leave the FOV.
	state := NewGameWithSeed(20, 10, 42)

	// Gather the initial explored set.
	initialExplored := make(map[Position]bool)
	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			if state.Map.Explored[y][x] {
				initialExplored[Position{X: x, Y: y}] = true
			}
		}
	}
	if len(initialExplored) == 0 {
		t.Fatal("no tiles explored in initial state")
	}

	// Step the player (move right, assuming floor).
	state2 := state.Step(ActionMoveRight)

	// Every tile that was explored before must still be explored.
	for pos := range initialExplored {
		if !state2.Map.Explored[pos.Y][pos.X] {
			t.Errorf("tile %+v was explored before step but not after", pos)
		}
	}
}

func TestFOVPlayerOriginAlwaysVisible(t *testing.T) {
	state := NewGameWithSeed(20, 10, 99)
	if !state.Map.Visible[state.Player.Pos.Y][state.Player.Pos.X] {
		t.Error("player origin should always be visible")
	}

	state2 := state.Step(ActionMoveDown)
	if !state2.Map.Visible[state2.Player.Pos.Y][state2.Player.Pos.X] {
		t.Error("player origin should be visible after step")
	}
}

func TestNewGameWithSeedInitializesFOV(t *testing.T) {
	state := NewGameWithSeed(20, 10, 777)

	if state.Map.Visible == nil {
		t.Fatal("Visible grid should be initialized by NewGameWithSeed")
	}
	if state.Map.Explored == nil {
		t.Fatal("Explored grid should be initialized by NewGameWithSeed")
	}

	// Player tile must be visible.
	pp := state.Player.Pos
	if !state.Map.Visible[pp.Y][pp.X] {
		t.Error("player spawn tile should be visible")
	}
	if !state.Map.Explored[pp.Y][pp.X] {
		t.Error("player spawn tile should be explored")
	}
}

func TestGenerateMapPlacesStairsDown(t *testing.T) {
	state := NewGameWithSeed(60, 30, 12345)
	foundStairs := false
	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			if state.Map.Tiles[y][x] == TileStairsDown {
				foundStairs = true
				break
			}
		}
		if foundStairs {
			break
		}
	}
	if !foundStairs {
		t.Error("expected generated map to contain TileStairsDown")
	}
}

func TestDescendingStairs(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	if state.Depth != 1 {
		t.Fatalf("initial depth = %d, want 1", state.Depth)
	}

	// Move player onto stairs tile
	var stairsPos Position
	found := false
	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			if state.Map.Tiles[y][x] == TileStairsDown {
				stairsPos = Position{X: x, Y: y}
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("no stairs found in test map")
	}

	// Teleport player directly onto stairs tile for test
	state.Player.Pos = stairsPos
	state.Player.Health = 85 // non-default HP to test persistence

	// Descend
	nextState := state.Step(ActionDescend)
	if nextState.Depth != 2 {
		t.Errorf("after ActionDescend depth = %d, want 2", nextState.Depth)
	}
	if nextState.Player.Health != 85 {
		t.Errorf("player HP after descend = %d, want 85", nextState.Player.Health)
	}

	lastLog := nextState.Log[len(nextState.Log)-1]
	if lastLog != "You descend to dungeon level 2." {
		t.Errorf("last log = %q, want 'You descend to dungeon level 2.'", lastLog)
	}

	// Verify FOV is re-calculated for new spawn position
	pp := nextState.Player.Pos
	if !nextState.Map.Visible[pp.Y][pp.X] {
		t.Error("player spawn position on level 2 should be visible in FOV")
	}
}

func TestActionDescendWithoutStairs(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	// Make sure player is standing on floor, not stairs
	state.Map.Tiles[state.Player.Pos.Y][state.Player.Pos.X] = TileFloor

	nextState := state.Step(ActionDescend)
	if nextState.Depth != 1 {
		t.Errorf("depth after ActionDescend without stairs = %d, want 1", nextState.Depth)
	}
	lastLog := nextState.Log[len(nextState.Log)-1]
	if lastLog != "There are no stairs here." {
		t.Errorf("log message = %q, want 'There are no stairs here.'", lastLog)
	}
}

func TestItemSpawning(t *testing.T) {
	state := NewGameWithSeed(60, 30, 12345)

	if len(state.Items) == 0 {
		t.Fatal("expected map items to be spawned in NewGameWithSeed")
	}

	for i, item := range state.Items {
		if item.Name == "" {
			t.Errorf("item %d has empty name", i)
		}
		if item.ItemType == ItemPotion && item.HealAmount <= 0 {
			t.Errorf("potion item %d has invalid HealAmount: %d", i, item.HealAmount)
		}
		if item.ItemType == ItemWeapon && item.DamageBonus <= 0 {
			t.Errorf("weapon item %d has invalid DamageBonus: %d", i, item.DamageBonus)
		}
		if state.Map.Tiles[item.Pos.Y][item.Pos.X] != TileFloor {
			t.Errorf("item %d at %+v is not on a floor tile", i, item.Pos)
		}
	}
}

func TestItemPickupAndInventory(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	pot := NewHealthPotion(100, Position{X: 2, Y: 2})
	state.Items = []Item{pot}
	state.Player.Pos = Position{X: 2, Y: 2}
	state.Player.Inventory = nil

	nextState := state.Step(ActionPickup)

	if len(nextState.Items) != 0 {
		t.Errorf("expected 0 items on ground after pickup, got %d", len(nextState.Items))
	}
	if len(nextState.Player.Inventory) != 1 {
		t.Fatalf("expected 1 item in inventory, got %d", len(nextState.Player.Inventory))
	}
	if got, want := nextState.Player.Inventory[0].Name, "Health Potion"; got != want {
		t.Errorf("inventory item name = %q, want %q", got, want)
	}
	lastLog := nextState.Log[len(nextState.Log)-1]
	if lastLog != "You pick up the Health Potion." {
		t.Errorf("log message = %q, want 'You pick up the Health Potion.'", lastLog)
	}
}

func TestActionPickupNothingHere(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	state.Items = nil
	state.Player.Pos = Position{X: 2, Y: 2}

	nextState := state.Step(ActionPickup)
	lastLog := nextState.Log[len(nextState.Log)-1]
	if lastLog != "There is nothing here to pick up." {
		t.Errorf("log message = %q, want 'There is nothing here to pick up.'", lastLog)
	}
}

func TestPotionRestoresHP(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	state.Player.Health = 50
	state.Player.MaxHealth = 100
	state.Player.Inventory = []Item{NewHealthPotion(1, Position{X: -1, Y: -1})}

	nextState := state.Step(ActionUseItem)

	if got, want := nextState.Player.Health, 75; got != want {
		t.Errorf("player HP = %d, want %d", got, want)
	}
	if len(nextState.Player.Inventory) != 0 {
		t.Errorf("expected empty inventory after consuming potion, got %d items", len(nextState.Player.Inventory))
	}
	lastLog := nextState.Log[len(nextState.Log)-1]
	if !strings.Contains(lastLog, "recover 25 HP") {
		t.Errorf("log message = %q, want recovery mention", lastLog)
	}
}

func TestPotionRestoresHPCappedAtMax(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	state.Player.Health = 90
	state.Player.MaxHealth = 100
	state.Player.Inventory = []Item{NewHealthPotion(1, Position{X: -1, Y: -1})}

	nextState := state.Step(ActionUseItem)

	if got, want := nextState.Player.Health, 100; got != want {
		t.Errorf("player HP = %d, want 100 (capped at MaxHealth)", got)
	}
}

func TestWeaponBonusDamageInCombat(t *testing.T) {
	state := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 2, Y: 2},
			Damage:    10,
			Health:    100,
			MaxHealth: 100,
		},
		Items: []Item{
			NewIronDagger(1, Position{X: 2, Y: 2}),
		},
		Entities: []Entity{
			NewOrc(2, Position{X: 2, Y: 1}), // Orc HP = 20
		},
	}

	// Pick up weapon
	state1 := state.Step(ActionPickup)

	if got, want := state1.Player.Damage, 13; got != want {
		t.Fatalf("player damage after pickup weapon = %d, want 13", got)
	}

	// Bump attack Orc
	state2 := state1.Step(ActionMoveUp)

	// Orc took 13 damage (20 -> 7)
	if len(state2.Entities) != 1 || state2.Entities[0].Health != 7 {
		t.Errorf("orc HP after attack = %d, want 7", state2.Entities[0].Health)
	}
}

func TestItemPersistenceAcrossStairs(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)

	var stairsPos Position
	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			if state.Map.Tiles[y][x] == TileStairsDown {
				stairsPos = Position{X: x, Y: y}
				break
			}
		}
	}

	state.Player.Pos = stairsPos
	state.Player.Inventory = []Item{NewHealthPotion(1, Position{X: -1, Y: -1}), NewIronDagger(2, Position{X: -1, Y: -1})}
	state.Player.Damage = 13

	nextState := state.Step(ActionDescend)

	if len(nextState.Player.Inventory) != 2 {
		t.Errorf("inventory len on floor 2 = %d, want 2", len(nextState.Player.Inventory))
	}
	if got, want := nextState.Player.Damage, 13; got != want {
		t.Errorf("player damage on floor 2 = %d, want 13", got)
	}
}

func TestTrollAndArcherHelpers(t *testing.T) {
	tr := NewTroll(1, Position{X: 2, Y: 3})
	if tr.Name != "Troll" || tr.Health != 40 || tr.Damage != 10 || tr.Rune != "T" {
		t.Errorf("unexpected troll stats: %+v", tr)
	}

	arc := NewArcher(2, Position{X: 4, Y: 5})
	if arc.Name != "Archer" || arc.Health != 15 || arc.Damage != 4 || !arc.IsRanged || arc.Range != 5 || arc.Rune != "A" {
		t.Errorf("unexpected archer stats: %+v", arc)
	}
}

func TestArcherRangedAttack(t *testing.T) {
	gm := GameMap{
		Width:  10,
		Height: 5,
		Tiles: [][]TileType{
			{TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall},
		},
	}

	playerPos := Position{X: 1, Y: 2}
	archerPos := Position{X: 4, Y: 2} // distance 3 (within range 5)

	gm.Explored = makeBoolGrid(10, 5)
	gm.ComputeFOV(playerPos, FOVRadius)

	state := GameState{
		Map: gm,
		Player: Entity{
			ID: 0, Name: "Player", IsPlayer: true, Pos: playerPos,
			Health: 100, MaxHealth: 100, Damage: 10,
		},
		Entities: []Entity{
			NewArcher(1, archerPos),
		},
	}

	// Player stays in place or moves down
	nextState := state.Step(ActionMoveDown)

	// Archer should stay at (4,2) and shoot player for 4 damage
	if got, want := nextState.Player.Health, 96; got != want {
		t.Errorf("player HP = %d, want 96 after archer shoot", got)
	}

	if len(nextState.Entities) != 1 || nextState.Entities[0].Pos != archerPos {
		t.Errorf("archer position changed, expected to hold position at %+v, got %+v", archerPos, nextState.Entities[0].Pos)
	}

	lastLog := nextState.Log[len(nextState.Log)-1]
	if !strings.Contains(lastLog, "Archer shoots Player for 4 damage.") {
		t.Errorf("log message = %q, want Archer shoot log", lastLog)
	}
}

func TestMonsterScalingByDepth(t *testing.T) {
	// Generate multiple games at depth 3 and verify Trolls / Archers can spawn
	foundTroll := false
	foundArcher := false

	for seed := int64(1); seed <= 20; seed++ {
		st := NewGameWithSeedAndDepth(60, 30, seed, 3)
		for _, e := range st.Entities {
			if e.Name == "Troll" {
				foundTroll = true
			}
			if e.Name == "Archer" {
				foundArcher = true
			}
		}
		if foundTroll && foundArcher {
			break
		}
	}

	if !foundTroll {
		t.Error("expected at least one Troll to spawn in depth 3 test runs")
	}
	if !foundArcher {
		t.Error("expected at least one Archer to spawn in depth 3 test runs")
	}
}
