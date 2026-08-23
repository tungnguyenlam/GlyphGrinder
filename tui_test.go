package main

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"glyphgrinder/internal/tuitest"
)

const tuiTestSeed int64 = 42

func newTUITestModel() model {
	m := initialModel(rand.New(rand.NewSource(tuiTestSeed)))
	m.glyphs = asciiGlyphs
	return m
}

// playerAt reports the grid position of the player glyph in a rendered view.
// Styling is stripped by looking for the raw player rune, which is safe as
// long as no other glyph reuses it (see AGENTS.md gotchas).
func playerAt(t *testing.T, lines []string) Position {
	t.Helper()
	for y, line := range lines {
		// Cells are styled, so count runes of the unstyled grid instead by
		// stripping ANSI escapes.
		plain := stripANSI(line)
		if x := strings.IndexRune(plain, '@'); x >= 0 {
			return Position{X: x, Y: y}
		}
	}
	t.Fatalf("player glyph not found in view:\n%s", strings.Join(lines, "\n"))
	return Position{}
}

func stripANSI(s string) string {
	var sb strings.Builder
	inEscape := false
	for _, r := range s {
		switch {
		case inEscape:
			if r == 'm' {
				inEscape = false
			}
		case r == '\x1b':
			inEscape = true
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func TestViewRendersFullGrid(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	state := d.Model().(model).state
	if state.Map.Width <= 80 || state.Map.Height <= 24 {
		t.Fatalf("production dungeon is %dx%d, want larger than an 80x24 terminal", state.Map.Width, state.Map.Height)
	}

	lines := d.Lines()
	if got, want := len(lines), state.Map.Height; got != want {
		t.Fatalf("view has %d rows, want %d", got, want)
	}
	for y, line := range lines {
		plain := []rune(stripANSI(line))
		if got, wantAtLeast := len(plain), state.Map.Width; got < wantAtLeast {
			t.Errorf("row %d has %d cells, want at least %d map cells", y, got, wantAtLeast)
			continue
		}
		for x, cell := range plain[:state.Map.Width] {
			if !strings.ContainsRune(" #.@g>", cell) {
				t.Errorf("map cell (%d,%d) = %q, want a dungeon glyph", x, y, cell)
			}
		}
	}
	if got, want := playerAt(t, lines), state.Player.Pos; got != want {
		t.Errorf("player rendered at %+v, want %+v", got, want)
	}
}

func TestResizeClipsViewportAndPreservesSidebar(t *testing.T) {
	state := openTestStateSized(60, 30)
	state.Player.Pos = Position{X: 30, Y: 15}
	state.Entities = nil
	state = state.refreshVisibility()
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Resize(50, 9)
	gotModel := d.Model().(model)
	if gotModel.windowWidth != 50 || gotModel.windowHeight != 9 {
		t.Fatalf("stored window size = %dx%d, want 50x9", gotModel.windowWidth, gotModel.windowHeight)
	}
	if got, want := gotModel.viewport(), (mapViewport{X: 22, Y: 11, Width: 16, Height: 9}); got != want {
		t.Fatalf("viewport = %+v, want %+v", got, want)
	}

	lines := d.Lines()
	if got, want := len(lines), 9; got != want {
		t.Fatalf("resized view has %d rows, want %d", got, want)
	}
	if got, want := playerAt(t, lines), (Position{X: 8, Y: 4}); got != want {
		t.Errorf("player rendered at %+v, want viewport position %+v", got, want)
	}
	plain := stripANSI(d.View())
	if got, want := strings.Index(strings.Split(plain, "\n")[0], "HP "), 18; got != want {
		t.Errorf("sidebar starts at column %d, want %d after 16 map cells and 2 spaces", got, want)
	}
	if !strings.Contains(plain, "Log:") || !strings.Contains(plain, "The fight begins.") {
		t.Errorf("resized view lost sidebar content:\n%s", plain)
	}
}

func TestViewportFollowsPlayerAndClampsAtMapEdges(t *testing.T) {
	centered := openTestStateSized(60, 30)
	centered.Player.Pos = Position{X: 30, Y: 15}
	centered.Entities = nil
	centered = centered.refreshVisibility()
	d := tuitest.New(t, model{state: centered, rng: rand.New(rand.NewSource(tuiTestSeed))})
	d.Resize(50, 9)

	d.Key("right")
	if got, want := d.Model().(model).state.Player.Pos, (Position{X: 31, Y: 15}); got != want {
		t.Fatalf("player world position after movement = %+v, want %+v", got, want)
	}
	if got, want := playerAt(t, d.Lines()), (Position{X: 8, Y: 4}); got != want {
		t.Errorf("camera did not follow player: rendered at %+v, want %+v", got, want)
	}

	tests := []struct {
		name       string
		world      Position
		wantScreen Position
		wantOrigin Position
	}{
		{name: "top left", world: Position{X: 1, Y: 1}, wantScreen: Position{X: 1, Y: 1}, wantOrigin: Position{}},
		{name: "bottom right", world: Position{X: 58, Y: 28}, wantScreen: Position{X: 14, Y: 7}, wantOrigin: Position{X: 44, Y: 21}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := openTestStateSized(60, 30)
			state.Player.Pos = tc.world
			state.Entities = nil
			state = state.refreshVisibility()
			d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})
			d.Resize(50, 9)

			viewport := d.Model().(model).viewport()
			if got := (Position{X: viewport.X, Y: viewport.Y}); got != tc.wantOrigin {
				t.Errorf("viewport origin = %+v, want %+v", got, tc.wantOrigin)
			}
			if got := playerAt(t, d.Lines()); got != tc.wantScreen {
				t.Errorf("player rendered at %+v, want %+v", got, tc.wantScreen)
			}
		})
	}
}

