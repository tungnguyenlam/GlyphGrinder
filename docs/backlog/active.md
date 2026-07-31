# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-31

## Current milestone

**Post-M7 Polish & Maintenance.** All milestones (M1–M7) are complete.

Serves `GOAL.md` directly: keep the finished roguelike correct, playable, and
easy for the next agent to extend without rediscovering finished work.

## Exact next action

**Pick next item from parking lot or declare maintenance complete.**

Remaining parking lot items: sound (likely non-goal), hunger clock, skill tree
expansion, persistent meta-progression (needs ADR), more hazard/item interaction
coverage (lava + regen race, multi-weapon inventory edge cases). The project is
in solid shape — all milestones complete, tree green, no known bugs. Next agent
can write coverage tests from the engineering section of `parking-lot.md` or
stop here.

## Recently completed (this session)

- [x] **Bugfix: gofmt + go vet failures.** Previous session left `game.go` and
  `main.go` not gofmt-clean. Also found and fixed two duplicate function
  declarations: `xpForLevel` (duplicate at end of file with a different formula,
  kept the first `100 * level * level` version near its caller `gainXP`) and
  `AvailableSkills` (identical copy at end of file). These prevented
  compilation.
- [x] **Bugfix: `TestBumpToAttackKillsMonster` stale assertion.** The test
  expected 2 log entries but the XP system (added later) now also logs "You gain
  10 XP." on kill. Updated to expect 3 log entries.
- [x] **Per-tile background colors for lava and water.** Added `LavaBg`
  (`#3D1500` dark red-orange) and `WaterBg` (`#001A2E` dark deep blue) to the
  `Palette` struct. Both lit and dimmed tile styles now include
  `.Background()`. Regression test `TestTileBackgroundColors` verifies all four
  variants (lit/dim × lava/water) contain ANSI `48;2;` background escapes and
  that floor/wall/door tiles do not. Golden frames regenerated.
- [x] **Palette token test coverage.** Added missing tokens (Scroll, Door, Lava,
  LavaBg, Water, WaterBg) to `TestDefaultPaletteTokens` in `colors_test.go`.
- [x] **Parking lot scrub.** Removed stale entries for screen shake, particle
  effects, golden-frame testing (all shipped last session), and per-tile
  background colors (shipped this session).

## Milestone plan

- [x] **M1** Playable core loop (dungeon generation, movement, combat, HUD, death/restart).
- [x] **M2** Visual pitch (Nerd Font & ASCII fallbacks, TrueColor palette, FOV shadowcasting, camera easing).
- [x] **M3** A run worth finishing (stairs & depth, items & inventory, Trolls & Archers with ranged AI, Amulet of Yendor win condition).
- [x] **M4** Ships to other people (`go install`, README rewrite, degradation checks, benchmarks, GitHub Actions CI).
- [x] **M5** Dungeon Mechanics & Magical Depth (Fireball/Teleport scrolls, interactive doors, status effects, seed replayability).
- [x] **M6** Visual Brilliance & Tactical Systems (Title/End screens, hazard tiles, item dropping, targeting UI, save/resume, render optimization).
- [x] **M7** Replayability, Headless Inspection & Class Archetypes (Input replay, `--dump-frame`, map reachability fuzzing, Warrior/Rogue/Mage archetypes).

## Acceptance criteria

- All milestones M1–M7 completed and logged in `docs/backlog/done.md`.
- `README.md` is complete and accurate with current features and controls.
- Performance benchmarks confirm sub-millisecond execution (< 0.065ms per View render).
- `./scripts/verify.sh` passes clean.
- Known rule bugs found in play (weapon drop bonus leak) are fixed with regression tests.

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-31  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, full test
suite pass (including `TestTileBackgroundColors`, `TestBumpToAttackKillsMonster`
with XP assertion, `TestDefaultPaletteTokens` with all 22 tokens), `go.mod`
tidy.
