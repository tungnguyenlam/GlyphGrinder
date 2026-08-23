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
		if got, want := len([]rune(stripANSI(line))), state.Map.Width; got != want {
			t.Errorf("row %d has %d cells, want %d", y, got, want)
		}
	}
	if got, want := playerAt(t, lines), state.Player.Pos; got != want {
		t.Errorf("player rendered at %+v, want %+v", got, want)
	}
}

func TestViewRendersMonstersAtStatePositions(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	state := d.Model().(model).state
	lines := d.Lines()

	for _, monster := range state.Entities {
		plain := []rune(stripANSI(lines[monster.Pos.Y]))
		if got := string(plain[monster.Pos.X]); got != monster.Rune {
			t.Errorf("monster %d rendered as %q at %+v, want %q", monster.ID, got, monster.Pos, monster.Rune)
		}
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
