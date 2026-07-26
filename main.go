package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state   GameState
	width   int      // terminal columns (from tea.WindowSizeMsg)
	height  int      // terminal rows  (from tea.WindowSizeMsg)
	palette Palette  // color palette for profile-aware rendering
	glyphs  GlyphSet // glyph set for rendering
}

// Default map dimensions for the larger dungeon.
const (
	defaultMapWidth  = 60
	defaultMapHeight = 30
)

func initialModel() model {
	return model{
		state:   NewGame(defaultMapWidth, defaultMapHeight),
		width:   80,
		height:  24,
		palette: DefaultPalette(),
		glyphs:  DetectGlyphSet(),
	}
}

func initialModelWithSeed(seed int64) model {
	return model{
		state:   NewGameWithSeed(defaultMapWidth, defaultMapHeight, seed),
		width:   80,
		height:  24,
		palette: DefaultPalette(),
		glyphs:  DetectGlyphSet(),
	}
}

func (m model) getPalette() Palette {
	if m.palette.FloorLit.TrueColor == "" && m.palette.FloorLit.ANSI256 == "" && m.palette.FloorLit.ANSI == "" {
		return DefaultPalette()
	}
	return m.palette
}

func (m model) getGlyphs() GlyphSet {
	if m.glyphs.Player == "" {
		return DetectGlyphSet()
	}
	return m.glyphs
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
		if m.state.Player.Health <= 0 {
			if k == "r" {
				m.state = NewGame(defaultMapWidth, defaultMapHeight)
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

// viewport computes the map sub-rectangle to render, centered on the player
// and clamped to map bounds. viewW and viewH are the number of map columns
// and rows that fit in the terminal (after reserving space for HUD and log).
func viewport(playerPos Position, mapW, mapH, viewW, viewH int) (x0, y0, x1, y1 int) {
	if viewW >= mapW {
		x0 = 0
		x1 = mapW
	} else {
		x0 = playerPos.X - viewW/2
		if x0 < 0 {
			x0 = 0
		}
		x1 = x0 + viewW
		if x1 > mapW {
			x1 = mapW
			x0 = x1 - viewW
		}
	}

	if viewH >= mapH {
		y0 = 0
		y1 = mapH
	} else {
		y0 = playerPos.Y - viewH/2
		if y0 < 0 {
			y0 = 0
		}
		y1 = y0 + viewH
		if y1 > mapH {
			y1 = mapH
			y0 = y1 - viewH
		}
	}

	return x0, y0, x1, y1
}

// reservedRows is the number of terminal rows reserved for non-map UI
// (1 HUD line + up to 5 log lines).
const reservedRows = 6

func (m model) View() string {
	if m.state.Map.Width == 0 || m.state.Map.Height == 0 {
		return ""
	}

	pal := m.getPalette()
	var sb strings.Builder

	// Render HUD / Status Bar
	hp := m.state.Player.Health
	if hp < 0 {
		hp = 0
	}
	hpStyle := lipgloss.NewStyle().Foreground(pal.HUDNormal).Bold(true)
	if hp == 0 {
		hpStyle = lipgloss.NewStyle().Foreground(pal.HUDWarning).Bold(true)
	}
	hudText := hpStyle.Render(fmt.Sprintf("HP: %d/%d", hp, m.state.Player.MaxHealth))

	if m.state.Player.Health <= 0 {
		gameOverStyle := lipgloss.NewStyle().Foreground(pal.HUDWarning).Bold(true)
		hudText += gameOverStyle.Render(" | *** GAME OVER *** (Press r to restart)")
	}
	sb.WriteString(hudText)
	sb.WriteString("\n")

	// Compute viewport — how much of the map fits in the terminal.
	viewW := m.width
	viewH := m.height - reservedRows
	if viewW <= 0 {
		viewW = m.state.Map.Width
	}
	if viewH <= 0 {
		viewH = m.state.Map.Height
	}
	x0, y0, x1, y1 := viewport(m.state.Player.Pos, m.state.Map.Width, m.state.Map.Height, viewW, viewH)

	gly := m.getGlyphs()

	// Render Map Grid (visible sub-rectangle only)
	playerStyle := lipgloss.NewStyle().Foreground(pal.Player).Bold(true)
	player := playerStyle.Render(gly.Player)
	floorStyle := lipgloss.NewStyle().Foreground(pal.FloorLit)
	floor := floorStyle.Render(gly.Floor)
	wallStyle := lipgloss.NewStyle().Foreground(pal.WallLit)
	wall := wallStyle.Render(gly.Wall)
	// Dimmed styles for explored-but-not-visible tiles (map memory).
	dimFloorStyle := lipgloss.NewStyle().Foreground(pal.FloorDim)
	dimFloor := dimFloorStyle.Render(gly.Floor)
	dimWallStyle := lipgloss.NewStyle().Foreground(pal.WallDim)
	dimWall := dimWallStyle.Render(gly.Wall)

	entityMap := make(map[Position]string, len(m.state.Entities))
	for _, e := range m.state.Entities {
		style := lipgloss.NewStyle().Foreground(ResolveEntityColor(e, pal)).Bold(true)
		glyph := ResolveEntityGlyph(e, gly)
		entityMap[e.Pos] = style.Render(glyph)
	}

	hasFOV := m.state.Map.Visible != nil && m.state.Map.Explored != nil
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
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
	logStyle := lipgloss.NewStyle().Foreground(pal.HUDLog)
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
