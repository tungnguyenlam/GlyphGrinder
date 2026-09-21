# Active Backlog

The only file you should need to read to answer "what do I do next".
Update it continuously — as soon as a sub-task lands or the plan changes.

**Last updated:** 2026-09-21 (in-game help fits a 9-row terminal)

## Current milestone

**M4 — Ships to other people.** Turn the complete three-depth roguelike into a
project people can install, understand, and trust on ordinary terminals.

Serves `GOAL.md` directly: M3 completed the game loop; M4 proves the single
binary installs cleanly, degrades gracefully, stays responsive, and remains
green outside the maintainer's machine.

## Exact next action

**Complete the external release checks after these commits are pushed** — local
M4 engineering is green, including `--dump-frame`, `--help`, in-game `?`
controls (compact overlay that fits 9-row terminals), wait/hjkl, low-HP
danger HUD, wrapped sidebar logs, remembered items, diagonal FOV, free
wall/OOB bumps, failed `p`/`e`/`>` log lines, distance-faded torch lighting,
and a `GlyphGrinder` make-build name. Remaining acceptance work cannot run
in this environment:

- Confirm the `Verify` GitHub Actions workflow passes on the first pushed
  commit and that its only verification entry point is `./scripts/verify.sh`.
- From outside this checkout, run
  `go install github.com/tungnguyenlam/GlyphGrinder@latest` after the public
  branch contains these commits and confirm `GlyphGrinder` launches.
- Run `make run` in a real terminal, capture a representative gameplay frame,
  add the image under `docs/assets/`, and embed it near the top of the README.

Why next: every local M4 check is green. These last observations inherently
require the remote repository or a real interactive terminal.

If the next session cannot push or open a TTY, audit `GOAL.md` against tests
and current code for another high-confidence local batch; park speculative
ideas in `docs/backlog/parking-lot.md`. A 2026-09-21 audit parked GOAL.md's
door glyphs (the dungeon has none) and restored the hunger-clock line; title/
death screens, camera easing, and seed replay stay in the parking lot.

## Milestone plan

- [x] **M4.1** Verify local and versioned `go install` paths; document install,
      controls, gameplay, and terminal requirements — local `GOBIN` install
      passes in `verify.sh`; versioned `@latest` awaits push (blocked).
- [x] **M4.2** Exercise startup degradation for ASCII/no Nerd Font, truecolor
      capability, tiny terminals, and mid-run resize without a TTY.
- [x] **M4.3** Benchmark the render path and enforce an evidence-based frame
      budget that keeps normal play smooth and quiet.
- [x] **M4.4** Run `./scripts/verify.sh` in CI on pushes and pull requests —
      workflow `.github/workflows/verify.yml:1-27` exists with push/PR, cache,
      per-ref cancel, 10-min timeout, sole canonical script; remote run awaits
      push (blocked).
- [x] **M4.5** `--dump-frame` renders one production frame without a TTY;
      `verify.sh` exercises the installed binary and `go run . --dump-frame`
      and keeps `TestViewRendersFullGrid` as the pure-View guard.
- [x] **M4.6** `--help`/`-h` print usage without a TTY; unknown arguments
      exit 2 with usage on stderr; `verify.sh` checks `--help` on the
      isolated install.
- [x] **M4.7** Letter keys (`WASD`, `p`, `e`, `r`, `q`) are case-insensitive
      so Caps Lock does not silently drop input.
- [x] **M4.8** Interactive `Run` errors still print the Bubble Tea TTY
      message to stderr; a missing TTY also hints at `--help` / `--dump-frame`.
- [x] **M4.9** `.` / space wait a turn; `hjkl` move; `make build` emits
      `./GlyphGrinder` to match `go install`.
- [x] **M4.10** `?` toggles in-game controls without spending a turn; sidebar
      shows `? help` when closed. Hunger-clock parking-lot line restored.
- [x] **M4.11** HP in `1..=10` uses danger styling; alive players keep a bar
      sliver so 1–9 HP is not an empty bar.
- [x] **M4.12** Sidebar log lines wrap to `sidebarWidth`; stored log text is
      unchanged.
- [x] **M4.13** Explored-but-unseen potions and swords render as faint memory
      glyphs; unexplored items and unseen monsters stay hidden.
- [x] **M4.14** Bresenham FOV cannot peek through a closed diagonal wall
      corner; a single adjacent wall still lets the glance land.
- [x] **M4.15** Walking into a wall or off the map does not spend a turn;
      wait (`.`/space) and bump-attacks still do.
- [x] **M4.16** Failed `p`/`e`/`>` append a reason and stay free; `?` help
      lists arrows with WASD/hjkl.
- [x] **M4.17** Visible floor/wall/stairs/items fade across four torch bands
      toward memory colors (55% at the FOV edge); actors and memory stay full
      semantic colors. Colorless profiles stay unstyled.
