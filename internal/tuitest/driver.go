// Package tuitest is a headless driver for Bubble Tea models.
//
// It exists so tests can exercise the real Update/View code without a TTY,
// without spawning a process, and without adding a test dependency. Tests
// follow an Observe -> Reason -> Act -> Synchronize loop: render the view,
// assert on it, send a key, render again.
package tuitest

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Driver wraps a tea.Model and replays messages through it in-process.
type Driver struct {
	t     *testing.T
	model tea.Model
	quit  bool
}

// New starts a driver around m and runs its Init command eagerly, so the
// model under test is in the same state it would be after program startup.
func New(t *testing.T, m tea.Model) *Driver {
	t.Helper()
	d := &Driver{t: t, model: m}
	if cmd := m.Init(); cmd != nil {
		d.dispatch(cmd)
	}
	return d
}

// Send delivers one message to Update and drains any command it returns.
func (d *Driver) Send(msg tea.Msg) *Driver {
	d.t.Helper()
	next, cmd := d.model.Update(msg)
	d.model = next
	if cmd != nil {
		d.dispatch(cmd)
	}
	return d
}

// Key sends a key press by its Bubble Tea name ("up", "w", "ctrl+c", ...).
func (d *Driver) Key(name string) *Driver {
	d.t.Helper()
	return d.Send(keyMsg(name))
}

// Keys sends several key presses in order.
func (d *Driver) Keys(names ...string) *Driver {
	d.t.Helper()
	for _, n := range names {
		d.Key(n)
	}
	return d
}

// Resize sends a window size message.
func (d *Driver) Resize(w, h int) *Driver {
	d.t.Helper()
	return d.Send(tea.WindowSizeMsg{Width: w, Height: h})
}

// View renders the current view.
func (d *Driver) View() string { return d.model.View() }

// Lines renders the current view split into rows, which is usually what grid
// assertions want.
func (d *Driver) Lines() []string {
	v := d.View()
	if v == "" {
		return nil
	}
	return strings.Split(v, "\n")
}

// Model returns the current model for type-asserted state assertions.
func (d *Driver) Model() tea.Model { return d.model }

// Quit reports whether a tea.Quit command has been produced so far.
func (d *Driver) Quit() bool { return d.quit }

// dispatch executes a command synchronously. Batched commands are flattened;
// tea.Quit is recorded rather than executed.
func (d *Driver) dispatch(cmd tea.Cmd) {
	msg := cmd()
	switch m := msg.(type) {
	case nil:
		return
	case tea.QuitMsg:
		d.quit = true
	case tea.BatchMsg:
		for _, c := range m {
			if c != nil {
				d.dispatch(c)
			}
		}
	default:
		d.Send(msg)
	}
}

func keyMsg(name string) tea.KeyMsg {
	switch name {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEscape}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)}
	}
}
