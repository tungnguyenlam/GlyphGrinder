# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M2 — It looks like the pitch.** Make GlyphGrinder visually distinct and responsive:
implement Field of View (FOV) & map memory, viewport/camera navigation, rich color palettes with Lip Gloss profile detection, Nerd Font glyph sets with ASCII fallback, and smooth tick-driven movement animation.

Serves `GOAL.md` directly: "A dungeon fades in... a torch pool of warm light slides across stone while the corridor behind you falls back into blue-grey memory... Nerd Font glyphs give every monster, door and potion a silhouette you recognise at a glance, and truecolor gives the dungeon depth."

## Exact next action

**M2.1 — Implement Field of View (FOV) and map memory in `GameMap`.**

- In `game.go`:
  - Add FOV / visibility tracking to `GameState` / `GameMap`: `Explored [][]bool` (tiles seen at least once) and `Visible [][]bool` (tiles currently lit by line-of-sight from player position).
  - Implement a pure FOV calculation function (e.g. shadowcasting or raycasting line-of-sight check up to a radius of ~6 tiles around player position).
  - Update `NewGameWithSeed` and `Step` to recalculate `Visible` and `Explored` tiles whenever player moves.
- In `main.go`:
  - Update `View()` rendering:
    - Currently visible tiles (`Visible[y][x] == true`): render bright entity/floor/wall colors.
    - Previously explored tiles (`Explored[y][x] == true && Visible[y][x] == false`): render dimmed floor/wall colors; do not render monsters on non-visible tiles.
    - Unexplored tiles (`Explored[y][x] == false`): render as blank space `" "`.
- Add unit tests in `game_test.go` and TUI tests in `tui_test.go` asserting FOV calculation, hiding unexplored tiles, dimming remembered tiles, and hiding non-visible monsters.

Why next: M2.1 is the foundational step of M2 (Lighting and visual pitch), establishing visibility and memory maps needed before camera, color ramps, and glyph sets.

## Milestone plan

- [ ] **M2.1** Field of view (FOV) & map memory — lit tiles bright, remembered tiles dimmed, unexplored tiles hidden.
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 20 test functions pass, `go.mod` tidy.
