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
			reachable := reachableFloors(g.Map, g.Player.Pos)
			floorCount := 0
			for y := 0; y < g.Map.Height; y++ {
				for x := 0; x < g.Map.Width; x++ {
					if g.Map.Tiles[y][x] == TileFloor {
						floorCount++
					}
				}
			}
			if got := len(reachable); got != floorCount {
				t.Errorf("seed %d: reachable floors = %d, want all %d", seed, got, floorCount)
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

func reachableFloors(m GameMap, start Position) map[Position]struct{} {
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
			if m.Tiles[next.Y][next.X] != TileFloor {
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

func mapSignature(m GameMap) string {
	signature := make([]byte, 0, m.Width*m.Height)
	for _, row := range m.Tiles {
		for _, tile := range row {
			signature = append(signature, byte(tile))
		}
	}
	return string(signature)
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

func combatTestState(targetHealth int) GameState {
	g := newTestGame(11)
	g.Entities = []Entity{
		{
			ID:        1,
			Pos:       Position{X: g.Player.Pos.X + 1, Y: g.Player.Pos.Y},
			Rune:      "g",
			Health:    targetHealth,
			MaxHealth: 20,
			Damage:    5,
		},
		{
			ID:        2,
			Pos:       Position{X: g.Player.Pos.X - 1, Y: g.Player.Pos.Y},
			Rune:      "g",
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
	const width, height = 7, 7
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
	return GameState{
		Map: GameMap{Width: width, Height: height, Tiles: tiles},
		Player: Entity{
			IsPlayer:  true,
			Rune:      "@",
			Health:    100,
			MaxHealth: 100,
			Damage:    10,
		},
		Log: log,
	}
}

func testMonster(id int, pos Position) Entity {
	return Entity{
		ID:        id,
		Pos:       pos,
		Rune:      "g",
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
