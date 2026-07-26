package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state GameState
}

func initialModel() model {
	// Initialize a 20x10 terminal grid as the map
	return model{state: NewGame(20, 10)}
}

func initialModelWithSeed(seed int64) model {
	return model{state: NewGameWithSeed(20, 10, seed)}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var act Action
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "w":
			act = ActionMoveUp
		case "down", "s":
			act = ActionMoveDown
		case "left", "a":
			act = ActionMoveLeft
		case "right", "d":
			act = ActionMoveRight
		}
		if act != ActionNone {
			m.state = m.state.Step(act)
		}
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

	var sb strings.Builder
	for y := 0; y < m.state.Map.Height; y++ {
		for x := 0; x < m.state.Map.Width; x++ {
			if x == m.state.Player.Pos.X && y == m.state.Player.Pos.Y {
				sb.WriteString(player)
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

	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
