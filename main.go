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
	state        GameState
	rng          *rand.Rand
	glyphs       glyphProfile
	motion       actorMotionState
	motionSeq    uint64
	windowWidth  int
	windowHeight int
}

type actorMotion struct {
	ID       int
	IsPlayer bool
	From     Position
	To       Position
}

type actorMotionState struct {
	Actors   []actorMotion
	Frame    int
	Sequence uint64
}

type animationTickMsg struct {
	Sequence uint64
}

type mapViewport struct {
	X      int
	Y      int
	Width  int
	Height int
}

const (
	productionDungeonWidth  = 96
	productionDungeonHeight = 48
	sidebarWidth            = 32
	mapSidebarGap           = 2
	animationFrames         = 3
	animationFrameDuration  = time.Second / 60
)

type glyphProfile struct {
	Player  string
	Monster string
	Ogre    string
	Bat     string
	Floor   string
	Wall    string
	Stairs  string
	Potion  string
	Sword   string
}

var (
	asciiGlyphs = glyphProfile{Player: "@", Monster: "g", Ogre: "O", Bat: "b", Floor: ".", Wall: "#", Stairs: ">", Potion: "!", Sword: "/"}
	richGlyphs  = glyphProfile{
		Player:  "\uf007", // Nerd Fonts: Font Awesome user
		Monster: "\uf6e2", // Nerd Fonts: Font Awesome ghost
		Ogre:    "\uf714", // Nerd Fonts: Font Awesome skull-crossbones
		Bat:     "\uf6d5", // Nerd Fonts: Font Awesome dragon
		Floor:   "·",
		Wall:    "█",
		Stairs:  "\uf063", // Nerd Fonts: Font Awesome arrow-down
		Potion:  "\uf0c3", // Nerd Fonts: Font Awesome flask
		Sword:   "\uf71c", // Nerd Fonts: Material Design sword-cross
	}
)

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
	Ogre        lipgloss.Color
	Bat         lipgloss.Color
	Potion      lipgloss.Color
	Sword       lipgloss.Color
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
	Ogre:        "#F59E0B",
	Bat:         "#A78BFA",
	Potion:      "#C4B5FD",
	Sword:       "#F8FAFC",
	UIText:      "#CBD5E1",
	UIAccent:    "#67E8F9",
	Health:      "#6EE7B7",
	HealthEmpty: "#334155",
	Danger:      "#FB7185",
}

type viewStyles struct {
	player       lipgloss.Style
	playerTrail  lipgloss.Style
	monster      lipgloss.Style
	monsterTrail lipgloss.Style
	ogre         lipgloss.Style
	ogreTrail    lipgloss.Style
	bat          lipgloss.Style
	batTrail     lipgloss.Style
	floor        lipgloss.Style
	wall         lipgloss.Style
	memoryFloor  lipgloss.Style
	memoryWall   lipgloss.Style
	stairs       lipgloss.Style
	memoryStairs lipgloss.Style
	potion       lipgloss.Style
	sword        lipgloss.Style
	uiText       lipgloss.Style
	uiAccent     lipgloss.Style
	health       lipgloss.Style
	healthEmpty  lipgloss.Style
	danger       lipgloss.Style
}

