# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M4 — Ships to other people.** Production-ready packaging, documentation, benchmarks, CI, and graceful degradation across terminal environments.

Serves `GOAL.md` directly: "And it just runs. `go install` or a single binary, no config file, no font wrangling beyond 'install a Nerd Font', graceful degradation to plainer glyphs and colors when the terminal can't do better."

## Exact next action

**M4.3 — Performance benchmarks (`BenchmarkView`, `BenchmarkStep`) holding frame budgets under 1ms.**

- Create `bench_test.go`:
  - `BenchmarkView`: measure `View()` execution time for 80x24 viewport with FOV, entities, and items.
  - `BenchmarkStep`: measure `Step()` execution time for player action + monster turns + FOV shadowcasting.
  - `BenchmarkComputeFOV`: measure `ComputeFOV()` shadowcasting algorithm speed.
- Verify sub-millisecond per-op execution (< 1ms per frame) to satisfy `GOAL.md` smooth battery-friendly performance goal.

Why next: M4.3 provides empirical performance benchmarks protecting render and step frame rates.

## Milestone plan

- [x] **M4.1** `go install` verification & README documentation rewrite.
- [x] **M4.2** Terminal capability & environment degradation checks (NO_COLOR, GLYPHGRINDER_ASCII, tiny window sizes).
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 65 test functions pass (including 3 new degradation, token resolution, and extreme viewport tests), `go.mod` tidy.


