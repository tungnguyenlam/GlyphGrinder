package main

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func newTestGame(seed int64) GameState {
	return NewGame(20, 10, rand.New(rand.NewSource(seed)))
}

func TestNewGameGeneratesWalledDungeon(t *testing.T) {
	g := newTestGame(1)

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
	if g.Map.Tiles[g.Player.Pos.Y][g.Player.Pos.X] != TileFloor {
		t.Error("dungeon should contain floor at the player position")
	}
}

func TestNewGamePlacesPlayerOnFloor(t *testing.T) {
	g := newTestGame(2)

	p := g.Player
	if !p.IsPlayer {
		t.Error("player entity should have IsPlayer set")
	}
	if g.Map.Tiles[p.Pos.Y][p.Pos.X] != TileFloor {
		t.Error("player must start on a floor tile")
	}
	if p.Health != p.MaxHealth {
		t.Errorf("player starts at %d/%d health, want full", p.Health, p.MaxHealth)
	}
	if got, want := g.Depth, 1; got != want {
		t.Errorf("starting depth = %d, want %d", got, want)
	}
}

func TestNewGamePlacesReachableStairsAwayFromPlayer(t *testing.T) {
	for _, seed := range []int64{1, 7, 42, 99} {
		t.Run("seed_"+strconv.FormatInt(seed, 10), func(t *testing.T) {
			g := newTestGame(seed)
			stairs := tilePositions(g.Map, TileStairs)
			if got, want := len(stairs), 1; got != want {
				t.Fatalf("stair count = %d, want %d", got, want)
			}
			if stairs[0] == g.Player.Pos {
				t.Fatal("stairs overlap the player's starting position")
			}
			if _, reachable := reachableTiles(g.Map, g.Player.Pos)[stairs[0]]; !reachable {
				t.Fatalf("stairs at %+v are not reachable from player at %+v", stairs[0], g.Player.Pos)
			}
		})
	}
}

func TestNewGamePlacesMonstersOnDistinctFloors(t *testing.T) {
	g := newTestGame(8)
	if len(g.Entities) == 0 {
		t.Fatal("new dungeon has no monsters")
	}

	occupied := map[Position]struct{}{g.Player.Pos: {}}
	for i, monster := range g.Entities {
		if got, want := monster.ID, i+1; got != want {
			t.Errorf("monster %d ID = %d, want %d", i, got, want)
		}
		if monster.IsPlayer {
			t.Errorf("monster %d is marked as player", monster.ID)
		}
		if monster.Health <= 0 || monster.Health != monster.MaxHealth || monster.Damage <= 0 {
			t.Errorf("monster %d has invalid combat stats: %+v", monster.ID, monster)
		}
		if g.Map.Tiles[monster.Pos.Y][monster.Pos.X] != TileFloor {
			t.Errorf("monster %d placed off floor at %+v", monster.ID, monster.Pos)
		}
		if _, exists := occupied[monster.Pos]; exists {
			t.Errorf("monster %d overlaps another actor at %+v", monster.ID, monster.Pos)
		}
		occupied[monster.Pos] = struct{}{}
	}
}

func TestNewGamePlacesPotionsOnDistinctUnoccupiedFloors(t *testing.T) {
	for _, seed := range []int64{3, 8, 21, 55} {
		t.Run("seed_"+strconv.FormatInt(seed, 10), func(t *testing.T) {
			g := newTestGame(seed)
			if got, want := len(g.Items), 2; got != want {
				t.Fatalf("item count = %d, want %d potions", got, want)
			}
			occupied := map[Position]struct{}{g.Player.Pos: {}}
			for _, monster := range g.Entities {
				occupied[monster.Pos] = struct{}{}
			}
			for _, item := range g.Items {
				if item.Type != ItemPotion {
					t.Errorf("item at %+v has type %d, want potion", item.Pos, item.Type)
				}
				if g.Map.Tiles[item.Pos.Y][item.Pos.X] != TileFloor {
					t.Errorf("potion placed off ordinary floor at %+v", item.Pos)
				}
				if _, exists := occupied[item.Pos]; exists {
					t.Errorf("potion overlaps actor or another item at %+v", item.Pos)
				}
				occupied[item.Pos] = struct{}{}
			}
		})
	}
}

