package main

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"glyphgrinder/internal/tuitest"
)

func init() {
	tickInterval = 0
}

// viewOrigin returns the top-left map coordinate of the rendered viewport
// for the given model. Tests use this to translate screen→map coordinates.
func viewOrigin(m model) (x0, y0 int) {
	vw := m.width
	vh := m.height - reservedRows
	if vw <= 0 {
		vw = m.state.Map.Width
	}
	if vh <= 0 {
		vh = m.state.Map.Height
	}
	camX, camY := m.getCamPos()
	x0, y0, _, _ = viewportCenter(camX, camY, m.state.Map.Width, m.state.Map.Height, vw, vh)
	return x0, y0
}

// playerAt reports the map-grid position of the player glyph in a rendered view.
// ox, oy is the camera origin (top-left map coordinate of the viewport).
// Line 0 is the HUD status bar, so map rows start at line index 1.
func playerAt(t *testing.T, lines []string, ox, oy int) Position {
	t.Helper()
	for lineIdx := 1; lineIdx < len(lines); lineIdx++ {
		plain := stripANSI(lines[lineIdx])
		x := strings.IndexRune(plain, '@')
		if x < 0 {
			x = strings.IndexRune(plain, '󰋋')
		}
		if x >= 0 {
			return Position{X: x + ox, Y: (lineIdx - 1) + oy}
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
	// With viewport, the rendered map may be clipped; just verify basics.
	// HUD is line 0, then viewport rows follow.
	if got := len(lines); got < 2 {
		t.Fatalf("view has %d total lines, want at least 2 (HUD + map)", got)
	}
	// Check HUD
	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "HP: 100/100") {
		t.Errorf("HUD line = %q, want 'HP: 100/100'", plainHUD)
	}

	// All rendered map rows should have the same width (the viewport width).
	vw := m.width
	if vw > m.state.Map.Width {
		vw = m.state.Map.Width
	}
	for i := 1; i < len(lines); i++ {
		plain := stripANSI(lines[i])
		if len(plain) == 0 {
			break // log lines follow
		}
		if got := len([]rune(plain)); got != vw {
			t.Errorf("map row %d has %d cells, want %d", i-1, got, vw)
		}
	}

	ox, oy := viewOrigin(m)
	wantPos := m.state.Player.Pos
	if got := playerAt(t, lines, ox, oy); got != wantPos {
		t.Errorf("player rendered at %+v, want %+v", got, wantPos)
	}
}

func TestMovementKeys(t *testing.T) {
	makeTestModel := func() model {
		return model{
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
					Health:    100,
					MaxHealth: 100,
				},
			},
		}
	}

	startPos := Position{X: 2, Y: 2}

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
			d := tuitest.New(t, makeTestModel())
			d.Keys(tc.keys...)
			mdl := d.Model().(model)
			ox, oy := viewOrigin(mdl)
			if got := playerAt(t, d.Lines(), ox, oy); got != tc.want {
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
	mdl := d.Model().(model)
	ox, oy := viewOrigin(mdl)
	posUp := playerAt(t, d.Lines(), ox, oy)
	if posUp.Y <= 0 || m.state.Map.Tiles[posUp.Y-1][posUp.X] != TileWall {
		t.Errorf("expected tile above player at %+v to be TileWall", posUp)
	}

	// Pressing left many times will hit a wall.
	for i := 0; i < 20; i++ {
		d.Key("left")
	}
	mdl = d.Model().(model)
	ox, oy = viewOrigin(mdl)
	posLeft := playerAt(t, d.Lines(), ox, oy)
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
	ox, oy := viewOrigin(m)
	vw := m.width
	vh := m.height - reservedRows

	visibleMonsterChecked := false
	for _, e := range m.state.Entities {
		// Only check monsters within the rendered viewport window AND visible in FOV
		if e.Pos.X >= ox && e.Pos.X < ox+vw && e.Pos.Y >= oy && e.Pos.Y < oy+vh {
			if m.state.Map.Visible != nil && m.state.Map.Visible[e.Pos.Y][e.Pos.X] {
				visibleMonsterChecked = true
				lineIdx := 1 + (e.Pos.Y - oy)
				if lineIdx >= len(lines) {
					t.Fatalf("monster pos Y %d out of bounds for lines len %d", e.Pos.Y, len(lines))
				}
				plainRow := []rune(stripANSI(lines[lineIdx]))
				colIdx := e.Pos.X - ox
				if colIdx >= len(plainRow) {
					t.Fatalf("monster pos X %d out of bounds for plain row len %d", e.Pos.X, len(plainRow))
				}
				gotRune := string(plainRow[colIdx])
				wantGlyph := ResolveEntityGlyph(e, m.getGlyphs())
				if gotRune != wantGlyph {
					t.Errorf("expected monster rune %q at %+v (screen row %d col %d), got %q", wantGlyph, e.Pos, lineIdx, colIdx, gotRune)
				}
			}
		}
	}

	if !visibleMonsterChecked {
		t.Log("Note: no monsters spawned within initial FOV of player for seed 12345")
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

	pos := playerAt(t, d.Lines(), 0, 0)
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
	pos := playerAt(t, lines, 0, 0)
	if pos != (Position{X: 2, Y: 1}) {
		t.Errorf("player position = %+v, want (2, 1)", pos)
	}

	// Map row Y=2 is line index 1+2=3
	plainRow := stripANSI(lines[3])
	wantOrcGlyph := ResolveEntityGlyph(NewOrc(1, Position{}), m.getGlyphs())
	if got := string([]rune(plainRow)[1]); got != wantOrcGlyph {
		t.Errorf("expected Orc %q at (1, 2), got %q in line: %q", wantOrcGlyph, got, plainRow)
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

func TestFOVUnexploredTilesRenderBlank(t *testing.T) {
	// Build a state with FOV initialized. Tiles outside FOV radius
	// should render as spaces.
	state := NewGameWithSeed(20, 10, 12345)
	m := model{state: state}
	d := tuitest.New(t, m)

	lines := d.Lines()
	// Map lines are 1..10 (after HUD).
	// At least some tiles far from the player should be unexplored (space).
	foundBlank := false
	for y := 1; y <= state.Map.Height; y++ {
		if y >= len(lines) {
			break
		}
		plain := stripANSI(lines[y])
		for _, r := range plain {
			if r == ' ' {
				foundBlank = true
				break
			}
		}
		if foundBlank {
			break
		}
	}
	if !foundBlank {
		t.Error("expected at least one blank (unexplored) tile in the initial view")
	}
}

func TestFOVHidesNonVisibleMonsters(t *testing.T) {
	// Place a monster outside the player's FOV. It should not be
	// rendered in the view.
	gm := GameMap{
		Width:  20,
		Height: 10,
		Tiles:  make([][]TileType, 10),
	}
	for y := 0; y < 10; y++ {
		gm.Tiles[y] = make([]TileType, 20)
		for x := 0; x < 20; x++ {
			if x == 0 || x == 19 || y == 0 || y == 9 {
				gm.Tiles[y][x] = TileWall
			} else {
				gm.Tiles[y][x] = TileFloor
			}
		}
	}
	// Put a wall to block sight between player and monster.
	gm.Tiles[4][10] = TileWall
	gm.Tiles[5][10] = TileWall
	gm.Tiles[3][10] = TileWall

	playerPos := Position{X: 5, Y: 5}
	monsterPos := Position{X: 15, Y: 5} // far away, behind walls

	gm.Explored = makeBoolGrid(20, 10)
	gm.ComputeFOV(playerPos, FOVRadius)

	// Verify the monster position is not visible.
	if gm.Visible[monsterPos.Y][monsterPos.X] {
		t.Skip("monster position is visible in this layout; test layout needs adjustment")
	}

	state := GameState{
		Map: gm,
		Player: Entity{
			ID: 0, Name: "Player", IsPlayer: true, Pos: playerPos,
			Rune: "@", Color: "#00FF00", Health: 100, MaxHealth: 100, Damage: 10,
		},
		Entities: []Entity{NewGoblin(1, monsterPos)},
	}

	m := model{state: state}
	d := tuitest.New(t, m)
	lines := d.Lines()

	// Monster at (15,5) → line index 1+5=6.
	if 6 >= len(lines) {
		t.Fatal("not enough lines for monster row")
	}
	plainRow := []rune(stripANSI(lines[6]))
	if 15 < len(plainRow) {
		got := string(plainRow[15])
		if got == "g" {
			t.Errorf("goblin at (15,5) should not be rendered when not visible, got %q", got)
		}
	}
}

func TestFOVExploredTilesShowDimmed(t *testing.T) {
	// After the player moves away from a tile, it should still render
	// (as dimmed) rather than blank space, because it was explored.
	state := NewGameWithSeed(20, 10, 12345)

	// Move the player a few steps to explore some tiles, then check
	// that the original position is still rendered (not blank).
	startPos := state.Player.Pos

	m := model{state: state}
	d := tuitest.New(t, m)

	// Move in one direction several times to shift the FOV.
	for i := 0; i < 3; i++ {
		d.Key("right")
	}

	lines := d.Lines()
	mdl := d.Model().(model)

	// Check the starting position: it should be explored.
	if !mdl.state.Map.Explored[startPos.Y][startPos.X] {
		t.Error("original player position should remain explored")
	}

	// If the start tile is no longer visible, it should render as '.'
	// or '#' (dimmed), not as ' ' (blank).
	lineIdx := 1 + startPos.Y
	if lineIdx < len(lines) {
		plain := []rune(stripANSI(lines[lineIdx]))
		if startPos.X < len(plain) {
			r := plain[startPos.X]
			if r == ' ' {
				t.Errorf("explored tile at start pos %+v rendered as blank, expected dimmed glyph", startPos)
			}
		}
	}
}

func TestViewportCalculation(t *testing.T) {
	cases := []struct {
		name                           string
		playerPos                      Position
		mapW, mapH                     int
		viewW, viewH                   int
		wantX0, wantY0, wantX1, wantY1 int
	}{
		{
			name:      "view larger than map",
			playerPos: Position{X: 10, Y: 5},
			mapW:      20, mapH: 10,
			viewW: 80, viewH: 24,
			wantX0: 0, wantY0: 0, wantX1: 20, wantY1: 10,
		},
		{
			name:      "centered in middle of map",
			playerPos: Position{X: 30, Y: 15},
			mapW:      60, mapH: 30,
			viewW: 20, viewH: 10,
			wantX0: 20, wantY0: 10, wantX1: 40, wantY1: 20,
		},
		{
			name:      "clamped to top-left border",
			playerPos: Position{X: 2, Y: 2},
			mapW:      60, mapH: 30,
			viewW: 20, viewH: 10,
			wantX0: 0, wantY0: 0, wantX1: 20, wantY1: 10,
		},
		{
			name:      "clamped to bottom-right border",
			playerPos: Position{X: 58, Y: 28},
			mapW:      60, mapH: 30,
			viewW: 20, viewH: 10,
			wantX0: 40, wantY0: 20, wantX1: 60, wantY1: 30,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x0, y0, x1, y1 := viewport(tc.playerPos, tc.mapW, tc.mapH, tc.viewW, tc.viewH)
			if x0 != tc.wantX0 || y0 != tc.wantY0 || x1 != tc.wantX1 || y1 != tc.wantY1 {
				t.Errorf("viewport(%+v, %d, %d, %d, %d) = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					tc.playerPos, tc.mapW, tc.mapH, tc.viewW, tc.viewH,
					x0, y0, x1, y1, tc.wantX0, tc.wantY0, tc.wantX1, tc.wantY1)
			}
		})
	}
}

func TestWindowSizeResizeViewport(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	// Resize window to small terminal: 40 cols, 15 rows (viewH = 15 - 6 = 9)
	d.Resize(40, 15)

	mdl := d.Model().(model)
	if mdl.width != 40 || mdl.height != 15 {
		t.Errorf("model dims = (%d, %d), want (40, 15)", mdl.width, mdl.height)
	}

	lines := d.Lines()
	// Line 0 is HUD. Map lines should be 9 rows long.
	mapRows := 0
	for i := 1; i < len(lines); i++ {
		plain := stripANSI(lines[i])
		if len(plain) == 0 {
			break
		}
		if got := len([]rune(plain)); got != 40 {
			t.Errorf("row %d width = %d, want 40", mapRows, got)
		}
		mapRows++
	}

	if mapRows != 9 {
		t.Errorf("rendered map rows = %d, want 9 (15 total height - 6 reserved)", mapRows)
	}
}

func TestCameraFollowsPlayer(t *testing.T) {
	// Create a large 60x30 map with an open floor area around (30, 15)
	tiles := make([][]TileType, 30)
	for y := 0; y < 30; y++ {
		tiles[y] = make([]TileType, 60)
		for x := 0; x < 60; x++ {
			if x == 0 || x == 59 || y == 0 || y == 29 {
				tiles[y][x] = TileWall
			} else {
				tiles[y][x] = TileFloor
			}
		}
	}

	gm := GameMap{
		Width:  60,
		Height: 30,
		Tiles:  tiles,
	}
	gm.Explored = makeBoolGrid(60, 30)
	gm.ComputeFOV(Position{X: 30, Y: 15}, FOVRadius)

	m := model{
		state: GameState{
			Map: gm,
			Player: Entity{
				ID: 0, Name: "Player", IsPlayer: true,
				Pos:  Position{X: 30, Y: 15},
				Rune: "@", Color: "#00FF00", Health: 100, MaxHealth: 100, Damage: 10,
			},
		},
		width:  20,
		height: 16, // viewW = 20, viewH = 16 - 6 = 10
	}

	d := tuitest.New(t, m)

	ox0, oy0 := viewOrigin(d.Model().(model))
	if ox0 != 20 || oy0 != 10 {
		t.Fatalf("initial camera origin = (%d, %d), want (20, 10)", ox0, oy0)
	}

	// Move player 5 steps right to (35, 15)
	for i := 0; i < 5; i++ {
		d.Key("right")
	}

	mdl := d.Model().(model)
	if mdl.state.Player.Pos != (Position{X: 35, Y: 15}) {
		t.Fatalf("player pos = %+v, want (35, 15)", mdl.state.Player.Pos)
	}

	ox1, oy1 := viewOrigin(mdl)
	if ox1 != 25 || oy1 != 10 {
		t.Errorf("camera origin after moving right = (%d, %d), want (25, 10)", ox1, oy1)
	}

	// Rendered view must still display player '@'
	lines := d.Lines()
	gotPos := playerAt(t, lines, ox1, oy1)
	if gotPos != (Position{X: 35, Y: 15}) {
		t.Errorf("playerAt returned %+v, want (35, 15)", gotPos)
	}
}

func makeAnimationTestModel() model {
	return model{
		state: GameState{
			Map: GameMap{
				Width:  10,
				Height: 10,
				Tiles: [][]TileType{
					{TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileFloor, TileWall},
					{TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall, TileWall},
				},
			},
			Player: Entity{
				ID:        0,
				Name:      "Player",
				IsPlayer:  true,
				Pos:       Position{X: 2, Y: 2},
				Rune:      "@",
				Color:     "#00FF00",
				Health:    100,
				MaxHealth: 100,
			},
		},
	}
}

func TestTickAnimationStateAndEasing(t *testing.T) {
	m := makeAnimationTestModel()
	if m.isAnimating() {
		t.Error("initial model should not be animating")
	}

	startPos := m.state.Player.Pos

	tickInterval = 16 * time.Millisecond
	defer func() { tickInterval = 0 }()

	nextMdl, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = nextMdl.(model)

	if cmd == nil {
		t.Fatal("Update(moveKey) should return a tick command when animating")
	}

	if !m.isAnimating() {
		t.Error("model should be animating after movement action")
	}

	wantPos := Position{X: startPos.X + 1, Y: startPos.Y}
	if m.state.Player.Pos != wantPos {
		t.Errorf("player position = %+v, want %+v", m.state.Player.Pos, wantPos)
	}

	frames := 0
	for m.isAnimating() {
		frames++
		if frames > 100 {
			t.Fatal("animation failed to converge within 100 ticks")
		}
		nextM, _ := m.Update(animTickMsg(time.Now()))
		m = nextM.(model)
	}

	camX, camY := m.getCamPos()
	if camX != float64(wantPos.X) || camY != float64(wantPos.Y) {
		t.Errorf("final camera pos = (%f, %f), want (%d, %d)", camX, camY, wantPos.X, wantPos.Y)
	}
}

func TestInputResponsivenessDuringTicks(t *testing.T) {
	tickInterval = 16 * time.Millisecond
	defer func() { tickInterval = 0 }()

	m := makeAnimationTestModel()
	startPos := m.state.Player.Pos

	nextM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = nextM.(model)

	nextM, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = nextM.(model)

	wantPos := Position{X: startPos.X + 1, Y: startPos.Y + 1}
	if m.state.Player.Pos != wantPos {
		t.Errorf("player position after rapid inputs = %+v, want %+v", m.state.Player.Pos, wantPos)
	}
}

func TestHUDDisplaysDepth(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	lines := d.Lines()
	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "Depth: 1") {
		t.Errorf("expected 'Depth: 1' in HUD line, got %q", plainHUD)
	}
}

func TestDescendLevelViaTUI(t *testing.T) {
	m := initialModelWithSeed(12345)
	var stairsPos Position
	found := false
	for y := 0; y < m.state.Map.Height; y++ {
		for x := 0; x < m.state.Map.Width; x++ {
			if m.state.Map.Tiles[y][x] == TileStairsDown {
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
		t.Fatal("no stairs in initial map")
	}

	m.state.Player.Pos = stairsPos
	d := tuitest.New(t, m)

	d.Key(">")

	mdl := d.Model().(model)
	if mdl.state.Depth != 2 {
		t.Errorf("model depth after descend = %d, want 2", mdl.state.Depth)
	}

	lines := d.Lines()
	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "Depth: 2") {
		t.Errorf("expected 'Depth: 2' in HUD line after descend, got %q", plainHUD)
	}
}

func TestViewRendersFloorItemsAndInventory(t *testing.T) {
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
				Health:    100,
				MaxHealth: 100,
			},
			Items: []Item{
				NewHealthPotion(1, Position{X: 2, Y: 1}),
			},
		},
	}

	d := tuitest.New(t, m)
	lines := d.Lines()
	plainHUD := stripANSI(lines[0])
	if strings.Contains(plainHUD, "Inv:") {
		t.Errorf("expected empty inventory in HUD initially, got %q", plainHUD)
	}

	// Move right onto potion tile (2, 1)
	d.Key("right")
	linesOnItem := d.Lines()
	fullText := stripANSI(strings.Join(linesOnItem, "\n"))
	if !strings.Contains(fullText, "You see a Health Potion here.") {
		t.Errorf("expected item prompt in log, got:\n%s", fullText)
	}

	// Press 'g' to pick up potion
	d.Key("g")
	linesAfterPickup := d.Lines()
	plainHUDAfter := stripANSI(linesAfterPickup[0])
	if !strings.Contains(plainHUDAfter, "Inv: [Health Potion]") {
		t.Errorf("expected 'Inv: [Health Potion]' in HUD line, got %q", plainHUDAfter)
	}
}

