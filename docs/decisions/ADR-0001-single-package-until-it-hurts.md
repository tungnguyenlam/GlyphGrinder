# ADR-0001: Single root package until a boundary actually hurts

**Date:** 2026-07-26
**Status:** Accepted

## Decision

All game code lives in one `package main` at the repository root (`game.go`,
`main.go`, and their tests). We do not pre-split into `internal/game`,
`internal/render`, `internal/input`, etc.

A new package is justified only when at least one of these is true:

- It needs a boundary a compiler can enforce (e.g. "this may not import
  Bubble Tea").
- It is genuinely reusable outside the game loop.
- The root package has grown past roughly a thousand lines and has an obvious
  seam that several agents keep tripping over.

`internal/tuitest` exists because it meets the first test: it is test
infrastructure that must not be reachable from game logic.

## Consequences

- Adding a file is cheap; no import-cycle puzzles, no premature interfaces, no
  `Get`/`Set` plumbing to cross a boundary that only exists on paper.
- Everything is package-private by default, so refactoring internals never
  breaks a "public API" — there isn't one.
- The root package will get crowded. That is accepted and is the intended
  signal: when navigating it becomes the bottleneck, split then, along the seam
  the crowding revealed, and give the new package its own `AGENTS.md`.
- Agents must not "tidy" code into new packages as a standalone task. Splitting
  is a deliberate decision recorded in a new ADR that supersedes this one, not
  a cleanup.
