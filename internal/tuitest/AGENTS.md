# internal/tuitest — boundaries

Headless driver for Bubble Tea models. Used by `tui_test.go` to observe →
reason → act → synchronize: render the view, assert on it, send a key, render
again. Workflow rules live in the root `AGENTS.md`; this file is boundaries
only.

- **No game imports.** This package must never import game types or the root
  package. It knows `tea.Model` and nothing about GlyphGrinder — tests do the
  type assertion via `Driver.Model()`.
- **No new dependencies, ever.** This exists instead of `teatest`/`pty` harnesses
  (ADR-0004). If it needs a real terminal to work, the design is wrong.
- **Everything is synchronous.** `dispatch` runs commands inline and flattens
  `tea.BatchMsg`; `tea.Quit` is recorded on the driver, not executed. No
  goroutines, no timers, no sleeps — a flaky driver poisons every UI test.
- **Test-only.** Every exported function takes `*testing.T` and calls
  `t.Helper()`. Nothing here may be reachable from `main`.
- **Assert on plain text.** Rendered output carries ANSI escapes; callers must
  strip them before measuring (see the styled-output notice). Don't add
  stripping here — the driver returns what `View` actually produced.
