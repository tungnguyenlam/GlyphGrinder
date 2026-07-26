package main

import "testing"

func TestNewGameBuildsWalledRoom(t *testing.T) {
	g := NewGame(20, 10)

	if got, want := len(g.Map.Tiles), 10; got != want {
		t.Fatalf("rows = %d, want %d", got, want)
	}
	for y, row := range g.Map.Tiles {
		if got, want := len(row), 20; got != want {
			t.Fatalf("row %d width = %d, want %d", y, got, want)
		}
	}

	for x := 0; x < g.Map.Width; x++ {
		if g.Map.Tiles[0][x] != TileWall || g.Map.Tiles[g.Map.Height-1][x] != TileWall {
			t.Errorf("column %d: expected wall on top and bottom border", x)
		}
	}
	for y := 0; y < g.Map.Height; y++ {
		if g.Map.Tiles[y][0] != TileWall || g.Map.Tiles[y][g.Map.Width-1] != TileWall {
			t.Errorf("row %d: expected wall on left and right border", y)
		}
	}

	// Player position in generated dungeon must be floor
	if g.Map.Tiles[g.Player.Pos.Y][g.Player.Pos.X] != TileFloor {
		t.Errorf("player spawn tile at %+v should be floor", g.Player.Pos)
	}
}

func TestNewGamePlacesPlayerOnFloor(t *testing.T) {
	g := NewGame(20, 10)

	p := g.Player
	if !p.IsPlayer {
		t.Error("player entity should have IsPlayer set")
	}
	if g.Map.Tiles[p.Pos.Y][p.Pos.X] != TileFloor {
		t.Errorf("player position %+v must be a floor tile", p.Pos)
	}
	if p.Health != p.MaxHealth {
		t.Errorf("player starts at %d/%d health, want full", p.Health, p.MaxHealth)
	}
}

func TestMonsterHelpers(t *testing.T) {
	gob := NewGoblin(1, Position{X: 2, Y: 3})
	if gob.ID != 1 || gob.Pos.X != 2 || gob.Pos.Y != 3 {
		t.Errorf("unexpected goblin pos/id: %+v", gob)
	}
	if gob.IsPlayer {
		t.Error("goblin should not be player")
	}
	if gob.Rune != "g" || gob.Health != 10 || gob.Damage != 3 {
		t.Errorf("unexpected goblin stats: %+v", gob)
	}

	orc := NewOrc(2, Position{X: 4, Y: 5})
	if orc.ID != 2 || orc.Pos.X != 4 || orc.Pos.Y != 5 {
		t.Errorf("unexpected orc pos/id: %+v", orc)
	}
	if orc.IsPlayer {
		t.Error("orc should not be player")
	}
	if orc.Rune != "o" || orc.Health != 20 || orc.Damage != 6 {
		t.Errorf("unexpected orc stats: %+v", orc)
	}
}

func TestNewGamePopulatesEntities(t *testing.T) {
	g := NewGameWithSeed(20, 10, 12345)

	if len(g.Entities) == 0 {
		t.Fatal("expected g.Entities to be populated with monsters")
	}

	seenIDs := make(map[int]bool)
	seenPos := make(map[Position]bool)

	seenIDs[g.Player.ID] = true
	seenPos[g.Player.Pos] = true

	for i, e := range g.Entities {
		if e.IsPlayer {
			t.Errorf("entity %d in Entities slice has IsPlayer=true", i)
		}
		if e.Health <= 0 || e.MaxHealth <= 0 || e.Damage <= 0 {
			t.Errorf("entity %d has invalid stats: %+v", i, e)
		}
		if e.Rune != "g" && e.Rune != "o" {
			t.Errorf("entity %d has unexpected rune %q", i, e.Rune)
		}
		if g.Map.Tiles[e.Pos.Y][e.Pos.X] != TileFloor {
			t.Errorf("entity %d at %+v is not on a floor tile", i, e.Pos)
		}
		if seenIDs[e.ID] {
			t.Errorf("duplicate entity ID %d", e.ID)
		}
		seenIDs[e.ID] = true
		if seenPos[e.Pos] {
			t.Errorf("duplicate entity position %+v", e.Pos)
		}
		seenPos[e.Pos] = true
	}
}
