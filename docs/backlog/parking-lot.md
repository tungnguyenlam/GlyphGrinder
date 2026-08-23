# Parking Lot

Uncommitted ideas. Nothing here is a promise, and deleting from this file is
cheap — do it freely when an idea stops being interesting or conflicts with
[GOAL.md](../../GOAL.md).

Anything vague that shows up in `active.md` gets demoted here rather than
lingering as a fuzzy reminder.

## Feel and presentation

- Screen shake / hit flash on damage.
- Particle-ish effects for spells and explosions using half-block characters.
- Smooth camera easing rather than a hard follow.
- A title screen and a death screen worth screenshotting.
- Per-tile background colors for water, lava, blood — needs a call on whether
  backgrounds fight readability.
- Sound? Almost certainly a non-goal (terminal bell only), noted so it stops
  getting re-proposed.

## Game design

- Hunger clock, or an explicit decision not to have one.
- Ranged attacks and line-of-sight targeting UI.
- Status effects (poison, slow, confusion).
- Character classes or a skill tree.
- Persistent meta-progression between runs — tension with permadeath in
  GOAL.md; would need an ADR.

## Engineering

- Save/resume mid-run (GOAL.md says permadeath, so this is only for
  quit-and-continue, not save-scumming).
- Replay from a seed plus an input log — would make bug reports reproducible
  and is cheap given ADR-0003's pure `View`.
- Golden-frame testing: snapshot `View` output and diff it.
- A `--dump-frame` flag so the real binary can be smoke-tested without a TTY
  (would resolve the TTY notice).
