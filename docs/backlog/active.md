# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M2 — It looks like the pitch.** Make GlyphGrinder visually distinct and responsive:
implement Field of View (FOV) & map memory, viewport/camera navigation, rich color palettes with Lip Gloss profile detection, Nerd Font glyph sets with ASCII fallback, and smooth tick-driven movement animation.

Serves `GOAL.md` directly: "A dungeon fades in... a torch pool of warm light slides across stone while the corridor behind you falls back into blue-grey memory... Nerd Font glyphs give every monster, door and potion a silhouette you recognise at a glance, and truecolor gives the dungeon depth."

## Exact next action

**M2.2 — Viewport/camera & terminal window resizing.**

- In `main.go`:
  - Add `width` and `height` fields to `model` for terminal dimensions (initially 80×24 or from first `tea.WindowSizeMsg`).
  - Handle `tea.WindowSizeMsg` in `Update` to capture terminal size.
  - In `View()`, compute a camera viewport centered on the player that clips the map to the terminal size. Only render the visible sub-rectangle of the map grid.
  - Support map sizes larger than the current 20×10 (e.g. 60×30). Update `initialModel()` to use the larger map.
- In `game.go`:
  - Ensure `NewGameWithSeed` and `GenerateMap` work correctly with larger map sizes.
- Add TUI tests using `d.Resize(w, h)` to verify camera follows the player and adapts to terminal size changes.

Why next: With FOV and map memory landed (M2.1), the next step is to support larger maps via a scrolling viewport, which is prerequisite for all subsequent visual polish.

## Milestone plan

- [x] **M2.1** Field of view (FOV) & map memory — lit tiles bright, remembered tiles dimmed, unexplored tiles hidden.
- [ ] **M2.2** Viewport/camera & terminal window resizing — support larger map sizes with camera centered on player and handling `tea.WindowSizeMsg`.
- [ ] **M2.3** Color palette & Lip Gloss profile-aware rendering — truecolor stone/floor/entity ramps with graceful degradation.
- [ ] **M2.4** Nerd Font glyph set with ASCII fallback.
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 29 test functions pass (9 new FOV tests), `go.mod` tidy.
