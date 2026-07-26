# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M3 — A run worth finishing.** Build depth and gameplay progression:
multiple dungeon levels with stairs, item pickup/inventory (potions & weapons), varied monster types/AI, and a victory condition.

Serves `GOAL.md` directly: "It is a real roguelike underneath the polish, not a tech demo: procedurally generated levels, permadeath, turn-based combat where a mistake costs you the run, items that change how you play."

## Exact next action

**M3.1 — Multiple dungeon levels & stairs.**

- In `game.go`:
  - Add `Depth int` to `GameState`.
  - Add `TileStairsDown` to `TileType` and place down stairs (`>`) in the furthest room during `GenerateMap`.
  - Handle descending stairs when player steps onto stairs or presses `>` / `enter`.
- In `main.go`:
  - Update HUD status bar to display `Depth: X`.
  - Add stair glyph support to `GlyphSet`.
- Add unit/TUI tests in `game_test.go` and `tui_test.go` verifying stair generation, descending level progression, depth HUD display, and FOV reset.

Why next: M3.1 establishes multi-floor dungeon progression, giving players a sense of descent and increasing challenge.

## Milestone plan

- [ ] **M3.1** Multiple dungeon levels & stairs (`>`) with depth HUD tracking.
- [ ] **M3.2** Items & Inventory (health potions `!`, weapons `/`, pickup `g`/`,`, drink `h`).
- [ ] **M3.3** Enhanced monster types & AI (Troll `T`, Archer `A`).
- [ ] **M3.4** Win condition & victory screen (Amulet of Yendor on Depth 5).

## Acceptance criteria for M3

- Player can descend through multiple dungeon floors via down stairs.
- Items generate on floors, can be picked up, viewed in inventory, and consumed.
- New monster types have distinct stats and behaviors.
- Descending to Depth 5 and retrieving the goal item triggers a victory state.
- All rules and interactions covered by `internal/tuitest` unit and TUI tests.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 41 test functions pass (2 new tick animation & camera easing tests), `go.mod` tidy.

