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

**Add a final depth and explicit victory flow.**

- Set a short, deterministic final depth and turn the final staircase into the
  run exit instead of generating another unbounded level.
- Add a distinct won state that freezes turns, renders a clear victory prompt,
  and restarts a fresh run with `r` just like death.
- Cover pre-final descent, final-depth victory, frozen post-win input, restart,
  and the complete headless key/render path.

Why next: depth, items, equipment, and enemy variety now create a progressing
run, but descent is still endless. A final exit closes the run arc required by
M3 and gives the player a reason to survive deeper floors.

## Milestone plan

- [x] **M3.1** Add dungeon depth, reachable stairs, and deterministic descent.
- [x] **M3.2** Add potions with pickup, inventory, and tactical use.
- [x] **M3.3** Add a weapon tier and equipment choice that changes combat.
- [x] **M3.4** Add at least two monster archetypes with distinct stats and
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

M3.2 deterministic health potions: placement, pickup, carried inventory,
turn-costing capped healing, restart/descent behavior, rich/ASCII rendering,
and full headless input/render coverage pass.

M3.3 iron sword equipment: deterministic placement/pickup, turn-costing equip,
15-damage combat, descent/restart state, rich/ASCII rendering, and headless
input/render coverage pass.

M3.4 goblin, ogre, and bat archetypes: deterministic mixed generation,
depth-scaled stats, one/alternate/two-action movement cadence, ID-ordered
collision and combat, named logs, distinct rich/ASCII glyphs and colors, and
headless visibility coverage pass.

Standing unattended prompt updated to execute substantial batches of related
work while retaining progressive tests, checkpoints, and reviewable commits.
