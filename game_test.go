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

	// Log should record hit, kill, and XP gain
	if len(newState.Log) != 3 {
		t.Fatalf("expected 3 log entries, got %d: %v", len(newState.Log), newState.Log)
	}
	if got, want := newState.Log[0], "Player hits Goblin for 10 damage."; got != want {
		t.Errorf("log[0] = %q, want %q", got, want)
	}
	if got, want := newState.Log[1], "Goblin dies."; got != want {
		t.Errorf("log[1] = %q, want %q", got, want)
	}
	if got, want := newState.Log[2], "You gain 10 XP."; got != want {
		t.Errorf("log[2] = %q, want %q", got, want)
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

func TestDepth5SpawnsAmuletWithoutStairs(t *testing.T) {
	st := NewGameWithSeedAndDepth(60, 30, 12345, 5)

	// Verify no stairs down exist on depth 5
	for y := 0; y < st.Map.Height; y++ {
		for x := 0; x < st.Map.Width; x++ {
			if st.Map.Tiles[y][x] == TileStairsDown {
				t.Error("expected no down stairs on depth 5")
			}
		}
	}

	// Verify Amulet of Yendor is spawned
	foundAmulet := false
	for _, it := range st.Items {
		if it.ItemType == ItemAmulet && it.Name == "Amulet of Yendor" {
			foundAmulet = true
			break
		}
	}

	if !foundAmulet {
		t.Error("expected Amulet of Yendor to spawn on depth 5")
	}

	if len(st.Log) == 0 || !strings.Contains(st.Log[0], "Amulet of Yendor is on this floor!") {
		t.Errorf("expected level 5 start log message, got %v", st.Log)
	}
}

func TestAmuletPickupTriggersVictory(t *testing.T) {
	st := NewGameWithSeedAndDepth(20, 10, 12345, 5)
	amuletPos := Position{X: 2, Y: 2}
	st.Items = []Item{NewAmuletOfYendor(100, amuletPos)}
	st.Player.Pos = amuletPos

	nextState := st.Step(ActionPickup)

	if !nextState.IsVictory {
		t.Error("expected IsVictory=true after picking up Amulet of Yendor")
	}

	fullLog := strings.Join(nextState.Log, "\n")
	if !strings.Contains(fullLog, "VICTORY IS YOURS!") {
		t.Errorf("expected victory message in log, got %q", fullLog)
	}

	// Further movement steps are ignored when IsVictory is true
	stepAfterVictory := nextState.Step(ActionMoveUp)
	if stepAfterVictory.Player.Pos != nextState.Player.Pos {
		t.Error("expected no player movement after victory state is reached")
	}
}

func TestFireballScrollAoE(t *testing.T) {
	gm := GameMap{
		Width:  15,
		Height: 15,
		Tiles:  make([][]TileType, 15),
	}
	for y := 0; y < 15; y++ {
		gm.Tiles[y] = make([]TileType, 15)
		for x := 0; x < 15; x++ {
			gm.Tiles[y][x] = TileFloor
		}
	}
	gm.Explored = makeBoolGrid(15, 15)

	playerPos := Position{X: 5, Y: 5}
	goblinPos := Position{X: 6, Y: 5}  // inside radius 3
	trollPos := Position{X: 12, Y: 12} // outside radius 3

	st := GameState{
		Map: gm,
		Player: Entity{
			ID: 0, Name: "Player", IsPlayer: true, Pos: playerPos,
			Health: 100, MaxHealth: 100,
			Inventory: []Item{
				NewFireballScroll(1, Position{X: -1, Y: -1}),
			},
		},
		Entities: []Entity{
			NewGoblin(1, goblinPos), // HP 10
			NewTroll(2, trollPos),   // HP 40
		},
	}

	nextState := st.Step(ActionUseItem)

	// Goblin at (6,5) has 10 HP and takes 15 damage from fireball, dying
	foundGoblin := false
	foundTroll := false
	for _, e := range nextState.Entities {
		if e.ID == 1 {
			foundGoblin = true
		}
		if e.ID == 2 {
			foundTroll = true
		}
	}

	if foundGoblin {
		t.Error("expected goblin in radius 3 to die from 15 damage Fireball")
	}
	if !foundTroll {
		t.Error("expected troll outside radius 3 to survive Fireball")
	}

	fullLog := strings.Join(nextState.Log, "\n")
	if !strings.Contains(fullLog, "Scroll of Fireball") {
		t.Errorf("expected Fireball scroll log, got %q", fullLog)
	}
}

func TestTeleportScroll(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	origPos := st.Player.Pos
	st.Player.Inventory = []Item{NewTeleportScroll(1, Position{X: -1, Y: -1})}

	nextState := st.Step(ActionUseItem)

	if nextState.Player.Pos == origPos {
		t.Errorf("expected player position to change after reading Teleport scroll, stayed at %+v", origPos)
	}

	tileAtNewPos := nextState.Map.Tiles[nextState.Player.Pos.Y][nextState.Player.Pos.X]
	if tileAtNewPos != TileFloor {
		t.Errorf("expected player to teleport onto a floor tile, got tile %v", tileAtNewPos)
	}

	fullLog := strings.Join(nextState.Log, "\n")
	if !strings.Contains(fullLog, "Scroll of Teleportation") {
		t.Errorf("expected Teleportation log message, got %q", fullLog)
	}
}

func TestDoorOpeningAndFOV(t *testing.T) {
	gm := GameMap{
		Width:  5,
		Height: 5,
		Tiles: [][]TileType{
			{TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileFloor, TileDoorClosed, TileFloor, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall},
			{TileWall, TileWall, TileWall, TileWall, TileWall},
		},
	}
	gm.Explored = makeBoolGrid(5, 5)
	playerPos := Position{X: 1, Y: 1}
	gm.ComputeFOV(playerPos, FOVRadius)

	// Tile (3,1) behind closed door (2,1) is not visible initially
	if gm.Visible[1][3] {
		t.Error("tile behind closed door should not be visible initially")
	}

	st := GameState{
		Map: gm,
		Player: Entity{
			ID: 0, Name: "Player", IsPlayer: true, Pos: playerPos,
			Health: 100, MaxHealth: 100,
		},
	}

	// Bump into closed door at (2,1)
	stOpened := st.Step(ActionMoveRight)

	if stOpened.Map.Tiles[1][2] != TileDoorOpen {
		t.Errorf("door at (2,1) = %v, want TileDoorOpen", stOpened.Map.Tiles[1][2])
	}
	if stOpened.Player.Pos != playerPos {
		t.Errorf("player position = %+v, want %+v (door opening takes 1 turn without moving)", stOpened.Player.Pos, playerPos)
	}

	// FOV now reveals tile (3,1) through opened door
	if !stOpened.Map.Visible[1][3] {
		t.Error("tile behind opened door should be visible after opening door")
	}

	// Next step moves through opened door onto (2,1)
	stMoved := stOpened.Step(ActionMoveRight)
	if stMoved.Player.Pos != (Position{X: 2, Y: 1}) {
		t.Errorf("player position after moving through open door = %+v, want (2,1)", stMoved.Player.Pos)
	}
}

func TestPoisonAndRegenStatusEffects(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	st.Player.Health = 50
	st.Player.MaxHealth = 100
	st.Player.Statuses = []ActiveStatus{
		{Type: StatusPoison, Duration: 2, Power: 2},
		{Type: StatusRegen, Duration: 2, Power: 5},
	}

	// Turn 1: 50 - 2 (poison) + 5 (regen) = 53 HP
	stTurn1 := st.Step(ActionMoveUp)

	if got, want := stTurn1.Player.Health, 53; got != want {
		t.Errorf("player HP after turn 1 = %d, want %d", got, want)
	}

	fullLog := strings.Join(stTurn1.Log, "\n")
	if !strings.Contains(fullLog, "poison damage") || !strings.Contains(fullLog, "soothing warmth") {
		t.Errorf("expected status tick logs, got %q", fullLog)
	}

	// Check status duration decremented from 2 to 1
	if len(stTurn1.Player.Statuses) != 2 || stTurn1.Player.Statuses[0].Duration != 1 {
		t.Errorf("expected status durations to decrement to 1, got %+v", stTurn1.Player.Statuses)
	}
}

func TestConfusedStatusMovement(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	st.Player.Statuses = []ActiveStatus{
		{Type: StatusConfused, Duration: 3, Power: 1},
	}

	nextState := st.Step(ActionMoveUp)

	fullLog := strings.Join(nextState.Log, "\n")
	if !strings.Contains(fullLog, "stumble about in confusion") {
		t.Errorf("expected confusion log, got %q", fullLog)
	}
}

func TestActionDropItem(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	item := NewHealthPotion(101, Position{X: -1, Y: -1})
	st.Player.Inventory = []Item{item}
	initialItemsCount := len(st.Items)

	stDropped := st.Step(ActionDropItem)

	if len(stDropped.Player.Inventory) != 0 {
		t.Errorf("expected empty inventory after dropping item, got %d", len(stDropped.Player.Inventory))
	}
	if len(stDropped.Items) != initialItemsCount+1 {
		t.Errorf("expected items count to increase by 1, got %d", len(stDropped.Items))
	}
	dropped := stDropped.Items[len(stDropped.Items)-1]
	if dropped.Pos != stDropped.Player.Pos {
		t.Errorf("dropped item pos = %+v, want player pos %+v", dropped.Pos, stDropped.Player.Pos)
	}
}

// Dropping a weapon must reverse the DamageBonus applied on pickup, otherwise
// drop+re-pickup stacks damage forever (Warrior start gear makes this easy to hit).
func TestDropWeaponRemovesDamageBonus(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	baseDamage := st.Player.Damage

	weapon := NewIronDagger(101, Position{X: -1, Y: -1})
	st.Player.Inventory = []Item{weapon}
	st.Player.Damage += weapon.DamageBonus

	stDropped := st.Step(ActionDropItem)
	if stDropped.Player.Damage != baseDamage {
		t.Errorf("after dropping weapon, Damage = %d, want base %d", stDropped.Player.Damage, baseDamage)
	}
	if len(stDropped.Player.Inventory) != 0 {
		t.Errorf("expected empty inventory after drop, got %d", len(stDropped.Player.Inventory))
	}

	// Re-pickup should restore exactly one bonus, not stack past base+bonus.
	stPicked := stDropped.Step(ActionPickup)
	want := baseDamage + weapon.DamageBonus
	if stPicked.Player.Damage != want {
		t.Errorf("after re-picking weapon, Damage = %d, want %d", stPicked.Player.Damage, want)
	}
}

// Warrior starts with an Iron Dagger already factored into Damage; dropping it
// must return them to bare-handed damage.
func TestWarriorDropStartingWeapon(t *testing.T) {
	st := NewGameWithSeedDepthAndClass(40, 20, 42, 1, ClassWarrior)
	if len(st.Player.Inventory) == 0 || st.Player.Inventory[0].ItemType != ItemWeapon {
		t.Fatalf("Warrior should start with a weapon in inventory")
	}
	bonus := st.Player.Inventory[0].DamageBonus
	armedDamage := st.Player.Damage

	stDropped := st.Step(ActionDropItem)
	want := armedDamage - bonus
	if stDropped.Player.Damage != want {
		t.Errorf("Warrior Damage after dropping starter weapon = %d, want %d", stDropped.Player.Damage, want)
	}
}

func TestLavaAndWaterHazardTiles(t *testing.T) {
	st := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileLava, TileWater, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 1, Y: 1},
			Health:    50,
			MaxHealth: 100,
			Statuses:  []ActiveStatus{{Type: StatusPoison, Duration: 5, Power: 2}},
		},
	}

	// Step right onto TileLava
	stLava := st.Step(ActionMoveRight)
	if stLava.Player.Health >= 50 {
		t.Errorf("expected HP to decrease after stepping on lava, got %d", stLava.Player.Health)
	}

	// Step right onto TileWater
	stWater := stLava.Step(ActionMoveRight)
	for _, s := range stWater.Player.Statuses {
		if s.Type == StatusPoison {
			t.Errorf("expected poison status to be cleansed by water")
		}
	}
}

