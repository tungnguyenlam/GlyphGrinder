package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"glyphgrinder/internal/tuitest"
)

// -update-golden regenerates files under testdata/. Usage:
//
//	go test -run Golden -update-golden
//
// Only this package defines the flag; other packages ignore it.
var updateGolden = flag.Bool("update-golden", false, "regenerate golden frames under testdata/")

// goldenModel builds a deterministic headless model: fixed seed, ASCII glyphs,
// default palette, fixed 80×24 terminal. Callers set screen / drive keys.
func goldenModel(seed int64) model {
	st := NewGameWithSeed(defaultMapWidth, defaultMapHeight, seed)
	return model{
		state:          st,
		width:          80,
		height:         24,
		palette:        DefaultPalette(),
		glyphs:         ASCIIGlyphs(),
		camX:           float64(st.Player.Pos.X),
		camY:           float64(st.Player.Pos.Y),
		camInitialized: true,
	}
}

// withStableRender locks Lip Gloss to TrueColor and clears the style cache so
// golden frames do not depend on the host terminal profile or prior tests.
func withStableRender(t *testing.T) {
	t.Helper()
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	globalRenderCache = nil
	t.Cleanup(func() {
		lipgloss.SetColorProfile(prev)
		globalRenderCache = nil
	})
}

func goldenPath(name string) string {
	return filepath.Join("testdata", name+".golden")
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := goldenPath(name)

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		t.Logf("updated golden %s (%d bytes)", path, len(got))
		return
	}

	wantBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v\nre-generate with: go test -run Golden -update-golden", path, err)
	}
	want := string(wantBytes)
	if got == want {
		return
	}

	// Diff-friendly failure: show plain-text delta first (ANSI noise is hard to read).
	gotPlain := stripANSI(got)
	wantPlain := stripANSI(want)
	if gotPlain != wantPlain {
		t.Errorf("golden %s plain-text mismatch:\n--- want ---\n%s\n--- got ---\n%s", name, wantPlain, gotPlain)
		return
	}
	t.Errorf("golden %s styled (ANSI) mismatch; plain text matches — color profile or style cache drift?\n  file: %s\n  want %d bytes, got %d bytes\n  re-generate with: go test -run Golden -update-golden",
		name, path, len(want), len(got))
}

func TestGoldenTitleScreen(t *testing.T) {
	withStableRender(t)

	m := goldenModel(42)
	m.screen = screenTitle
	d := tuitest.New(t, m)

	got := d.View()
	if !strings.Contains(stripANSI(got), "SELECT CLASS ARCHETYPE TO BEGIN DESCENT") {
		t.Fatalf("title golden source frame missing expected subtitle:\n%s", stripANSI(got))
	}
	assertGolden(t, "title_screen", got)
}

func TestGoldenInGameWarrior(t *testing.T) {
	withStableRender(t)

	// Title → Warrior (keeps seed 42) → fixed walk so the camera settles on a
	// known frame without relying on random combat outcomes.
	m := goldenModel(42)
	m.screen = screenTitle
	d := tuitest.New(t, m)
	d.Keys("1", "w", "w", "d", "d")

	got := d.View()
	plain := stripANSI(got)
	if !strings.Contains(plain, "HP:") {
		t.Fatalf("in-game golden source missing HUD:\n%s", plain)
	}
	if !strings.Contains(plain, "Seed: 42") {
		t.Fatalf("in-game golden source missing Seed: 42:\n%s", plain)
	}
	if !strings.Contains(plain, "Warrior") && !strings.Contains(plain, "Iron Dagger") {
		// Inventory line carries the class starter weapon; either is fine signal.
		if !strings.Contains(plain, "Inv:") {
			t.Fatalf("in-game golden source missing Warrior inventory signal:\n%s", plain)
		}
	}

	// Camera must have finished easing or the golden would be frame-flaky.
	final := d.Model().(model)
	if final.isAnimating() || final.flashTicks > 0 {
		t.Fatalf("frame still animating (animTicks=%d flashTicks=%d cam=(%f,%f) player=%+v)",
			final.animTicks, final.flashTicks, final.camX, final.camY, final.state.Player.Pos)
	}

	assertGolden(t, "ingame_warrior_seed42", got)
}