func newViewStyles(p colorPalette) viewStyles {
	return viewStyles{
		player:       lipgloss.NewStyle().Foreground(p.Player).Bold(true),
		playerTrail:  lipgloss.NewStyle().Foreground(p.Player).Faint(true),
		monster:      lipgloss.NewStyle().Foreground(p.Monster).Bold(true),
		monsterTrail: lipgloss.NewStyle().Foreground(p.Monster).Faint(true),
		ogre:         lipgloss.NewStyle().Foreground(p.Ogre).Bold(true),
		ogreTrail:    lipgloss.NewStyle().Foreground(p.Ogre).Faint(true),
		bat:          lipgloss.NewStyle().Foreground(p.Bat).Bold(true),
		batTrail:     lipgloss.NewStyle().Foreground(p.Bat).Faint(true),
		floor:        lipgloss.NewStyle().Foreground(p.LitFloor).Background(p.LitFloorBG),
		wall:         lipgloss.NewStyle().Foreground(p.LitWall).Background(p.LitWallBG),
		memoryFloor:  lipgloss.NewStyle().Foreground(p.MemoryFloor).Background(p.MemoryBG),
		memoryWall:   lipgloss.NewStyle().Foreground(p.MemoryWall).Background(p.MemoryBG),
		stairs:       lipgloss.NewStyle().Foreground(p.UIAccent).Background(p.LitFloorBG).Bold(true),
		memoryStairs: lipgloss.NewStyle().Foreground(p.MemoryWall).Background(p.MemoryBG),
		potion:       lipgloss.NewStyle().Foreground(p.Potion).Background(p.LitFloorBG).Bold(true),
		sword:        lipgloss.NewStyle().Foreground(p.Sword).Background(p.LitFloorBG).Bold(true),
		uiText:       lipgloss.NewStyle().Foreground(p.UIText),
		uiAccent:     lipgloss.NewStyle().Foreground(p.UIAccent).Bold(true),
		health:       lipgloss.NewStyle().Foreground(p.Health),
		healthEmpty:  lipgloss.NewStyle().Foreground(p.HealthEmpty),
		danger:       lipgloss.NewStyle().Foreground(p.Danger).Bold(true),
	}
}

func initialModel(rng *rand.Rand) model {
	return model{
		state:  NewGame(productionDungeonWidth, productionDungeonHeight, rng),
		rng:    rng,
		glyphs: glyphProfileForEnvironment(os.Getenv),
	}
}

func glyphProfileForEnvironment(getenv func(string) string) glyphProfile {
	if strings.EqualFold(strings.TrimSpace(getenv("TERM")), "dumb") {
		return asciiGlyphs
	}

	locale := getenv("LC_ALL")
	if locale == "" {
		locale = getenv("LC_CTYPE")
	}
	if locale == "" {
		locale = getenv("LANG")
	}
	locale = strings.ToLower(locale)
	if strings.Contains(locale, "utf-8") || strings.Contains(locale, "utf8") {
		return richGlyphs
	}
	return asciiGlyphs
}

func (p glyphProfile) withASCIIFallback() glyphProfile {
	if lipgloss.Width(p.Player) != 1 ||
		lipgloss.Width(p.Monster) != 1 ||
		lipgloss.Width(p.Ogre) != 1 ||
		lipgloss.Width(p.Bat) != 1 ||
		lipgloss.Width(p.Floor) != 1 ||
		lipgloss.Width(p.Wall) != 1 ||
		lipgloss.Width(p.Stairs) != 1 ||
		lipgloss.Width(p.Potion) != 1 ||
		lipgloss.Width(p.Sword) != 1 {
		return asciiGlyphs
	}
	return p
}

func (p glyphProfile) monsterGlyph(archetype MonsterArchetype) string {
	switch archetype {
	case MonsterOgre:
		return p.Ogre
	case MonsterBat:
		return p.Bat
	default:
		return p.Monster
	}
}

