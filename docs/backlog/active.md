# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**Post-M7 Polish & Maintenance.** All milestones (M1–M7) are complete.

Serves `GOAL.md` directly: keep the finished roguelike correct, playable, and
easy for the next agent to extend without rediscovering finished work.

## Exact next action

**Golden-frame regression harness (parking-lot engineering item).**

1. Add a small helper (e.g. in `tui_test.go` or a new `golden_test.go`) that
   builds a fixed-seed model, forces ASCII glyphs + a stable color profile, and
   writes/compares `View()` output against a checked-in golden file under
   `testdata/`.
2. Start with one golden: title screen, and one in-game frame (seed `42`, depth
   1, Warrior, after a fixed key sequence). Use an update flag
   (`-update-golden`) only in that test file, not a new dependency.
3. Run `./scripts/verify.sh` and record the result below.

Why next: protects View/palette/glyph regressions without a TTY; parking lot
still lists this as open engineering work after M7.

## Recently completed (this session)

- [x] **Bugfix: dropping a weapon no longer keeps its DamageBonus.**
  `ActionDropItem` now subtracts `DamageBonus` when the dropped item is a
  weapon (clamped at 0). Regression tests:
  `TestDropWeaponRemovesDamageBonus`, `TestWarriorDropStartingWeapon`.
  Without the fix, Warriors (and anyone who drop+re-picked a dagger) stacked
  permanent damage.
- [x] **Parking lot scrub.** Removed M1–M7 items that already shipped so the
  lot is only uncommitted ideas again.

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
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, full test
suite pass (including `TestDropWeaponRemovesDamageBonus` and
`TestWarriorDropStartingWeapon`), `go.mod` tidy.
