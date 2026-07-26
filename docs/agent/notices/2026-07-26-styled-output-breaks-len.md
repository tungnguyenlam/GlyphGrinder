# Styled output means `len(line)` is not the column count

**Status:** Active
**Scope:** `main.go` (`View`), `tui_test.go`, anything measuring or slicing rendered output
**Related:** [ADR-0003](../../decisions/ADR-0003-view-is-pure.md), [notice: TUI needs a TTY](2026-07-26-tui-needs-a-tty.md)

## Why It Matters

`View` renders every cell through Lip Gloss (`playerStyle.Render`,
`wallStyle.Render`, ...), so each visible character arrives wrapped in ANSI
escape sequences. A 20-column row is far more than 20 bytes and more than 20
runes. Naive assertions and layout math silently do the wrong thing:

- `len(line) == 20` fails even when the row is correct.
- `strings.Index(line, "@")` returns a byte offset, not a column.
- Truncating or padding a styled string can cut an escape sequence in half and
  bleed color into the rest of the frame.

This also applies to alignment: whether the terminal has truecolor changes how
much escape text Lip Gloss emits, so byte lengths are not even stable across
environments.

## Required Behavior

- Strip ANSI escapes before asserting on positions or widths in tests — use the
  `stripANSI` helper in `tui_test.go`, and count with `[]rune`, not `len`.
- For display width of multi-byte or wide glyphs (Nerd Font icons, CJK), use
  `lipgloss.Width` rather than counting runes.
- Build layout from unstyled strings and apply styling last, so widths are
  computed on plain text.
- Never slice or pad an already-styled string.

## Revisit When

Rendering moves to a cell-buffer model (a `[][]Cell` composed then styled once
at the end) where views are measured before any escape codes exist.