func TestTinyResizeKeepsViewportValid(t *testing.T) {
	state := openTestStateSized(60, 30)
	state.Player.Pos = Position{X: 30, Y: 15}
	state.Entities = nil
	state = state.refreshVisibility()
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Resize(1, 1)
	if got, want := d.Model().(model).viewport(), (mapViewport{X: 30, Y: 15, Width: 1, Height: 1}); got != want {
		t.Fatalf("tiny viewport = %+v, want %+v", got, want)
	}
	if got, want := playerAt(t, d.Lines()), (Position{}); got != want {
		t.Errorf("player rendered at %+v in tiny viewport, want %+v", got, want)
	}
}

func TestViewRendersOnlyVisibleMonsters(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	state := d.Model().(model).state
	lines := d.Lines()

	for _, monster := range state.Entities {
		plain := []rune(stripANSI(lines[monster.Pos.Y]))
		got := string(plain[monster.Pos.X])
		if state.Map.Visible[monster.Pos.Y][monster.Pos.X] && got != asciiGlyphs.Monster {
			t.Errorf("visible monster %d rendered as %q at %+v, want %q", monster.ID, got, monster.Pos, asciiGlyphs.Monster)
		}
		if !state.Map.Visible[monster.Pos.Y][monster.Pos.X] && got == asciiGlyphs.Monster {
			t.Errorf("hidden monster %d rendered at %+v", monster.ID, monster.Pos)
		}
	}
}

