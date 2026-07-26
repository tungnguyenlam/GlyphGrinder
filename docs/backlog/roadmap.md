# Roadmap

Committed future work, in order. A milestone moves from here into
`active.md` when it starts. Nothing here is a promise about *when*.

## M1 — Playable core loop *(completed 2026-07-27 — see [done.md](done.md))*

Generated dungeon, monsters, bump-to-attack combat, death and restart.

## M2 — It looks like the pitch *(completed 2026-07-27 — see [done.md](done.md))*

Nerd Font glyphs, truecolor palette, FOV shadowcasting, movement animation, viewport camera.

## M3 — A run worth finishing *(completed 2026-07-27 — see [done.md](done.md))*

Multiple dungeon floors, items & inventory (potions & weapons), varied monster types (Trolls & Archers) with ranged AI, and Amulet of Yendor win condition.

## M4 — Ships to other people *(in progress — see [active.md](active.md))*

- `go install` verified from a clean machine; README rewritten as install +
  controls + screenshot.
- Startup degradation checks: no Nerd Font, no truecolor, tiny terminal, resize
  mid-run.
- Benchmark the render path and hold a frame budget (GOAL.md: smooth, on
  battery, no fan).
- CI running `./scripts/verify.sh` on push.
