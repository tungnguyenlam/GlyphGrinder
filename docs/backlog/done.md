# Done

Completed milestones, newest first. This is the project's memory of what has
already been built and tried — check here before "discovering" that something
is missing.

## 2026-07-26 — M0: Agent-centric scaffolding

Set the repo up so agents can pick up work across sessions without
re-discovering context.

- Added `AGENTS.md` (state, startup routine, gotchas, handoff rules),
  `GOAL.md`, `docs/agent/{index,continuity}.md`, four notices, four ADRs, this
  backlog, `prompt/improve.md`, and a thin `CLAUDE.md` pointer.
- Added `internal/tuitest` — a headless Bubble Tea driver (no new dependencies)
  so `Init`/`Update`/`View` can be tested without a TTY.
- Added the first tests: `game_test.go` (map borders, player placement) and
  `tui_test.go` (grid render, all movement keys, wall collision, quit keys).
  Repo previously had none.
- Added the canonical verification command `./scripts/verify.sh` (`make
  verify`): gofmt check, `go vet`, build, headless smoke frame, full test
  suite, `go mod tidy` check. Plus a `Makefile` with `run`/`build`/`test`/`fmt`.
- Fixed `go.mod`: `bubbletea` and `lipgloss` were marked `// indirect` despite
  being imported directly. `verify.sh` now fails if tidiness regresses.

Verified: `./scripts/verify.sh` — VERIFY OK.

## 2026-02-27 — Initial prototype (pre-dating this log)

Bubble Tea program rendering a fixed 20x10 walled room with an `@` movable by
arrows or WASD, quit on `q`/`ctrl+c`. Flat `GameState`/`GameMap`/`Entity`
model in `game.go`. No tests, no build tooling.
