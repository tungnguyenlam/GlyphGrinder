# GlyphGrinder

A fast, turn-based Go roguelike that makes the terminal stop looking like a
terminal. Explore procedural dungeons through a moving pool of light, build a
small survival kit, and escape before depth-scaled monsters end the run.

A real starting frame, rendered by the game's own view code in a 62x20
ASCII-profile terminal (`@` you, `#` walls, `.` floor, `/` an iron sword):

```
                              HP [##########] 100/100
                              Depth 1
                              Goal reach depth 3
                              Potions 0 (p)
                              Weapon fists, 10 dmg
                              Log:
           .                  ? help
           #      #
         .##.####.#.
         ...........
         .........
        ......@.#
         .......#
         .......#
         /...###
          ..
           .
```

A UTF-8 terminal with truecolor and a Nerd Font shows the same scene with
warm torchlight falloff, colored glyphs, and recognizable silhouettes — no
configuration needed either way.

## Install

GlyphGrinder requires Go 1.25 or newer. Install the latest version directly
from GitHub:

```sh
go install github.com/tungnguyenlam/GlyphGrinder@latest
GlyphGrinder
```

Go places the command in `GOBIN` (or the default Go bin directory). That
directory must be on `PATH` for the second command to work.

To play from a source checkout instead:

```sh
git clone https://github.com/tungnguyenlam/GlyphGrinder.git
cd GlyphGrinder
make run
```

## Play

Survive three increasingly dangerous dungeon depths. Walk onto `>` and press
`>` to descend — that costs a turn, so standing next to a monster is risky.
The staircase on depth 3 is the exit and wins the run. Potions heal but
consume a turn, and the iron sword also costs a turn to equip.

| Key | Action |
| --- | --- |
| Arrow keys, `WASD`, or `hjkl` | Move; bump a monster to attack |
| `.` or space | Wait a turn |
| `?` | Toggle in-game controls |
| `p` | Drink a carried health potion |
| `e` | Equip a carried iron sword |
| `>` | Descend or escape while standing on stairs (costs a turn) |
| `r` | Start a fresh run after death or victory |
| `q` or `Ctrl+C` | Quit |

Letter keys (`WASD`, `hjkl`, `p`, `e`, `r`, `q`) are case-insensitive.

Goblins pursue every turn, high-health ogres act every other turn, and fragile
bats can take two actions. Enemy health and damage rise with dungeon depth.

## Terminal support

The game needs an interactive terminal (TTY) to play. `--help` and
`--dump-frame` run without one. A UTF-8 terminal with truecolor and a Nerd Font
gives the intended rich glyphs and palette, but no configuration is required:
non-UTF-8 locales and `TERM=dumb` automatically use one-cell ASCII glyphs. The
camera adapts to terminal resizing and remains usable in small windows.

If the command is not found after installation, add Go's bin directory to
`PATH`; `go env GOBIN` shows an explicitly configured location, while the
default is the `bin` directory under `go env GOPATH`.

## Develop

```sh
./scripts/verify.sh    # or: make verify — the one verification command
```

Runs the formatting check, `go vet`, build, an isolated `go install` plus
`--dump-frame` and `--help` of the installed binary, a degradation matrix that
renders that binary under `TERM=dumb`/POSIX (pure ASCII) and a UTF-8 locale
(rich glyphs), a matching `go run . --dump-frame` smoke render, the full test
suite under the race detector, and a `go mod tidy` check. No test requires a
TTY. The suite includes seeded full-run simulations that play whole games end
to end — a shortest-path bot that proves the dungeon stays escapable, and a
wandering bot that spends turns the way a distracted player does — and fails
if a run ever soft-locks, breaks the world's invariants, or the dungeon stops
being escapable.

[GOAL.md](GOAL.md) describes the product direction and
[docs/backlog/active.md](docs/backlog/active.md) tracks the current engineering
slice.

## Working on this with an AI agent

Start at [AGENTS.md](AGENTS.md). The cross-session workflow is
[docs/agent/continuity.md](docs/agent/continuity.md) and the doc index is
[docs/agent/index.md](docs/agent/index.md).

## License

GlyphGrinder is free software under the terms of the GNU General Public
License, version 3; see [LICENSE](LICENSE).