func TestNewGameIsDeterministicForSeed(t *testing.T) {
	first := newTestGame(42)
	second := newTestGame(42)

	if !reflect.DeepEqual(first, second) {
		t.Fatal("same seed generated different game states")
	}
}

func TestGeneratedDungeonFloorsAreReachable(t *testing.T) {
	for _, seed := range []int64{1, 7, 42, 99} {
		t.Run("seed_"+strconv.FormatInt(seed, 10), func(t *testing.T) {
			g := newTestGame(seed)
			reachable := reachableTiles(g.Map, g.Player.Pos)
			walkableCount := 0
			for y := 0; y < g.Map.Height; y++ {
				for x := 0; x < g.Map.Width; x++ {
					if g.Map.Tiles[y][x] != TileWall {
						walkableCount++
					}
				}
			}
			if got := len(reachable); got != walkableCount {
				t.Errorf("seed %d: reachable tiles = %d, want all %d", seed, got, walkableCount)
			}
		})
	}
}

func TestDifferentSeedsVaryDungeon(t *testing.T) {
	maps := make(map[string]struct{})
	for _, seed := range []int64{1, 2, 3, 4} {
		maps[mapSignature(newTestGame(seed).Map)] = struct{}{}
	}
	if len(maps) < 2 {
		t.Fatal("different seeds all generated the same map")
	}
}

func TestVisibilityStopsBehindWalls(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 3, Y: 3}
	g.Map.Tiles[2][3] = TileWall
	g.Map.Visible = nil
	g.Map.Explored = nil
	g = g.refreshVisibility()

	if !g.Map.Visible[2][3] {
		t.Error("blocking wall should itself be visible")
	}
	if g.Map.Visible[1][3] {
		t.Error("tile directly behind wall should not be visible")
	}
	if g.Map.Explored[1][3] {
		t.Error("never-seen tile behind wall should not be explored")
	}
}

func TestVisibilityRefreshPreservesExplorationMemory(t *testing.T) {
	g := openTestStateSized(17, 7)
	g.Player.Pos = Position{X: 2, Y: 3}
	g.Map.Visible = nil
	g.Map.Explored = nil
	g = g.refreshVisibility()
	remembered := Position{X: 2, Y: 2}
	priorExplored := g.Map.Explored

	g.Player.Pos = Position{X: 14, Y: 3}
	got := g.refreshVisibility()

	if got.Map.Visible[remembered.Y][remembered.X] {
		t.Error("old tile remained visible outside the new radius")
	}
	if !got.Map.Explored[remembered.Y][remembered.X] {
		t.Error("old visible tile was not remembered as explored")
	}
	if priorExplored[3][14] {
		t.Error("visibility refresh mutated prior exploration backing grid")
	}
}

func TestStepRefreshesVisibilityAfterPlayerMoves(t *testing.T) {
	g := openTestStateSized(12, 7)
	g.Player.Pos = Position{X: 2, Y: 3}
	g.Entities = nil
	g.Map.Visible = nil
	g.Map.Explored = nil
	g = g.refreshVisibility()
	newlyVisible := Position{X: 9, Y: 3}
	if g.Map.Visible[newlyVisible.Y][newlyVisible.X] {
		t.Fatal("fixture tile is already visible before movement")
	}

	got := g.Step(ActionMoveRight)

	if !got.Map.Visible[newlyVisible.Y][newlyVisible.X] {
		t.Error("player movement did not refresh visibility")
	}
	if g.Map.Visible[newlyVisible.Y][newlyVisible.X] {
		t.Error("Step mutated the prior state's visibility grid")
	}
}

func reachableTiles(m GameMap, start Position) map[Position]struct{} {
	reachable := map[Position]struct{}{start: {}}
	queue := []Position{start}
	directions := []Position{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, direction := range directions {
			next := Position{X: current.X + direction.X, Y: current.Y + direction.Y}
			if next.X < 0 || next.X >= m.Width || next.Y < 0 || next.Y >= m.Height {
				continue
			}
			if m.Tiles[next.Y][next.X] == TileWall {
				continue
			}
			if _, seen := reachable[next]; seen {
				continue
			}
			reachable[next] = struct{}{}
			queue = append(queue, next)
		}
	}
	return reachable
}