func TestGlyphProfileSelection(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want glyphProfile
	}{
		{name: "UTF-8 locale", env: map[string]string{"LANG": "en_US.UTF-8"}, want: richGlyphs},
		{name: "compact UTF8 locale", env: map[string]string{"LC_CTYPE": "C.UTF8"}, want: richGlyphs},
		{name: "plain locale", env: map[string]string{"LANG": "C"}, want: asciiGlyphs},
		{name: "missing locale", env: map[string]string{}, want: asciiGlyphs},
		{name: "dumb terminal overrides Unicode", env: map[string]string{"TERM": "dumb", "LANG": "en_US.UTF-8"}, want: asciiGlyphs},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getenv := func(key string) string { return tc.env[key] }
			if got := glyphProfileForEnvironment(getenv); got != tc.want {
				t.Errorf("profile = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestGlyphProfilesStayOneCellWide(t *testing.T) {
	for name, profile := range map[string]glyphProfile{"ASCII": asciiGlyphs, "rich": richGlyphs} {
		for role, glyph := range map[string]string{
			"player":  profile.Player,
			"monster": profile.Monster,
			"floor":   profile.Floor,
			"wall":    profile.Wall,
			"stairs":  profile.Stairs,
		} {
			if got := lipgloss.Width(glyph); got != 1 {
				t.Errorf("%s %s glyph %q is %d cells wide, want 1", name, role, glyph, got)
			}
		}
	}
}

func TestIncompleteGlyphProfileFallsBackToASCII(t *testing.T) {
	if got := (glyphProfile{Player: "@"}).withASCIIFallback(); got != asciiGlyphs {
		t.Errorf("incomplete profile resolved to %+v, want ASCII %+v", got, asciiGlyphs)
	}
}

func TestViewRendersRichGlyphSemantics(t *testing.T) {
	state := openTestState()
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Entities = []Entity{testMonster(1, Position{X: 4, Y: 3})}
	state.Map.Tiles[2][3] = TileStairs
	state = state.refreshVisibility()
	d := tuitest.New(t, model{
		state:  state,
		rng:    rand.New(rand.NewSource(tuiTestSeed)),
		glyphs: richGlyphs,
	})
	lines := d.Lines()

	tests := []struct {
		name string
		pos  Position
		want string
	}{
		{name: "player", pos: Position{X: 3, Y: 3}, want: richGlyphs.Player},
		{name: "monster", pos: Position{X: 4, Y: 3}, want: richGlyphs.Monster},
		{name: "floor", pos: Position{X: 2, Y: 3}, want: richGlyphs.Floor},
		{name: "wall", pos: Position{X: 0, Y: 3}, want: richGlyphs.Wall},
		{name: "stairs", pos: Position{X: 3, Y: 2}, want: richGlyphs.Stairs},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := renderedCell(t, lines, tc.pos); got != tc.want {
				t.Errorf("rendered glyph = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestViewRendersASCIIStairs(t *testing.T) {
	state := openTestState()
	stairs := Position{X: 3, Y: 2}
	state.Map.Tiles[stairs.Y][stairs.X] = TileStairs
	state = state.refreshVisibility()
	d := tuitest.New(t, model{
		state:  state,
		rng:    rand.New(rand.NewSource(tuiTestSeed)),
		glyphs: asciiGlyphs,
	})

	if got := renderedCell(t, d.Lines(), stairs); got != asciiGlyphs.Stairs {
		t.Errorf("rendered stair glyph = %q, want %q", got, asciiGlyphs.Stairs)
	}
}

func TestViewDistinguishesHiddenAndRememberedTiles(t *testing.T) {
	state := openTestStateSized(11, 7)
	state.Player.Pos = Position{X: 2, Y: 3}
	state.Map.Tiles[3][3] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 4, Y: 3})}
	state.Map.Visible = nil
	state.Map.Explored = nil
	state = state.refreshVisibility()
	remembered := Position{X: 9, Y: 3}
	state.Map.Explored[remembered.Y][remembered.X] = true
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})
	lines := d.Lines()

	if got := renderedCell(t, lines, Position{X: 4, Y: 3}); got != " " {
		t.Errorf("unseen monster tile rendered as %q, want blank", got)
	}
	if got := renderedCell(t, lines, remembered); got != "." {
		t.Errorf("remembered floor rendered as %q, want dim floor glyph", got)
	}
}

func renderedCell(t *testing.T, lines []string, pos Position) string {
	t.Helper()
	if pos.Y < 0 || pos.Y >= len(lines) {
		t.Fatalf("row %d outside rendered view", pos.Y)
	}
	plain := []rune(stripANSI(lines[pos.Y]))
	if pos.X < 0 || pos.X >= len(plain) {
		t.Fatalf("column %d outside rendered row %d", pos.X, pos.Y)
	}
	return string(plain[pos.X])
}

func TestViewRendersHealthAndRecentLog(t *testing.T) {
	state := openTestState()
	state.Player.Health = 65
	state.Log = append(state.Log, "A test event.")
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})
	plain := stripANSI(d.View())

	if !strings.Contains(plain, "HP [######....] 65/100") {
		t.Errorf("view does not contain health bar:\n%s", plain)
	}
	if !strings.Contains(plain, "Depth 1") {
		t.Errorf("view does not contain dungeon depth:\n%s", plain)
	}
	if !strings.Contains(plain, "A test event.") {
		t.Errorf("view does not contain recent log entry:\n%s", plain)
	}
}

