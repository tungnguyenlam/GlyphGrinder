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

**Add deterministic health potions with pickup, inventory, and tactical use.**

- Place potions deterministically on ordinary floor cells without overlapping
  the player, stairs, monsters, or one another.
- Pick a potion up by moving onto it, carry potion inventory across dungeon
  depths, and clear inventory on a fresh run.
- Map `p` to consume one potion only while injured, heal without exceeding max
  health, spend a turn so nearby monsters can respond, and log the result.
- Render potions through rich and ASCII glyph profiles, show their inventory
  count in the sidebar, and cover placement/pickup/use plus headless input and
  rendering.

Why next: descent now establishes a run, but health only decreases. Potions add
the first build resource and a meaningful risk/reward decision about when to
spend a combat turn recovering.

## Milestone plan

- [x] **M3.1** Add dungeon depth, reachable stairs, and deterministic descent.
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

M3.1 reachable stairs and deterministic descent: state and headless coverage
pass; rich/ASCII stairs, depth sidebar, preserved health, fresh visibility,
and cleared camera motion verified.

Standing unattended prompt updated to execute substantial batches of related
work while retaining progressive tests, checkpoints, and reviewable commits.
