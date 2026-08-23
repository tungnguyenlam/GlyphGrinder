# ADR-0002: Flat game state, not an ECS

**Date:** 2026-07-26
**Status:** Accepted

## Decision

`GameState` is a plain struct holding a `GameMap`, the `Player` entity, a slice
of other `Entity` values, and a message `Log`. Entities are concrete structs
with fixed rule fields (`Pos`, `Health`, `MaxHealth`, `Damage`) — not component
bags, not interfaces, not an entity-component-system. Presentation choices such
as glyphs and colors live on the Bubble Tea model/rendering side instead of in
game entities.

The state travels by value: Bubble Tea's `Update` receives a `model` by value,
mutates its copy, and returns it. We keep it that way rather than switching to
pointer receivers.

## Consequences

- Adding a capability means adding a field to `Entity` or a slice to
  `GameState`. This is boring and obvious, which is the point — an agent with
  no prior context can find and change game state in one file.
- Every entity carries every field, including ones it doesn't use. Accepted:
  the entity count in a single dungeon level is small enough that the memory
  cost is irrelevant.
- Value semantics mean a mutation is only kept if the mutated copy is returned
  from `Update`. Forgetting `return m, ...` silently discards the change; this
  is the most likely class of bug in this codebase, so game-state changes need
  a test that asserts the change survived a round trip through `Update`.
- Slices inside the struct (`Map.Tiles`, `Entities`, `Log`) are *not* deep
  copied by assignment. Two copies of a `GameState` share their tile grid.
  Anything that needs an independent snapshot (undo, replay, AI lookahead) must
  clone explicitly.
- If per-entity behavior later needs real polymorphism, that is a new ADR
  superseding this one — not an incremental drift toward components.