func TestPaletteStylesKeepSemanticRolesDistinct(t *testing.T) {
	styles := newViewStyles(defaultPalette)
	semanticStyles := []lipgloss.Style{
		styles.floor,
		styles.memoryFloor,
		styles.player,
		styles.monster,
	}
	for i := range semanticStyles {
		if got := stripANSI(semanticStyles[i].Render(".")); got != "." {
			t.Errorf("style %d changed plain glyph to %q", i, got)
		}
		for j := 0; j < i; j++ {
			sameForeground := reflect.DeepEqual(semanticStyles[i].GetForeground(), semanticStyles[j].GetForeground())
			sameBackground := reflect.DeepEqual(semanticStyles[i].GetBackground(), semanticStyles[j].GetBackground())
			sameBold := semanticStyles[i].GetBold() == semanticStyles[j].GetBold()
			if sameForeground && sameBackground && sameBold {
				t.Errorf("semantic styles %d and %d render identically", i, j)
			}
		}
	}
}

func TestMovementKeys(t *testing.T) {
	cases := []struct {
		name    string
		keys    []string
		actions []Action
	}{
		{"arrow up", []string{"up"}, []Action{ActionMoveUp}},
		{"arrow down", []string{"down"}, []Action{ActionMoveDown}},
		{"arrow left", []string{"left"}, []Action{ActionMoveLeft}},
		{"arrow right", []string{"right"}, []Action{ActionMoveRight}},
		{"wasd", []string{"w", "a"}, []Action{ActionMoveUp, ActionMoveLeft}},
		{"round trip", []string{"w", "s", "a", "d"}, []Action{ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := tuitest.New(t, newTUITestModel())
			want := d.Model().(model).state
			for _, action := range tc.actions {
				want = want.Step(action)
			}
			d.Keys(tc.keys...)
			if got := playerAt(t, d.Lines()); got != want.Player.Pos {
				t.Errorf("after %v player at %+v, want %+v", tc.keys, got, want.Player.Pos)
			}
		})
	}
}

func TestMovementAnimationAdvancesOnlyOnTicks(t *testing.T) {
	state := openTestState()
	state.Entities = nil
	m := model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed)), glyphs: asciiGlyphs}
	from := state.Player.Pos
	to := Position{X: from.X + 1, Y: from.Y}

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = next.(model)
	if cmd == nil {
		t.Fatal("movement did not schedule an animation tick")
	}
	if got := m.state.Player.Pos; got != to {
		t.Fatalf("resolved player position = %+v, want %+v", got, to)
	}
	if got := playerAt(t, strings.Split(m.View(), "\n")); got != from {
		t.Errorf("frame 0 rendered player at %+v, want prior position %+v", got, from)
	}

	next, cmd = m.Update(animationTickMsg{Sequence: m.motion.Sequence})
	m = next.(model)
	if cmd == nil {
		t.Fatal("first animation frame did not schedule the next tick")
	}
	lines := strings.Split(m.View(), "\n")
	if got := renderedCell(t, lines, from); got != asciiGlyphs.Player {
		t.Errorf("frame 1 source trail = %q, want %q", got, asciiGlyphs.Player)
	}
	if got := renderedCell(t, lines, to); got != asciiGlyphs.Player {
		t.Errorf("frame 1 destination = %q, want %q", got, asciiGlyphs.Player)
	}

	next, cmd = m.Update(animationTickMsg{Sequence: m.motion.Sequence})
	m = next.(model)
	if cmd == nil {
		t.Fatal("second animation frame did not schedule the settling tick")
	}
	if got := playerAt(t, strings.Split(m.View(), "\n")); got != to {
		t.Errorf("frame 2 rendered player at %+v, want destination %+v", got, to)
	}

	next, cmd = m.Update(animationTickMsg{Sequence: m.motion.Sequence})
	m = next.(model)
	if cmd != nil || len(m.motion.Actors) != 0 {
		t.Errorf("settled animation = %+v, cmd nil = %v", m.motion, cmd == nil)
	}
}