func TestTurnCountKillsAndDamageDealtTracking(t *testing.T) {
	st := GameState{
		Map: GameMap{
			Width:  5,
			Height: 5,
			Tiles: [][]TileType{
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileFloor, TileFloor, TileFloor, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
				{TileWall, TileWall, TileWall, TileWall, TileWall},
			},
		},
		Player: Entity{
			ID:        0,
			Name:      "Player",
			IsPlayer:  true,
			Pos:       Position{X: 1, Y: 1},
			Damage:    15,
			Health:    100,
			MaxHealth: 100,
		},
		Entities: []Entity{
			{ID: 1, Name: "Goblin", Pos: Position{X: 2, Y: 1}, Health: 10, MaxHealth: 10, Damage: 2},
		},
	}

	stAttack := st.Step(ActionMoveRight)
	if stAttack.TurnCount != 1 {
		t.Errorf("TurnCount = %d, want 1", stAttack.TurnCount)
	}
	if stAttack.DamageDealt != 15 {
		t.Errorf("DamageDealt = %d, want 15", stAttack.DamageDealt)
	}
	if stAttack.Kills != 1 {
		t.Errorf("Kills = %d, want 1", stAttack.Kills)
	}
}

func TestClassArchetypesInitialization(t *testing.T) {
	stWarrior := NewGameWithSeedDepthAndClass(60, 30, 12345, 1, ClassWarrior)
	if stWarrior.Player.Health != 120 || len(stWarrior.Player.Inventory) != 1 || stWarrior.Player.Name != "Warrior" {
		t.Errorf("warrior state mismatch: %+v", stWarrior.Player)
	}

	stRogue := NewGameWithSeedDepthAndClass(60, 30, 12345, 1, ClassRogue)
	if stRogue.Player.Health != 90 || len(stRogue.Player.Inventory) != 2 || stRogue.Player.Name != "Rogue" {
		t.Errorf("rogue state mismatch: %+v", stRogue.Player)
	}

	stMage := NewGameWithSeedDepthAndClass(60, 30, 12345, 1, ClassMage)
	if stMage.Player.Health != 80 || len(stMage.Player.Inventory) != 2 || stMage.Player.Name != "Mage" {
		t.Errorf("mage state mismatch: %+v", stMage.Player)
	}
}

