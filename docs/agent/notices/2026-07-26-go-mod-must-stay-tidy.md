# `go.mod` must stay tidy — direct deps were mislabelled

**Status:** Active
**Scope:** `go.mod`, `go.sum`, `scripts/verify.sh`
**Related:** [ADR-0004](../../decisions/ADR-0004-no-new-dependencies.md)

## Why It Matters

Until 2026-07-26 every requirement in `go.mod` — including
`github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss`,
which `main.go` imports directly — was marked `// indirect`. That happens when
deps are added by hand or by `go mod download` without a `go mod tidy`, and it
is actively misleading: an agent reading `go.mod` to answer "what does this
project actually depend on?" sees a flat list of twenty transitive packages
with no signal about which two are load-bearing.

It also hides dependency growth. If everything is `// indirect`, a newly added
direct dependency looks exactly like a transitive one.

## Required Behavior

- Run `go mod tidy` after any change to imports. `./scripts/verify.sh` fails
  the build if `go.mod` would change, so an untidy file will not survive a
  handoff.
- The `require` block without `// indirect` is the honest list of direct
  dependencies. Keep it short (see ADR-0004) and treat any addition to it as a
  decision worth stating in the commit message.
- `verify.sh` restores `go.mod`/`go.sum` after its tidiness check; do not
  "simplify" that step into a bare `go mod tidy` that mutates the tree during
  verification.

## Revisit When

The tidiness check is enforced somewhere earlier (a pre-commit hook or CI job)
and has held for long enough that the mislabelling failure mode is no longer a
live risk.