func TestPotionUseViaTUI(t *testing.T) {
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
				Health:    60,
				MaxHealth: 100,
				Inventory: []Item{NewHealthPotion(1, Position{X: -1, Y: -1})},
			},
		},
	}

	d := tuitest.New(t, m)
	linesBefore := d.Lines()
	plainHUDBefore := stripANSI(linesBefore[0])
	if !strings.Contains(plainHUDBefore, "HP: 60/100") {
		t.Errorf("expected HP: 60/100 before drinking potion, got %q", plainHUDBefore)
	}

	// Press 'h' to drink potion
	d.Key("h")
	linesAfter := d.Lines()
	plainHUDAfter := stripANSI(linesAfter[0])
	if !strings.Contains(plainHUDAfter, "HP: 85/100") {
		t.Errorf("expected HP: 85/100 after drinking potion, got %q", plainHUDAfter)
	}

	fullText := stripANSI(strings.Join(linesAfter, "\n"))
	if !strings.Contains(fullText, "You drink the Health Potion and recover 25 HP.") {
		t.Errorf("expected potion log message, got:\n%s", fullText)
	}
}

func TestArcherRangedAttackViaTUI(t *testing.T) {
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
	archerPos := Position{X: 4, Y: 2}
	gm.Explored = makeBoolGrid(10, 5)
	gm.ComputeFOV(playerPos, FOVRadius)

	m := model{
		state: GameState{
			Map: gm,
			Player: Entity{
				ID: 0, Name: "Player", IsPlayer: true, Pos: playerPos,
				Rune: "@", Color: "#00FF00", Health: 100, MaxHealth: 100, Damage: 10,
			},
			Entities: []Entity{
				NewArcher(1, archerPos),
			},
		},
	}

	d := tuitest.New(t, m)

	// Step down to (1, 3). Archer at (4, 2) shoots player from range.
	d.Key("down")

	lines := d.Lines()
	fullText := stripANSI(strings.Join(lines, "\n"))
	if !strings.Contains(fullText, "Archer shoots Player for 4 damage.") {
		t.Errorf("expected archer shoot message in TUI log, got:\n%s", fullText)
	}

	plainHUD := stripANSI(lines[0])
	if !strings.Contains(plainHUD, "HP: 96/100") {
		t.Errorf("expected HP: 96/100 after archer ranged hit, got %q", plainHUD)
	}
}

