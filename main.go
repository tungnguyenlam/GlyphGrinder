package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state GameState
	rng   *rand.Rand
}

func initialModel(rng *rand.Rand) model {
	return model{state: NewGame(20, 10, rng), rng: rng}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		action := ActionNone
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			if m.state.GameOver {
				m.state = NewGame(m.state.Map.Width, m.state.Map.Height, m.rng)
			}
			return m, nil
		case "up", "w":
			action = ActionMoveUp
		case "down", "s":
			action = ActionMoveDown
		case "left", "a":
			action = ActionMoveLeft
		case "right", "d":
			action = ActionMoveRight
		}
		m.state = m.state.Step(action)
	}
	return m, nil
}

func (m model) View() string {
	if m.state.Map.Width == 0 || m.state.Map.Height == 0 {
		return ""
	}

	playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.state.Player.Color)).Bold(true)
	player := playerStyle.Render(m.state.Player.Rune)
	floorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	floor := floorStyle.Render(".")
	wallStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	wall := wallStyle.Render("#")
	memoryFloorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#222222"))
	memoryFloor := memoryFloorStyle.Render(".")
	memoryWallStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	memoryWall := memoryWallStyle.Render("#")
	entityGlyphs := make(map[Position]string, len(m.state.Entities))
	for _, entity := range m.state.Entities {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(entity.Color)).Bold(true)
		entityGlyphs[entity.Pos] = style.Render(entity.Rune)
	}

	var sb strings.Builder
	for y := 0; y < m.state.Map.Height; y++ {
		for x := 0; x < m.state.Map.Width; x++ {
			visible := m.state.Map.Visible[y][x]
			explored := m.state.Map.Explored[y][x]
			if !explored {
				sb.WriteByte(' ')
			} else if x == m.state.Player.Pos.X && y == m.state.Player.Pos.Y {
				sb.WriteString(player)
			} else if entity, ok := entityGlyphs[Position{X: x, Y: y}]; visible && ok {
				sb.WriteString(entity)
			} else if !visible && m.state.Map.Tiles[y][x] == TileWall {
				sb.WriteString(memoryWall)
			} else if !visible {
				sb.WriteString(memoryFloor)
			} else {
				// Render floor tiles or walls based on the map state
				if m.state.Map.Tiles[y][x] == TileWall {
					sb.WriteString(wall)
				} else {
					sb.WriteString(floor)
				}
			}
		}
		if y < m.state.Map.Height-1 {
			sb.WriteString("\n")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, sb.String(), "  ", renderSidebar(m.state))
}

func renderSidebar(state GameState) string {
	const (
		healthBarWidth  = 10
		visibleLogLines = 4
	)
	filled := 0
	if state.Player.MaxHealth > 0 {
		filled = state.Player.Health * healthBarWidth / state.Player.MaxHealth
	}
	filled = max(0, min(healthBarWidth, filled))

	var sb strings.Builder
	fmt.Fprintf(
		&sb,
		"HP [%s%s] %d/%d\nLog:",
		strings.Repeat("#", filled),
		strings.Repeat(".", healthBarWidth-filled),
		state.Player.Health,
		state.Player.MaxHealth,
	)
	logStart := max(0, len(state.Log)-visibleLogLines)
	for _, entry := range state.Log[logStart:] {
		sb.WriteByte('\n')
		sb.WriteString(entry)
	}
	if state.GameOver {
		sb.WriteString("\nGAME OVER\nPress r to restart")
	}
	return sb.String()
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(rng), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
