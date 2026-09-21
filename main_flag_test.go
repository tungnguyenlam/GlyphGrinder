package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    cliCommand
		wantErr string
	}{
		{name: "no args", args: nil, want: cliPlay},
		{name: "empty slice", args: []string{}, want: cliPlay},
		{name: "dump-frame alone", args: []string{"--dump-frame"}, want: cliDumpFrame},
		{name: "help long", args: []string{"--help"}, want: cliHelp},
		{name: "help short", args: []string{"-h"}, want: cliHelp},
		{name: "help wins over dump-frame", args: []string{"--dump-frame", "--help"}, want: cliHelp},
		{name: "help before dump-frame still help", args: []string{"--help", "--dump-frame"}, want: cliHelp},
		{name: "short help wins over dump-frame", args: []string{"-h", "--dump-frame"}, want: cliHelp},
		{name: "unknown alone", args: []string{"--foo"}, wantErr: "unknown argument: --foo"},
		{name: "unknown after dump-frame", args: []string{"--dump-frame", "extra"}, wantErr: "unknown argument: extra"},
		{name: "unknown before dump-frame", args: []string{"other", "--dump-frame"}, wantErr: "unknown argument: other"},
		{name: "single dash dump-frame is unknown", args: []string{"-dump-frame"}, wantErr: "unknown argument: -dump-frame"},
		{name: "underscore dump-frame is unknown", args: []string{"--dump_frame"}, wantErr: "unknown argument: --dump_frame"},
		{name: "missing dashes dump-frame is unknown", args: []string{"dump-frame"}, wantErr: "unknown argument: dump-frame"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseArgs(tc.args)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("parseArgs(%q) = %v, nil; want error %q", tc.args, got, tc.wantErr)
				}
				if got := err.Error(); got != tc.wantErr {
					t.Errorf("error = %q, want %q", got, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseArgs(%q) error = %v, want %v", tc.args, err, tc.want)
			}
			if got != tc.want {
				t.Errorf("parseArgs(%q) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestWriteUsageMentionsHeadlessFlags(t *testing.T) {
	var buf bytes.Buffer
	writeUsage(&buf)
	got := buf.String()
	for _, want := range []string{"--dump-frame", "--help", "-h", "TTY"} {
		if !strings.Contains(got, want) {
			t.Errorf("usage missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "tea.NewProgram") {
		t.Errorf("usage leaked implementation detail")
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

func TestFormatRunErrorHintsOnMissingTTY(t *testing.T) {
	ttyErr := fmt.Errorf("could not open a new TTY: open /dev/tty: device not configured")
	got := formatRunError(ttyErr)
	if !strings.HasPrefix(got, "glyphgrinder: could not open a new TTY:") {
		t.Errorf("TTY error lost the original message:\n%s", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("TTY error missing trailing newline: %q", got)
	}
	if !strings.Contains(got, "--help") || !strings.Contains(got, "--dump-frame") {
		t.Errorf("TTY error missing headless-flag hint:\n%s", got)
	}

	other := formatRunError(fmt.Errorf("boom"))
	if got, want := other, "glyphgrinder: boom\n"; got != want {
		t.Errorf("non-TTY error = %q, want %q", got, want)
	}
	if strings.Contains(other, "--help") {
		t.Errorf("non-TTY error should not hint at flags: %q", other)
	}
}