func tilePositions(m GameMap, want TileType) []Position {
	var positions []Position
	for y, row := range m.Tiles {
		for x, tile := range row {
			if tile == want {
				positions = append(positions, Position{X: x, Y: y})
			}
		}
	}
	return positions
}

func mapSignature(m GameMap) string {
	signature := make([]byte, 0, m.Width*m.Height)
	for _, row := range m.Tiles {
		for _, tile := range row {
			signature = append(signature, byte(tile))
		}
	}
	return string(signature)
}

func TestDescendRequiresPlayerToStandOnStairs(t *testing.T) {
	g := newTestGame(12)
	before := g

	got := g.Descend(rand.New(rand.NewSource(99)))

	if !reflect.DeepEqual(got, before) {
		t.Errorf("descent away from stairs changed state:\ngot  %+v\nwant %+v", got, before)
	}
}

func TestDescendGeneratesNextLevelAndPreservesPlayerStats(t *testing.T) {
	rng := rand.New(rand.NewSource(18))
	g := NewGame(20, 10, rng)
	g.Player.Pos = tilePositions(g.Map, TileStairs)[0]
	g.Player.Health = 37
	g.Player.MaxHealth = 120
	g.Player.Damage = 14
	g.Potions = 2
	oldMap := mapSignature(g.Map)

	got := g.Descend(rng)

	if got.Depth != 2 {
		t.Errorf("depth after descent = %d, want 2", got.Depth)
	}
	if got.Player.Health != 37 || got.Player.MaxHealth != 120 || got.Player.Damage != 14 {
		t.Errorf("player stats after descent = %+v, want health 37/120 and damage 14", got.Player)
	}
	if got.Potions != 2 {
		t.Errorf("potion inventory after descent = %d, want 2", got.Potions)
	}
	if got.Map.Width != g.Map.Width || got.Map.Height != g.Map.Height {
		t.Errorf("next map dimensions = %dx%d, want %dx%d", got.Map.Width, got.Map.Height, g.Map.Width, g.Map.Height)
	}
	if got.Map.Tiles[got.Player.Pos.Y][got.Player.Pos.X] != TileFloor {
		t.Errorf("player starts next depth on tile %d, want floor", got.Map.Tiles[got.Player.Pos.Y][got.Player.Pos.X])
	}
	stairs := tilePositions(got.Map, TileStairs)
	if len(stairs) != 1 {
		t.Fatalf("next depth stair count = %d, want 1", len(stairs))
	}
	if _, reachable := reachableTiles(got.Map, got.Player.Pos)[stairs[0]]; !reachable {
		t.Fatalf("next stairs at %+v are unreachable", stairs[0])
	}
	if got.Map.Explored == nil || !got.Map.Explored[got.Player.Pos.Y][got.Player.Pos.X] {
		t.Error("next depth visibility was not initialized around the player")
	}
	if &got.Map.Visible[0][0] == &g.Map.Visible[0][0] || &got.Map.Explored[0][0] == &g.Map.Explored[0][0] {
		t.Error("descent retained visibility grids from the prior depth")
	}
	if oldMap == mapSignature(got.Map) {
		t.Error("descent generated the same dungeon layout")
	}
	if gotLog, want := got.Log[len(got.Log)-1], "You descend to depth 2."; gotLog != want {
		t.Errorf("last log entry = %q, want %q", gotLog, want)
	}
}

func TestDescentSequenceIsDeterministicForSeed(t *testing.T) {
	descend := func(seed int64) GameState {
		rng := rand.New(rand.NewSource(seed))
		g := NewGame(20, 10, rng)
		g.Player.Pos = tilePositions(g.Map, TileStairs)[0]
		return g.Descend(rng)
	}

	first := descend(73)
	second := descend(73)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("same seeded run generated different second levels")
	}
}

