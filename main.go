package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
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

type screenMode uint8

const (
	screenPlaying screenMode = iota
	screenTitle
	screenGameOver
	screenVictory
	screenTargeting
)

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
	flashTicks     int      // screen hit flash ticks
	screen         screenMode
	targetCursor   Position
}

// Default map dimensions for the larger dungeon.
const (
	defaultMapWidth  = 60
	defaultMapHeight = 30
)

func initialModel() model {
	var st GameState
	if seedStr := os.Getenv("GLYPHGRINDER_SEED"); seedStr != "" {
		if sVal, err := strconv.ParseInt(seedStr, 10, 64); err == nil {
			st = NewGameWithSeed(defaultMapWidth, defaultMapHeight, sVal)
		} else {
			st = NewGame(defaultMapWidth, defaultMapHeight)
		}
	} else if loaded, ok, _ := LoadAndRemoveSaveGame(DefaultSaveFilePath); ok {
		st = loaded
		st.Log = append(st.Log, "Resumed previous run from save file.")
	} else {
		st = NewGame(defaultMapWidth, defaultMapHeight)
	}
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

func initialTitleModel() model {
	m := initialModel()
	m.screen = screenTitle
	return m
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

// continuousCameraEasing calculates continuous camera center coordinates for smooth visual tracking.
func continuousCameraEasing(curX, curY, targetX, targetY float64, animTicks int) (float64, float64) {
	if animTicks <= 0 {
		return targetX, targetY
	}
	t := float64(maxAnimTicks-animTicks) / float64(maxAnimTicks)
	easeT := 1 - math.Pow(1-t, 3)

	newX := curX + (targetX-curX)*easeT
	newY := curY + (targetY-curY)*easeT
	return newX, newY
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
		if m.flashTicks > 0 {
			m.flashTicks--
		}

		if m.isAnimating() || m.flashTicks > 0 {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" {
			if m.state.Player.Health > 0 && !m.state.IsVictory {
				_ = SaveGame(m.state, DefaultSaveFilePath)
			}
			return m, tea.Quit
		}

		if m.screen == screenTitle {
			switch k {
			case "1":
				m.state = NewGameWithSeedDepthAndClass(defaultMapWidth, defaultMapHeight, m.state.Seed, 1, ClassWarrior)
				m.screen = screenPlaying
			case "2":
				m.state = NewGameWithSeedDepthAndClass(defaultMapWidth, defaultMapHeight, m.state.Seed, 1, ClassRogue)
				m.screen = screenPlaying
			case "3":
				m.state = NewGameWithSeedDepthAndClass(defaultMapWidth, defaultMapHeight, m.state.Seed, 1, ClassMage)
				m.screen = screenPlaying
			case " ", "enter", "w", "a", "s", "d", "up", "down", "left", "right":
				m.screen = screenPlaying
			}
			return m, nil
		}

		if m.screen == screenTargeting {
			switch k {
			case "up", "w":
				if m.targetCursor.Y > 0 {
					m.targetCursor.Y--
				}
			case "down", "s":
				if m.targetCursor.Y < m.state.Map.Height-1 {
					m.targetCursor.Y++
				}
			case "left", "a":
				if m.targetCursor.X > 0 {
					m.targetCursor.X--
				}
			case "right", "d":
				if m.targetCursor.X < m.state.Map.Width-1 {
					m.targetCursor.X++
				}
			case "esc":
				m.screen = screenPlaying
			case "enter", "f", " ":
				m.screen = screenPlaying
				oldHealth := m.state.Player.Health
				m.state = m.state.Step(ActionUseItem)
				if m.state.Player.Health < oldHealth {
					m.flashTicks = 4
				}
				return m, tickCmd()
			}
			return m, nil
		}

		if m.state.Player.Health <= 0 || m.state.IsVictory {
			if k == "r" {
				m.state = NewGame(defaultMapWidth, defaultMapHeight)
				m.camX = float64(m.state.Player.Pos.X)
				m.camY = float64(m.state.Player.Pos.Y)
				m.camInitialized = true
				m.animTicks = 0
				m.screen = screenPlaying
			}
			return m, nil
		}

		if k == "t" {
			m.screen = screenTargeting
			m.targetCursor = m.state.Player.Pos
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
		case "x", "D":
			act = ActionDropItem
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
			oldHealth := m.state.Player.Health
			m.state = m.state.Step(act)
			if m.state.Player.Health < oldHealth {
				m.flashTicks = 4
			}
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

func renderTitleScreen(pal Palette, gly GlyphSet, seed int64) string {
	titleStyle := lipgloss.NewStyle().Foreground(pal.Stairs).Bold(true)
	subStyle := lipgloss.NewStyle().Foreground(pal.Player).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(pal.Weapon)
	descStyle := lipgloss.NewStyle().Foreground(pal.HUDLog)
	promptStyle := lipgloss.NewStyle().Foreground(pal.HUDNormal).Bold(true)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  ____ _                    ____    _         _           "))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render(" / ___| |_   _ _ __ | |__  / ___| _ __(_)_ __   __| | ___ _ __ "))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("| |  _| | | | | '_ \\| '_ \\| |  _ | '__| | '_ \\ / _` |/ _ \\ '__|"))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("| |_| | | |_| | |_) | | | | |_| || |  | | | | | (_| |  __/ |   "))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render(" \\____|_|\\__, | .__/|_| |_|\\____||_|  |_|_| |_|\\__,_|\\___|_|   "))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("         |___/|_|                                              "))
	sb.WriteString("\n\n")

	sb.WriteString(subStyle.Render("      --- SELECT CLASS ARCHETYPE TO BEGIN DESCENT ---"))
	sb.WriteString("\n\n")

	sb.WriteString(keyStyle.Render("  [1]") + descStyle.Render(" Warrior — 120 HP, High Melee Damage, Starts with Iron Dagger\n"))
	sb.WriteString(keyStyle.Render("  [2]") + descStyle.Render(" Rogue   —  90 HP, High Burst Damage, Starts with Teleport Scroll\n"))
	sb.WriteString(keyStyle.Render("  [3]") + descStyle.Render(" Mage    —  80 HP, Starts with Fireball Scroll & Regen Potion\n\n"))

	sb.WriteString(fmt.Sprintf("                 Dungeon Seed: %d\n\n", seed))

	sb.WriteString(keyStyle.Render("  [W/A/S/D / Arrows]") + descStyle.Render(" Move & Bump Attack\n"))
	sb.WriteString(keyStyle.Render("  [G / ,]") + descStyle.Render(" Pick Up Item    ") + keyStyle.Render("[X / D]") + descStyle.Render(" Drop Item\n"))
	sb.WriteString(keyStyle.Render("  [H / 1-9]") + descStyle.Render(" Use Item       ") + keyStyle.Render("[T]") + descStyle.Render(" Target Ranged / Scroll\n"))
	sb.WriteString(keyStyle.Render("  [> / Enter]") + descStyle.Render(" Descend Stairs ") + keyStyle.Render("[Q]") + descStyle.Render(" Quit Game\n\n"))

	sb.WriteString(promptStyle.Render("   === PRESS [1], [2], [3] OR [SPACE] TO BEGIN DESCENT ==="))
	sb.WriteString("\n")
	return sb.String()
}

type renderCache struct {
	pal Palette
	gly GlyphSet

	Player        string
	Floor         string
	Wall          string
	Stairs        string
	Lava          string
	Water         string
	DoorClosed    string
	DoorOpen      string
	DimFloor      string
	DimWall       string
	DimStairs     string
	DimLava       string
	DimWater      string
	DimDoorClosed string
	DimDoorOpen   string
	TargetReticle string

	HUDNormalStyle  lipgloss.Style
	HUDWarningStyle lipgloss.Style
	HUDLogStyle     lipgloss.Style

	GoblinGlyph string
	GoblinStyle lipgloss.Style
	OrcGlyph    string
	OrcStyle    lipgloss.Style
	TrollGlyph  string
	TrollStyle  lipgloss.Style
	ArcherGlyph string
	ArcherStyle lipgloss.Style

	PotionGlyph string
	PotionStyle lipgloss.Style
	WeaponGlyph string
	WeaponStyle lipgloss.Style
	AmuletGlyph string
	AmuletStyle lipgloss.Style
	ScrollGlyph string
	ScrollStyle lipgloss.Style
}

var globalRenderCache *renderCache

func getRenderCache(pal Palette, gly GlyphSet) *renderCache {
	if globalRenderCache != nil && globalRenderCache.pal == pal && globalRenderCache.gly == gly {
		return globalRenderCache
	}

	c := &renderCache{
		pal: pal,
		gly: gly,

		Player:        lipgloss.NewStyle().Foreground(pal.Player).Bold(true).Render(gly.Player),
		Floor:         lipgloss.NewStyle().Foreground(pal.FloorLit).Render(gly.Floor),
		Wall:          lipgloss.NewStyle().Foreground(pal.WallLit).Render(gly.Wall),
		Stairs:        lipgloss.NewStyle().Foreground(pal.Stairs).Bold(true).Render(gly.StairsDown),
		Lava:          lipgloss.NewStyle().Foreground(pal.Lava).Bold(true).Render(gly.Lava),
		Water:         lipgloss.NewStyle().Foreground(pal.Water).Bold(true).Render(gly.Water),
		DoorClosed:    lipgloss.NewStyle().Foreground(pal.Door).Render(gly.DoorClosed),
		DoorOpen:      lipgloss.NewStyle().Foreground(pal.Door).Render(gly.DoorOpen),
		DimFloor:      lipgloss.NewStyle().Foreground(pal.FloorDim).Render(gly.Floor),
		DimWall:       lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.Wall),
		DimStairs:     lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.StairsDown),
		DimLava:       lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.Lava),
		DimWater:      lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.Water),
		DimDoorClosed: lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.DoorClosed),
		DimDoorOpen:   lipgloss.NewStyle().Foreground(pal.WallDim).Render(gly.DoorOpen),
		TargetReticle: lipgloss.NewStyle().Foreground(pal.HUDWarning).Bold(true).Render("X"),

		HUDNormalStyle:  lipgloss.NewStyle().Foreground(pal.HUDNormal).Bold(true),
		HUDWarningStyle: lipgloss.NewStyle().Foreground(pal.HUDWarning).Bold(true),
		HUDLogStyle:     lipgloss.NewStyle().Foreground(pal.HUDLog),

		GoblinGlyph: gly.Goblin,
		GoblinStyle: lipgloss.NewStyle().Foreground(pal.Goblin).Bold(true),
		OrcGlyph:    gly.Orc,
		OrcStyle:    lipgloss.NewStyle().Foreground(pal.Orc).Bold(true),
		TrollGlyph:  gly.Troll,
		TrollStyle:  lipgloss.NewStyle().Foreground(pal.Troll).Bold(true),
		ArcherGlyph: gly.Archer,
		ArcherStyle: lipgloss.NewStyle().Foreground(pal.Archer).Bold(true),

		PotionGlyph: gly.Potion,
		PotionStyle: lipgloss.NewStyle().Foreground(pal.Potion).Bold(true),
		WeaponGlyph: gly.Weapon,
		WeaponStyle: lipgloss.NewStyle().Foreground(pal.Weapon).Bold(true),
		AmuletGlyph: gly.Amulet,
		AmuletStyle: lipgloss.NewStyle().Foreground(pal.Amulet).Bold(true),
		ScrollGlyph: gly.Scroll,
		ScrollStyle: lipgloss.NewStyle().Foreground(pal.Scroll).Bold(true),
	}
	globalRenderCache = c
	return c
}

func (m model) View() string {
	if m.state.Map.Width == 0 || m.state.Map.Height == 0 {
		return ""
	}

	pal := m.getPalette()
	gly := m.getGlyphs()

	if m.screen == screenTitle {
		return renderTitleScreen(pal, gly, m.state.Seed)
	}

	cache := getRenderCache(pal, gly)
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
	hpStyle := cache.HUDNormalStyle
	if m.flashTicks > 0 || hp == 0 {
		hpStyle = cache.HUDWarningStyle
	}
	hudStr := fmt.Sprintf("HP: %d/%d | Depth: %d | Seed: %d", hp, m.state.Player.MaxHealth, depth, m.state.Seed)
	if len(m.state.Player.Inventory) > 0 {
		var invNames []string
		for _, item := range m.state.Player.Inventory {
			invNames = append(invNames, item.Name)
		}
		hudStr += fmt.Sprintf(" | Inv: [%s]", strings.Join(invNames, ", "))
	}
	hudText := hpStyle.Render(hudStr)

	if m.screen == screenTargeting {
		targetBannerStyle := cache.WeaponStyle
		hudText += targetBannerStyle.Render(" | [TARGETING MODE: Use Arrows, Enter to fire, Esc to cancel]")
	} else if m.state.IsVictory {
		victoryStyle := lipgloss.NewStyle().Foreground(pal.Stairs).Bold(true)
		hudText += victoryStyle.Render(fmt.Sprintf(" | *** VICTORY! Slain: %d | Turns: %d *** (Press r to restart)", m.state.Kills, m.state.TurnCount))
	} else if m.state.Player.Health <= 0 {
		gameOverStyle := cache.HUDWarningStyle
		hudText += gameOverStyle.Render(fmt.Sprintf(" | *** GAME OVER *** Slain: %d | Turns: %d (Press r to restart)", m.state.Kills, m.state.TurnCount))
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

	entityMap := make(map[Position]string, len(m.state.Entities))
	for _, e := range m.state.Entities {
		switch e.Name {
		case "Goblin":
			entityMap[e.Pos] = cache.GoblinStyle.Render(cache.GoblinGlyph)
		case "Orc":
			entityMap[e.Pos] = cache.OrcStyle.Render(cache.OrcGlyph)
		case "Troll":
			entityMap[e.Pos] = cache.TrollStyle.Render(cache.TrollGlyph)
		case "Archer":
			entityMap[e.Pos] = cache.ArcherStyle.Render(cache.ArcherGlyph)
		default:
			style := lipgloss.NewStyle().Foreground(ResolveEntityColor(e, pal)).Bold(true)
			glyph := ResolveEntityGlyph(e, gly)
			entityMap[e.Pos] = style.Render(glyph)
		}
	}

	itemMap := make(map[Position]string, len(m.state.Items))
	for _, it := range m.state.Items {
		switch it.ItemType {
		case ItemPotion:
			itemMap[it.Pos] = cache.PotionStyle.Render(cache.PotionGlyph)
		case ItemWeapon:
			itemMap[it.Pos] = cache.WeaponStyle.Render(cache.WeaponGlyph)
		case ItemAmulet:
			itemMap[it.Pos] = cache.AmuletStyle.Render(cache.AmuletGlyph)
		case ItemScroll:
			itemMap[it.Pos] = cache.ScrollStyle.Render(cache.ScrollGlyph)
		default:
			style := lipgloss.NewStyle().Foreground(ResolveItemColor(it, pal)).Bold(true)
			glyph := ResolveItemGlyph(it, gly)
			itemMap[it.Pos] = style.Render(glyph)
		}
	}

	hasFOV := m.state.Map.Visible != nil && m.state.Map.Explored != nil
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			pos := Position{X: x, Y: y}
			visible := !hasFOV || m.state.Map.Visible[y][x]
			explored := !hasFOV || m.state.Map.Explored[y][x]

			if m.screen == screenTargeting && pos == m.targetCursor {
				sb.WriteString(cache.TargetReticle)
			} else if !explored {
				// Tile has never been seen — render as blank.
				sb.WriteString(" ")
			} else if visible {
				// Currently in line-of-sight — full brightness.
				if pos == m.state.Player.Pos {
					sb.WriteString(cache.Player)
				} else if renderedEntity, found := entityMap[pos]; found {
					sb.WriteString(renderedEntity)
				} else if renderedItem, found := itemMap[pos]; found {
					sb.WriteString(renderedItem)
				} else {
					switch m.state.Map.Tiles[y][x] {
					case TileWall:
						sb.WriteString(cache.Wall)
					case TileStairsDown:
						sb.WriteString(cache.Stairs)
					case TileDoorClosed:
						sb.WriteString(cache.DoorClosed)
					case TileDoorOpen:
						sb.WriteString(cache.DoorOpen)
					case TileLava:
						sb.WriteString(cache.Lava)
					case TileWater:
						sb.WriteString(cache.Water)
					default:
						sb.WriteString(cache.Floor)
					}
				}
			} else {
				// Explored but not visible — dimmed, no monsters or items.
				switch m.state.Map.Tiles[y][x] {
				case TileWall:
					sb.WriteString(cache.DimWall)
				case TileStairsDown:
					sb.WriteString(cache.DimStairs)
				case TileDoorClosed:
					sb.WriteString(cache.DimDoorClosed)
				case TileDoorOpen:
					sb.WriteString(cache.DimDoorOpen)
				case TileLava:
					sb.WriteString(cache.DimLava)
				case TileWater:
					sb.WriteString(cache.DimWater)
				default:
					sb.WriteString(cache.DimFloor)
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
	logStyle := cache.HUDLogStyle
	for i := start; i < logCount; i++ {
		sb.WriteString(logStyle.Render(m.state.Log[i]))
		if i < logCount-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-dump-frame" || arg == "--dump-frame" {
			m := initialModel()
			fmt.Print(m.View())
			os.Exit(0)
		}
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
