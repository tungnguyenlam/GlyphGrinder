# Done

Completed milestones, newest first. This is the project's memory of what has
already been built and tried — check here before "discovering" that something
is missing.

## 2026-07-27 — M2: It looks like the pitch

Transformed GlyphGrinder into a visually distinct, responsive terminal roguelike.

- M2.1: Field of view (FOV) shadowcasting & map memory (lit tiles bright, explored dimmed, unexplored hidden).
- M2.2: Viewport camera centered on player and terminal window resize handling (`tea.WindowSizeMsg`).
- M2.3: Rich color palette system with Lip Gloss profile-aware rendering (TrueColor, ANSI256, ANSI fallbacks).
- M2.4: Nerd Font glyph set detection with ASCII fallback.
- M2.5: Tick-driven movement animation with continuous camera easing interpolation (~60 FPS ticks in main, instant 0ms ticks in unit tests).

Verified: `./scripts/verify.sh` — VERIFY OK.

## 2026-07-27 — M1: Playable core loop

Transformed the initial static 20x10 square room into a playable roguelike loop.

- M1.1: Factored movement into `tryMove(dx, dy)`.
- M1.2: Decoupled input handling in `Update` from turn processing via `Step(Action)`.
- M1.3: Dungeon map generation with rooms and corridors carved into `GameMap`, seeded deterministically for testability.
- M1.4: Entity system populated with monsters (Goblins and Orcs) rendered on the map.
- M1.5: Bump-to-attack combat dealing damage, removing dead entities, and appending combat events to `GameState.Log`.
- M1.6: Turn-based monster AI stepping towards the player and attacking when adjacent.
- M1.7: Rendered HUD status bar (`HP: X/Y`), message log display below the map, GAME OVER banner on player death, and single-keypress restart (`r`).

Verified: `./scripts/verify.sh` — VERIFY OK.

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
