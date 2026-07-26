package main

import (
	"encoding/json"
	"os"
)

// DefaultSaveFilePath is the default local file path used for mid-run auto-saving.
const DefaultSaveFilePath = ".glyphgrinder_save.json"

// SaveGame serializes the current GameState to a JSON file on disk.
func SaveGame(st GameState, path string) error {
	if path == "" {
		path = DefaultSaveFilePath
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadGame reads a saved GameState from a JSON file on disk.
func LoadGame(path string) (GameState, error) {
	if path == "" {
		path = DefaultSaveFilePath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return GameState{}, err
	}
	var st GameState
	if err := json.Unmarshal(data, &st); err != nil {
		return GameState{}, err
	}
	return st, nil
}

// LoadAndRemoveSaveGame checks if a save file exists, loads the state, and deletes
// the save file to strictly uphold the roguelike permadeath contract (no save-scumming).
func LoadAndRemoveSaveGame(path string) (GameState, bool, error) {
	if path == "" {
		path = DefaultSaveFilePath
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return GameState{}, false, nil
	}
	st, err := LoadGame(path)
	if err != nil {
		_ = os.Remove(path)
		return GameState{}, false, err
	}
	_ = os.Remove(path)
	return st, true, nil
}
