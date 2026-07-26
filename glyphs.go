package main

import "os"

// GlyphSet holds the display symbols for all game map tiles and entities.
// It supports both Nerd Font unicode symbols and clean ASCII fallbacks.
type GlyphSet struct {
	Player     string
	Goblin     string
	Orc        string
	Troll      string
	Archer     string
	Floor      string
	Wall       string
	StairsDown string
	Potion     string
	Weapon     string
	Amulet     string
	Scroll     string
	DoorClosed string
	DoorOpen   string
}

// ASCIIGlyphs returns the standard ASCII fallback glyph set.
func ASCIIGlyphs() GlyphSet {
	return GlyphSet{
		Player:     "@",
		Goblin:     "g",
		Orc:        "o",
		Troll:      "T",
		Archer:     "A",
		Floor:      ".",
		Wall:       "#",
		StairsDown: ">",
		Potion:     "!",
		Weapon:     "/",
		Amulet:     "*",
		Scroll:     "?",
		DoorClosed: "+",
		DoorOpen:   "/",
	}
}

// NerdFontGlyphs returns the enriched Nerd Font glyph set.
func NerdFontGlyphs() GlyphSet {
	return GlyphSet{
		Player:     "󰋋", // U+F02CB (nf-md-account / knight)
		Goblin:     "󰆧", // U+F01A7 (nf-md-ghost)
		Orc:        "󰌆", // U+F0306 (nf-md-skull)
		Troll:      "󰇄", // U+F01C4 (troll / monster)
		Archer:     "󰓤", // U+F04E4 (archer / bow)
		Floor:      "·", // U+00B7 (middle dot)
		Wall:       "▓", // U+2593 (dark shade block)
		StairsDown: "󰌑", // U+F0311 (stairs down)
		Potion:     "󰏗", // U+F03D7 (potion)
		Weapon:     "󰓥", // U+F04E5 (sword)
		Amulet:     "󰇮", // U+F01EE (amulet / gem)
		Scroll:     "󰈙", // U+F0219 (scroll / document)
		DoorClosed: "󰌝", // U+F031D (closed door)
		DoorOpen:   "󰌟", // U+F031F (open door)
	}
}

// ResolveEntityGlyph maps an entity to its GlyphSet token or falls back to e.Rune.
func ResolveEntityGlyph(e Entity, g GlyphSet) string {
	if e.IsPlayer {
		return g.Player
	}
	switch e.Name {
	case "Goblin":
		return g.Goblin
	case "Orc":
		return g.Orc
	case "Troll":
		return g.Troll
	case "Archer":
		return g.Archer
	default:
		if e.Rune != "" {
			return e.Rune
		}
		return "?"
	}
}

// ResolveItemGlyph maps an item to its GlyphSet token.
func ResolveItemGlyph(it Item, g GlyphSet) string {
	switch it.ItemType {
	case ItemPotion:
		return g.Potion
	case ItemWeapon:
		return g.Weapon
	case ItemAmulet:
		return g.Amulet
	case ItemScroll:
		return g.Scroll
	default:
		return "?"
	}
}

// DetectGlyphSetFromEnv checks environment variables to determine whether to use
// Nerd Font glyphs or fallback ASCII glyphs.
func DetectGlyphSetFromEnv(getenv func(string) string) GlyphSet {
	if getenv == nil {
		getenv = os.Getenv
	}

	// Explicit override via environment variables
	val := getenv("GLYPHGRINDER_NERD_FONTS")
	if val == "" {
		val = getenv("GLYPHGRINDER_GLYPHS")
	}
	if val == "" {
		val = getenv("NERD_FONTS")
	}

	if val != "" {
		switch val {
		case "0", "false", "ascii", "off", "no":
			return ASCIIGlyphs()
		case "1", "true", "nerd", "nerdfont", "on", "yes":
			return NerdFontGlyphs()
		}
	}

	// Explicit ASCII flags
	if getenv("GLYPHGRINDER_ASCII") == "1" || getenv("GLYPHGRINDER_ASCII") == "true" || getenv("NO_UNICODE") == "1" {
		return ASCIIGlyphs()
	}

	// Dumb terminal check
	if getenv("TERM") == "dumb" {
		return ASCIIGlyphs()
	}

	// Default to Nerd Font glyphs on modern terminals
	return NerdFontGlyphs()
}

// DetectGlyphSet auto-detects the appropriate GlyphSet from OS environment.
func DetectGlyphSet() GlyphSet {
	return DetectGlyphSetFromEnv(os.Getenv)
}
