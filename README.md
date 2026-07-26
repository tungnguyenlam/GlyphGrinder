# GlyphGrinder
A Go-based roguelike proving terminal games can be beautiful. Features fluid animations, Nerd Font icons, and rich truecolor graphics powered by Bubble Tea.

> Early days: today it's a walled room you can move an `@` around in. See
> [GOAL.md](GOAL.md) for where it's going and
> [docs/backlog/active.md](docs/backlog/active.md) for what's next.

## Play

```sh
make run          # needs a real terminal
```

Arrows or WASD to move, `q` to quit.

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
