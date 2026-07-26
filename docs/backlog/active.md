# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M2 — It looks like the pitch.** Make GlyphGrinder visually distinct and responsive:
implement Field of View (FOV) & map memory, viewport/camera navigation, rich color palettes with Lip Gloss profile detection, Nerd Font glyph sets with ASCII fallback, and smooth tick-driven movement animation.

Serves `GOAL.md` directly: "A dungeon fades in... a torch pool of warm light slides across stone while the corridor behind you falls back into blue-grey memory... Nerd Font glyphs give every monster, door and potion a silhouette you recognise at a glance, and truecolor gives the dungeon depth."

## Exact next action

**M2.5 — Tick-driven movement animation.**

- In `main.go` / `game.go`:
  - Introduce smooth interpolation / tick-driven animation state for entity movements or camera slides, ensuring turn-based logic remains discrete while rendering provides smooth transitions.
  - Implement tick messages (`tea.Tick`) to advance animation frames smoothly when player or entities move.
- Add unit/TUI tests verifying frame progression and input responsiveness during ticks.

Why next: M2.5 completes M2 ("It looks like the pitch"), fulfilling the `GOAL.md` vision where "the `@` you control moves with weight — steps ease, the camera drifts to follow... nothing snaps; things move."

## Milestone plan

- [x] **M2.1** Field of view (FOV) & map memory — lit tiles bright, remembered tiles dimmed, unexplored tiles hidden.
- [x] **M2.2** Viewport/camera & terminal window resizing — support larger map sizes with camera centered on player and handling `tea.WindowSizeMsg`.
- [x] **M2.3** Color palette & Lip Gloss profile-aware rendering — truecolor stone/floor/entity ramps with graceful degradation.
- [x] **M2.4** Nerd Font glyph set with ASCII fallback.
- [ ] **M2.5** Tick-driven movement animation.

## Acceptance criteria for M2

- Unexplored dungeon tiles are hidden until within line-of-sight of player.
- Previously seen tiles remain visible in a dimmed memory state.
- Camera follows player smoothly across larger maps and handles terminal resize events (`tea.WindowSizeMsg`).
- Color palette scales gracefully depending on terminal color capabilities.
- All visual and gameplay rules covered by `internal/tuitest` unit and TUI tests.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 39 test functions pass (4 new glyph set & environment detection tests), `go.mod` tidy.
