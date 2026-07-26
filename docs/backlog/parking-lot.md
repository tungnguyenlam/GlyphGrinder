# Parking Lot

Uncommitted ideas. Nothing here is a promise, and deleting from this file is
cheap — do it freely when an idea stops being interesting or conflicts with
[GOAL.md](../../GOAL.md).

Anything vague that shows up in `active.md` gets demoted here rather than
lingering as a fuzzy reminder.

Items already shipped in M1–M7 (title/end screens, hit flash, status effects,
ranged targeting, classes, save/resume, seed+input replay, `--dump-frame`, FOV,
camera easing, map reachability fuzz, View style cache) were removed on
2026-07-27 so this list stays actionable.

## Feel and presentation

- Screen shake on heavy hits (hit flash already exists; shake would sell impact).
- Particle-ish effects for spells and explosions using half-block characters.
- Per-tile background colors for water, lava, blood — needs a call on whether
  backgrounds fight readability.
- Sound? Almost certainly a non-goal (terminal bell only), noted so it stops
  getting re-proposed.

## Game design

- Hunger clock, or an explicit decision not to have one.
- Skill tree beyond the three class archetypes (Warrior / Rogue / Mage).
- Persistent meta-progression between runs — tension with permadeath in
  GOAL.md; would need an ADR.

## Engineering

- Golden-frame testing: snapshot `View` output and diff it across glyph/color
  modes.
- More hazard/item interaction coverage (e.g. lava + regen race, multi-weapon
  inventory edge cases).
