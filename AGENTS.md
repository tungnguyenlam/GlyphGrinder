# AGENTS.md

Entry point for AI agents working in this repository. Read this first, then
follow the Startup Routine. Keep this file concise — durable detail belongs in
`docs/agent/notices/` and `docs/decisions/`, both indexed from
`docs/agent/index.md`.

## Current State

GlyphGrinder is a terminal roguelike written in Go, built on Bubble Tea
(`github.com/charmbracelet/bubbletea`) and Lip Gloss for styling. It is very
early: a single `package main` at the repo root holds the whole game — `game.go`
has the flat `GameState`/`GameMap`/`Entity` model, `NewGame`, and turn resolution
in `GameState.Step`; `main.go` maps Bubble Tea input to game actions and renders
the result. Right now the game boots into a fixed 20x10 walled room where you
move an `@` with arrows or WASD and quit with `q`; there are no monsters, no
combat, no level generation, and no animation yet.

The north star is [GOAL.md](GOAL.md) — read it before planning work, and check
every task against it.

<!-- TEMPORARY-PRIORITY START (added 2026-07-26; delete this whole block when the
     "Playable core loop" milestone in docs/backlog/active.md is done)
     Right now everything is about getting from "a square you can walk around"
     to a playable core loop: real map generation, monsters, and turn-based
     combat. Prefer work that moves that forward. Visual polish (animation,
     lighting, Nerd Font glyph sets) is deliberately deferred until there is a
     game to look at — park polish ideas in docs/backlog/parking-lot.md instead
     of building them now.
TEMPORARY-PRIORITY END -->

## Startup Routine

1. Read [`docs/backlog/active.md`](docs/backlog/active.md) — the current
   milestone and the exact next action.
2. Read [`docs/agent/index.md`](docs/agent/index.md) — the index of every
   notice, ADR and package doc.
3. Read [`docs/agent/continuity.md`](docs/agent/continuity.md) — the
   cross-session workflow you are expected to follow.
4. Scope with search before reading broadly: `rg 'TileWall' -n`,
   `rg 'func \(m model\)' -n`. Do not read the whole repo; read the area you
   are changing.
5. Read the nearest subtree `AGENTS.md` (e.g.
   [`internal/tuitest/AGENTS.md`](internal/tuitest/AGENTS.md)) before editing
   files under it.
6. Update `docs/backlog/active.md` **continuously** — as soon as a sub-task
   lands or the plan changes, not at the end of the session.

## Instruction Maintenance

You may edit this file and any subtree `AGENTS.md` when doing so makes future
sessions safer or faster. Rules:

- Keep this file short. Push detail into `docs/agent/notices/` (small, dated,
  expiring) or `docs/decisions/` (architecture-shaping ADRs).
- A new notice or ADR is not done until it is linked from
  `docs/agent/index.md`.
- Delete stale instructions in the same change that makes them stale. A gotcha
  you just fixed should be removed from Current Gotchas by that same commit.
- Do not duplicate the workflow into subtree files; subtree files are
  boundaries and invariants only.

## Common Tasks

| If you are... | Read |
| --- | --- |
| Changing game rules, map or entity state | `game.go`, [ADR-0002](docs/decisions/ADR-0002-flat-game-state.md) |
| Changing input handling or rendering | `main.go`, [ADR-0003](docs/decisions/ADR-0003-view-is-pure.md) |
| Writing or fixing a test that drives the TUI | [`internal/tuitest/AGENTS.md`](internal/tuitest/AGENTS.md) |
| Trying to run the game to check something | [notice: TUI needs a TTY](docs/agent/notices/2026-07-26-tui-needs-a-tty.md) |
| Deciding where new code should live | [ADR-0001](docs/decisions/ADR-0001-single-package-until-it-hurts.md) |
| Picking up work with no other instruction | [`prompt/improve.md`](prompt/improve.md) |
| Verifying anything | `./scripts/verify.sh` (`make verify`) |

## Current Gotchas

Sharp edges found during real work. Prune entries that stop being true.

- **You cannot run the game to check your work.** `tea.WithAltScreen()` in
  `main()` requires a real TTY; piping input gives
  `could not open a new TTY: open /dev/tty: device not configured` and exit 1.
  Use `internal/tuitest` instead. Ask the user to run `make run` if a change
  genuinely needs human eyes.
- **Bounds checks in `GameState.Step` are load-bearing but redundant.** `NewGame`
  always walls the border, so the `Pos.Y > 0`-style checks never fire today.
  They will matter the moment map generation stops guaranteeing a wall ring —
  do not delete them without also guaranteeing the invariant.
- **`GameState.Log` is written by combat but not rendered yet.** Inspect it
  through state/headless tests until M1.7 adds the visible log; do not mistake
  the absent UI for missing combat events.
- **`View` returns styled cells**, so a rendered line's `len()` is not its
  column count. Tests must strip ANSI escapes before asserting positions — see
  `stripANSI` in `tui_test.go`.
- **Everything is `package main` at the repo root**, including tests. Only
  `internal/tuitest` is a separate package. See ADR-0001 before adding a new
  one.

## Handoff Rules

- **Unfinished work** → `docs/backlog/active.md` (current milestone, exact next
  action, blockers). Uncommitted ideas → `docs/backlog/parking-lot.md`.
  Committed future work → `docs/backlog/roadmap.md`. Finished milestones →
  `docs/backlog/done.md`.
- **Durable warnings** → a dated notice in `docs/agent/notices/`, linked from
  `docs/agent/index.md`. Short-lived ones can live in Current Gotchas above.
- **Architecture decisions** → an ADR in `docs/decisions/`.
- **The one canonical verification command** is `./scripts/verify.sh`
  (equivalently `make verify`). Run it before every handoff and record the
  result in `docs/backlog/active.md`. Do not invent ad hoc verification steps;
  add the check to that script instead.
