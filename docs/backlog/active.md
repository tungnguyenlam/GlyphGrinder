# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M5 — Dungeon Mechanics & Magical Depth.** Expand tactical gameplay with consumable magic scrolls, interactive doors, status effects, and seed replayability.

Serves `GOAL.md` directly: "It is a real roguelike underneath the polish... items that change how you play... turn-based combat where a mistake costs you the run."

## Exact next action

**M5.2 — Interactive Doors (Closed Doors `+` & Open Doors `/`).**

- In `game.go`:
  - Add `TileDoorClosed` and `TileDoorOpen` to `TileType`.
  - Place `TileDoorClosed` at room entrances during corridor carving in `GenerateMap`.
  - In `Step(Action)`: Bumping `TileDoorClosed` transforms it to `TileDoorOpen`, logs `"You open the door."`, consumes 1 turn, and reveals room via FOV.
  - Update `ComputeFOV` shadowcasting: `TileDoorClosed` blocks vision; `TileDoorOpen` allows vision.
- In `glyphs.go` & `colors.go`:
  - Add `DoorClosed` (`+`, `󰌝`) and `DoorOpen` (`/`, `󰌟`) to `GlyphSet` and wood color (`#D7CCC8`) to `Palette`.
- Add unit and TUI tests in `game_test.go` and `tui_test.go` verifying door opening, turn consumption, and FOV line of sight.

Why next: M5.2 adds classic roguelike room boundary mechanics and tactical FOV reveals when opening doors.

## Milestone plan

- [x] **M5.1** Magic Scrolls (Fireball AoE damage & Teleportation scrolls).
- [ ] **M5.2** Interactive Doors (Carved closed doors `+` / open doors `/` between rooms).
- [ ] **M5.3** Status Effects (Poison, Regeneration, Confusion with turn duration).
- [ ] **M5.4** Seed Replayability & Input Verification (Deterministic seed options and golden frame tests).

## Acceptance criteria for M5

- Magic scrolls generate, pick up to inventory, and execute spell effects (Fireball AoE damage, Teleportation).
- Doors spawn between rooms/corridors and open upon player interaction.
- Status effects apply over turns and affect player/monster combat.
- All rules and interactions covered by unit and TUI tests.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 68 test functions pass (including 3 new Fireball AoE, Teleport, and TUI scroll tests) + 3 performance benchmarks pass, `go.mod` tidy.


