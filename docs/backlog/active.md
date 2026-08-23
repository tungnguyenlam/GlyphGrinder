# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-08-23

## Current milestone

**M3 — A run worth finishing.** Turn the polished single-level combat loop into
a progressing run with dungeon levels, items, varied enemies, and a win
condition.

Serves `GOAL.md` directly: M2 delivered the visual pitch; M3 now supplies the
longer tactical arc and build-changing choices required of a real roguelike.

## Exact next action

**Add stairs and deterministic descent to deeper generated levels.**

- Add a down-stair tile on a reachable floor away from the player and track the
  current dungeon depth in `GameState`.
- Map `>` to a descend action that works only while standing on stairs, creates
  the next seeded dungeon, preserves the run's player health, and resets
  visibility/camera motion safely.
- Render stairs through both rich and ASCII glyph profiles and show depth in
  the sidebar.
- Cover deterministic stair placement, invalid descent, successful descent,
  preserved stats, and the full headless input/render path.

Why next: a run currently ends on the first map with no objective. Descent is
the smallest complete vertical slice that establishes progression for items,
enemy variety, difficulty, and an eventual win condition.

## Milestone plan

- [ ] **M3.1** Add dungeon depth, reachable stairs, and deterministic descent.
- [ ] **M3.2** Add potions with pickup, inventory, and tactical use.
- [ ] **M3.3** Add a weapon tier and equipment choice that changes combat.
- [ ] **M3.4** Add at least two monster archetypes with distinct stats and
      movement behavior, scaling by depth.
- [ ] **M3.5** Add a final depth and explicit win state with one-key restart.

## Acceptance criteria for M3

- A run spans multiple increasingly difficult dungeon levels and ends in an
  explicit win or death.
- Potions and weapons create at least two meaningful tactical choices during a
  run; inventory/use/equip behavior is deterministic and visible in the UI.
- At least three monster archetypes are recognizable by glyph/color and behave
  differently enough to change positioning decisions.
- Progression, items, enemies, death, restart, and victory have state-level and
  headless coverage; no test requires a TTY.
- `./scripts/verify.sh` passes.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-08-23  —  VERIFY OK
```

M2.5 tick-driven player/monster motion and M2 completion: gofmt clean, `go vet`
clean, builds, headless smoke frame renders, all tests pass, `go.mod` tidy.

Standing unattended prompt updated to execute substantial batches of related
work while retaining progressive tests, checkpoints, and reviewable commits.
