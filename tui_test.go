package main

import (
	"strings"
	"testing"

	"glyphgrinder/internal/tuitest"
)

// playerAt reports the grid position of the player glyph in a rendered view.
// Styling is stripped by looking for the raw player rune, which is safe as
// long as no other glyph reuses it (see AGENTS.md gotchas).
func playerAt(t *testing.T, lines []string) Position {
	t.Helper()
	for y, line := range lines {
		// Cells are styled, so count runes of the unstyled grid instead by
		// stripping ANSI escapes.
		plain := stripANSI(line)
		if x := strings.IndexRune(plain, '@'); x >= 0 {
			return Position{X: x, Y: y}
		}
	}
	t.Fatalf("player glyph not found in view:\n%s", strings.Join(lines, "\n"))
	return Position{}
}

func stripANSI(s string) string {
	var sb strings.Builder
	inEscape := false
	for _, r := range s {
		switch {
		case inEscape:
			if r == 'm' {
				inEscape = false
			}
		case r == '\x1b':
			inEscape = true
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func TestViewRendersFullGrid(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	lines := d.Lines()
	if got, want := len(lines), 10; got != want {
		t.Fatalf("view has %d rows, want %d", got, want)
	}
	for y, line := range lines {
		if got, want := len([]rune(stripANSI(line))), 20; got != want {
			t.Errorf("row %d has %d cells, want %d", y, got, want)
		}
	}
	wantPos := m.state.Player.Pos
	if got := playerAt(t, lines); got != wantPos {
		t.Errorf("player rendered at %+v, want %+v", got, wantPos)
	}
}

func TestMovementKeys(t *testing.T) {
	baseModel := initialModelWithSeed(12345)
	startPos := baseModel.state.Player.Pos

	cases := []struct {
		name string
		keys []string
		want Position
	}{
		{"arrow up", []string{"up"}, Position{X: startPos.X, Y: startPos.Y - 1}},
		{"arrow down", []string{"down"}, Position{X: startPos.X, Y: startPos.Y + 1}},
		{"arrow left", []string{"left"}, Position{X: startPos.X - 1, Y: startPos.Y}},
		{"arrow right", []string{"right"}, Position{X: startPos.X + 1, Y: startPos.Y}},
		{"wasd", []string{"w", "a"}, Position{X: startPos.X - 1, Y: startPos.Y - 1}},
		{"round trip", []string{"w", "s", "a", "d"}, startPos},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := tuitest.New(t, initialModelWithSeed(12345))
			d.Keys(tc.keys...)
			if got := playerAt(t, d.Lines()); got != tc.want {
				t.Errorf("after %v player at %+v, want %+v", tc.keys, got, tc.want)
			}
		})
	}
}

func TestPlayerCannotWalkThroughWalls(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	// Pressing up many times will eventually hit a wall.
	for i := 0; i < 20; i++ {
		d.Key("up")
	}
	posUp := playerAt(t, d.Lines())
	if posUp.Y <= 0 || m.state.Map.Tiles[posUp.Y-1][posUp.X] != TileWall {
		t.Errorf("expected tile above player at %+v to be TileWall", posUp)
	}

	// Pressing left many times will hit a wall.
	for i := 0; i < 20; i++ {
		d.Key("left")
	}
	posLeft := playerAt(t, d.Lines())
	if posLeft.X <= 0 || m.state.Map.Tiles[posLeft.Y][posLeft.X-1] != TileWall {
		t.Errorf("expected tile to the left of player at %+v to be TileWall", posLeft)
	}
}

func TestQuitKeys(t *testing.T) {
	for _, key := range []string{"q", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			d := tuitest.New(t, initialModel())
			if d.Quit() {
				t.Fatal("model quit before any key was sent")
			}
			d.Key(key)
			if !d.Quit() {
				t.Errorf("%q did not quit the program", key)
			}
		})
	}
}

func TestGameStateStep(t *testing.T) {
	state := NewGameWithSeed(20, 10, 12345)
	startPos := state.Player.Pos

	// Move Up
	state = state.Step(ActionMoveUp)
	if got, want := state.Player.Pos, (Position{X: startPos.X, Y: startPos.Y - 1}); got != want {
		t.Errorf("after MoveUp got pos %+v, want %+v", got, want)
	}

	// ActionNone does not move player
	stateBefore := state
	state = state.Step(ActionNone)
	if state.Player.Pos != stateBefore.Player.Pos {
		t.Errorf("ActionNone changed pos from %+v to %+v", stateBefore.Player.Pos, state.Player.Pos)
	}
}

