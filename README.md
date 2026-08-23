# GlyphGrinder

A Go terminal roguelike proving terminal games can be beautiful. Explore a
large procedural dungeon through torch-like field of view, fight pursuing
monsters, and try to survive with truecolor depth and a player-following camera.

It is still early; [GOAL.md](GOAL.md) describes the intended destination and
[docs/backlog/active.md](docs/backlog/active.md) tracks the next slice.

## Play

```sh
make run          # needs a real terminal
```

Arrows or WASD move, `r` restarts after death, and `q` quits. UTF-8 terminals
use richer one-cell glyphs (best with a Nerd Font); non-UTF-8 and `TERM=dumb`
environments fall back automatically to `@`, `g`, `.`, and `#`.

## Develop

```sh
./scripts/verify.sh    # or: make verify — the one verification command
```

Runs the gofmt check, `go vet`, build, a headless smoke render, the test suite,
and a `go mod tidy` check.

## Working on this with an AI agent

Start at [AGENTS.md](AGENTS.md). The cross-session workflow is
[docs/agent/continuity.md](docs/agent/continuity.md) and the doc index is
[docs/agent/index.md](docs/agent/index.md).
