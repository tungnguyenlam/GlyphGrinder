# CLAUDE.md

Pointer file. The source of truth is:

- [`AGENTS.md`](AGENTS.md) — read this first (state, startup routine, gotchas).
- [`GOAL.md`](GOAL.md) — the north star every task is checked against.
- [`prompt/improve.md`](prompt/improve.md) — the standing prompt for an
  unattended work session.

Claude-specific notes only:

- Verification is `./scripts/verify.sh` (or `make verify`). Run it after each
  change, not just at the end.
- Don't try to run the game to check your work — it needs a TTY and will exit 1.
  Use `internal/tuitest`.
- Prefer `rg` over reading files wholesale; the root package is small but
  growing on purpose (ADR-0001).
