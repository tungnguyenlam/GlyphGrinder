# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M3 — A run worth finishing.** Build depth and gameplay progression:
multiple dungeon levels with stairs, item pickup/inventory (potions & weapons), varied monster types/AI, and a victory condition.

Serves `GOAL.md` directly: "It is a real roguelike underneath the polish, not a tech demo: procedurally generated levels, permadeath, turn-based combat where a mistake costs you the run, items that change how you play."

## Exact next action

**M3.3 — Enhanced monster types & AI (Troll `T`, Archer `A`).**

- In `game.go`:
  - Add `NewTroll(id int, pos Position)` (HP 40, Damage 10) and `NewArcher(id int, pos Position)` (HP 15, Damage 4, ranged attack range 5).
  - Scale monster pool by floor depth in `NewGameWithSeedAndDepth` (Depth 1: Goblin/Orc; Depth 2+: Trolls; Depth 3+: Archers).
  - Update monster AI in `runMonsterTurns`: Archers fire ranged attacks if player is within range 5 and line of sight; Trolls move/attack in melee.
- In `glyphs.go` & `colors.go`:
  - Add `Troll` (`T`, `󰇄`) and `Archer` (`A`, `󰓤`) to `GlyphSet` and colors to `Palette`.
- Add unit and TUI tests in `game_test.go` and `tui_test.go` verifying depth monster scaling, Troll high-damage combat, and Archer ranged attacks.

Why next: M3.3 adds enemy variety and tactical choices, elevating combat depth before M3.4 win condition.

## Milestone plan

- [x] **M3.1** Multiple dungeon levels & stairs (`>`) with depth HUD tracking.
- [x] **M3.2** Items & Inventory (health potions `!`, weapons `/`, pickup `g`/`,`, drink `h`).
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

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 55 test functions pass (9 new item spawning, pickup, potion HP restoration, weapon damage bonus, inventory persistence, and TUI item tests), `go.mod` tidy.