func TestStepPicksUpPotion(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 2, Y: 3}
	g.Entities = nil
	g.Items = []Item{{Type: ItemPotion, Pos: Position{X: 3, Y: 3}}}
	priorItems := append([]Item(nil), g.Items...)

	got := g.Step(ActionMoveRight)

	if got.Player.Pos != (Position{X: 3, Y: 3}) {
		t.Errorf("player position = %+v, want potion cell {3 3}", got.Player.Pos)
	}
	if got.Potions != 1 || len(got.Items) != 0 {
		t.Errorf("after pickup potions = %d, ground items = %+v; want 1 and none", got.Potions, got.Items)
	}
	if !reflect.DeepEqual(g.Items, priorItems) {
		t.Errorf("pickup mutated prior items: got %+v, want %+v", g.Items, priorItems)
	}
	if gotLog, want := got.Log[len(got.Log)-1], "You pick up a health potion."; gotLog != want {
		t.Errorf("last log entry = %q, want %q", gotLog, want)
	}
}

func TestUsePotionHealsAndSpendsTurn(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 3, Y: 3}
	g.Player.Health = 70
	g.Potions = 2
	g.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}

	got := g.Step(ActionUsePotion)

	if got.Potions != 1 {
		t.Errorf("potions after use = %d, want 1", got.Potions)
	}
	if got.Player.Health != 90 {
		t.Errorf("health after healing 25 and taking 5 damage = %d, want 90", got.Player.Health)
	}
	if gotLog, want := got.Log, []string{
		"The fight begins.",
		"You drink a potion and recover 25 health.",
		"Monster 1 hits you for 5 damage.",
	}; !reflect.DeepEqual(gotLog, want) {
		t.Errorf("log after potion turn = %q, want %q", gotLog, want)
	}
}

func TestUsePotionCapsHealingAtMaxHealth(t *testing.T) {
	g := openTestState()
	g.Player.Health = 90
	g.Potions = 1
	g.Entities = nil

	got := g.Step(ActionUsePotion)

	if got.Player.Health != 100 || got.Potions != 0 {
		t.Errorf("after potion health = %d and inventory = %d, want 100 and 0", got.Player.Health, got.Potions)
	}
	if gotLog, want := got.Log[len(got.Log)-1], "You drink a potion and recover 10 health."; gotLog != want {
		t.Errorf("last log entry = %q, want %q", gotLog, want)
	}
}

func TestUsePotionIsFreeWhenUnavailableOrAtFullHealth(t *testing.T) {
	tests := []struct {
		name    string
		health  int
		potions int
	}{
		{name: "no inventory", health: 50, potions: 0},
		{name: "full health", health: 100, potions: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := openTestState()
			g.Player.Health = tc.health
			g.Potions = tc.potions
			g.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}

			got := g.Step(ActionUsePotion)

			if !reflect.DeepEqual(got, g) {
				t.Errorf("invalid potion use changed state:\ngot  %+v\nwant %+v", got, g)
			}
		})
	}
}

func TestStepMovesPlayer(t *testing.T) {
	cases := []struct {
		name   string
		action Action
		delta  Position
	}{
		{name: "up", action: ActionMoveUp, delta: Position{Y: -1}},
		{name: "down", action: ActionMoveDown, delta: Position{Y: 1}},
		{name: "left", action: ActionMoveLeft, delta: Position{X: -1}},
		{name: "right", action: ActionMoveRight, delta: Position{X: 1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := newTestGame(3)
			start := g.Player.Pos
			got := g.Step(tc.action)
			want := Position{X: start.X + tc.delta.X, Y: start.Y + tc.delta.Y}
			if got.Player.Pos != want {
				t.Errorf("player pos = %+v, want %+v", got.Player.Pos, want)
			}
		})
	}
}

func TestStepRejectsBlockedMovement(t *testing.T) {
	t.Run("wall", func(t *testing.T) {
		g := newTestGame(4)
		floor, action := floorBesideWall(t, g.Map)
		g.Player.Pos = floor

		got := g.Step(action)
		if got.Player.Pos != g.Player.Pos {
			t.Errorf("player pos = %+v, want %+v", got.Player.Pos, g.Player.Pos)
		}
	})

	t.Run("map bounds", func(t *testing.T) {
		g := newTestGame(5)
		g.Player.Pos = Position{X: 0, Y: 0}

		got := g.Step(ActionMoveLeft)
		if got.Player.Pos != g.Player.Pos {
			t.Errorf("player pos = %+v, want %+v", got.Player.Pos, g.Player.Pos)
		}
	})
}