- [x] **M4.18** Open `?` help omits log/goal, uses compact control lines, and
      clips to `windowHeight` so a 50x9 terminal stays 9 rows.

## Acceptance criteria for M4

- The README gives a working install command, complete controls, the run goal,
  terminal expectations, a genuine gameplay screenshot, and the canonical
  contributor workflow.
- Installation, startup fallback, tiny/mid-run resize behavior, and the render
  frame budget have automated coverage that does not require a TTY.
- CI runs the same canonical verification command used locally.
- `./scripts/verify.sh` passes.

## Blockers

- The workflow and versioned `go install ...@latest` cannot see unpushed local
  commits. The standing prompt forbids pushing or publishing.
- A genuine README screenshot needs `make run` in a real TTY; the agent
  environment cannot open `/dev/tty` (see the active TTY notice).

## Last verification

```
./scripts/verify.sh   —  2026-09-21  —  VERIFY OK (go1.27.0 darwin/arm64, /opt/homebrew/bin/go)
```

Tiny-terminal help batch (2026-09-21): open `?` omits the log and goal lines,
pairs `p`/`e` and `>`/`r`, puts `q` on the `? hide` row, and `clipLines` trims
the sidebar to `windowHeight`. A 50x9 window stays 9 rows with help open.
Proven by `TestInGameHelpFitsTinyTerminal` and
`TestResizeClipsViewportAndPreservesSidebar`.

Torch falloff batch (2026-09-21): visible floor, wall, stairs, and items use
four Euclidean torch bands from the animated player cell, blending up to 55%
toward memory colors at the FOV edge so the pool stays distinct from
exploration memory. Player, monsters, and unseen memory glyphs stay at full
semantic colors. Colorless terminals stay unstyled. Proven by
`TestTorchBandIncreasesWithDistance`, `TestLerpHexMixesChannels`, and
`TestViewTorchFallsOffWithDistance`.

Failed-action log batch (2026-09-21): off-stairs `>` logs `There are no stairs
here.`; empty/full-health `p` logs `You have no potions.` /
`You are already at full health.`; failed `e` logs
`You have no sword to equip.` / `You already wield the sword.` None of these
spend a monster turn. In-game help lists `arrows/WASD/hjkl`. Extra `>` after
victory stays a no-op. Proven by `TestDescendRequiresPlayerToStandOnStairs`,
`TestDescendKeyOffStairsLogsWithoutChangingDepth`,
`TestUsePotionIsFreeWhenUnavailableOrAtFullHealth`,
`TestEquipSwordIsFreeWhenUnavailableOrAlreadyEquipped`,
`TestFailedPotionAndEquipKeysLogWithoutSpendingATurn`, and
`TestInGameHelpTogglesWithoutSpendingATurn`.

Wall-bump batch (2026-09-21): `resolvePlayerMove` returns `acted=false` on
walls and map edges, so a missed corridor is free. Wait (`.`/space) and
bump-attacks still spend a turn. Monster cadence, lethal-hit, and TUI tests
that used a wall as a wait now call `ActionWait` / `.`. Proven by
`TestStepBlockedMovementDoesNotSpendATurn`, `TestWallBumpDoesNotSpendATurn`,
`TestStepWaitAdvancesMonsterTurnsWithoutMoving`, and
`TestStepBumpAttackDamagesMonster`.

Diagonal FOV batch (2026-09-21): a simultaneous Bresenham X/Y step is blocked
when both cardinal neighbors are walls, so a torch cannot see through a closed
corner. Cardinal occlusion and “the wall itself is visible” are unchanged. A
single adjacent wall still lets the glance land. Proven by
`TestVisibilityDoesNotPeekThroughDiagonalWallGap`,
`TestVisibilityDoesNotPeekThroughDistantDiagonalWallGap`,
`TestVisibilityCanGlancePastSingleWallCorner`, and
`TestVisibilityStopsBehindWalls`.

Remembered items batch (2026-09-21): explored-but-unseen potions/swords use
faint memory styles on `MemoryBG`; unexplored tiles stay blank; monsters stay
hidden when unseen. Proven by `TestViewRemembersSeenItems`.

Sidebar log wrap batch (2026-09-21): `wrapLogLine` word-wraps HUD log text to
32 cells (hard-break only for overlong tokens). `GameState.Log` stays intact.
Proven by `TestWrapLogLineFitsSidebarWidth` and
`TestSidebarLogLinesFitReservedWidth` (potion recover message).

Low-HP HUD batch (2026-09-21): health in `1..=10` styles the `HP` label, filled
bar, and `n/max` fraction with `styles.danger`; alive players keep at least one
`#` so 1–9 HP is not an empty bar. 0 HP stays empty plus GAME OVER. Proven by
`TestViewSignalsDangerInLastTenHitPoints`.

