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

**Render HUD & message log in `View`, and handle game-over state.**

- In `main.go`, update `View()`:
  - Render a status bar displaying player HP (e.g. `HP: 94/100`).
  - Render the last 3–5 entries of `m.state.Log` below the map.
  - If `m.state.Player.Health <= 0`, render a GAME OVER header/banner indicating game over and restart prompt ("Press r to restart").
- In `main.go`, update `Update()`:
  - If player is dead (`Player.Health <= 0`), listen for `"r"` key to restart the game by re-initializing the state (`m.state = NewGame(20, 10)`).
- Add unit tests in `game_test.go` and TUI tests in `tui_test.go` asserting HUD/log rendering, game-over screen rendering, and restart key functionality.

Why next: M1.7 is the final sub-task of M1 (Playable core loop). Finishing M1.7 satisfies all M1 acceptance criteria.

## Milestone plan

- [x] **M1.1** Factor movement into `tryMove(dx, dy)` — *completed 2026-07-27*.
- [x] **M1.2** Move turn resolution out of the key switch: `Update` handles
      input, a `Step(action Action)` method on `GameState` advances the world
      one turn. — *completed 2026-07-27*.
- [x] **M1.3** Map generation — rooms plus corridors into `GameMap`, seeded
      from an explicit RNG stored in the model so tests are deterministic
      (ADR-0003). Replace the hard-coded 20x10 room in `initialModel`. — *completed 2026-07-27*.
- [x] **M1.4** Populate `GameState.Entities` with monsters at generation time
      and render them in `View`. — *completed 2026-07-27*.
- [x] **M1.5** Bump-to-attack: moving into an occupied tile deals
      `Entity.Damage`, reduces `Health`, removes the entity at zero, and
      appends a line to `GameState.Log`. — *completed 2026-07-27*.
- [x] **M1.6** Monster turns: each monster steps toward the player after the
      player's turn resolves. — *completed 2026-07-27*.
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 18 test
functions pass, `go.mod` tidy.
