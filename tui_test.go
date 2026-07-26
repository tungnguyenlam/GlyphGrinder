package main

import (
	"strings"
	"testing"

	"glyphgrinder/internal/tuitest"
)

// playerAt reports the grid position of the player glyph in a rendered view.
// Line 0 is the HUD status bar, so map rows start at line index 1.
func playerAt(t *testing.T, lines []string) Position {
	t.Helper()
	for lineIdx := 1; lineIdx < len(lines); lineIdx++ {
		plain := stripANSI(lines[lineIdx])
		if x := strings.IndexRune(plain, '@'); x >= 0 {
			return Position{X: x, Y: lineIdx - 1}
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
	// Line 0 is HUD; lines 1..10 are map rows (10 rows total)
	if got := len(lines); got < 11 {
		t.Fatalf("view has %d total lines, want at least 11 (HUD + 10 map rows)", got)
	}
	// Check HUD
	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "HP: 100/100") {
		t.Errorf("HUD line = %q, want 'HP: 100/100'", plainHUD)
	}

	mapLines := lines[1:11]
	for y, line := range mapLines {
		if got, want := len([]rune(stripANSI(line))), 20; got != want {
			t.Errorf("map row %d has %d cells, want %d", y, got, want)
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
		lineIdx := 1 + e.Pos.Y
		if lineIdx >= len(lines) {
			t.Fatalf("monster pos Y %d out of bounds for lines len %d", e.Pos.Y, len(lines))
		}
		plainRow := []rune(stripANSI(lines[lineIdx]))
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

	// Map row Y=2 is line index 1+2=3
	plainRow := stripANSI(lines[3])
	if got := string([]rune(plainRow)[1]); got != "o" {
		t.Errorf("expected Orc 'o' at (1, 2), got %q in line: %q", got, plainRow)
	}
}

func TestViewRendersHUDAndLog(t *testing.T) {
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

	// Initially HUD has HP: 100/100, log is empty
	lines := d.Lines()
	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "HP: 100/100") {
		t.Errorf("expected HP: 100/100 in HUD line, got %q", plainHUD)
	}

	// Bump attack goblin (kills 10 HP goblin)
	d.Key("up")

	linesAfter := d.Lines()
	fullText := stripANSI(strings.Join(linesAfter, "\n"))
	if !strings.Contains(fullText, "Player hits Goblin for 10 damage.") {
		t.Errorf("expected hit message in rendered log, got:\n%s", fullText)
	}
	if !strings.Contains(fullText, "Goblin dies.") {
		t.Errorf("expected death message in rendered log, got:\n%s", fullText)
	}
}

func TestGameOverAndRestartViaTUI(t *testing.T) {
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
				Health:    5, // Low health
				MaxHealth: 100,
			},
			Entities: []Entity{
				NewOrc(1, Position{X: 2, Y: 1}), // Orc deals 6 damage
			},
		},
	}

	d := tuitest.New(t, m)

	// Player attacks Orc, Orc counter-attacks for 6 damage, killing Player (5 HP -> 0 HP)
	d.Key("up")

	linesDead := d.Lines()
	plainHUD := stripANSI(linesDead[0])
	if !strings.Contains(plainHUD, "HP: 0/100") {
		t.Errorf("expected HP: 0/100 in dead HUD, got %q", plainHUD)
	}
	if !strings.Contains(plainHUD, "GAME OVER") {
		t.Errorf("expected GAME OVER in dead HUD, got %q", plainHUD)
	}
	if !strings.Contains(plainHUD, "Press r to restart") {
		t.Errorf("expected restart prompt in dead HUD, got %q", plainHUD)
	}

	fullDeadText := stripANSI(strings.Join(linesDead, "\n"))
	if !strings.Contains(fullDeadText, "Player dies.") {
		t.Errorf("expected 'Player dies.' in log, got:\n%s", fullDeadText)
	}

	// Movement key while dead does not act or revive
	d.Key("up")
	linesStillDead := d.Lines()
	if !strings.Contains(stripANSI(linesStillDead[0]), "GAME OVER") {
		t.Errorf("expected still GAME OVER after movement key while dead")
	}

	// Press 'r' to restart game
	d.Key("r")
	linesRestarted := d.Lines()
	plainHUDRestarted := stripANSI(linesRestarted[0])
	if !strings.Contains(plainHUDRestarted, "HP: 100/100") {
		t.Errorf("expected HP: 100/100 after restart, got %q", plainHUDRestarted)
	}
	if strings.Contains(plainHUDRestarted, "GAME OVER") {
		t.Errorf("did not expect GAME OVER after restart")
	}
}
