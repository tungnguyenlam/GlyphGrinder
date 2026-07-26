package main

import (
	"encoding/json"
	"os"
)

// ReplayData holds the initial seed and complete action sequence to reproduce a run.
type ReplayData struct {
	Seed    int64    `json:"seed"`
	Actions []Action `json:"actions"`
}

// ReplayRun plays back a sequence of actions from an initial seed and returns the final GameState.
func ReplayRun(width, height int, seed int64, actions []Action) GameState {
	st := NewGameWithSeed(width, height, seed)
	for _, act := range actions {
		st = st.Step(act)
	}
	return st
}

// SaveReplay writes the seed and action sequence to a JSON file on disk.
func SaveReplay(seed int64, actions []Action, path string) error {
	data := ReplayData{
		Seed:    seed,
		Actions: actions,
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// LoadReplay reads a seed and action sequence from a JSON file on disk.
func LoadReplay(path string) (ReplayData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ReplayData{}, err
	}
	var data ReplayData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return ReplayData{}, err
	}
	return data, nil
}