func TestFireballSpawnsParticles(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	// Give player a fireball scroll
	st.Player.Inventory = []Item{NewFireballScroll(1, Position{X: -1, Y: -1})}
	st.Player.Pos = Position{X: 10, Y: 5}

	// Add some monsters in fireball radius
	st.Entities = []Entity{
		NewGoblin(1, Position{X: 11, Y: 5}), // adjacent - will be hit
		NewOrc(2, Position{X: 13, Y: 5}),    // within radius 3 - will be hit
		NewTroll(3, Position{X: 15, Y: 8}),  // outside radius - safe
	}

	stAfter := st.Step(ActionUseItem)

	// Fireball should spawn particles
	if len(stAfter.Particles) == 0 {
		t.Fatal("expected particles after casting Fireball, got none")
	}
	// Should have center particles + ring particles (12 ring positions + 4 center = ~16)
	// Some may be on walls/out of bounds or on walls, so expect at least 8
	if len(stAfter.Particles) < 8 {
		t.Errorf("expected at least 10 particles, got %d", len(stAfter.Particles))
	}
	// All particles should have fire colors
	fireColors := map[string]bool{"#FF4500": true, "#FF6600": true, "#FF8800": true, "#FFAA00": true, "#FFCC00": true, "#FFFF00": true}
	for _, p := range stAfter.Particles {
		if !fireColors[p.Color] {
			t.Errorf("particle has non-fire color: %s", p.Color)
		}
		if p.Lifetime <= 0 || p.MaxLife <= 0 {
			t.Errorf("particle has invalid lifetime: %+v", p)
		}
	}
}

