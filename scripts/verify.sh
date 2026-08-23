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
	rm -f "$install_dir/GlyphGrinder"
	rmdir "$install_dir" 2>/dev/null || true
}
trap cleanup_install EXIT
GOBIN="$install_dir" go install .
if [[ ! -x "$install_dir/GlyphGrinder" ]]; then
	echo "go install did not produce the GlyphGrinder command" >&2
	exit 1
fi
echo "go install produced an executable GlyphGrinder command"
cleanup_install
trap - EXIT

step "smoke test (renders one frame headlessly)"
# The real binary needs a TTY and cannot run in CI or under an agent; the
# headless driver in internal/tuitest is the smoke test. See
# docs/agent/notices/2026-07-26-tui-needs-a-tty.md
go test -count=1 -run 'TestViewRendersFullGrid' .

step "full test suite"
go test -count=1 ./...

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
