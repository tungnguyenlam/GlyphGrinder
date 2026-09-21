package main

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

func TestHasDumpFrameArg(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "no args", args: nil, want: false},
		{name: "empty slice", args: []string{}, want: false},
		{name: "flag alone", args: []string{"--dump-frame"}, want: true},
		{name: "flag with others", args: []string{"--dump-frame", "extra"}, want: true},
		{name: "flag after other", args: []string{"other", "--dump-frame"}, want: true},
		{name: "single dash not accepted", args: []string{"-dump-frame"}, want: false},
		{name: "underscore not accepted", args: []string{"--dump_frame"}, want: false},
		{name: "missing dash", args: []string{"dump-frame"}, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasDumpFrameArg(tc.args); got != tc.want {
				t.Errorf("hasDumpFrameArg(%q) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestDumpFrameRendersHeadlesslyWithoutTTY(t *testing.T) {
	// Exercise the same dumpFrame helper the binary uses for --dump-frame.
	// Glyphs follow the process environment (ASCII or Nerd Font); HUD text
	// is ASCII either way. The subprocess half is exercised by
	// scripts/verify.sh via both the installed binary and `go run . --dump-frame`.
	var buf bytes.Buffer
	dumpFrame(&buf, rand.New(rand.NewSource(9)))
	plain := stripANSI(buf.String())
	if !strings.Contains(plain, "HP ") {
		t.Errorf("dumped frame missing HP HUD")
	}
	if !strings.Contains(plain, "Depth 1") {
		t.Errorf("dumped frame missing depth HUD")
	}
	lines := strings.Split(strings.TrimSuffix(plain, "\n"), "\n")
	if got, want := len(lines), productionDungeonHeight; got != want {
		t.Errorf("dumped frame has %d rows, want %d", got, want)
	}
}
