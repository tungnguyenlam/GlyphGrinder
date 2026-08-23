# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-08-23

## Current milestone

**M2 — It looks like the pitch.** Turn the now-playable roguelike into a
terminal presentation with visibility, depth, motion, and camera behavior that
serves the “that’s a terminal?” goal without sacrificing deterministic rules.

Serves `GOAL.md` directly: M1 supplied the real roguelike underneath; M2 now
targets the dungeon lighting, remembered space, smooth movement, recognizable
glyphs, and terminal-aware framing promised by the vision.

## Exact next action

**Add field of view and exploration memory.**

- Track visible and explored cells in game state, initialized at generation and
  refreshed only when the player changes position; `View` remains pure.
- Use a small radius and wall-blocked line of sight so nearby walls are visible,
  tiles behind them are hidden, and previously seen cells remain remembered.
- Render unseen cells blank, remembered cells dim, and visible cells normally;
  hide monsters outside current visibility.
- Add state tests for occlusion, movement refresh, and memory plus headless
  tests for hidden/remembered rendering.

Why next: visibility is the largest missing piece of the north-star dungeon
look and establishes the state/rendering seam that the palette will style next.

## Milestone plan

- [ ] **M2.1** Field of view and exploration memory with wall occlusion;
      invisible monsters stay hidden and explored terrain renders dimly.
- [ ] **M2.2** Replace ad hoc colors with a cohesive truecolor terrain/entity
      palette that degrades through Lip Gloss color profiles.
- [ ] **M2.3** Add a resize-aware viewport and player-following camera, then
      grow generated maps beyond the initial 20x10 frame.
- [ ] **M2.4** Add recognizable terrain/entity glyph profiles with a dependable
      ASCII fallback.
- [ ] **M2.5** Animate player and monster movement from tick messages without
      moving rule resolution into `View`.

## Acceptance criteria for M2

- The dungeon distinguishes visible, remembered, and unexplored space; unseen
  monsters do not render.
- Terrain and actors use a coherent color/glyph system with a documented ASCII
  fallback path.
- A dungeon larger than the terminal is navigable through a resize-aware camera
  that follows the player.
- Player and monster steps animate through modelled tick state while input and
  rule resolution remain deterministic and `View` stays pure.
- Every behavior above has headless coverage; no test requires a TTY.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-08-23  —  VERIFY OK
```

M1.7 health/log UI, death state, and restart (completing M1): gofmt clean,
`go vet` clean, builds, headless smoke frame renders, all tests pass, `go.mod`
tidy.
