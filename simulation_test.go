package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// The unit tests drive fixed seeds through individual rules. This file plays
// whole runs instead: a deterministic beeline bot walks production-size
// dungeons from the starting floor to the final stairs across many seeds,
// checking the structural invariants a player would rely on after every turn.
// It is the shipping canary that the engine never panics, never soft-locks,
// and still produces at least one escapable run.

const (
	simulationSeedCount  = 200
	simulationMaxSteps   = 1500
	simulationPotionAtHP = 35
)

type runOutcome int

const (
	runDied runOutcome = iota
	runWon
)

// TestSimulatedRunsTerminateWithInvariants plays seeded full runs and asserts
// that every run ends in death or victory inside the step bound, that the
// world stays structurally sound after every turn, and that at least one of
// the fixed seeds wins — a winnability canary. It is deterministic: fixed
// seeds, fixed bot, so a failure means the engine changed, not that luck did.
func TestSimulatedRunsTerminateWithInvariants(t *testing.T) {
	wins := 0
	for seed := int64(1); seed <= simulationSeedCount; seed++ {
		seed := seed
		t.Run("seed_"+strconv.FormatInt(seed, 10), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			g := NewGame(productionDungeonWidth, productionDungeonHeight, rng)
			for step := 0; step < simulationMaxSteps; step++ {
				checkRunInvariants(t, g, seed, step)
				if g.GameOver || g.Won {
					if g.Won {
						wins++
					}
					return
				}
				next := playSimulationTurn(t, g, rng, seed, step)
				// MonsterTurns is per-depth (ogre cadence restarts on each
				// floor), so it may only reset together with a descent.
				if next.MonsterTurns < g.MonsterTurns && next.Depth == g.Depth {
					t.Fatalf("seed %d step %d: MonsterTurns went backwards %d -> %d", seed, step, g.MonsterTurns, next.MonsterTurns)
				}
				g = next
			}
			t.Fatalf("seed %d: run did not terminate (win or death) within %d steps at depth %d", seed, simulationMaxSteps, g.Depth)
		})
	}
	if wins == 0 {
		t.Fatalf("no seeded run reached victory; the game is no longer escapable for these %d seeds", simulationSeedCount)
	}
	t.Logf("%d of %d seeded runs escaped the dungeon", wins, simulationSeedCount)
}

// playSimulationTurn is the deterministic bot: drink when hurt, equip the
// sword when carried, use the stairs when standing on them, otherwise walk one
// cardinal step along the shortest path to the stairs (bump-attacking any
// monster on that step). Every turn spends time or is terminal, so runs always
// end: each bump kills a monster in finitely many hits, and the path length
// shrinks on every unblocked step.
func playSimulationTurn(t *testing.T, g GameState, rng *rand.Rand, seed int64, step int) GameState {
	t.Helper()
	if g.Player.Health <= simulationPotionAtHP && g.Potions > 0 {
		return g.Step(ActionUsePotion)
	}
	if g.HasSword && !g.Equipped {
		return g.Step(ActionEquipWeapon)
	}
	if g.Map.Tiles[g.Player.Pos.Y][g.Player.Pos.X] == TileStairs {
		return g.useStairs(rng)
	}
	stairs := tilePositions(g.Map, TileStairs)
	if len(stairs) != 1 {
		t.Fatalf("seed %d step %d: stair count = %d, want 1", seed, step, len(stairs))
	}
	next, ok := nextCardinalStep(g.Map, g.Player.Pos, stairs[0])
	if !ok {
		t.Fatalf("seed %d step %d: no path from %+v to stairs %+v", seed, step, g.Player.Pos, stairs[0])
	}
	delta := Position{X: next.X - g.Player.Pos.X, Y: next.Y - g.Player.Pos.Y}
	return g.Step(actionForDelta(delta))
}

