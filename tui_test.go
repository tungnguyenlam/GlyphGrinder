package main

import (
	"math/rand"
	"strings"
	"testing"

	"glyphgrinder/internal/tuitest"
)

const tuiTestSeed int64 = 42

func newTUITestModel() model {
	return initialModel(rand.New(rand.NewSource(tuiTestSeed)))
}

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
	d := tuitest.New(t, newTUITestModel())
	state := d.Model().(model).state

	lines := d.Lines()
	if got, want := len(lines), state.Map.Height; got != want {
		t.Fatalf("view has %d rows, want %d", got, want)
	}
	for y, line := range lines {
		plain := []rune(stripANSI(line))
		if got, wantAtLeast := len(plain), state.Map.Width; got < wantAtLeast {
			t.Errorf("row %d has %d cells, want at least %d map cells", y, got, wantAtLeast)
			continue
		}
		for x, cell := range plain[:state.Map.Width] {
			if !strings.ContainsRune(" #.@g", cell) {
				t.Errorf("map cell (%d,%d) = %q, want a dungeon glyph", x, y, cell)
			}
		}
	}
	if got, want := playerAt(t, lines), state.Player.Pos; got != want {
		t.Errorf("player rendered at %+v, want %+v", got, want)
	}
}

func TestViewRendersOnlyVisibleMonsters(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	state := d.Model().(model).state
	lines := d.Lines()

	for _, monster := range state.Entities {
		plain := []rune(stripANSI(lines[monster.Pos.Y]))
		got := string(plain[monster.Pos.X])
		if state.Map.Visible[monster.Pos.Y][monster.Pos.X] && got != monster.Rune {
			t.Errorf("visible monster %d rendered as %q at %+v, want %q", monster.ID, got, monster.Pos, monster.Rune)
		}
		if !state.Map.Visible[monster.Pos.Y][monster.Pos.X] && got == monster.Rune {
			t.Errorf("hidden monster %d rendered at %+v", monster.ID, monster.Pos)
		}
	}
}

func TestViewDistinguishesHiddenAndRememberedTiles(t *testing.T) {
	state := openTestStateSized(11, 7)
	state.Player.Pos = Position{X: 2, Y: 3}
	state.Map.Tiles[3][3] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 4, Y: 3})}
	state.Map.Visible = nil
	state.Map.Explored = nil
	state = state.refreshVisibility()
	remembered := Position{X: 9, Y: 3}
	state.Map.Explored[remembered.Y][remembered.X] = true
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})
	lines := d.Lines()

	if got := renderedCell(t, lines, Position{X: 4, Y: 3}); got != " " {
		t.Errorf("unseen monster tile rendered as %q, want blank", got)
	}
	if got := renderedCell(t, lines, remembered); got != "." {
		t.Errorf("remembered floor rendered as %q, want dim floor glyph", got)
	}
}

func renderedCell(t *testing.T, lines []string, pos Position) string {
	t.Helper()
	if pos.Y < 0 || pos.Y >= len(lines) {
		t.Fatalf("row %d outside rendered view", pos.Y)
	}
	plain := []rune(stripANSI(lines[pos.Y]))
	if pos.X < 0 || pos.X >= len(plain) {
		t.Fatalf("column %d outside rendered row %d", pos.X, pos.Y)
	}
	return string(plain[pos.X])
}

func TestViewRendersHealthAndRecentLog(t *testing.T) {
	state := openTestState()
	state.Player.Health = 65
	state.Log = append(state.Log, "A test event.")
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})
	plain := stripANSI(d.View())

	if !strings.Contains(plain, "HP [######....] 65/100") {
		t.Errorf("view does not contain health bar:\n%s", plain)
	}
	if !strings.Contains(plain, "A test event.") {
		t.Errorf("view does not contain recent log entry:\n%s", plain)
	}
}

