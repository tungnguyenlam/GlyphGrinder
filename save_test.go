package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadGame(t *testing.T) {
	tempDir := t.TempDir()
	savePath := filepath.Join(tempDir, "test_save.json")

	st := NewGameWithSeed(20, 10, 12345)
	st.Player.Health = 77
	st.Player.Inventory = []Item{NewHealthPotion(1, Position{X: -1, Y: -1})}
	st.TurnCount = 42

	err := SaveGame(st, savePath)
	if err != nil {
		t.Fatalf("failed to save game: %v", err)
	}

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Fatalf("save file was not created at %s", savePath)
	}

	loaded, ok, err := LoadAndRemoveSaveGame(savePath)
	if err != nil {
		t.Fatalf("failed to load game: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true when loading save game")
	}

	if loaded.Seed != st.Seed {
		t.Errorf("loaded seed = %d, want %d", loaded.Seed, st.Seed)
	}
	if loaded.Player.Health != 77 {
		t.Errorf("loaded HP = %d, want 77", loaded.Player.Health)
	}
	if loaded.TurnCount != 42 {
		t.Errorf("loaded TurnCount = %d, want 42", loaded.TurnCount)
	}
	if len(loaded.Player.Inventory) != 1 {
		t.Errorf("loaded inventory length = %d, want 1", len(loaded.Player.Inventory))
	}

	// Verify save file was deleted after loading (permadeath requirement)
	if _, err := os.Stat(savePath); !os.IsNotExist(err) {
		t.Errorf("expected save file to be deleted after loading, but it still exists")
	}
}

func TestLoadNonExistentSaveGame(t *testing.T) {
	tempDir := t.TempDir()
	savePath := filepath.Join(tempDir, "non_existent.json")

	_, ok, err := LoadAndRemoveSaveGame(savePath)
	if err != nil {
		t.Fatalf("unexpected error loading non-existent save: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for non-existent save file")
	}
}