// nextCardinalStep BFS-es over non-wall cells from `from` and returns the
// first cardinal step toward `to`, mirroring how the player and monsters
// actually move.
func nextCardinalStep(m GameMap, from, to Position) (Position, bool) {
	if from == to {
		return Position{}, false
	}
	parent := make(map[Position]Position, m.Width*m.Height)
	parent[from] = from
	queue := []Position{from}
	directions := [...]Position{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, direction := range directions {
			next := Position{X: current.X + direction.X, Y: current.Y + direction.Y}
			if next.X < 0 || next.X >= m.Width || next.Y < 0 || next.Y >= m.Height {
				continue
			}
			if m.Tiles[next.Y][next.X] == TileWall {
				continue
			}
			if _, seen := parent[next]; seen {
				continue
			}
			parent[next] = current
			if next == to {
				step := next
				for parent[step] != from {
					step = parent[step]
				}
				return step, true
			}
			queue = append(queue, next)
		}
	}
	return Position{}, false
}

func actionForDelta(delta Position) Action {
	switch {
	case delta.Y < 0:
		return ActionMoveUp
	case delta.Y > 0:
		return ActionMoveDown
	case delta.X < 0:
		return ActionMoveLeft
	default:
		return ActionMoveRight
	}
}

// checkRunInvariants asserts the structural promises every game rule keeps:
// actors and items live on in-bounds non-wall tiles, never overlap, carry
// sane stats, and the terminal states mean what they say.
func checkRunInvariants(t *testing.T, g GameState, seed int64, step int) {
	t.Helper()
	fail := func(format string, args ...any) {
		t.Helper()
		t.Fatalf("seed %d step %d: %s", seed, step, fmt.Sprintf(format, args...))
	}

	inBounds := func(p Position) bool {
		return p.X >= 0 && p.X < g.Map.Width && p.Y >= 0 && p.Y < g.Map.Height
	}
	assertOpenTile := func(what string, p Position) {
		if !inBounds(p) {
			fail("%s out of bounds at %+v", what, p)
		}
		if g.Map.Tiles[p.Y][p.X] == TileWall {
			fail("%s inside a wall at %+v", what, p)
		}
	}

	assertOpenTile("player", g.Player.Pos)
	if g.Player.Health < 0 || g.Player.Health > g.Player.MaxHealth {
		fail("player health %d outside 0..%d", g.Player.Health, g.Player.MaxHealth)
	}
	if g.Depth < 1 || g.Depth > finalDepth {
		fail("depth %d outside 1..%d", g.Depth, finalDepth)
	}
	if g.Potions < 0 {
		fail("negative potions %d", g.Potions)
	}
	if g.GameOver && g.Player.Health != 0 {
		fail("game over with %d health", g.Player.Health)
	}
	if g.Won && (g.Depth != finalDepth || g.GameOver) {
		fail("won state at depth %d with game over %v", g.Depth, g.GameOver)
	}

	occupied := map[Position]bool{g.Player.Pos: true}
	ids := make(map[int]bool, len(g.Entities))
	for _, monster := range g.Entities {
		assertOpenTile("monster", monster.Pos)
		if occupied[monster.Pos] {
			fail("monster %d overlaps an actor at %+v", monster.ID, monster.Pos)
		}
		occupied[monster.Pos] = true
		if ids[monster.ID] {
			fail("duplicate monster ID %d", monster.ID)
		}
		ids[monster.ID] = true
		if monster.Health <= 0 || monster.Health > monster.MaxHealth || monster.Damage <= 0 {
			fail("monster %d has invalid stats: %+v", monster.ID, monster)
		}
	}

	items := make(map[Position]bool, len(g.Items))
	for _, item := range g.Items {
		assertOpenTile("item", item.Pos)
		if items[item.Pos] {
			fail("two items share %+v", item.Pos)
		}
		items[item.Pos] = true
	}
	if len(g.Items) > 3 {
		fail("%d items on one floor, want at most 3", len(g.Items))
	}
}
