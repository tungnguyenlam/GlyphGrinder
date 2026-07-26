# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M3 — A run worth finishing.** Build depth and gameplay progression:
multiple dungeon levels with stairs, item pickup/inventory (potions & weapons), varied monster types/AI, and a victory condition.

Serves `GOAL.md` directly: "It is a real roguelike underneath the polish, not a tech demo: procedurally generated levels, permadeath, turn-based combat where a mistake costs you the run, items that change how you play."

## Exact next action

**M3.4 — Win condition & victory screen (Amulet of Yendor on Depth 5).**

- In `game.go`:
  - Add `ItemAmulet` item type ("Amulet of Yendor", `*`).
  - On Depth 5 (`depth == 5`), spawn the Amulet of Yendor instead of down stairs.
  - Picking up the Amulet sets `IsVictory bool` on `GameState` and logs victory message.
- In `main.go`:
  - Add `Amulet` (`*`, `󰇮`) to `GlyphSet` and gold color to `Palette`.
  - Update `View()` to display victory HUD banner when `IsVictory` is true.
  - Ignore action keys in victory state (except `r` restart and `q` quit).
- Add unit and TUI tests in `game_test.go` and `tui_test.go` verifying Depth 5 Amulet spawn, victory trigger on pickup, victory screen rendering, and restart key.

Why next: M3.4 completes M3 "A run worth finishing" by providing a satisfying goal and end-game state.

## Milestone plan

- [x] **M3.1** Multiple dungeon levels & stairs (`>`) with depth HUD tracking.
- [x] **M3.2** Items & Inventory (health potions `!`, weapons `/`, pickup `g`/`,`, drink `h`).
- [x] **M3.3** Enhanced monster types & AI (Troll `T`, Archer `A`).
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 59 test functions pass (4 new Troll, Archer ranged attack, depth monster scaling, and TUI archer tests), `go.mod` tidy.


