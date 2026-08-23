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
	entityGlyphs := make(map[Position]string, len(m.state.Entities))
	for _, entity := range m.state.Entities {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(entity.Color)).Bold(true)
		entityGlyphs[entity.Pos] = style.Render(entity.Rune)
	}

	var sb strings.Builder
	for y := 0; y < m.state.Map.Height; y++ {
		for x := 0; x < m.state.Map.Width; x++ {
			if x == m.state.Player.Pos.X && y == m.state.Player.Pos.Y {
				sb.WriteString(player)
			} else if entity, ok := entityGlyphs[Position{X: x, Y: y}]; ok {
				sb.WriteString(entity)
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(rng), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