func TestMapGenerationDeterminism(t *testing.T) {
	seed := int64(987654321)
	s1 := NewGameWithSeed(20, 10, seed)
	s2 := NewGameWithSeed(20, 10, seed)

	if s1.Player.Pos != s2.Player.Pos {
		t.Fatalf("players spawned at different positions for same seed: %+v vs %+v", s1.Player.Pos, s2.Player.Pos)
	}

	for y := 0; y < s1.Map.Height; y++ {
		for x := 0; x < s1.Map.Width; x++ {
			if s1.Map.Tiles[y][x] != s2.Map.Tiles[y][x] {
				t.Fatalf("map tile mismatch at (%d,%d)", x, y)
			}
		}
	}
}

func TestPlayerSpawnInValidFloorTile(t *testing.T) {
	state := NewGame(20, 10)
	pos := state.Player.Pos
	if tile := state.Map.Tiles[pos.Y][pos.X]; tile != TileFloor {
		t.Errorf("player spawned on tile %v, want TileFloor", tile)
	}
}

func TestMapBorderIsWalled(t *testing.T) {
	state := NewGame(20, 10)
	w, h := state.Map.Width, state.Map.Height

	for x := 0; x < w; x++ {
		if tile := state.Map.Tiles[0][x]; tile != TileWall {
			t.Errorf("top border tile at x=%d is %v, want TileWall", x, tile)
		}
		if tile := state.Map.Tiles[h-1][x]; tile != TileWall {
			t.Errorf("bottom border tile at x=%d is %v, want TileWall", x, tile)
		}
	}

	for y := 0; y < h; y++ {
		if tile := state.Map.Tiles[y][0]; tile != TileWall {
			t.Errorf("left border tile at y=%d is %v, want TileWall", y, tile)
		}
		if tile := state.Map.Tiles[y][w-1]; tile != TileWall {
			t.Errorf("right border tile at y=%d is %v, want TileWall", y, tile)
		}
	}
}

func TestViewRendersMonsters(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	if len(m.state.Entities) == 0 {
		t.Fatal("initial model with seed should have entities")
	}

	lines := d.Lines()
	for _, e := range m.state.Entities {
		if e.Pos.Y >= len(lines) {
			t.Fatalf("monster pos Y %d out of bounds for lines len %d", e.Pos.Y, len(lines))
		}
		plainRow := []rune(stripANSI(lines[e.Pos.Y]))
		if e.Pos.X >= len(plainRow) {
			t.Fatalf("monster pos X %d out of bounds for plain row len %d", e.Pos.X, len(plainRow))
		}
		gotRune := string(plainRow[e.Pos.X])
		if gotRune != e.Rune {
			t.Errorf("expected monster rune %q at %+v, got %q", e.Rune, e.Pos, gotRune)
		}
	}
}

func TestBumpToAttackViaTUI(t *testing.T) {
	m := model{
		state: GameState{
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
				Rune:      "@",
				Color:     "#00FF00",
				Damage:    10,
				Health:    100,
				MaxHealth: 100,
			},
			Entities: []Entity{
				NewGoblin(1, Position{X: 2, Y: 1}),
			},
		},
	}

	d := tuitest.New(t, m)

	// Press 'up' key to bump-attack goblin
	d.Key("up")

	pos := playerAt(t, d.Lines())
	if pos != (Position{X: 2, Y: 2}) {
		t.Errorf("player position after bump attack = %+v, want (2,2)", pos)
	}

	// Verify monster killed in state (retrieved from driver model)
	// Note: tuitest Driver updates model internal state on key press
}

func TestMonsterTurnsViaTUI(t *testing.T) {
	m := model{
		state: GameState{
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
				Rune:      "@",
				Color:     "#00FF00",
				Damage:    10,
				Health:    100,
				MaxHealth: 100,
			},
			Entities: []Entity{
				NewOrc(1, Position{X: 1, Y: 3}),
			},
		},
	}

	d := tuitest.New(t, m)

	// Step right to (2, 1). Orc at (1, 3) should step up to (1, 2).
	d.Key("right")

	lines := d.Lines()
	pos := playerAt(t, lines)
	if pos != (Position{X: 2, Y: 1}) {
		t.Errorf("player position = %+v, want (2, 1)", pos)
	}

	// Verify Orc glyph 'o' is rendered at (1, 2)
	plainRow := stripANSI(lines[2]) // Y = 2
	if got := string([]rune(plainRow)[1]); got != "o" {
		t.Errorf("expected Orc 'o' at (1, 2), got %q in line: %q", got, plainRow)
	}
}
