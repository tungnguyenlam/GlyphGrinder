# GlyphGrinder

A Go terminal roguelike built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

GlyphGrinder proves that terminal games can be fluid, responsive, and visually striking. Featuring procedurally generated dungeons, smooth continuous camera easing, dynamic line-of-sight shadowcasting, truecolor themes, and tactical turn-based combat.

```
HP: 100/100 | Depth: 1 | Inv: [Health Potion, Iron Dagger]
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓················▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓····󰋋···········▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓·······󰏗········▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓··········󰆧·····▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓·▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
You see a Health Potion here. (Press g or , to pick up)
```

---

## Features

- **Procedural Dungeon Generation**: Rooms connected by carved corridors across 5 dungeon levels with hazard tiles (Lava & Water).
- **Dynamic Field of View**: Recursive shadowcasting FOV algorithm (lit tiles torch-lit, explored tiles slate blue memory, unexplored hidden).
- **Camera Interpolation & Smooth Animations**: ~60 FPS tick-driven sub-tile camera interpolation easing to follow the player seamlessly without flicker.
- **Title Screen & Stat Cards**: Stylized ASCII title screen and detailed end-game run summaries (turns survived, monsters slain, damage dealt, dungeon seed).
- **Profile-Aware Visuals**: Curated TrueColor palette (`#00FF87`, `#FF55FF`, `#00E5FF`, `#FF3B30`, `#D32F2F`, `#FF4500` Lava, `#00BFFF` Water) with automatic Lip Gloss degradation for ANSI256 and 16-color ANSI terminals.
- **Nerd Font & ASCII Support**: Auto-detects Nerd Font unicode symbols (`󰋋` player, `󰆧` goblin, `󰌆` orc, `󰇄` troll, `󰓤` archer, `󰏗` potion, `󰓥` sword, `󰇮` amulet, `󰈸` lava, `󰌠` water) with clean `@`, `g`, `o`, `T`, `A`, `!`, `/`, `*`, `~` ASCII fallbacks.
- **Items, Inventory & Drop**: Pick up and drop items (`x`/`D`), manage inventory in status bar, drink health & regen potions, read fireball & teleport scrolls, and equip weapons.
- **Targeting UI & Magical Scrolls**: Line-of-sight Ranged Targeting mode (`t`) for precision magic scroll casting.
- **Mid-Run Save & Resume**: Auto-saves on quit (`q`), auto-resumes on launch, and deletes save file upon loading to strictly enforce permadeath.
- **Varied Monster AI & Ranged Combat**: Heavy-hitting Trolls, fast Goblins, ferocious Orcs, and Archers that snipe from up to 5 tiles away in line of sight.
- **Win Condition**: Descend to Depth 5, retrieve the legendary **Amulet of Yendor**, and claim victory.

---

## Controls

| Key | Action |
| --- | --- |
| `W` / `A` / `S` / `D` or Arrows | Move / Melee bump-attack |
| `g` or `,` | Pick up item underfoot into inventory |
| `x` or `D` | Drop last item from inventory |
| `h` | Drink health potion from inventory |
| `1` – `9` | Use / equip item from inventory slot 1–9 |
| `t` | Line-of-sight Ranged / Scroll targeting mode |
| `>` or `Enter` | Descend down stairs |
| `r` | Restart run (after Game Over or Victory) |
| `q` or `Ctrl+C` | Auto-save & quit game |

---

## Quick Start

### Installation via `go install`

```sh
go install github.com/tungnguyenlam/GlyphGrinder@latest
glyphgrinder
```

### Running from source

```sh
make run          # requires a terminal emulator with TTY
```

*Note:* GlyphGrinder works best in terminals with Nerd Font support and 24-bit TrueColor (e.g. Alacritty, Kitty, WezTerm, iTerm2, tmux). To force ASCII mode:

```sh
GLYPHGRINDER_ASCII=1 make run
```

---

## Architecture & Verification

GlyphGrinder follows a strict single-package root model (`package main`) backed by zero third-party dependencies outside the standard library and Charm stack:

- `game.go`: Pure, value-semantic game engine (`GameState`, `GameMap`, `Entity`, `Item`, `Step(Action)`).
- `main.go`: Pure Bubble Tea model (`Init`, `Update`, `View`) with decoupled rendering and camera interpolation.
- `colors.go` & `glyphs.go`: CompleteColor profiles and Nerd Font detection.
- `internal/tuitest`: Headless TUI test driver for deterministic unit and integration tests.

### Running Verification

```sh
./scripts/verify.sh     # or: make verify
```

Runs `gofmt` check, `go vet` static analysis, compilation, headless smoke frame render, the entire test suite, and `go mod tidy` check.

---

## License

MIT License. See [GOAL.md](GOAL.md) for project vision and design non-goals.