func TestVictoryStateAndRestartViaTUI(t *testing.T) {
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
				ID: 0, Name: "Player", IsPlayer: true, Pos: Position{X: 2, Y: 2},
				Rune: "@", Color: "#00FF00", Health: 100, MaxHealth: 100,
			},
			Items: []Item{
				NewAmuletOfYendor(1, Position{X: 2, Y: 2}),
			},
		},
	}

	d := tuitest.New(t, m)

	// Press 'g' to pick up Amulet of Yendor
	d.Key("g")

	linesWon := d.Lines()
	plainHUDWon := stripANSI(linesWon[0])
	if !strings.Contains(plainHUDWon, "VICTORY!") {
		t.Errorf("expected VICTORY! in HUD line after picking up Amulet, got %q", plainHUDWon)
	}
	if !strings.Contains(plainHUDWon, "Press r to restart") {
		t.Errorf("expected restart prompt in victory HUD line, got %q", plainHUDWon)
	}

	// Movement key while won does not act or change state
	d.Key("up")
	linesStillWon := d.Lines()
	if !strings.Contains(stripANSI(linesStillWon[0]), "VICTORY!") {
		t.Errorf("expected still VICTORY after movement key while won")
	}

	// Press 'r' to restart game
	d.Key("r")
	linesRestarted := d.Lines()
	plainHUDRestarted := stripANSI(linesRestarted[0])
	if strings.Contains(plainHUDRestarted, "VICTORY!") {
		t.Errorf("did not expect VICTORY after restart")
	}
	if !strings.Contains(plainHUDRestarted, "HP: 100/100") {
		t.Errorf("expected fresh HP: 100/100 after restart, got %q", plainHUDRestarted)
	}
}

