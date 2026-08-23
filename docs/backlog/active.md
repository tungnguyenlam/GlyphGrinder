# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-08-23

## Current milestone

**M4 — Ships to other people.** Turn the complete three-depth roguelike into a
project people can install, understand, and trust on ordinary terminals.

Serves `GOAL.md` directly: M3 completed the game loop; M4 proves the single
binary installs cleanly, degrades gracefully, stays responsive, and remains
green outside the maintainer's machine.

## Exact next action

**Run the canonical verification suite in GitHub Actions.**

- Add a minimal workflow for pushes and pull requests that checks out the
  repository, installs the Go version declared by `go.mod`, caches modules, and
  runs only `./scripts/verify.sh`.
- Give the workflow read-only repository permissions and cancel superseded
  runs on the same ref.
- Validate the workflow structure locally as far as the environment permits,
  then run the canonical suite once more.

Why next: installability, fallback, and performance are now verified locally.
CI is the remaining M4 acceptance criterion and should invoke the exact same
entry point rather than grow a second verification definition.

## Milestone plan

- [x] **M4.1** Verify local and versioned `go install` paths; document install,
      controls, gameplay, and terminal requirements.
- [x] **M4.2** Exercise startup degradation for ASCII/no Nerd Font, truecolor
      capability, tiny terminals, and mid-run resize without a TTY.
- [x] **M4.3** Benchmark the render path and enforce an evidence-based frame
      budget that keeps normal play smooth and quiet.
- [ ] **M4.4** Run `./scripts/verify.sh` in CI on pushes and pull requests.

## Acceptance criteria for M4

- The README gives a working install command, complete controls, the run goal,
  terminal expectations, and the canonical contributor workflow.
- Installation, startup fallback, tiny/mid-run resize behavior, and the render
  frame budget have automated coverage that does not require a TTY.
- CI runs the same canonical verification command used locally.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-08-23  —  VERIFY OK
```

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
