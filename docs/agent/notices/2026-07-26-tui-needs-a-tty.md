# The game binary needs a TTY for interactive play

**Status:** Active — interactive mode still requires a TTY; `--help` and `--dump-frame` cover headless invocations
**Scope:** `main.go` (`main`), `scripts/verify.sh`, any attempt to "just run it and see"
**Related:** [`internal/tuitest/AGENTS.md`](../../../internal/tuitest/AGENTS.md), [ADR-0003](../../decisions/ADR-0003-view-is-pure.md)

## Why It Matters

`main()` starts the interactive program with `tea.WithAltScreen()`, which needs
a real terminal. In an agent session, CI, or any piped shell there is no
controlling TTY, so the interactive binary fails immediately:

```
$ printf 'q' | ./GlyphGrinder
glyphgrinder: could not open a new TTY: open /dev/tty: device not configured
glyphgrinder: play needs a TTY; try --help or --dump-frame
$ echo $?
1
```

(The original Bubble Tea message still goes to stderr with a trailing newline;
a second line points at `--help` / `--dump-frame`. The exit status is still 1.)

An agent that tries `go run .` to check a change will read that exit-1 as "my
change broke the game" and start debugging something that was never broken.

## Required Behavior

- Never use bare `go run .` / the built binary as an interactive verification
  step. It will fail for environmental reasons in every automated context.
- For a headless smoke test of the real binary, use the explicit flag:
  `go run . --dump-frame` or `GlyphGrinder --dump-frame` — it builds the
  initial model and prints one `View()` frame to stdout without calling
  `tea.NewProgram`. `GlyphGrinder --help` (or `-h`) prints usage and exits 0.
  Unknown arguments print usage to stderr and exit 2. `./scripts/verify.sh`
  exercises `--dump-frame` and `--help` on the isolated `go install` binary.
- Verify interactive UI behavior through `internal/tuitest`, which drives
  `Init`, `Update` and `View` in-process with no terminal involved. See
  `tui_test.go` for the pattern.
- If a change genuinely needs human eyes (colors, animation smoothness, glyph
  rendering in a specific font), stop and ask the user to run `make run`; state
  precisely what they should look for.

## Revisit When

The interactive program gains a true non-alt-screen interactive fallback, or
`--dump-frame` is removed. The current flag intentionally only renders one
static frame; it does not replace `make run` for visual verification.
