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

**Add recognizable terrain/entity glyph profiles with an automatic ASCII
fallback.**

- Centralize the player, monster, floor, and wall glyphs into explicit rich and
  ASCII profiles rather than scattering literals through state and rendering.
- Select the rich profile only when the process environment indicates a
  Unicode-capable terminal; keep plain ASCII dependable in other environments.
- Preserve one-cell alignment for every profile and restart, and add headless
  tests for selection and rendered glyph semantics.

Why next: the viewport now delivers a large, camera-followed dungeon; distinct
silhouettes and a reliable fallback are the remaining static visual layer
before motion work begins.

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

M2.3 resize-aware viewport and 96x48 production dungeon: gofmt clean, `go vet`
clean, builds, headless smoke frame renders, all tests pass, `go.mod` tidy.

Standing unattended prompt updated to execute substantial batches of related
work while retaining progressive tests, checkpoints, and reviewable commits.