func TestStepIgnoresNoAction(t *testing.T) {
	g := newTestGame(6)

	got := g.Step(ActionNone)
	if got.Player.Pos != g.Player.Pos {
		t.Errorf("player pos = %+v, want %+v", got.Player.Pos, g.Player.Pos)
	}
}

func TestStepBumpAttackDamagesMonster(t *testing.T) {
	g := combatTestState(20)
	start := g.Player.Pos
	untargeted := g.Entities[1]
	priorLogBacking := g.Log[:cap(g.Log)]

	got := g.Step(ActionMoveRight)

	if got.Player.Pos != start {
		t.Errorf("player moved to %+v during attack, want %+v", got.Player.Pos, start)
	}
	if got.Entities[0].Health != 10 {
		t.Errorf("target health = %d, want 10", got.Entities[0].Health)
	}
	if got.Entities[1] != untargeted {
		t.Errorf("untargeted monster changed: got %+v, want %+v", got.Entities[1], untargeted)
	}
	if gotLog, want := got.Log, []string{
		"The fight begins.",
		"You hit monster 1 for 10 damage.",
		"Monster 1 hits you for 5 damage.",
		"Monster 2 hits you for 5 damage.",
	}; !reflect.DeepEqual(gotLog, want) {
		t.Errorf("log = %q, want %q", gotLog, want)
	}
	if g.Entities[0].Health != 20 {
		t.Errorf("Step mutated prior state target health to %d, want 20", g.Entities[0].Health)
	}
	if priorLogBacking[1] != "" {
		t.Errorf("Step mutated prior log backing array with %q", priorLogBacking[1])
	}
}

func TestStepBumpAttackKillsMonster(t *testing.T) {
	g := combatTestState(10)

	got := g.Step(ActionMoveRight)

	if len(got.Entities) != 1 || got.Entities[0].ID != 2 {
		t.Fatalf("entities after kill = %+v, want only monster 2", got.Entities)
	}
	if gotLog, want := got.Log, []string{
		"The fight begins.",
		"You kill monster 1.",
		"Monster 2 hits you for 5 damage.",
	}; !reflect.DeepEqual(gotLog, want) {
		t.Errorf("log = %q, want %q", gotLog, want)
	}
	if len(g.Entities) != 2 || g.Entities[0].Health != 10 {
		t.Errorf("Step mutated prior entities: %+v", g.Entities)
	}
}

func TestMonsterTurnPursuesPlayer(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 2, Y: 3}
	g.Entities = []Entity{testMonster(1, Position{X: 5, Y: 3})}

	got := g.Step(ActionMoveUp)

	if got.Player.Pos != (Position{X: 2, Y: 2}) {
		t.Errorf("player pos = %+v, want {2 2}", got.Player.Pos)
	}
	if got.Entities[0].Pos != (Position{X: 4, Y: 3}) {
		t.Errorf("monster pos = %+v, want {4 3}", got.Entities[0].Pos)
	}
	if g.Entities[0].Pos != (Position{X: 5, Y: 3}) {
		t.Errorf("Step mutated prior monster position to %+v", g.Entities[0].Pos)
	}
}

func TestMonsterTurnAvoidsWallsAndActors(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 1, Y: 3}
	g.Entities = []Entity{
		testMonster(1, Position{X: 4, Y: 3}),
		testMonster(2, Position{X: 3, Y: 3}),
	}
	g.Map.Tiles[2][4] = TileWall

	got := g.Step(ActionMoveUp)

	if got.Entities[0].Pos != (Position{X: 4, Y: 3}) {
		t.Errorf("blocked monster moved to %+v, want {4 3}", got.Entities[0].Pos)
	}
	if got.Entities[1].Pos != (Position{X: 2, Y: 3}) {
		t.Errorf("unblocked monster moved to %+v, want {2 3}", got.Entities[1].Pos)
	}
}

func TestMonsterTurnAttacksAfterBlockedPlayerAction(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 3, Y: 3}
	g.Map.Tiles[3][4] = TileWall
	g.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}
	priorLogBacking := g.Log[:cap(g.Log)]

	got := g.Step(ActionMoveRight)

	if got.Player.Health != 95 {
		t.Errorf("player health = %d, want 95", got.Player.Health)
	}
	if g.Player.Health != 100 {
		t.Errorf("Step mutated prior player health to %d", g.Player.Health)
	}
	if gotLog, want := got.Log, []string{"The fight begins.", "Monster 1 hits you for 5 damage."}; !reflect.DeepEqual(gotLog, want) {
		t.Errorf("log = %q, want %q", gotLog, want)
	}
	if priorLogBacking[1] != "" {
		t.Errorf("Step mutated prior log backing array with %q", priorLogBacking[1])
	}
}