In-game help batch (2026-09-21): `?` toggles `model.showHelp` without calling
`Step`; other non-quit keys dismiss the overlay. Sidebar keeps a one-line
`? help` cue (8 rows with the default log, still fits the 50x9 resize test)
and expands to WASD/hjkl/wait/item/stairs/restart/quit. Hunger-clock parking
item restored after `ea61434` replaced it.

Wait/hjkl/binary-name batch (2026-09-21): `ActionWait` spends a monster turn
without moving; `.` and space map to it; `ActionNone` still freezes time.
Vim `hjkl` aliases the move keys. `make build` writes `./GlyphGrinder`;
`make clean` also removes leftover `glyphgrinder`. README controls updated.

TTY-error hint batch (2026-09-21): `formatRunError` keeps the original
`glyphgrinder: could not open a new TTY: ...` line and adds
`play needs a TTY; try --help or --dump-frame` only when the error text
contains `could not open a new TTY`. Non-TTY run errors stay a single line.
Proven by `TestFormatRunErrorHintsOnMissingTTY`. TTY notice example updated.

Caps Lock keys batch (2026-09-21): `Update` lowercases single-character keys
before matching, so `WASD`/`p`/`e`/`r`/`q` work with Caps Lock. Covered by
`TestMovementKeys` uppercase cases, `TestUppercaseLetterKeysMatchLowercaseActions`,
and `TestQuitKeys` `"Q"`. README notes letter keys are case-insensitive.

Headless CLI help batch (2026-09-21): `parseArgs` accepts `--help`/`-h`
(stdout usage, exit 0) and `--dump-frame`; unknown arguments print usage to
stderr and exit 2. `verify.sh` checks `--help` on the isolated install.
Proven: `--help` exit 0, `--not-a-flag` exit 2 with `unknown argument` on
stderr.

`--dump-frame` batch (2026-09-21): `main.dumpFrame` prints one `View()` frame
without `tea.NewProgram`; `parseArgs` recognizes `--dump-frame`; `verify.sh` checks
the isolated `go install` binary and `go run . --dump-frame` for `HP ` / `Depth 1`
after stripping CSI; HUD grep is ANSI-safe; README Develop and the TTY notice
document the flag. Parking-lot `--dump-frame` item removed.

Honest-failure batch (2026-08-23): `verify.sh` no longer claims
`all files gofmt-clean` when `gofmt` is missing or errors — it checks `go` and
`gofmt` exist up front and trusts `gofmt`'s exit status. Proven by running the
script with an empty toolchain PATH (fails naming `go`), a `go`-only PATH
(fails naming `gofmt`), and a syntax-error probe file (fails with the parse
error; previously passed vacuously). Startup errors now go to stderr with a
trailing newline.

Session audit (2026-08-23): game logic, value-state slice isolation, the
69-test suite, CI workflow (exec bit on `verify.sh` is `100755`, canonical
script is the only entry point), Makefile, README, and agent docs re-checked —
no further provable gaps. Parking lot pruned of two answered items (View
allocation profiling, map-generation fuzzing); TTY notice example refreshed
for the stderr error format.

M3.5 three-depth victory flow: pre-final descent, final-stair escape without
regeneration, frozen won state, visible goal/victory UI, and one-key fresh-run
restart pass at state and headless levels.

M4.1 public GitHub module identity, isolated executable `go install` check,
complete installation/gameplay/control/terminal README, and path fallout all
pass the canonical suite.

M4.2 injectable auto-detected Lip Gloss renderer, truecolor/colorless output
coverage, existing locale/`TERM` glyph fallback, 1x1 startup, and resize during
active animation all pass headlessly.

M4.3 baseline on Apple M4 (`go test -run '^$' -bench
'^BenchmarkViewProductionTrueColor$' -benchmem -count=5 .`): full 96x48 frames
97.9–101.0 us/op, 86,215 B/op, 893 allocs/op; clipped 80x24 frames 75.4–79.9
us/op, 47,679 B/op, 889 allocs/op. This is about 0.6% of the 16.7 ms animation
interval, so no optimization was justified; the averaged production-frame
budget guard passes 20 consecutive runs.

M4.4 read-only GitHub Actions workflow: push/PR triggers, cached `go.mod` Go
toolchain, per-ref cancellation, ten-minute timeout, and the canonical script
as its sole verification command. YAML parses locally and the same script is
green; remote execution awaits a push.

Audit 2026-08-26: re-ran `./scripts/verify.sh` with `/opt/homebrew/bin/go`
(go1.27) — VERIFY OK. Found stale `AGENTS.md` Current State and the missing
`--dump-frame` path; both are now landed.
