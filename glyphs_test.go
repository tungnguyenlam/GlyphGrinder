package main

import (
	"strings"
	"testing"

	"glyphgrinder/internal/tuitest"
)

func TestGlyphSetTokens(t *testing.T) {
	ascii := ASCIIGlyphs()
	if ascii.Player != "@" || ascii.Goblin != "g" || ascii.Orc != "o" || ascii.Troll != "T" || ascii.Archer != "A" || ascii.Floor != "." || ascii.Wall != "#" || ascii.Potion != "!" || ascii.Weapon != "/" || ascii.Amulet != "*" {
		t.Errorf("unexpected ASCIIGlyphs tokens: %+v", ascii)
	}

	nerd := NerdFontGlyphs()
	if nerd.Player != "󰋋" || nerd.Goblin != "󰆧" || nerd.Orc != "󰌆" || nerd.Troll != "󰇄" || nerd.Archer != "󰓤" || nerd.Floor != "·" || nerd.Wall != "▓" || nerd.Potion != "󰏗" || nerd.Weapon != "󰓥" || nerd.Amulet != "󰇮" {
		t.Errorf("unexpected NerdFontGlyphs tokens: %+v", nerd)
	}
}

func TestResolveEntityGlyph(t *testing.T) {
	ascii := ASCIIGlyphs()
	nerd := NerdFontGlyphs()

	player := Entity{IsPlayer: true, Name: "Player", Rune: "@"}
	goblin := Entity{IsPlayer: false, Name: "Goblin", Rune: "g"}
	orc := Entity{IsPlayer: false, Name: "Orc", Rune: "o"}
	troll := Entity{IsPlayer: false, Name: "Troll", Rune: "T"}
	archer := Entity{IsPlayer: false, Name: "Archer", Rune: "A"}
	custom := Entity{IsPlayer: false, Name: "Dragon", Rune: "D"}

	// ASCII mode
	if got := ResolveEntityGlyph(player, ascii); got != "@" {
		t.Errorf("ResolveEntityGlyph(player, ascii) = %q, want '@'", got)
	}
	if got := ResolveEntityGlyph(goblin, ascii); got != "g" {
		t.Errorf("ResolveEntityGlyph(goblin, ascii) = %q, want 'g'", got)
	}
	if got := ResolveEntityGlyph(orc, ascii); got != "o" {
		t.Errorf("ResolveEntityGlyph(orc, ascii) = %q, want 'o'", got)
	}
	if got := ResolveEntityGlyph(troll, ascii); got != "T" {
		t.Errorf("ResolveEntityGlyph(troll, ascii) = %q, want 'T'", got)
	}
	if got := ResolveEntityGlyph(archer, ascii); got != "A" {
		t.Errorf("ResolveEntityGlyph(archer, ascii) = %q, want 'A'", got)
	}
	if got := ResolveEntityGlyph(custom, ascii); got != "D" {
		t.Errorf("ResolveEntityGlyph(custom, ascii) = %q, want 'D'", got)
	}

	// Nerd Font mode
	if got := ResolveEntityGlyph(player, nerd); got != "󰋋" {
		t.Errorf("ResolveEntityGlyph(player, nerd) = %q, want '󰋋'", got)
	}
	if got := ResolveEntityGlyph(goblin, nerd); got != "󰆧" {
		t.Errorf("ResolveEntityGlyph(goblin, nerd) = %q, want '󰆧'", got)
	}
	if got := ResolveEntityGlyph(orc, nerd); got != "󰌆" {
		t.Errorf("ResolveEntityGlyph(orc, nerd) = %q, want '󰌆'", got)
	}
	if got := ResolveEntityGlyph(troll, nerd); got != "󰇄" {
		t.Errorf("ResolveEntityGlyph(troll, nerd) = %q, want '󰇄'", got)
	}
	if got := ResolveEntityGlyph(archer, nerd); got != "󰓤" {
		t.Errorf("ResolveEntityGlyph(archer, nerd) = %q, want '󰓤'", got)
	}
	if got := ResolveEntityGlyph(custom, nerd); got != "D" {
		t.Errorf("ResolveEntityGlyph(custom, nerd) = %q, want 'D'", got)
	}
}

