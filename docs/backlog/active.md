# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**M5 — Dungeon Mechanics & Magical Depth.** Expand tactical gameplay with consumable magic scrolls, interactive doors, status effects, and seed replayability.

Serves `GOAL.md` directly: "It is a real roguelike underneath the polish... items that change how you play... turn-based combat where a mistake costs you the run."

## Exact next action

**M5.4 — Seed Replayability & Input Verification.**

- In `main.go`:
  - Display run seed in HUD status bar (`Seed: X`) for run sharing and bug reproduction.
  - Support setting initial seed via `GLYPHGRINDER_SEED` environment variable or CLI flag.
- Add unit and TUI tests in `tui_test.go` verifying deterministic state replayability across identical seed and keypress sequences.

Why next: M5.4 makes runs 100% reproducible for bug reports, competitive seed sharing, and regression testing.

## Milestone plan

- [x] **M5.1** Magic Scrolls (Fireball AoE damage & Teleportation scrolls).
- [x] **M5.2** Interactive Doors (Carved closed doors `+` / open doors `/` between rooms).
- [x] **M5.3** Status Effects (Poison, Regeneration, Confusion with turn duration).
- [ ] **M5.4** Seed Replayability & Input Verification (Deterministic seed options and golden frame tests).

## Acceptance criteria for M5

- Magic scrolls generate, pick up to inventory, and execute spell effects (Fireball AoE damage, Teleportation).
- Doors spawn between rooms/corridors and open upon player interaction.
- Status effects apply over turns and affect player/monster combat.
- Seed display and deterministic input replayability verified by unit and TUI tests.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 72 test functions pass (including 2 new status effect tests) + 3 performance benchmarks pass, `go.mod` tidy.


