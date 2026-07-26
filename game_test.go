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