func TestResolveItemGlyph(t *testing.T) {
	ascii := ASCIIGlyphs()
	nerd := NerdFontGlyphs()

	pot := Item{ItemType: ItemPotion}
	weap := Item{ItemType: ItemWeapon}
	amu := Item{ItemType: ItemAmulet}

	if got := ResolveItemGlyph(pot, ascii); got != "!" {
		t.Errorf("potion ascii = %q, want '!'", got)
	}
	if got := ResolveItemGlyph(weap, ascii); got != "/" {
		t.Errorf("weapon ascii = %q, want '/'", got)
	}
	if got := ResolveItemGlyph(amu, ascii); got != "*" {
		t.Errorf("amulet ascii = %q, want '*'", got)
	}

	if got := ResolveItemGlyph(pot, nerd); got != "󰏗" {
		t.Errorf("potion nerd = %q, want '󰏗'", got)
	}
	if got := ResolveItemGlyph(weap, nerd); got != "󰓥" {
		t.Errorf("weapon nerd = %q, want '󰓥'", got)
	}
	if got := ResolveItemGlyph(amu, nerd); got != "󰇮" {
		t.Errorf("amulet nerd = %q, want '󰇮'", got)
	}
}

func TestDetectGlyphSetFromEnv(t *testing.T) {
	mockEnv := func(env map[string]string) func(string) string {
		return func(key string) string {
			return env[key]
		}
	}

	cases := []struct {
		name string
		env  map[string]string
		want GlyphSet
	}{
		{
			name: "explicit GLYPHGRINDER_NERD_FONTS=0",
			env:  map[string]string{"GLYPHGRINDER_NERD_FONTS": "0"},
			want: ASCIIGlyphs(),
		},
		{
			name: "explicit GLYPHGRINDER_NERD_FONTS=1",
			env:  map[string]string{"GLYPHGRINDER_NERD_FONTS": "1"},
			want: NerdFontGlyphs(),
		},
		{
			name: "explicit GLYPHGRINDER_GLYPHS=ascii",
			env:  map[string]string{"GLYPHGRINDER_GLYPHS": "ascii"},
			want: ASCIIGlyphs(),
		},
		{
			name: "explicit NERD_FONTS=true",
			env:  map[string]string{"NERD_FONTS": "true"},
			want: NerdFontGlyphs(),
		},
		{
			name: "GLYPHGRINDER_ASCII=1",
			env:  map[string]string{"GLYPHGRINDER_ASCII": "1"},
			want: ASCIIGlyphs(),
		},
		{
			name: "NO_UNICODE=1",
			env:  map[string]string{"NO_UNICODE": "1"},
			want: ASCIIGlyphs(),
		},
		{
			name: "TERM=dumb",
			env:  map[string]string{"TERM": "dumb"},
			want: ASCIIGlyphs(),
		},
		{
			name: "default env",
			env:  map[string]string{},
			want: NerdFontGlyphs(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DetectGlyphSetFromEnv(mockEnv(tc.env))
			if got != tc.want {
				t.Errorf("DetectGlyphSetFromEnv() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestTUIModelRenderingModes(t *testing.T) {
	modes := []struct {
		name       string
		glyphs     GlyphSet
		wantPlayer string
		wantWall   string
	}{
		{"ASCII Mode", ASCIIGlyphs(), "@", "#"},
		{"Nerd Font Mode", NerdFontGlyphs(), "󰋋", "▓"},
	}

	for _, tc := range modes {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModelWithSeed(12345)
			m.glyphs = tc.glyphs
			d := tuitest.New(t, m)

			lines := d.Lines()
			fullText := stripANSI(strings.Join(lines, "\n"))

			if !strings.Contains(fullText, tc.wantPlayer) {
				t.Errorf("rendered view in %s missing player glyph %q", tc.name, tc.wantPlayer)
			}
			if !strings.Contains(fullText, tc.wantWall) {
				t.Errorf("rendered view in %s missing wall glyph %q", tc.name, tc.wantWall)
			}
		})
	}
}