func (s viewStyles) monsterStyle(archetype MonsterArchetype, trail bool) lipgloss.Style {
	switch archetype {
	case MonsterOgre:
		if trail {
			return s.ogreTrail
		}
		return s.ogre
	case MonsterBat:
		if trail {
			return s.batTrail
		}
		return s.bat
	default:
		if trail {
			return s.monsterTrail
		}
		return s.monster
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case animationTickMsg:
		if len(m.motion.Actors) == 0 || msg.Sequence != m.motion.Sequence {
			return m, nil
		}
		m.motion.Frame++
		if m.motion.Frame >= animationFrames {
			m.motion = actorMotionState{}
			return m, nil
		}
		return m, animationTick(m.motion.Sequence)
	case tea.KeyMsg:
		action := ActionNone
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			if m.state.GameOver || m.state.Won {
				m.state = NewGame(m.state.Map.Width, m.state.Map.Height, m.rng)
				m.motion = actorMotionState{}
			}
			return m, nil
		case ">":
			previousDepth := m.state.Depth
			previousWon := m.state.Won
			m.state = m.state.Descend(m.rng)
			if m.state.Depth != previousDepth || m.state.Won != previousWon {
				m.motionSeq++
				m.motion = actorMotionState{}
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
		case "p":
			action = ActionUsePotion
		case "e":
			action = ActionEquipWeapon
		}
		if action == ActionNone {
			return m, nil
		}
		previous := m.state
		m.state = m.state.Step(action)
		m.motionSeq++
		m.motion = newActorMotion(previous, m.state, m.motionSeq)
		if len(m.motion.Actors) > 0 {
			return m, animationTick(m.motion.Sequence)
		}
	}
	return m, nil
}

func animationTick(sequence uint64) tea.Cmd {
	return tea.Tick(animationFrameDuration, func(time.Time) tea.Msg {
		return animationTickMsg{Sequence: sequence}
	})
}

func newActorMotion(previous, next GameState, sequence uint64) actorMotionState {
	actors := make([]actorMotion, 0, len(next.Entities)+1)
	if previous.Player.Pos != next.Player.Pos {
		actors = append(actors, actorMotion{IsPlayer: true, From: previous.Player.Pos, To: next.Player.Pos})
	}

	previousMonsters := make(map[int]Position, len(previous.Entities))
	for _, entity := range previous.Entities {
		previousMonsters[entity.ID] = entity.Pos
	}
	for _, entity := range next.Entities {
		from, ok := previousMonsters[entity.ID]
		if ok && from != entity.Pos {
			actors = append(actors, actorMotion{ID: entity.ID, From: from, To: entity.Pos})
		}
	}
	return actorMotionState{Actors: actors, Sequence: sequence}
}

func (m model) animatedPosition(id int, isPlayer bool, resolved Position) Position {
	for _, actor := range m.motion.Actors {
		if actor.ID != id || actor.IsPlayer != isPlayer {
			continue
		}
		if m.motion.Frame == 0 {
			return actor.From
		}
		return actor.To
	}
	return resolved
}

func (m model) viewport() mapViewport {
	if m.state.Map.Width <= 0 || m.state.Map.Height <= 0 {
		return mapViewport{}
	}

	width := m.state.Map.Width
	if m.windowWidth > 0 {
		width = min(width, max(1, m.windowWidth-sidebarWidth-mapSidebarGap))
	}
	height := m.state.Map.Height
	if m.windowHeight > 0 {
		height = min(height, max(1, m.windowHeight))
	}

	playerPos := m.animatedPosition(0, true, m.state.Player.Pos)
	x := playerPos.X - width/2
	y := playerPos.Y - height/2
	x = max(0, min(m.state.Map.Width-width, x))
	y = max(0, min(m.state.Map.Height-height, y))
	return mapViewport{X: x, Y: y, Width: width, Height: height}
}

