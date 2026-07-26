# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M4 — Ships to other people.** Production-ready packaging, documentation, benchmarks, CI, and graceful degradation across terminal environments.

Serves `GOAL.md` directly: "And it just runs. `go install` or a single binary, no config file, no font wrangling beyond 'install a Nerd Font', graceful degradation to plainer glyphs and colors when the terminal can't do better."

## Exact next action

**M4.2 — Terminal capability & environment degradation checks.**

- In `colors.go`, `glyphs.go`, `colors_test.go`, `glyphs_test.go`:
  - Test environment variable overrides (`NO_COLOR`, `GLYPHGRINDER_ASCII`, `TERM=dumb`, `GLYPHGRINDER_NERD_FONTS`).
  - Add tests in `colors_test.go` and `glyphs_test.go` ensuring complete coverage for ASCII glyph fallbacks and palette rendering under basic terminals.
  - Add TUI tests in `tui_test.go` testing extreme viewport sizes (e.g. 10x5, 5x5, 0x0) ensuring no panics and clean clipping.

Why next: M4.2 ensures GlyphGrinder degrades gracefully on any terminal environment without crashing.

## Milestone plan

- [x] **M4.1** `go install` verification & README documentation rewrite.
- [ ] **M4.2** Terminal capability & environment degradation checks (NO_COLOR, GLYPHGRINDER_ASCII, tiny window sizes).
- [ ] **M4.3** Performance benchmarks (`BenchmarkView`, `BenchmarkStep`) holding frame budgets under 1ms.
- [ ] **M4.4** GitHub Actions CI workflow running `./scripts/verify.sh`.

## Acceptance criteria for M4

- Binary builds and installs via `go install`.
- `README.md` provides clear controls, installation guide, and design summary.
- Terminal degradation modes (ASCII fallback, monochrome, small windows) tested and verified.
- Benchmarks verify sub-millisecond view render time.
- CI configuration in `.github/workflows/ci.yml` passes `./scripts/verify.sh`.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 62 test functions pass, `go.mod` tidy.


