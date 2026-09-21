# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-09-21 (`--dump-frame` landed; next batch is headless CLI help)

## Current milestone

**M4 — Ships to other people.** Turn the complete three-depth roguelike into a
project people can install, understand, and trust on ordinary terminals.

Serves `GOAL.md` directly: M3 completed the game loop; M4 proves the single
binary installs cleanly, degrades gracefully, stays responsive, and remains
green outside the maintainer's machine.

## Exact next action

**Current batch — make the binary self-describing without a TTY:**

1. Parse `--help`/`-h` in `main.go` so those flags print usage to stdout and
   exit 0 without calling `tea.NewProgram`.
2. Reject unknown arguments with a stderr usage line and a non-zero status
   instead of falling through into the TTY-only interactive path.
3. Cover both paths in `main_flag_test.go`; keep `--dump-frame` working as the
   production-frame smoke test.
4. Mention `--help` and `--dump-frame` in the README Develop (or Install)
   section so a first-run TTY failure is not the only documented non-play
   invocation.

Why this batch next: `--dump-frame` made the real binary headless-verifiable,
but `GlyphGrinder --help` and typos still open alt-screen and die with
`could not open a new TTY`. That is a first-impression M4 gap on ordinary
machines without a controlling terminal (pipes, CI, `go install` smoke).
External release checks remain blocked (see Blockers) and cannot be the next
executable action.

After this batch, re-evaluate blocked external release checks:

- Confirm the `Verify` GitHub Actions workflow passes on the first pushed
  commit and that its only verification entry point is `./scripts/verify.sh`.
- From outside this checkout, run
  `go install github.com/tungnguyenlam/GlyphGrinder@latest` after the public
  branch contains these commits and confirm `GlyphGrinder` launches.
- Run `make run` in a real terminal, capture a representative gameplay frame,
  add the image under `docs/assets/`, and embed it near the top of the README.

## Milestone plan

- [x] **M4.1** Verify local and versioned `go install` paths; document install,
      controls, gameplay, and terminal requirements — local `GOBIN` install
      passes in `verify.sh`; versioned `@latest` awaits push (blocked).
- [x] **M4.2** Exercise startup degradation for ASCII/no Nerd Font, truecolor
      capability, tiny terminals, and mid-run resize without a TTY.
- [x] **M4.3** Benchmark the render path and enforce an evidence-based frame
      budget that keeps normal play smooth and quiet.
- [x] **M4.4** Run `./scripts/verify.sh` in CI on pushes and pull requests —
      workflow `.github/workflows/verify.yml:1-27` exists with push/PR, cache,
      per-ref cancel, 10-min timeout, sole canonical script; remote run awaits
      push (blocked).
- [x] **M4.5** `--dump-frame` renders one production frame without a TTY;
      `verify.sh` exercises the installed binary and `go run . --dump-frame`
      and keeps `TestViewRendersFullGrid` as the pure-View guard.

## Acceptance criteria for M4

- The README gives a working install command, complete controls, the run goal,
  terminal expectations, a genuine gameplay screenshot, and the canonical
  contributor workflow.
- Installation, startup fallback, tiny/mid-run resize behavior, and the render
  frame budget have automated coverage that does not require a TTY.
- CI runs the same canonical verification command used locally.
- `./scripts/verify.sh` passes.

## Blockers

- The workflow and versioned `go install ...@latest` cannot see unpushed local
  commits. The standing prompt forbids pushing or publishing.
- A genuine README screenshot needs `make run` in a real TTY; the agent
  environment cannot open `/dev/tty` (see the active TTY notice).

## Last verification

```
./scripts/verify.sh   —  2026-09-21  —  VERIFY OK (go1.27.0 darwin/arm64, /opt/homebrew/bin/go)
```

`--dump-frame` batch (2026-09-21): `main.dumpFrame` prints one `View()` frame
without `tea.NewProgram`; `hasDumpFrameArg` is unit-tested; `verify.sh` checks
the isolated `go install` binary and `go run . --dump-frame` for `HP ` / `Depth 1`
after stripping CSI; HUD grep is ANSI-safe; README Develop and the TTY notice
document the flag. Parking-lot `--dump-frame` item removed.

Honest-failure batch (2026-08-23): `verify.sh` no longer claims
`all files gofmt-clean` when `gofmt` is missing or errors — it checks `go` and
`gofmt` exist up front and trusts `gofmt`'s exit status. Proven by running the
script with an empty toolchain PATH (fails naming `go`), a `go`-only PATH
(fails naming `gofmt`), and a syntax-error probe file (fails with the parse
error; previously passed vacuously). Startup errors now go to stderr with a
trailing newline.

Session audit (2026-08-23): game logic, value-state slice isolation, the
69-test suite, CI workflow (exec bit on `verify.sh` is `100755`, canonical
script is the only entry point), Makefile, README, and agent docs re-checked —
no further provable gaps. Parking lot pruned of two answered items (View
allocation profiling, map-generation fuzzing); TTY notice example refreshed
for the stderr error format.

M3.5 three-depth victory flow: pre-final descent, final-stair escape without
regeneration, frozen won state, visible goal/victory UI, and one-key fresh-run
restart pass at state and headless levels.

M4.1 public GitHub module identity, isolated executable `go install` check,
complete installation/gameplay/control/terminal README, and path fallout all
pass the canonical suite.

M4.2 injectable auto-detected Lip Gloss renderer, truecolor/colorless output
coverage, existing locale/`TERM` glyph fallback, 1x1 startup, and resize during
active animation all pass headlessly.

M4.3 baseline on Apple M4 (`go test -run '^$' -bench
'^BenchmarkViewProductionTrueColor$' -benchmem -count=5 .`): full 96x48 frames
97.9–101.0 us/op, 86,215 B/op, 893 allocs/op; clipped 80x24 frames 75.4–79.9
us/op, 47,679 B/op, 889 allocs/op. This is about 0.6% of the 16.7 ms animation
interval, so no optimization was justified; the averaged production-frame
budget guard passes 20 consecutive runs.

M4.4 read-only GitHub Actions workflow: push/PR triggers, cached `go.mod` Go
toolchain, per-ref cancellation, ten-minute timeout, and the canonical script
as its sole verification command. YAML parses locally and the same script is
green; remote execution awaits a push.

Audit 2026-08-26: re-ran `./scripts/verify.sh` with `/opt/homebrew/bin/go`
(go1.27) — VERIFY OK. Found stale `AGENTS.md` Current State and the missing
`--dump-frame` path; both are now landed.
