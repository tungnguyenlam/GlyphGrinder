package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tickInterval = 16 * time.Millisecond

const (
	camSpeed     = 0.35
	maxAnimTicks = 4
)

type animTickMsg time.Time

func tickCmd() tea.Cmd {
	if tickInterval == 0 {
		return func() tea.Msg {
			return animTickMsg(time.Now())
		}
	}
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

type model struct {
	state          GameState
	width          int      // terminal columns (from tea.WindowSizeMsg)
	height         int      // terminal rows  (from tea.WindowSizeMsg)
	palette        Palette  // color palette for profile-aware rendering
	glyphs         GlyphSet // glyph set for rendering
	camX           float64  // visual camera center X coordinate
	camY           float64  // visual camera center Y coordinate
	camInitialized bool     // whether camera position has been initialized
	animTicks      int      // remaining animation ticks for current move step
}

// Default map dimensions for the larger dungeon.
const (
	defaultMapWidth  = 60
	defaultMapHeight = 30
)

func initialModel() model {
	st := NewGame(defaultMapWidth, defaultMapHeight)
	return model{
		state:          st,
		width:          80,
		height:         24,
		palette:        DefaultPalette(),
		glyphs:         DetectGlyphSet(),
		camX:           float64(st.Player.Pos.X),
		camY:           float64(st.Player.Pos.Y),
		camInitialized: true,
	}
}

func initialModelWithSeed(seed int64) model {
	st := NewGameWithSeed(defaultMapWidth, defaultMapHeight, seed)
	return model{
		state:          st,
		width:          80,
		height:         24,
		palette:        DefaultPalette(),
		glyphs:         DetectGlyphSet(),
		camX:           float64(st.Player.Pos.X),
		camY:           float64(st.Player.Pos.Y),
		camInitialized: true,
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

func (m model) getCamPos() (float64, float64) {
	if !m.camInitialized {
		return float64(m.state.Player.Pos.X), float64(m.state.Player.Pos.Y)
	}
	return m.camX, m.camY
}

func (m model) isAnimating() bool {
	if !m.camInitialized {
		return false
	}
	targetX := float64(m.state.Player.Pos.X)
	targetY := float64(m.state.Player.Pos.Y)
	return m.animTicks > 0 || m.camX != targetX || m.camY != targetY
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
	case animTickMsg:
		if !m.camInitialized {
			m.camX = float64(m.state.Player.Pos.X)
			m.camY = float64(m.state.Player.Pos.Y)
			m.camInitialized = true
		}
		targetX := float64(m.state.Player.Pos.X)
		targetY := float64(m.state.Player.Pos.Y)

		dx := targetX - m.camX
		dy := targetY - m.camY

		if math.Abs(dx) < 0.05 {
			m.camX = targetX
		} else {
			m.camX += dx * camSpeed
		}

		if math.Abs(dy) < 0.05 {
			m.camY = targetY
		} else {
			m.camY += dy * camSpeed
		}

		if m.animTicks > 0 {
			m.animTicks--
		}

		if m.isAnimating() {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
		if m.state.Player.Health <= 0 {
			if k == "r" {
				m.state = NewGame(defaultMapWidth, defaultMapHeight)
				m.camX = float64(m.state.Player.Pos.X)
				m.camY = float64(m.state.Player.Pos.Y)
				m.camInitialized = true
				m.animTicks = 0
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
		case ">", "enter":
			act = ActionDescend
		case "g", ",":
			act = ActionPickup
		case "h":
			act = ActionUseItem
		case "1":
			act = ActionUseItem1
		case "2":
			act = ActionUseItem2
		case "3":
			act = ActionUseItem3
		case "4":
			act = ActionUseItem4
		case "5":
			act = ActionUseItem5
		case "6":
			act = ActionUseItem6
		case "7":
			act = ActionUseItem7
		case "8":
			act = ActionUseItem8
		case "9":
			act = ActionUseItem9
		}
		if act != ActionNone {
			if !m.camInitialized {
				m.camX = float64(m.state.Player.Pos.X)
				m.camY = float64(m.state.Player.Pos.Y)
				m.camInitialized = true
			}
			oldDepth := m.state.Depth
			m.state = m.state.Step(act)
			if m.state.Depth != oldDepth {
				m.camX = float64(m.state.Player.Pos.X)
				m.camY = float64(m.state.Player.Pos.Y)
				m.animTicks = 0
			} else {
				m.animTicks = maxAnimTicks
			}
			return m, tickCmd()
		}
	}
	return m, nil
}

// viewport computes the map sub-rectangle to render, centered on playerPos.
func viewport(playerPos Position, mapW, mapH, viewW, viewH int) (x0, y0, x1, y1 int) {
	return viewportCenter(float64(playerPos.X), float64(playerPos.Y), mapW, mapH, viewW, viewH)
}

// viewportCenter computes the map sub-rectangle centered on continuous coordinates (centerX, centerY).
func viewportCenter(centerX, centerY float64, mapW, mapH, viewW, viewH int) (x0, y0, x1, y1 int) {
	cx := int(math.Round(centerX))
	cy := int(math.Round(centerY))

	if viewW >= mapW {
		x0 = 0
		x1 = mapW
	} else {
		x0 = cx - viewW/2
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
		y0 = cy - viewH/2
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
	depth := m.state.Depth
	if depth < 1 {
		depth = 1
	}
	hpStyle := lipgloss.NewStyle().Foreground(pal.HUDNormal).Bold(true)
	if hp == 0 {
		hpStyle = lipgloss.NewStyle().Foreground(pal.HUDWarning).Bold(true)
	}
	hudStr := fmt.Sprintf("HP: %d/%d | Depth: %d", hp, m.state.Player.MaxHealth, depth)
	if len(m.state.Player.Inventory) > 0 {
		var invNames []string
		for _, item := range m.state.Player.Inventory {
			invNames = append(invNames, item.Name)
		}
		hudStr += fmt.Sprintf(" | Inv: [%s]", strings.Join(invNames, ", "))
	}
	hudText := hpStyle.Render(hudStr)

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
	camX, camY := m.getCamPos()
	x0, y0, x1, y1 := viewportCenter(camX, camY, m.state.Map.Width, m.state.Map.Height, viewW, viewH)

	gly := m.getGlyphs()

	// Render Map Grid (visible sub-rectangle only)
	playerStyle := lipgloss.NewStyle().Foreground(pal.Player).Bold(true)
	player := playerStyle.Render(gly.Player)
	floorStyle := lipgloss.NewStyle().Foreground(pal.FloorLit)
	floor := floorStyle.Render(gly.Floor)
	wallStyle := lipgloss.NewStyle().Foreground(pal.WallLit)
	wall := wallStyle.Render(gly.Wall)
	stairsStyle := lipgloss.NewStyle().Foreground(pal.Stairs).Bold(true)
	stairs := stairsStyle.Render(gly.StairsDown)
	// Dimmed styles for explored-but-not-visible tiles (map memory).
	dimFloorStyle := lipgloss.NewStyle().Foreground(pal.FloorDim)
	dimFloor := dimFloorStyle.Render(gly.Floor)
	dimWallStyle := lipgloss.NewStyle().Foreground(pal.WallDim)
	dimWall := dimWallStyle.Render(gly.Wall)
	dimStairsStyle := lipgloss.NewStyle().Foreground(pal.WallDim)
	dimStairs := dimStairsStyle.Render(gly.StairsDown)

	entityMap := make(map[Position]string, len(m.state.Entities))
	for _, e := range m.state.Entities {
		style := lipgloss.NewStyle().Foreground(ResolveEntityColor(e, pal)).Bold(true)
		glyph := ResolveEntityGlyph(e, gly)
		entityMap[e.Pos] = style.Render(glyph)
	}

	itemMap := make(map[Position]string, len(m.state.Items))
	for _, it := range m.state.Items {
		style := lipgloss.NewStyle().Foreground(ResolveItemColor(it, pal)).Bold(true)
		glyph := ResolveItemGlyph(it, gly)
		itemMap[it.Pos] = style.Render(glyph)
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
				} else if renderedItem, found := itemMap[pos]; found {
					sb.WriteString(renderedItem)
				} else {
					switch m.state.Map.Tiles[y][x] {
					case TileWall:
						sb.WriteString(wall)
					case TileStairsDown:
						sb.WriteString(stairs)
					default:
						sb.WriteString(floor)
					}
				}
			} else {
				// Explored but not visible — dimmed, no monsters or items.
				switch m.state.Map.Tiles[y][x] {
				case TileWall:
					sb.WriteString(dimWall)
				case TileStairsDown:
					sb.WriteString(dimStairs)
				default:
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
