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

**Replace ad hoc colors with a semantic dungeon palette.**

- Define one semantic palette for lit stone/floor, remembered terrain, player,
  monsters, health, log text, and danger/game-over UI instead of scattered
  literals.
- Apply foreground/background depth consistently while relying on Lip Gloss's
  terminal color-profile conversion for graceful degradation.
- Keep colors out of rule decisions and preserve plain-text glyph output; add
  focused rendering tests for semantic style distinctions without snapshotting
  entire ANSI frames.

Why next: M2.1 now exposes lit and remembered states, but their current greys
and actor colors are still disconnected literals rather than a designed ramp.

## Milestone plan

- [x] **M2.1** Field of view and exploration memory with wall occlusion;
      invisible monsters stay hidden and explored terrain renders dimly. —
      *completed 2026-08-23*.
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

M2.1 field of view and exploration memory: gofmt clean, `go vet` clean, builds,
headless smoke frame renders, all tests pass, `go.mod` tidy.
