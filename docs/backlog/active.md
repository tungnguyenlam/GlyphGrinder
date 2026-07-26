# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M1 — Playable core loop.** Turn the current fixed 20x10 walled room into a
small but real roguelike level: a generated dungeon, at least one monster that
takes a turn, and bump-to-attack combat that can kill and be killed.

Serves `GOAL.md` directly: right now there is nothing to play, and the vision's
"real roguelike underneath the polish" has no polish worth adding until the
loop exists. Visual work (animation, lighting, Nerd Font glyph sets) is
deliberately deferred to M2 — park ideas in `parking-lot.md`.

## Exact next action

**Implement dungeon map generation (rooms + corridors) in `game.go` with deterministic seed support.**

- Define generator function/struct for creating `GameMap` with rooms connected by corridors.
- Support deterministic RNG seeding stored in model/state (ADR-0003).
- Update `initialModel` in `main.go` to generate dungeon maps.
- Ensure spawn position places player in a valid floor tile inside a room.
- Add tests in `tui_test.go` verifying map generation and deterministic seed behavior.

Why next: dungeon map generation (M1.3) replaces the static 20x10 room with real levels needed for monster placement and combat (M1.4 - M1.6).

## Milestone plan

- [x] **M1.1** Factor movement into `tryMove(dx, dy)` — *completed 2026-07-27*.
- [x] **M1.2** Move turn resolution out of the key switch: `Update` handles
      input, a `Step(action Action)` method on `GameState` advances the world
      one turn. — *completed 2026-07-27*.
- [ ] **M1.3** Map generation — rooms plus corridors into `GameMap`, seeded
      from an explicit RNG stored in the model so tests are deterministic
      (ADR-0003). Replace the hard-coded 20x10 room in `initialModel`.
- [ ] **M1.4** Populate `GameState.Entities` with monsters at generation time
      and render them in `View` (currently `Entities` is never read or
      written).
- [ ] **M1.5** Bump-to-attack: moving into an occupied tile deals
      `Entity.Damage`, reduces `Health`, removes the entity at zero, and
      appends a line to `GameState.Log`.
- [ ] **M1.6** Monster turns: each monster steps toward the player after the
      player's turn resolves.
- [ ] **M1.7** Render the log and a health bar alongside the map; handle player
      death with a game-over state and restart on one keypress.

## Acceptance criteria for M1

- Launching produces a different dungeon each run; the same seed produces the
  same dungeon.
- The player can walk into a monster, damage it, kill it, and see it in the log.
- A monster can kill the player, and the game says so instead of hanging.
- `GameState.Entities` and `GameState.Log` are both read and written by real
  gameplay code.
- Every rule above is covered by a `internal/tuitest`-driven test; no test
  requires a TTY.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 4 test
functions (11 cases) pass, `go.mod` tidy.
