#!/usr/bin/env bash
# The canonical verification command for GlyphGrinder.
#
#   ./scripts/verify.sh     (or: make verify)
#
# Every doc in this repo points here. If you need a new check, add it to this
# script rather than inventing an ad hoc command in a doc or a PR description.
set -euo pipefail

cd "$(dirname "$0")/.."

step() { printf '\n\033[1;36m==> %s\033[0m\n' "$1"; }

# HUD strings survive Lip Gloss styling because CSI SGR wraps the label, not
# the characters inside it — except "Depth 1", where the number sits outside
# the styled "Depth". Strip escapes so a color-capable renderer cannot split
# that check.
frame_has_hud() {
	local plain
	plain="$(sed $'s/\x1b\\[[0-9;]*m//g' "$1")"
	grep -q "HP " <<<"$plain" && grep -q "Depth 1" <<<"$plain"
}

# A verification that cannot run its tools must fail loudly, never pass
# vacuously: check the toolchain exists before any step relies on it.
for tool in go gofmt; do
	if ! command -v "$tool" >/dev/null 2>&1; then
		echo "verify.sh: required tool '$tool' not found on PATH; install Go and retry" >&2
		exit 1
	fi
done

step "gofmt (formatting check)"
# gofmt exits non-zero on syntax errors without listing the offending file, so
# trust its exit status, not just its output. No `|| true` on this call.
if ! unformatted="$(gofmt -l .)"; then
	echo "gofmt failed; fix the syntax errors it reported and re-run" >&2
	exit 1
fi
unformatted="$(printf '%s' "$unformatted" | grep -v '^vendor/' || true)"
if [[ -n "$unformatted" ]]; then
	echo "These files are not gofmt-clean:" >&2
	echo "$unformatted" >&2
	echo "Run: gofmt -w ." >&2
	exit 1
fi
echo "all files gofmt-clean"

step "go vet (static analysis)"
go vet ./...

step "go build (compiles)"
go build ./...

step "go install (isolated binary)"
install_dir="$(mktemp -d "${TMPDIR:-/tmp}/glyphgrinder-install.XXXXXX")"
cleanup_install() {
	rm -f "$install_dir/GlyphGrinder" "$install_dir/frame.txt" "$install_dir/help.txt" "$install_dir/ascii.txt" "$install_dir/rich.txt"
	rmdir "$install_dir" 2>/dev/null || true
}
trap cleanup_install EXIT
GOBIN="$install_dir" go install .
if [[ ! -x "$install_dir/GlyphGrinder" ]]; then
	echo "go install did not produce the GlyphGrinder command" >&2
	exit 1
fi
echo "go install produced an executable GlyphGrinder command"
if ! "$install_dir/GlyphGrinder" --dump-frame >"$install_dir/frame.txt"; then
	echo "GlyphGrinder --dump-frame failed" >&2
	exit 1
fi
if ! frame_has_hud "$install_dir/frame.txt"; then
	echo "GlyphGrinder --dump-frame output missing expected HUD" >&2
	cat "$install_dir/frame.txt" >&2 || true
	exit 1
fi
echo "GlyphGrinder --dump-frame produced a valid frame"
if ! "$install_dir/GlyphGrinder" --help >"$install_dir/help.txt"; then
	echo "GlyphGrinder --help failed" >&2
	exit 1
fi
if ! grep -q -- "--dump-frame" "$install_dir/help.txt" || ! grep -q -- "--help" "$install_dir/help.txt"; then
	echo "GlyphGrinder --help missing expected flags" >&2
	cat "$install_dir/help.txt" >&2 || true
	exit 1
fi
echo "GlyphGrinder --help produced usage"

step "degradation matrix (installed binary)"
# The real binary must pick its glyph profile from the environment without a
# TTY: TERM=dumb and a POSIX locale fall back to one-cell ASCII, while a UTF-8
# locale selects the Nerd Font glyphs. This is the startup-degradation check
# M4 promises; glyphs never depend on the detected color profile.
if ! TERM=dumb LC_ALL=C LANG=C "$install_dir/GlyphGrinder" --dump-frame >"$install_dir/ascii.txt"; then
	echo "TERM=dumb GlyphGrinder --dump-frame failed" >&2
	exit 1
fi
if LC_ALL=C grep -q '[^ -~]' "$install_dir/ascii.txt"; then
	echo "TERM=dumb frame is not pure ASCII:" >&2
	LC_ALL=C grep -n '[^ -~]' "$install_dir/ascii.txt" | head -3 >&2 || true
	exit 1
fi
if ! TERM=xterm-256color LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 "$install_dir/GlyphGrinder" --dump-frame >"$install_dir/rich.txt"; then
	echo "UTF-8 locale GlyphGrinder --dump-frame failed" >&2
	exit 1
fi
if ! LC_ALL=C grep -q '[^ -~]' "$install_dir/rich.txt"; then
	echo "UTF-8 locale frame has no rich glyphs" >&2
	exit 1
fi
echo "TERM=dumb/POSIX renders ASCII; UTF-8 renders rich glyphs"
cleanup_install
trap - EXIT

step "smoke test (--dump-frame and headless View)"
# The real binary renders headlessly via --dump-frame so CI and agents can
# smoke-test it without a TTY; View is also exercised directly via tuitest.
# See docs/agent/notices/2026-07-26-tui-needs-a-tty.md
dump_file="$(mktemp "${TMPDIR:-/tmp}/glyphgrinder-dump.XXXXXX")"
if ! go run . --dump-frame >"$dump_file"; then
	echo "go run --dump-frame failed" >&2
	rm -f "$dump_file"
	exit 1
fi
if ! frame_has_hud "$dump_file"; then
	echo "go run --dump-frame output missing expected HUD" >&2
	cat "$dump_file" >&2 || true
	rm -f "$dump_file"
	exit 1
fi
echo "go run --dump-frame produced a valid frame"
rm -f "$dump_file"
go test -count=1 -run 'TestViewRendersFullGrid' .

step "full test suite"
# The race detector guards the value-state slice isolation contract (ADR-0002)
# and the synchronous tuitest driver (internal/tuitest/AGENTS.md): tests drive
# the model from a single goroutine and resolve animation ticks without real
# delays, so -race stays deterministic and fast.
go test -count=1 -race ./...

step "go.mod tidiness"
cp go.mod go.mod.verifybak
cp go.sum go.sum.verifybak
trap 'mv go.mod.verifybak go.mod; mv go.sum.verifybak go.sum' EXIT
go mod tidy
if ! diff -q go.mod go.mod.verifybak >/dev/null; then
	echo "go.mod is not tidy. Run: go mod tidy" >&2
	diff go.mod.verifybak go.mod || true
	exit 1
fi
echo "go.mod is tidy"

printf '\n\033[1;32mVERIFY OK\033[0m\n'
