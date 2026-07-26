# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-07-27

## Current milestone

**Post-M7 Polish & Maintenance.** All milestones (M1 Playable Core Loop, M2 Visual Pitch & FOV, M3 Deep Progression & Items, M4 Shipping & CI, M5 Dungeon Mechanics & Magical Depth, M6 Visual Brilliance & Tactical Systems, M7 Replayability, Headless Inspection & Class Archetypes) are complete.

Serves `GOAL.md` directly: GlyphGrinder is a complete, highly polished, fast, responsive, and fully verified Go terminal roguelike.

## Exact next action

**Maintenance Mode / Post-Launch Polish.**

- Run `./scripts/verify.sh` to confirm green build and test status.
- Consult `docs/backlog/parking-lot.md` for post-launch ideas when starting new feature work.

Why next: Milestones M1 through M7 have been completed and verified against `GOAL.md`.

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

## Blockers

None. No decision is waiting on the user.

## Last verification

```
./scripts/verify.sh   —  2026-07-27  —  VERIFY OK
```

gofmt clean, `go vet` clean, builds, headless smoke frame renders, 85 test functions pass + 3 performance benchmarks pass, `go.mod` tidy.