func TestMonsterTurnsResolveInStableIDOrder(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 3, Y: 3}
	g.Map.Tiles[3][4] = TileWall
	g.Entities = []Entity{
		testMonster(2, Position{X: 3, Y: 4}),
		testMonster(1, Position{X: 3, Y: 2}),
	}

	got := g.Step(ActionMoveRight)

	want := []string{
		"The fight begins.",
		"Monster 1 hits you for 5 damage.",
		"Monster 2 hits you for 5 damage.",
	}
	if !reflect.DeepEqual(got.Log, want) {
		t.Errorf("log = %q, want stable ID order %q", got.Log, want)
	}
}

func TestMonsterAttackEndsRunAndFreezesState(t *testing.T) {
	g := openTestState()
	g.Player.Pos = Position{X: 3, Y: 3}
	g.Player.Health = 5
	g.Map.Tiles[3][4] = TileWall
	g.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}

	dead := g.Step(ActionMoveRight)

	if !dead.GameOver {
		t.Fatal("lethal monster attack did not mark game over")
	}
	if dead.Player.Health != 0 {
		t.Errorf("player health = %d, want 0", dead.Player.Health)
	}
	if got, want := dead.Log[len(dead.Log)-1], "You die."; got != want {
		t.Errorf("final log entry = %q, want %q", got, want)
	}
	if after := dead.Step(ActionMoveLeft); !reflect.DeepEqual(after, dead) {
		t.Errorf("game-over state changed on another turn:\ngot  %+v\nwant %+v", after, dead)
	}
}

func combatTestState(targetHealth int) GameState {
	g := newTestGame(11)
	g.Entities = []Entity{
		{
			ID:        1,
			Pos:       Position{X: g.Player.Pos.X + 1, Y: g.Player.Pos.Y},
			Health:    targetHealth,
			MaxHealth: 20,
			Damage:    5,
		},
		{
			ID:        2,
			Pos:       Position{X: g.Player.Pos.X - 1, Y: g.Player.Pos.Y},
			Health:    20,
			MaxHealth: 20,
			Damage:    5,
		},
	}
	g.Log = make([]string, 1, 3)
	g.Log[0] = "The fight begins."
	return g
}

func openTestState() GameState {
	return openTestStateSized(7, 7)
}

func openTestStateSized(width, height int) GameState {
	tiles := make([][]TileType, height)
	for y := range tiles {
		tiles[y] = make([]TileType, width)
		for x := range tiles[y] {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				tiles[y][x] = TileWall
			}
		}
	}
	log := make([]string, 1, 3)
	log[0] = "The fight begins."
	g := GameState{
		Map: GameMap{Width: width, Height: height, Tiles: tiles},
		Player: Entity{
			IsPlayer:  true,
			Pos:       Position{X: width / 2, Y: height / 2},
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
		Log:   log,
		Depth: 1,
	}
	return g.refreshVisibility()
}

func testMonster(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Pos:       pos,
		Health:    20,
		MaxHealth: 20,
		Damage:    5,
	}
}

func floorBesideWall(t *testing.T, m GameMap) (Position, Action) {
	t.Helper()
	directions := []struct {
		delta  Position
		action Action
	}{
		{delta: Position{Y: -1}, action: ActionMoveUp},
		{delta: Position{Y: 1}, action: ActionMoveDown},
		{delta: Position{X: -1}, action: ActionMoveLeft},
		{delta: Position{X: 1}, action: ActionMoveRight},
	}
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			if m.Tiles[y][x] != TileFloor {
				continue
			}
			for _, direction := range directions {
				if m.Tiles[y+direction.delta.Y][x+direction.delta.X] == TileWall {
					return Position{X: x, Y: y}, direction.action
				}
			}
		}
	}
	t.Fatal("generated dungeon has no floor beside a wall")
	return Position{}, ActionNone
}
