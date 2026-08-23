package main

import (
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var benchmarkViewSink string

func TestProductionViewMeetsAnimationFrameBudget(t *testing.T) {
	const (
		trueColorProfile = 0 // termenv.TrueColor, exposed through Lip Gloss.
		warmupFrames     = 10
		measuredFrames   = 100
	)
	renderer := lipgloss.NewRenderer(io.Discard)
	renderer.SetColorProfile(trueColorProfile)
	m := model{
		state:    NewGame(productionDungeonWidth, productionDungeonHeight, rand.New(rand.NewSource(tuiTestSeed))),
		renderer: renderer,
		glyphs:   richGlyphs,
	}

	for range warmupFrames {
		benchmarkViewSink = m.View()
	}
	started := time.Now()
	for range measuredFrames {
		benchmarkViewSink = m.View()
	}
	average := time.Since(started) / measuredFrames
	if average >= animationFrameDuration {
		t.Errorf("average production render = %s, want under animation frame budget %s", average, animationFrameDuration)
	}
}

func BenchmarkViewProductionTrueColor(b *testing.B) {
	const trueColorProfile = 0 // termenv.TrueColor, exposed through Lip Gloss.
	renderer := lipgloss.NewRenderer(io.Discard)
	renderer.SetColorProfile(trueColorProfile)
	state := NewGame(productionDungeonWidth, productionDungeonHeight, rand.New(rand.NewSource(tuiTestSeed)))

	benchmarks := []struct {
		name   string
		width  int
		height int
	}{
		{name: "full_map"},
		{name: "clipped_80x24", width: 80, height: 24},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			m := model{
				state:        state,
				renderer:     renderer,
				glyphs:       richGlyphs,
				windowWidth:  benchmark.width,
				windowHeight: benchmark.height,
			}
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				benchmarkViewSink = m.View()
			}
		})
	}
}