func TestAnimationTracksMonstersAndCamera(t *testing.T) {
	state := openTestStateSized(60, 30)
	state.Player.Pos = Position{X: 30, Y: 15}
	state.Entities = []Entity{testMonster(1, Position{X: 26, Y: 13})}
	state = state.refreshVisibility()
	m := model{
		state:        state,
		rng:          rand.New(rand.NewSource(tuiTestSeed)),
		glyphs:       asciiGlyphs,
		windowWidth:  50,
		windowHeight: 9,
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = next.(model)
	if got, want := m.viewport().X, 22; got != want {
		t.Errorf("camera origin on frame 0 = %d, want %d", got, want)
	}
	lines := strings.Split(m.View(), "\n")
	if got, want := playerAt(t, lines), (Position{X: 8, Y: 4}); got != want {
		t.Errorf("animated player at %+v, want %+v", got, want)
	}
	if got := renderedCell(t, lines, Position{X: 4, Y: 2}); got != asciiGlyphs.Monster {
		t.Errorf("frame 0 monster glyph = %q, want %q at prior cell", got, asciiGlyphs.Monster)
	}

	next, _ = m.Update(animationTickMsg{Sequence: m.motion.Sequence})
	m = next.(model)
	next, _ = m.Update(animationTickMsg{Sequence: m.motion.Sequence})
	m = next.(model)
	if got, want := m.viewport().X, 23; got != want {
		t.Errorf("camera origin on frame 2 = %d, want %d", got, want)
	}
	lines = strings.Split(m.View(), "\n")
	if got, want := playerAt(t, lines), (Position{X: 8, Y: 4}); got != want {
		t.Errorf("settling player at %+v, want %+v", got, want)
	}
	monster := m.state.Entities[0]
	monsterScreen := Position{X: monster.Pos.X - m.viewport().X, Y: monster.Pos.Y - m.viewport().Y}
	if got := renderedCell(t, lines, monsterScreen); got != asciiGlyphs.Monster {
		t.Errorf("frame 2 monster glyph = %q, want %q at destination", got, asciiGlyphs.Monster)
	}
}

func TestMovementInputReplacesInFlightAnimation(t *testing.T) {
	state := openTestStateSized(9, 7)
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Entities = nil
	state = state.refreshVisibility()
	m := model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed)), glyphs: asciiGlyphs}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = next.(model)
	staleSequence := m.motion.Sequence
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = next.(model)

	if got, want := m.state.Player.Pos, (Position{X: 5, Y: 3}); got != want {
		t.Fatalf("resolved position after rapid input = %+v, want %+v", got, want)
	}
	if got, want := playerAt(t, strings.Split(m.View(), "\n")), (Position{X: 4, Y: 3}); got != want {
		t.Errorf("replacement animation starts at %+v, want last resolved cell %+v", got, want)
	}
	next, cmd := m.Update(animationTickMsg{Sequence: staleSequence})
	m = next.(model)
	if cmd != nil || m.motion.Frame != 0 {
		t.Errorf("stale tick advanced replacement animation to frame %d", m.motion.Frame)
	}
}

func TestDescendKeyRegeneratesAndRendersNextDepth(t *testing.T) {
	rng := rand.New(rand.NewSource(tuiTestSeed))
	state := NewGame(20, 10, rng)
	state.Player.Pos = tilePositions(state.Map, TileStairs)[0]
	state.Player.Health = 63
	state = state.refreshVisibility()
	m := model{
		state:  state,
		rng:    rng,
		glyphs: asciiGlyphs,
		motion: actorMotionState{
			Actors:   []actorMotion{{IsPlayer: true, From: Position{X: 1, Y: 1}, To: state.Player.Pos}},
			Frame:    1,
			Sequence: 4,
		},
		motionSeq:    4,
		windowWidth:  50,
		windowHeight: 9,
	}
	d := tuitest.New(t, m)

	d.Key(">")

	got := d.Model().(model)
	if got.state.Depth != 2 {
		t.Fatalf("depth after > = %d, want 2", got.state.Depth)
	}
	if got.state.Player.Health != 63 {
		t.Errorf("health after > = %d, want 63", got.state.Player.Health)
	}
	if len(got.motion.Actors) != 0 || got.motion.Frame != 0 {
		t.Errorf("descent retained actor/camera motion: %+v", got.motion)
	}
	viewport := got.viewport()
	wantScreen := Position{X: got.state.Player.Pos.X - viewport.X, Y: got.state.Player.Pos.Y - viewport.Y}
	if rendered := playerAt(t, d.Lines()); rendered != wantScreen {
		t.Errorf("next-depth player rendered at %+v, want %+v", rendered, wantScreen)
	}
	plain := stripANSI(d.View())
	if !strings.Contains(plain, "HP [######....] 63/100") || !strings.Contains(plain, "Depth 2") {
		t.Errorf("next-depth view lost preserved health or depth:\n%s", plain)
	}
}

