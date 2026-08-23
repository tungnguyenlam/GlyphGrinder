# Roadmap

Committed future work, in order. A milestone moves from here into
`active.md` when it starts. Nothing here is a promise about *when*.

## M3 — A run worth finishing *(in progress — see [active.md](active.md))*

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
