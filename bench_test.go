package main

import (
	"testing"
)

func BenchmarkView(b *testing.B) {
	m := initialModelWithSeed(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

func BenchmarkStep(b *testing.B) {
	state := NewGameWithSeed(60, 30, 12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = state.Step(ActionMoveRight)
	}
}

func BenchmarkComputeFOV(b *testing.B) {
	state := NewGameWithSeed(60, 30, 12345)
	origin := state.Player.Pos
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state.Map.ComputeFOV(origin, FOVRadius)
	}
}