func TestExtremeViewportDimensions(t *testing.T) {
	m := initialModelWithSeed(12345)
	d := tuitest.New(t, m)

	tinyDimensions := []struct {
		w, h int
	}{
		{0, 0},
		{5, 5},
		{10, 5},
		{1, 1},
		{200, 100},
	}

	for _, dim := range tinyDimensions {
		d.Resize(dim.w, dim.h)
		output := d.View()
		if dim.w > 0 && dim.h > 0 && output == "" {
			t.Errorf("View() empty for dimensions (%d, %d)", dim.w, dim.h)
		}
	}
}

func TestScrollUsageViaTUI(t *testing.T) {
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
				ID: 0, Name: "Player", IsPlayer: true, Pos: Position{X: 2, Y: 2},
				Rune: "@", Color: "#00FF00", Health: 100, MaxHealth: 100,
				Inventory: []Item{
					NewFireballScroll(1, Position{X: -1, Y: -1}),
				},
			},
		},
	}

	d := tuitest.New(t, m)

	// Press 'h' to use fireball scroll
	d.Key("h")

	lines := d.Lines()
	fullText := stripANSI(strings.Join(lines, "\n"))

	if !strings.Contains(fullText, "Scroll of Fireball") {
		t.Errorf("expected Fireball scroll log message in TUI view, got:\n%s", fullText)
	}
}
