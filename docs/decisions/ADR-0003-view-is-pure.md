# ADR-0003: `View` renders, `Update` decides

**Date:** 2026-07-26
**Status:** Accepted

## Decision

`model.View()` is pure: it reads model state and returns a string. It must not
mutate state, advance turns, roll dice, read the clock, or perform I/O.

Everything that changes the world happens in `model.Update` in response to a
message. Time-based effects (animation frames, monster turns on a timer) arrive
as messages from a `tea.Cmd`, not from `View` noticing that time has passed.

## Consequences

- `View` can be called any number of times, in any order, with no side effects
  — which is what makes the headless driver in `internal/tuitest` a valid way
  to observe game state, and what makes tests deterministic.
- Animation must be modelled as state ("this entity is 0.4 of the way through a
  step") plus a tick message, rather than as a `View` that draws something
  different each call. This is more work up front and is the price of testable
  rendering; it also means an animation can be asserted on frame by frame.
- Randomness lives in `Update`. A test can therefore control outcomes by
  controlling the seed in the model, without threading a RNG through rendering.
- `View` may allocate freely (it builds strings), but it must not be the place
  where a bug's *effect* lives — if the display is wrong, the state was wrong
  or the rendering was wrong, and those two are never entangled.
