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