func (m model) View() string {
	if m.state.Map.Width == 0 || m.state.Map.Height == 0 {
		return ""
	}
	viewport := m.viewport()
	glyphs := m.glyphs.withASCIIFallback()

	styles := newViewStyles(defaultPalette)
	player := styles.player.Render(glyphs.Player)
	floor := styles.floor.Render(glyphs.Floor)
	wall := styles.wall.Render(glyphs.Wall)
	memoryFloor := styles.memoryFloor.Render(glyphs.Floor)
	memoryWall := styles.memoryWall.Render(glyphs.Wall)
	stairs := styles.stairs.Render(glyphs.Stairs)
	memoryStairs := styles.memoryStairs.Render(glyphs.Stairs)
	potion := styles.potion.Render(glyphs.Potion)
	sword := styles.sword.Render(glyphs.Sword)
	entityGlyphs := make(map[Position]string, len(m.state.Entities))
	monsterArchetypes := make(map[int]MonsterArchetype, len(m.state.Entities))
	for _, entity := range m.state.Entities {
		pos := m.animatedPosition(entity.ID, false, entity.Pos)
		entityGlyphs[pos] = styles.monsterStyle(entity.Archetype, false).Render(glyphs.monsterGlyph(entity.Archetype))
		monsterArchetypes[entity.ID] = entity.Archetype
	}
	itemGlyphs := make(map[Position]string, len(m.state.Items))
	for _, item := range m.state.Items {
		switch item.Type {
		case ItemPotion:
			itemGlyphs[item.Pos] = potion
		case ItemSword:
			itemGlyphs[item.Pos] = sword
		}
	}
	playerPos := m.animatedPosition(0, true, m.state.Player.Pos)
	type trailGlyph struct {
		glyph    string
		isPlayer bool
	}
	trailGlyphs := make(map[Position]trailGlyph, len(m.motion.Actors))
	if m.motion.Frame > 0 && m.motion.Frame < animationFrames-1 {
		for _, actor := range m.motion.Actors {
			archetype := monsterArchetypes[actor.ID]
			glyph := styles.monsterStyle(archetype, true).Render(glyphs.monsterGlyph(archetype))
			if actor.IsPlayer {
				glyph = styles.playerTrail.Render(glyphs.Player)
			}
			trailGlyphs[actor.From] = trailGlyph{glyph: glyph, isPlayer: actor.IsPlayer}
		}
	}

	var sb strings.Builder
	for y := viewport.Y; y < viewport.Y+viewport.Height; y++ {
		for x := viewport.X; x < viewport.X+viewport.Width; x++ {
			visible := m.state.Map.Visible[y][x]
			explored := m.state.Map.Explored[y][x]
			if !explored {
				sb.WriteByte(' ')
			} else if x == playerPos.X && y == playerPos.Y {
				sb.WriteString(player)
			} else if entity, ok := entityGlyphs[Position{X: x, Y: y}]; visible && ok {
				sb.WriteString(entity)
			} else if trail, ok := trailGlyphs[Position{X: x, Y: y}]; ok && (trail.isPlayer || visible) {
				sb.WriteString(trail.glyph)
			} else if item, ok := itemGlyphs[Position{X: x, Y: y}]; visible && ok {
				sb.WriteString(item)
			} else if !visible {
				switch m.state.Map.Tiles[y][x] {
				case TileWall:
					sb.WriteString(memoryWall)
				case TileStairs:
					sb.WriteString(memoryStairs)
				default:
					sb.WriteString(memoryFloor)
				}
			} else {
				switch m.state.Map.Tiles[y][x] {
				case TileWall:
					sb.WriteString(wall)
				case TileStairs:
					sb.WriteString(stairs)
				default:
					sb.WriteString(floor)
				}
			}
		}
		if y < viewport.Y+viewport.Height-1 {
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
	fmt.Fprintf(&sb, "%s %d\n", styles.uiAccent.Render("Depth"), state.Depth)
	if state.Depth < finalDepth {
		fmt.Fprintf(&sb, "%s reach depth %d\n", styles.uiAccent.Render("Goal"), finalDepth)
	} else {
		fmt.Fprintf(&sb, "%s escape via >\n", styles.uiAccent.Render("Goal"))
	}
	fmt.Fprintf(&sb, "%s %d %s\n", styles.uiAccent.Render("Potions"), state.Potions, styles.uiText.Render("(p)"))
	weapon := "fists"
	if state.HasSword {
		weapon = "sword (e)"
	}
	if state.Equipped {
		weapon = "sword"
	}
	fmt.Fprintf(&sb, "%s %s, %d dmg\n", styles.uiAccent.Render("Weapon"), styles.uiText.Render(weapon), state.Player.Damage)
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
	} else if state.Won {
		sb.WriteByte('\n')
		sb.WriteString(styles.health.Render("VICTORY"))
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
