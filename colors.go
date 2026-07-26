package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Palette contains CompleteColor definitions for all game visual tokens.
// Each CompleteColor provides explicit values for TrueColor (24-bit),
// ANSI256 (8-bit), and ANSI (4-bit 16-color) terminal profiles, ensuring
// graceful degradation across terminal capabilities without color loss.
type Palette struct {
	Player     lipgloss.CompleteColor
	Goblin     lipgloss.CompleteColor
	Orc        lipgloss.CompleteColor
	Troll      lipgloss.CompleteColor
	Archer     lipgloss.CompleteColor
	FloorLit   lipgloss.CompleteColor
	WallLit    lipgloss.CompleteColor
	FloorDim   lipgloss.CompleteColor
	WallDim    lipgloss.CompleteColor
	Stairs     lipgloss.CompleteColor
	Potion     lipgloss.CompleteColor
	Weapon     lipgloss.CompleteColor
	HUDNormal  lipgloss.CompleteColor
	HUDWarning lipgloss.CompleteColor
	HUDLog     lipgloss.CompleteColor
}

// DefaultPalette returns the curated GlyphGrinder color theme.
// Lit tiles use warm torchlit tones (#4A3E2D floor, #998574 stone wall).
// Dimmed tiles use cool slate blue-grey memory tones (#1E2530 floor, #3B485E wall).
func DefaultPalette() Palette {
	return Palette{
		Player: lipgloss.CompleteColor{
			TrueColor: "#00FF87",
			ANSI256:   "46",
			ANSI:      "2",
		},
		Goblin: lipgloss.CompleteColor{
			TrueColor: "#85E82A",
			ANSI256:   "118",
			ANSI:      "10",
		},
		Orc: lipgloss.CompleteColor{
			TrueColor: "#FF3B30",
			ANSI256:   "196",
			ANSI:      "1",
		},
		Troll: lipgloss.CompleteColor{
			TrueColor: "#D32F2F",
			ANSI256:   "160",
			ANSI:      "1",
		},
		Archer: lipgloss.CompleteColor{
			TrueColor: "#FF9800",
			ANSI256:   "208",
			ANSI:      "3",
		},
		FloorLit: lipgloss.CompleteColor{
			TrueColor: "#4A3E2D",
			ANSI256:   "238",
			ANSI:      "8",
		},
		WallLit: lipgloss.CompleteColor{
			TrueColor: "#998574",
			ANSI256:   "245",
			ANSI:      "7",
		},
		FloorDim: lipgloss.CompleteColor{
			TrueColor: "#1E2530",
			ANSI256:   "235",
			ANSI:      "0",
		},
		WallDim: lipgloss.CompleteColor{
			TrueColor: "#3B485E",
			ANSI256:   "240",
			ANSI:      "8",
		},
		Stairs: lipgloss.CompleteColor{
			TrueColor: "#FFD700",
			ANSI256:   "220",
			ANSI:      "3",
		},
		Potion: lipgloss.CompleteColor{
			TrueColor: "#FF55FF",
			ANSI256:   "207",
			ANSI:      "5",
		},
		Weapon: lipgloss.CompleteColor{
			TrueColor: "#00E5FF",
			ANSI256:   "45",
			ANSI:      "6",
		},
		HUDNormal: lipgloss.CompleteColor{
			TrueColor: "#00FF87",
			ANSI256:   "46",
			ANSI:      "2",
		},
		HUDWarning: lipgloss.CompleteColor{
			TrueColor: "#FF3B30",
			ANSI256:   "196",
			ANSI:      "1",
		},
		HUDLog: lipgloss.CompleteColor{
			TrueColor: "#8E8E93",
			ANSI256:   "245",
			ANSI:      "7",
		},
	}
}

// ResolveEntityColor maps an entity to its palette token or falls back to lipgloss.Color(e.Color).
func ResolveEntityColor(e Entity, p Palette) lipgloss.TerminalColor {
	if e.IsPlayer {
		return p.Player
	}
	switch e.Name {
	case "Goblin":
		return p.Goblin
	case "Orc":
		return p.Orc
	case "Troll":
		return p.Troll
	case "Archer":
		return p.Archer
	default:
		return lipgloss.Color(e.Color)
	}
}

// ResolveItemColor maps an item to its palette token.
func ResolveItemColor(it Item, p Palette) lipgloss.TerminalColor {
	switch it.ItemType {
	case ItemPotion:
		return p.Potion
	case ItemWeapon:
		return p.Weapon
	default:
		return p.HUDNormal
	}
}
