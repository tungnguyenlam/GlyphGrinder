# ADR-0004: Standard library plus the Charm stack, nothing else by default

**Date:** 2026-07-26
**Status:** Accepted

## Decision

GlyphGrinder's direct dependencies are `charmbracelet/bubbletea` and
`charmbracelet/lipgloss`. Everything else — map generation, pathfinding, field
of view, RNG, easing curves, save serialization, test assertions — is written
against the standard library.

Adding a direct dependency requires an ADR superseding this one, stating what
it replaces and why writing it ourselves is worse.

## Consequences

- `go install` stays a single fast build with a small, auditable dependency
  tree, which serves the "it just runs" part of `GOAL.md`.
- We write our own roguelike algorithms. For this genre they are small,
  well-documented, and tuning them by hand is part of making the game feel
  good — a general-purpose library would be adapted more than used.
- No assertion library: tests use `testing` with `got`/`want` comparisons. No
  `teatest` either — `internal/tuitest` is ~100 lines and drives the model
  directly, which is faster and less brittle than a pty-backed harness.
- Charm's own transitive dependencies (termenv, ansi, runewidth, ...) are
  accepted as the cost of Bubble Tea and are not counted against this rule.
  They must stay marked `// indirect` in `go.mod`; see
  [notice: go.mod must stay tidy](../agent/notices/2026-07-26-go-mod-must-stay-tidy.md).
- If a needed algorithm turns out to be genuinely hard to get right (unicode
  width handling, for instance), prefer a Charm-ecosystem package over a novel
  one before writing an ADR.
