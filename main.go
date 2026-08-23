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

type colorPalette struct {
	LitFloor    lipgloss.Color
	LitFloorBG  lipgloss.Color
	LitWall     lipgloss.Color
	LitWallBG   lipgloss.Color
	MemoryFloor lipgloss.Color
	MemoryWall  lipgloss.Color
	MemoryBG    lipgloss.Color
	Player      lipgloss.Color
	Monster     lipgloss.Color
	UIText      lipgloss.Color
	UIAccent    lipgloss.Color
	Health      lipgloss.Color
	HealthEmpty lipgloss.Color
	Danger      lipgloss.Color
}

var defaultPalette = colorPalette{
	LitFloor:    "#64748B",
	LitFloorBG:  "#111827",
	LitWall:     "#B8C2D1",
	LitWallBG:   "#1E293B",
	MemoryFloor: "#303846",
	MemoryWall:  "#46505F",
	MemoryBG:    "#0B1018",
	Player:      "#FFD166",
	Monster:     "#EF6A6A",
	UIText:      "#CBD5E1",
	UIAccent:    "#67E8F9",
	Health:      "#6EE7B7",
	HealthEmpty: "#334155",
	Danger:      "#FB7185",
}

type viewStyles struct {
	player      lipgloss.Style
	monster     lipgloss.Style
	floor       lipgloss.Style
	wall        lipgloss.Style
	memoryFloor lipgloss.Style
	memoryWall  lipgloss.Style
	uiText      lipgloss.Style
	uiAccent    lipgloss.Style
	health      lipgloss.Style
	healthEmpty lipgloss.Style
	danger      lipgloss.Style
}

func newViewStyles(p colorPalette) viewStyles {
	return viewStyles{
		player:      lipgloss.NewStyle().Foreground(p.Player).Bold(true),
		monster:     lipgloss.NewStyle().Foreground(p.Monster).Bold(true),
		floor:       lipgloss.NewStyle().Foreground(p.LitFloor).Background(p.LitFloorBG),
		wall:        lipgloss.NewStyle().Foreground(p.LitWall).Background(p.LitWallBG),
		memoryFloor: lipgloss.NewStyle().Foreground(p.MemoryFloor).Background(p.MemoryBG),
		memoryWall:  lipgloss.NewStyle().Foreground(p.MemoryWall).Background(p.MemoryBG),
		uiText:      lipgloss.NewStyle().Foreground(p.UIText),
		uiAccent:    lipgloss.NewStyle().Foreground(p.UIAccent).Bold(true),
		health:      lipgloss.NewStyle().Foreground(p.Health),
		healthEmpty: lipgloss.NewStyle().Foreground(p.HealthEmpty),
		danger:      lipgloss.NewStyle().Foreground(p.Danger).Bold(true),
	}
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

	styles := newViewStyles(defaultPalette)
	player := styles.player.Render(m.state.Player.Rune)
	floor := styles.floor.Render(".")
	wall := styles.wall.Render("#")
	memoryFloor := styles.memoryFloor.Render(".")
	memoryWall := styles.memoryWall.Render("#")
	entityGlyphs := make(map[Position]string, len(m.state.Entities))
	for _, entity := range m.state.Entities {
		entityGlyphs[entity.Pos] = styles.monster.Render(entity.Rune)
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

	return lipgloss.JoinHorizontal(lipgloss.Top, sb.String(), "  ", renderSidebar(m.state, styles))
}

func renderSidebar(state GameState, styles viewStyles) string {
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
	sb.WriteString(styles.uiAccent.Render("HP "))
	sb.WriteByte('[')
	sb.WriteString(styles.health.Render(strings.Repeat("#", filled)))
	sb.WriteString(styles.healthEmpty.Render(strings.Repeat(".", healthBarWidth-filled)))
	fmt.Fprintf(&sb, "] %d/%d\n", state.Player.Health, state.Player.MaxHealth)
	sb.WriteString(styles.uiAccent.Render("Log:"))
	logStart := max(0, len(state.Log)-visibleLogLines)
	for _, entry := range state.Log[logStart:] {
		sb.WriteByte('\n')
		sb.WriteString(styles.uiText.Render(entry))
	}
	if state.GameOver {
		sb.WriteByte('\n')
		sb.WriteString(styles.danger.Render("GAME OVER"))
		sb.WriteByte('\n')
		sb.WriteString(styles.uiAccent.Render("Press r to restart"))
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
