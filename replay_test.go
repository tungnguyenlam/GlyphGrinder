package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestReplayRunDeterminism(t *testing.T) {
	seed := int64(888123)
	actions := []Action{
		ActionMoveUp,
		ActionMoveRight,
		ActionMoveDown,
		ActionMoveLeft,
		ActionPickup,
		ActionUseItem,
	}

	stOriginal := NewGameWithSeed(60, 30, seed)
	for _, act := range actions {
		stOriginal = stOriginal.Step(act)
	}

	stReplayed := ReplayRun(60, 30, seed, actions)

	if stOriginal.Player.Pos != stReplayed.Player.Pos {
		t.Errorf("player pos = %+v, want %+v", stReplayed.Player.Pos, stOriginal.Player.Pos)
	}
	if stOriginal.Player.Health != stReplayed.Player.Health {
		t.Errorf("player HP = %d, want %d", stReplayed.Player.Health, stOriginal.Player.Health)
	}
	if stOriginal.TurnCount != stReplayed.TurnCount {
		t.Errorf("TurnCount = %d, want %d", stReplayed.TurnCount, stOriginal.TurnCount)
	}
}

func TestSaveAndLoadReplayFile(t *testing.T) {
	tempDir := t.TempDir()
	replayPath := filepath.Join(tempDir, "test_replay.json")

	seed := int64(54321)
	actions := []Action{ActionMoveUp, ActionMoveRight, ActionDescend}

	err := SaveReplay(seed, actions, replayPath)
	if err != nil {
		t.Fatalf("failed to save replay: %v", err)
	}

	data, err := LoadReplay(replayPath)
	if err != nil {
		t.Fatalf("failed to load replay: %v", err)
	}

	if data.Seed != seed {
		t.Errorf("loaded seed = %d, want %d", data.Seed, seed)
	}
	if !reflect.DeepEqual(data.Actions, actions) {
		t.Errorf("loaded actions = %+v, want %+v", data.Actions, actions)
	}
}