func TestPlayerCannotWalkThroughWalls(t *testing.T) {
	d := tuitest.New(t, newTUITestModel())
	want := d.Model().(model).state

	// Repeated moves must stop wherever the generated dungeon blocks the path.
	for i := 0; i < 20; i++ {
		d.Key("up")
		want = want.Step(ActionMoveUp)
	}
	if got := playerAt(t, d.Lines()); got != want.Player.Pos {
		t.Errorf("player stopped at %+v, want %+v", got, want.Player.Pos)
	}

	for i := 0; i < 20; i++ {
		d.Key("left")
		want = want.Step(ActionMoveLeft)
	}
	if got := playerAt(t, d.Lines()); got != want.Player.Pos {
		t.Errorf("player stopped at %+v, want %+v", got, want.Player.Pos)
	}
}

func TestBumpAttackSurvivesUpdateRoundTrip(t *testing.T) {
	state := combatTestState(20)
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	got := d.Model().(model).state
	if got.Entities[0].Health != 10 {
		t.Errorf("target health after key = %d, want 10", got.Entities[0].Health)
	}
	if got.Player.Pos != state.Player.Pos {
		t.Errorf("player moved to %+v during attack, want %+v", got.Player.Pos, state.Player.Pos)
	}
	if gotLog, want := got.Log, "You hit monster 1 for 10 damage."; len(gotLog) < 2 || gotLog[1] != want {
		t.Errorf("log after key = %q, want player event %q first", gotLog, want)
	}

	d.Key("right")
	got = d.Model().(model).state
	if len(got.Entities) != 1 || got.Entities[0].ID != 2 {
		t.Errorf("entities after second key = %+v, want only monster 2", got.Entities)
	}
}

func TestMonsterAttackSurvivesUpdateRoundTrip(t *testing.T) {
	state := openTestState()
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Map.Tiles[3][4] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	got := d.Model().(model).state
	if got.Player.Health != 95 {
		t.Errorf("player health after key = %d, want 95", got.Player.Health)
	}
	if got.Entities[0].Pos != (Position{X: 3, Y: 2}) {
		t.Errorf("attacking monster moved to %+v, want {3 2}", got.Entities[0].Pos)
	}
	if gotLog, want := got.Log, "Monster 1 hits you for 5 damage."; len(gotLog) != 2 || gotLog[1] != want {
		t.Errorf("log after key = %q, want final entry %q", gotLog, want)
	}
}

func TestGameOverAndRestartSurviveUpdateRoundTrip(t *testing.T) {
	state := openTestState()
	state.Player.Pos = Position{X: 3, Y: 3}
	state.Player.Health = 5
	state.Map.Tiles[3][4] = TileWall
	state.Entities = []Entity{testMonster(1, Position{X: 3, Y: 2})}
	d := tuitest.New(t, model{state: state, rng: rand.New(rand.NewSource(tuiTestSeed))})

	d.Key("right")
	deadModel := d.Model().(model)
	dead := deadModel.state
	if !dead.GameOver {
		t.Fatal("lethal input did not produce game over")
	}
	if len(deadModel.motion.Actors) != 0 {
		t.Errorf("lethal stationary turn started motion: %+v", deadModel.motion)
	}
	plain := stripANSI(d.View())
	if !strings.Contains(plain, "GAME OVER") || !strings.Contains(plain, "Press r to restart") {
		t.Errorf("game-over view is missing restart prompt:\n%s", plain)
	}

	d.Key("r")
	restartedModel := d.Model().(model)
	restarted := restartedModel.state
	if restarted.GameOver {
		t.Fatal("restart left game in game-over state")
	}
	if restarted.Player.Health != restarted.Player.MaxHealth {
		t.Errorf("restarted health = %d/%d, want full", restarted.Player.Health, restarted.Player.MaxHealth)
	}
	if len(restarted.Entities) == 0 {
		t.Error("restarted dungeon has no monsters")
	}
	if len(restartedModel.motion.Actors) != 0 {
		t.Errorf("restart retained actor motion: %+v", restartedModel.motion)
	}
}

func TestQuitKeys(t *testing.T) {
	for _, key := range []string{"q", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			d := tuitest.New(t, newTUITestModel())
			if d.Quit() {
				t.Fatal("model quit before any key was sent")
			}
			d.Key(key)
			if !d.Quit() {
				t.Errorf("%q did not quit the program", key)
			}
		})
	}
}
