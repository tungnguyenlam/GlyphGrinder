package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"glyphgrinder/internal/tuitest"
)

func TestDefaultPaletteTokens(t *testing.T) {
	pal := DefaultPalette()

	tokens := map[string]lipgloss.CompleteColor{
		"Player":     pal.Player,
		"Goblin":     pal.Goblin,
		"Orc":        pal.Orc,
		"Troll":      pal.Troll,
		"Archer":     pal.Archer,
		"FloorLit":   pal.FloorLit,
		"WallLit":    pal.WallLit,
		"FloorDim":   pal.FloorDim,
		"WallDim":    pal.WallDim,
		"Stairs":     pal.Stairs,
		"Potion":     pal.Potion,
		"Weapon":     pal.Weapon,
		"Amulet":     pal.Amulet,
		"HUDNormal":  pal.HUDNormal,
		"HUDWarning": pal.HUDWarning,
		"HUDLog":     pal.HUDLog,
	}

	for name, tok := range tokens {
		if tok.TrueColor == "" {
			t.Errorf("token %s missing TrueColor", name)
		}
		if tok.ANSI256 == "" {
			t.Errorf("token %s missing ANSI256", name)
		}
		if tok.ANSI == "" {
			t.Errorf("token %s missing ANSI", name)
		}
	}
}

func TestResolveEntityColor(t *testing.T) {
	pal := DefaultPalette()

	player := Entity{IsPlayer: true, Color: "#FFFFFF"}
	if got := ResolveEntityColor(player, pal); got != pal.Player {
		t.Errorf("player color = %v, want %v", got, pal.Player)
	}

	goblin := Entity{Name: "Goblin", Color: "#FFFFFF"}
	if got := ResolveEntityColor(goblin, pal); got != pal.Goblin {
		t.Errorf("goblin color = %v, want %v", got, pal.Goblin)
	}

	orc := Entity{Name: "Orc", Color: "#FFFFFF"}
	if got := ResolveEntityColor(orc, pal); got != pal.Orc {
		t.Errorf("orc color = %v, want %v", got, pal.Orc)
	}

	troll := Entity{Name: "Troll", Color: "#FFFFFF"}
	if got := ResolveEntityColor(troll, pal); got != pal.Troll {
		t.Errorf("troll color = %v, want %v", got, pal.Troll)
	}

	archer := Entity{Name: "Archer", Color: "#FFFFFF"}
	if got := ResolveEntityColor(archer, pal); got != pal.Archer {
		t.Errorf("archer color = %v, want %v", got, pal.Archer)
	}

	custom := Entity{Name: "Dragon", Color: "#FF00FF"}
	if got := ResolveEntityColor(custom, pal); got != lipgloss.Color("#FF00FF") {
		t.Errorf("custom entity color = %v, want %v", got, lipgloss.Color("#FF00FF"))
	}
}

func TestResolveItemColor(t *testing.T) {
	pal := DefaultPalette()

	pot := Item{ItemType: ItemPotion}
	if got := ResolveItemColor(pot, pal); got != pal.Potion {
		t.Errorf("potion color = %v, want %v", got, pal.Potion)
	}

	weap := Item{ItemType: ItemWeapon}
	if got := ResolveItemColor(weap, pal); got != pal.Weapon {
		t.Errorf("weapon color = %v, want %v", got, pal.Weapon)
	}

	amu := Item{ItemType: ItemAmulet}
	if got := ResolveItemColor(amu, pal); got != pal.Amulet {
		t.Errorf("amulet color = %v, want %v", got, pal.Amulet)
	}
}

func TestColorProfileRendering(t *testing.T) {
	// Restore default profile after test
	defer lipgloss.SetColorProfile(lipgloss.ColorProfile())

	profiles := []struct {
		name    string
		profile termenv.Profile
	}{
		{"TrueColor", termenv.TrueColor},
		{"ANSI256", termenv.ANSI256},
		{"ANSI", termenv.ANSI},
		{"Ascii", termenv.Ascii},
	}

	for _, tc := range profiles {
		t.Run(tc.name, func(t *testing.T) {
			lipgloss.SetColorProfile(tc.profile)

			m := initialModelWithSeed(12345)
			d := tuitest.New(t, m)

			view := d.View()
			if view == "" {
				t.Fatalf("View() output is empty for profile %s", tc.name)
			}

			lines := d.Lines()
			if len(lines) < 2 {
				t.Fatalf("View() lines len = %d, want >= 2 for profile %s", len(lines), tc.name)
			}

			// Verify HUD status line is plain text equivalent
			hudPlain := stripANSI(lines[0])
			if !strings.Contains(hudPlain, "HP: 100/100") {
				t.Errorf("HUD line = %q, want 'HP: 100/100' for profile %s", hudPlain, tc.name)
			}
		})
	}
}
