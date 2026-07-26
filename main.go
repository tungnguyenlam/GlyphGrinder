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
		k := msg.String()
		if k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
		if m.state.Player.Health <= 0 {
			if k == "r" {
				m.state = NewGame(20, 10)
			}
			return m, nil
		}

		var act Action
		switch k {
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

	var sb strings.Builder

	// Render HUD / Status Bar
	hp := m.state.Player.Health
	if hp < 0 {
		hp = 0
	}
	hpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	if hp == 0 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	}
	hudText := hpStyle.Render(fmt.Sprintf("HP: %d/%d", hp, m.state.Player.MaxHealth))

	if m.state.Player.Health <= 0 {
		gameOverStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
		hudText += gameOverStyle.Render(" | *** GAME OVER *** (Press r to restart)")
	}
	sb.WriteString(hudText)
	sb.WriteString("\n")

	// Render Map Grid
	playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.state.Player.Color)).Bold(true)
	player := playerStyle.Render(m.state.Player.Rune)
	floorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	floor := floorStyle.Render(".")
	wallStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	wall := wallStyle.Render("#")
	// Dimmed styles for explored-but-not-visible tiles (map memory).
	dimFloorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#222222"))
	dimFloor := dimFloorStyle.Render(".")
	dimWallStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	dimWall := dimWallStyle.Render("#")

	entityMap := make(map[Position]string, len(m.state.Entities))
	for _, e := range m.state.Entities {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(e.Color)).Bold(true)
		entityMap[e.Pos] = style.Render(e.Rune)
	}

	hasFOV := m.state.Map.Visible != nil && m.state.Map.Explored != nil
	for y := 0; y < m.state.Map.Height; y++ {
		for x := 0; x < m.state.Map.Width; x++ {
			pos := Position{X: x, Y: y}
			visible := !hasFOV || m.state.Map.Visible[y][x]
			explored := !hasFOV || m.state.Map.Explored[y][x]

			if !explored {
				// Tile has never been seen — render as blank.
				sb.WriteString(" ")
			} else if visible {
				// Currently in line-of-sight — full brightness.
				if pos == m.state.Player.Pos {
					sb.WriteString(player)
				} else if renderedEntity, found := entityMap[pos]; found {
					sb.WriteString(renderedEntity)
				} else {
					if m.state.Map.Tiles[y][x] == TileWall {
						sb.WriteString(wall)
					} else {
						sb.WriteString(floor)
					}
				}
			} else {
				// Explored but not visible — dimmed, no monsters.
				if m.state.Map.Tiles[y][x] == TileWall {
					sb.WriteString(dimWall)
				} else {
					sb.WriteString(dimFloor)
				}
			}
		}
		sb.WriteString("\n")
	}

	// Render Message Log (last 5 entries)
	logCount := len(m.state.Log)
	start := 0
	if logCount > 5 {
		start = logCount - 5
	}
	logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	for i := start; i < logCount; i++ {
		sb.WriteString(logStyle.Render(m.state.Log[i]))
		if i < logCount-1 {
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
