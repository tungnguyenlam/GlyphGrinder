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

**Animate resolved player and monster movement through tick-driven render
state.**

- Capture actor positions before and after each resolved turn in `Update`, then
  store short-lived interpolation state without changing deterministic
  `GameState` rule positions.
- Schedule fixed-rate Bubble Tea tick messages while motion is active and ease
  render coordinates toward their resolved cells; never read time or mutate in
  `View`.
- Keep input responsive and define how another key during motion settles or
  replaces the current transition before resolving the next turn.
- Add headless frame-by-frame tests for player/monster movement, camera
  interaction, rest state, and restart/death behavior.

Why next: visibility, palette, camera, and glyph silhouettes now provide the
static presentation promised by `GOAL.md`; motion is the last M2 acceptance
criterion.

## Milestone plan

- [x] **M2.1** Field of view and exploration memory with wall occlusion;
      invisible monsters stay hidden and explored terrain renders dimly. —
      *completed 2026-08-23*.
- [x] **M2.2** Replace ad hoc colors with a cohesive truecolor terrain/entity
      palette that degrades through Lip Gloss color profiles. — *completed
      2026-08-23*.
- [x] **M2.3** Add a resize-aware viewport and player-following camera, then
      grow generated maps beyond the initial 20x10 frame. — *completed
      2026-08-23*.
- [x] **M2.4** Add recognizable terrain/entity glyph profiles with a dependable
      ASCII fallback. — *completed 2026-08-23*.
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

M2.4 environment-selected rich/ASCII glyph profiles: gofmt clean, `go vet`
clean, builds, headless smoke frame renders, all tests pass, `go.mod` tidy.

Standing unattended prompt updated to execute substantial batches of related
work while retaining progressive tests, checkpoints, and reviewable commits.
