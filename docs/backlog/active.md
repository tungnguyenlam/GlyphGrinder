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

**Make the repository installable and document the complete play path.**

- Adopt the public GitHub module identity from `origin` so documented
  `go install github.com/tungnguyenlam/GlyphGrinder@latest` installs the right
  command when a version is available.
- Add an isolated local `go install .` check to `./scripts/verify.sh` without
  leaving a binary in the repository or the user's normal Go bin directory.
- Rewrite the README around installation, all controls, the three-depth win
  condition, terminal fallback behavior, and contributor verification.
- Cover any path/name fallout and run the canonical verification suite.

Why next: the game now has a full run, but its module identity is local-only,
installation is unverified, and the README still describes the M1-era loop.

## Milestone plan

- [ ] **M4.1** Verify local and versioned `go install` paths; document install,
      controls, gameplay, and terminal requirements.
- [ ] **M4.2** Exercise startup degradation for ASCII/no Nerd Font, truecolor
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
