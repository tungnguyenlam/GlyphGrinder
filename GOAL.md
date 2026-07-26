# GOAL

GlyphGrinder is a roguelike that makes someone say "*that's* a terminal?"

## What it should feel like

You launch it from a shell and the terminal stops looking like a terminal. A
dungeon fades in. The `@` you control moves with weight — steps ease, the
camera drifts to follow, a torch pool of warm light slides across stone while
the corridor behind you falls back into blue-grey memory. Nothing snaps; things
*move*. Nerd Font glyphs give every monster, door and potion a silhouette you
recognise at a glance, and truecolor gives the dungeon depth instead of the
sixteen-color flatness people expect from a TUI.

It is fast and quiet. Keys respond instantly — no input lag, no flicker, no
tearing when the frame redraws. It runs at a smooth frame rate in a normal
terminal on a laptop, on battery, without a fan spinning up.

It is a real roguelike underneath the polish, not a tech demo: procedurally
generated levels, permadeath, turn-based combat where a mistake costs you the
run, items that change how you play. A run should be tense in the last ten
hit points and re-startable in one keypress.

And it just runs. `go install` or a single binary, no config file, no font
wrangling beyond "install a Nerd Font", graceful degradation to plainer glyphs
and colors when the terminal can't do better.

## Checking work against this

Every task should answer: **does this move us toward or away from GOAL.md?**

Toward: it makes the game more playable, more beautiful, more responsive, or
easier to run.

Away: it adds configuration surface, adds a dependency the user has to install,
costs frame time without a visual payoff, or builds engine generality no
current feature needs.

## Explicit non-goals

- Not a graphical game. No SDL, no web frontend, no mouse-first UI.
- Not an engine or a framework for other people's roguelikes.
- No networking, no multiplayer, no cloud saves, no telemetry.
- No plugin/scripting system.
