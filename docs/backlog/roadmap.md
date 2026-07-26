# Roadmap

Committed future work, in order. A milestone moves from here into
`active.md` when it starts. Nothing here is a promise about *when*.

## M1 — Playable core loop *(in progress — see [active.md](active.md))*

Generated dungeon, monsters, bump-to-attack combat, death and restart.

## M2 — It looks like the pitch

The point at which the README's claim stops being aspirational.

- Nerd Font glyph set for entities and terrain, with an ASCII fallback when the
  terminal or font can't do it.
- A real palette: truecolor stone/floor/entity ramps instead of the current
  three hard-coded greys, with graceful degradation via `lipgloss`'s color
  profile detection.
- Field of view and memory: lit tiles bright, remembered tiles dimmed,
  unexplored tiles black.
- Movement animation — entities interpolate between tiles over a few frames,
  driven by tick messages (ADR-0003), never by `View` reading the clock.
- A viewport/camera so the map can be larger than the terminal, following the
  player and reacting to `tea.WindowSizeMsg` (currently ignored entirely).

## M3 — A run worth finishing

- Multiple dungeon levels with stairs and increasing difficulty.
- Items: pick up, inventory, use/equip; at minimum potions and one weapon tier.
- More than one monster type, with different stats and movement behavior.
- A win condition — a reason to descend.

## M4 — Ships to other people

- `go install` verified from a clean machine; README rewritten as install +
  controls + screenshot.
- Startup degradation checks: no Nerd Font, no truecolor, tiny terminal, resize
  mid-run.
- Benchmark the render path and hold a frame budget (GOAL.md: smooth, on
  battery, no fan).
- CI running `./scripts/verify.sh` on push.