func TestTeleportSpawnsParticles(t *testing.T) {
	st := NewGameWithSeed(20, 10, 12345)
	oldPos := st.Player.Pos
	st.Player.Inventory = []Item{NewTeleportScroll(1, Position{X: -1, Y: -1})}

	stAfter := st.Step(ActionUseItem)

	// Teleport should spawn particles at both old and new position
	if len(stAfter.Particles) == 0 {
		t.Fatal("expected particles after casting Teleport, got none")
	}
	if stAfter.Player.Pos == oldPos {
		t.Fatal("expected player position to change after teleport")
	}
	// Should have particles at both positions (8 around each + 3 at each = ~22)
	if len(stAfter.Particles) < 15 {
		t.Errorf("expected at least 15 particles, got %d", len(stAfter.Particles))
	}
	// All particles should have teleport colors
	teleportColors := map[string]bool{"#E040FB": true, "#AA00FF": true, "#6600CC": true, "#FF00FF": true, "#8800FF": true}
	for _, p := range stAfter.Particles {
		if !teleportColors[p.Color] {
			t.Errorf("particle has non-teleport color: %s", p.Color)
		}
	}
}

func TestParticlesDecayAndExpire(t *testing.T) {
	st := GameState{
		Map: GameMap{
			Width:  20,
			Height: 10,
			Tiles:  make([][]TileType, 10),
		},
		Player: Entity{
			ID: 0, Name: "Player", IsPlayer: true, Pos: Position{X: 10, Y: 5},
			Health: 100, MaxHealth: 100, Damage: 10,
		},
		Particles: []Particle{
			{Pos: Position{X: 10, Y: 5}, Glyph: "░", Color: "#FF4500", Lifetime: 1, MaxLife: 1},
			{Pos: Position{X: 11, Y: 5}, Glyph: "▒", Color: "#FF6600", Lifetime: 2, MaxLife: 2},
			{Pos: Position{X: 12, Y: 5}, Glyph: "▓", Color: "#FF8800", Lifetime: 3, MaxLife: 3},
		},
	}
	for y := range st.Map.Tiles {
		st.Map.Tiles[y] = make([]TileType, 20)
		for x := range st.Map.Tiles[y] {
			if x == 0 || x == 19 || y == 0 || y == 9 {
				st.Map.Tiles[y][x] = TileWall
			} else {
				st.Map.Tiles[y][x] = TileFloor
			}
		}
	}

	// First step - particles with lifetime 1 should expire
	st = st.Step(ActionMoveUp)
	if len(st.Particles) != 2 {
		t.Errorf("after 1 step, expected 2 particles (lifetime 2 and 3), got %d", len(st.Particles))
	}

	// Second step - particles with lifetime 2 should expire
	st = st.Step(ActionMoveUp)
	if len(st.Particles) != 1 {
		t.Errorf("after 2 steps, expected 1 particle (lifetime 3), got %d", len(st.Particles))
	}

	// Third step - last particle should expire
	st = st.Step(ActionMoveUp)
	if len(st.Particles) != 0 {
		t.Errorf("after 3 steps, expected 0 particles, got %d", len(st.Particles))
	}
}
