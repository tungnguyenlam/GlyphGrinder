# Agent Doc Index

Table of contents only. Content lives in the linked files. **Anything added to
`notices/` or `decisions/` must be linked here in the same change** — likewise
when one is resolved or deleted.

## Start here

| Doc | What it is |
| --- | --- |
| [`../../AGENTS.md`](../../AGENTS.md) | Entry point: current state, startup routine, gotchas, handoff rules |
| [`../../GOAL.md`](../../GOAL.md) | North-star vision — what the finished game should feel like |
| [`../backlog/active.md`](../backlog/active.md) | Current milestone and the exact next action |
| [`continuity.md`](continuity.md) | The cross-session workflow: pickup, documentation, handoff |
| [`../../prompt/improve.md`](../../prompt/improve.md) | Standing prompt for an unattended autonomous session |

## Verification

One command, referenced by every doc here:

```sh
./scripts/verify.sh     # or: make verify
```

Runs gofmt check → `go vet` → `go build` → headless smoke test → full test
suite → `go mod tidy` check. Add new checks to that script, never to a doc.

## Backlog

| File | Purpose |
| --- | --- |
| [`../backlog/active.md`](../backlog/active.md) | The only place "what do I do next" needs to be answered from |
| [`../backlog/roadmap.md`](../backlog/roadmap.md) | Committed future work, in order |
| [`../backlog/parking-lot.md`](../backlog/parking-lot.md) | Uncommitted ideas — no promise they happen |
| [`../backlog/done.md`](../backlog/done.md) | Log of completed milestones |

## Notices

Small, dated, durable context with an expiry condition.

| Notice | Summary |
| --- | --- |
| [2026-07-26-tui-needs-a-tty](notices/2026-07-26-tui-needs-a-tty.md) | Resolved via `-dump-frame` flag. Binary can now dump frames headlessly. |
| [2026-07-26-styled-output-breaks-len](notices/2026-07-26-styled-output-breaks-len.md) | Lip Gloss escapes mean `len(line)` ≠ column count; strip ANSI before asserting. |
| [2026-07-26-go-mod-must-stay-tidy](notices/2026-07-26-go-mod-must-stay-tidy.md) | Direct deps were mislabelled `// indirect`; `verify.sh` now enforces tidiness. |

## Decisions (ADRs)

| ADR | Decision |
| --- | --- |
| [ADR-0001](../decisions/ADR-0001-single-package-until-it-hurts.md) | Stay in one root `package main` until a real boundary appears |
| [ADR-0002](../decisions/ADR-0002-flat-game-state.md) | Game state is a flat struct passed by value, not an ECS |
| [ADR-0003](../decisions/ADR-0003-view-is-pure.md) | `View` is pure rendering; all mutation happens in `Update` |
| [ADR-0004](../decisions/ADR-0004-no-new-dependencies.md) | Default to the standard library plus the existing Charm stack |

## Package map

| Package / file | Doc |
| --- | --- |
| `glyphgrinder` (root `package main`) — game model + Bubble Tea program | [`../../AGENTS.md`](../../AGENTS.md), [ADR-0001](../decisions/ADR-0001-single-package-until-it-hurts.md) |
| `internal/tuitest` — headless Bubble Tea driver used by all UI tests | [`../../internal/tuitest/AGENTS.md`](../../internal/tuitest/AGENTS.md) |
