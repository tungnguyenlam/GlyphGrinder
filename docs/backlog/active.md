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

**Measure and hold the render path to an animation-frame budget.**

- Add repeatable benchmarks for full 96x48 truecolor rendering and a typical
  clipped viewport, reporting allocations as well as time.
- Establish the baseline against the 16.7 ms animation-frame interval and
  identify the dominant measured cost before changing code.
- Optimize only if the measurements justify it, preserving byte-for-byte plain
  rendering semantics and state purity.
- Add a conservative automated frame-budget guard to canonical verification
  and record the before/after evidence in this backlog.

Why next: terminal fallback is now explicit and tested. GOAL.md also promises
smooth, quiet rendering, so the hot path needs evidence and a regression bound
before CI makes the suite authoritative.

## Milestone plan

- [x] **M4.1** Verify local and versioned `go install` paths; document install,
      controls, gameplay, and terminal requirements.
- [x] **M4.2** Exercise startup degradation for ASCII/no Nerd Font, truecolor
      capability, tiny terminals, and mid-run resize without a TTY.
- [ ] **M4.3** Benchmark the render path and enforce an evidence-based frame
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
