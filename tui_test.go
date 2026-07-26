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
	d := tuitest.New(t, initialModel())

	lines := d.Lines()
	if got, want := len(lines), 10; got != want {
		t.Fatalf("view has %d rows, want %d", got, want)
	}
	for y, line := range lines {
		if got, want := len([]rune(stripANSI(line))), 20; got != want {
			t.Errorf("row %d has %d cells, want %d", y, got, want)
		}
	}
	if got, want := playerAt(t, lines), (Position{X: 10, Y: 5}); got != want {
		t.Errorf("player rendered at %+v, want %+v", got, want)
	}
}

func TestMovementKeys(t *testing.T) {
	cases := []struct {
		name string
		keys []string
		want Position
	}{
		{"arrow up", []string{"up"}, Position{X: 10, Y: 4}},
		{"arrow down", []string{"down"}, Position{X: 10, Y: 6}},
		{"arrow left", []string{"left"}, Position{X: 9, Y: 5}},
		{"arrow right", []string{"right"}, Position{X: 11, Y: 5}},
		{"wasd", []string{"w", "a"}, Position{X: 9, Y: 4}},
		{"round trip", []string{"w", "s", "a", "d"}, Position{X: 10, Y: 5}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := tuitest.New(t, initialModel())
			d.Keys(tc.keys...)
			if got := playerAt(t, d.Lines()); got != tc.want {
				t.Errorf("after %v player at %+v, want %+v", tc.keys, got, tc.want)
			}
		})
	}
}

func TestPlayerCannotWalkThroughWalls(t *testing.T) {
	d := tuitest.New(t, initialModel())

	// Ten presses is more than enough to reach any border of the 20x10 room.
	for i := 0; i < 20; i++ {
		d.Key("up")
	}
	if got := playerAt(t, d.Lines()).Y; got != 1 {
		t.Errorf("player stopped at y=%d, want 1 (just inside the top wall)", got)
	}

	for i := 0; i < 20; i++ {
		d.Key("left")
	}
	if got := playerAt(t, d.Lines()).X; got != 1 {
		t.Errorf("player stopped at x=%d, want 1 (just inside the left wall)", got)
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