func TestMovementKeys(t *testing.T) {
	cases := []struct {
		name    string
		keys    []string
		actions []Action
	}{
		{"arrow up", []string{"up"}, []Action{ActionMoveUp}},
		{"arrow down", []string{"down"}, []Action{ActionMoveDown}},
		{"arrow left", []string{"left"}, []Action{ActionMoveLeft}},
		{"arrow right", []string{"right"}, []Action{ActionMoveRight}},
		{"wasd", []string{"w", "a"}, []Action{ActionMoveUp, ActionMoveLeft}},
		{"round trip", []string{"w", "s", "a", "d"}, []Action{ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := tuitest.New(t, newTUITestModel())
			want := d.Model().(model).state
			for _, action := range tc.actions {
				want = want.Step(action)
			}
			d.Keys(tc.keys...)
			if got := playerAt(t, d.Lines()); got != want.Player.Pos {
				t.Errorf("after %v player at %+v, want %+v", tc.keys, got, want.Player.Pos)
			}
		})
	}
}

func TestPlayerCannotWalkThroughWalls(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	want := d.Model().(model).state

	// Repeated moves must stop wherever the generated dungeon blocks the path.
	for i := 0; i < 20; i++ {
		d.Key("up")
		want = want.Step(ActionMoveUp)
	}
	if got := playerAt(t, d.Lines()); got != want.Player.Pos {
		t.Errorf("player stopped at %+v, want %+v", got, want.Player.Pos)
	}

	for i := 0; i < 20; i++ {
		d.Key("left")
		want = want.Step(ActionMoveLeft)
	}
	if got := playerAt(t, d.Lines()); got != want.Player.Pos {
		t.Errorf("player stopped at %+v, want %+v", got, want.Player.Pos)
	}
}

func TestBumpAttackSurvivesUpdateRoundTrip(t *testing.T) {
	state := combatTestState(20)
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	got := d.Model().(model).state
	if got.Entities[0].Health != 10 {
		t.Errorf("target health after key = %d, want 10", got.Entities[0].Health)
	}
	if got.Player.Pos != state.Player.Pos {
		t.Errorf("player moved to %+v during attack, want %+v", got.Player.Pos, state.Player.Pos)
	}
	if gotLog, want := got.Log, "You hit monster 1 for 10 damage."; len(gotLog) < 2 || gotLog[1] != want {
		t.Errorf("log after key = %q, want player event %q first", gotLog, want)
	}

	d.Key("right")
	got = d.Model().(model).state
	if len(got.Entities) != 1 || got.Entities[0].ID != 2 {
		t.Errorf("entities after second key = %+v, want only monster 2", got.Entities)
	}
}

func TestMonsterAttackSurvivesUpdateRoundTrip(t *testing.T) {
	state := openTestState()
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Map.Tiles[3][4] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	got := d.Model().(model).state
	if got.Player.Health != 95 {
		t.Errorf("player health after key = %d, want 95", got.Player.Health)
	}
	if got.Entities[0].Pos != (Position{X: 3, Y: 2}) {
		t.Errorf("attacking monster moved to %+v, want {3 2}", got.Entities[0].Pos)
	}
	if gotLog, want := got.Log, "Monster 1 hits you for 5 damage."; len(gotLog) != 2 || gotLog[1] != want {
		t.Errorf("log after key = %q, want final entry %q", gotLog, want)
	}
}

func TestGameOverAndRestartSurviveUpdateRoundTrip(t *testing.T) {
	state := openTestState()
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Player.Health = 5
	state.Map.Tiles[3][4] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	dead := d.Model().(model).state
	if !dead.GameOver {
		t.Fatal("lethal input did not produce game over")
	}
	plain := stripANSI(d.View())
	if !strings.Contains(plain, "GAME OVER") || !strings.Contains(plain, "Press r to restart") {
		t.Errorf("game-over view is missing restart prompt:\n%s", plain)
	}

	d.Key("r")
	restarted := d.Model().(model).state
	if restarted.GameOver {
		t.Fatal("restart left game in game-over state")
	}
	if restarted.Player.Health != restarted.Player.MaxHealth {
		t.Errorf("restarted health = %d/%d, want full", restarted.Player.Health, restarted.Player.MaxHealth)
	}
	if len(restarted.Entities) == 0 {
		t.Error("restarted dungeon has no monsters")
	}
}

func TestQuitKeys(t *testing.T) {
	for _, key := range []string{"q", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			d := tuitest.New(t, newTUITestModel())
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
