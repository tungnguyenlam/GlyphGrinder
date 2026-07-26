# The game binary cannot run headlessly

**Status:** Active
**Scope:** `main.go` (`main`), `scripts/verify.sh`, any attempt to "just run it and see"
**Related:** [`internal/tuitest/AGENTS.md`](../../../internal/tuitest/AGENTS.md), [ADR-0003](../../decisions/ADR-0003-view-is-pure.md)

## Why It Matters

`main()` starts the program with `tea.WithAltScreen()`, which needs a real
terminal. In an agent session, CI, or any piped shell there is no controlling
TTY, so the binary fails immediately:

```
$ printf 'q' | ./glyphgrinder
Error: could not open a new TTY: open /dev/tty: device not configured
$ echo $?
1
```

An agent that tries `go run .` to check a change will read that exit-1 as "my
change broke the game" and start debugging something that was never broken.

## Required Behavior

- Never use `go run .` / the built binary as a verification step. It will fail
  for environmental reasons in every automated context.
- Verify UI behavior through `internal/tuitest`, which drives `Init`, `Update`
  and `View` in-process with no terminal involved. See `tui_test.go` for the
  pattern.
- `./scripts/verify.sh` uses the headless render test as its smoke step for
  exactly this reason — keep it that way.
- If a change genuinely needs human eyes (colors, animation smoothness, glyph
  rendering in a specific font), stop and ask the user to run `make run`; state
  precisely what they should look for.

## Revisit When

The program grows a headless or non-alt-screen mode (e.g. a `--dump-frame`
flag, or a test entry point that builds the model without `tea.NewProgram`)
that lets an automated run render a frame end-to-end.
