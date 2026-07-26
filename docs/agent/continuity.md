# Continuity: Working Across Sessions

Agent sessions do not share memory. This file is the workflow that lets the
next session start where this one stopped, without re-discovering the codebase
or re-arguing settled decisions.

## New session pickup

Do these in order, before touching code:

1. Read [`../../AGENTS.md`](../../AGENTS.md) — current state, gotchas, handoff
   rules.
2. Read [`../backlog/active.md`](../backlog/active.md) — milestone, exact next
   action, blockers, last verification result.
3. Read [`index.md`](index.md) — scan the notice and ADR tables; open any that
   touch the area you are about to change.
4. Read [`../../GOAL.md`](../../GOAL.md) if you are planning work rather than
   executing an already-specified next action.
5. Scope with search, not with reading: `rg -n '<symbol>'`,
   `rg -n 'func ' game.go main.go`. Read only the files you will change plus
   their direct callers.
6. Read the nearest subtree `AGENTS.md` before editing anything under it.
7. Run `./scripts/verify.sh` **before** your first edit, so you know whether a
   later failure is yours.

## Continuous documentation

Update docs as work happens, not at the end. The session may be cut off at any
point; whatever is written down is what survives.

- A sub-task lands → tick it in `active.md` and set the new next action, in the
  same commit as the code.
- The plan changes → rewrite the plan in `active.md` immediately, including
  *why* it changed. A stale plan is worse than no plan.
- You discover a trap → write a notice (below) as soon as you understand it,
  not after you have worked around it.
- You make a call about structure → write an ADR before building on it.
- A gotcha in `AGENTS.md` stops being true → delete it in the change that fixed
  it.

## Notices

A notice is for context that is **too important to lose but too small for an
ADR**: a trap, a subtle invariant, a temporary compatibility rule, a
cross-package contract, an environment limitation.

Write one when the answer to "would the next agent waste an hour rediscovering
this?" is yes. Do *not* write one for things the code already says clearly, or
for a one-off bug that is now fixed.

- Location: `docs/agent/notices/`
- Filename: `YYYY-MM-DD-short-topic.md` (date it was written, kebab-case topic)
- Link it from [`index.md`](index.md) **in the same change** — adding,
  resolving or deleting.

Template — use exactly these sections:

```markdown
# <Title>

**Status:** Active | Resolved (YYYY-MM-DD) | Superseded by <link>
**Scope:** <files, packages or workflows this applies to>
**Related:** <links to ADRs, other notices, backlog items — or "None">

## Why It Matters
<what goes wrong if an agent doesn't know this — concrete>

## Required Behavior
<what to do instead — specific, checkable>

## Revisit When
<the condition that makes this notice obsolete>
```

When a notice's Revisit-When condition is met: set Status to Resolved with the
date and delete the body if it is now misleading, or delete the file outright.
Either way, update `index.md`.

## Backlog

`docs/backlog/active.md` is the only file that should ever need to be read to
answer "what do I do next". It must always answer, explicitly:

- **Current milestone** — one sentence, plus how it serves `GOAL.md`.
- **Exact next action** — a specific, executable step ("factor the four
  movement branches in `main.go:27-47` into `tryMove(dx, dy int)`"), not a
  theme ("improve movement").
- **Acceptance criteria** — how the next agent knows the milestone is done.
- **Blockers** — anything needing a human, an external decision, or missing
  information.
- **Last verification** — the exact command run, when, and its result.

Keep vague reminders out. Anything vague either becomes a concrete next action
or gets demoted to `parking-lot.md`.

The other three files:

- `roadmap.md` — committed future work, ordered. Milestones move from here into
  `active.md`.
- `parking-lot.md` — ideas we are *not* committing to. Free to grow; nothing
  here is a promise. Deleting from here is cheap and encouraged.
- `done.md` — append-only log of completed milestones, newest first, with the
  date and what changed. This is the project's memory of what was already
  tried.

## Subtree `AGENTS.md`

A package deserves its own file when it has **local rules a general agent would
otherwise violate**: an architectural boundary ("nothing here may import
Bubble Tea"), a test requirement, an invariant callers depend on.

It does not deserve one when the file would just restate the root workflow,
describe what the code obviously does, or list every function. That is noise
and it goes stale.

Format: 3–6 bullets, each a hard boundary or invariant. No workflow, no
verification instructions, no duplication of the root file.

## Handoff checklist

Every single time work pauses — end of session, context running out, blocked,
or handing to a human:

1. `docs/backlog/active.md` updated: milestone, next action, blockers.
2. Any new trap written up as a notice and linked from `index.md`.
3. Any structural decision written up as an ADR and linked from `index.md`.
4. `AGENTS.md` gotchas added or pruned to match reality.
5. Run `./scripts/verify.sh`.
6. Record that run in `active.md` under **Last verification** — the command,
   the date, and the actual result. If it fails, say so and say exactly what
   failed; never hand off with an unrecorded or fabricated result.
7. If the tree is green and the work is coherent, commit it.
