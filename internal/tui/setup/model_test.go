package setup

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel() Model {
	return NewModel(
		"opencode",
		func(path string) error { return nil }, // mock register: always success
		func() string { return "/global/path" },
		func() string { return "./local/path" },
	)
}

func key(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func keySpecial(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func TestInitialScreen(t *testing.T) {
	m := newTestModel()
	if m.screen != ScreenScope {
		t.Errorf("initial screen = %v, want ScreenScope", m.screen)
	}
	if m.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", m.cursor)
	}
}

func TestCursorDown(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(key("j"))
	m = next.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor after j = %d, want 1", m.cursor)
	}
}

func TestCursorDownAtBottom(t *testing.T) {
	m := newTestModel()
	m.cursor = 1
	next, _ := m.Update(key("j"))
	m = next.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor stays at 1, got %d", m.cursor)
	}
}

func TestCursorUp(t *testing.T) {
	m := newTestModel()
	m.cursor = 1
	next, _ := m.Update(key("k"))
	m = next.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor after k = %d, want 0", m.cursor)
	}
}

func TestCursorUpAtTop(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(key("k"))
	m = next.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor stays at 0, got %d", m.cursor)
	}
}

func TestEnterOnScopeGlobalTransitionsToConfirm(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(keySpecial(tea.KeyEnter))
	m = next.(Model)
	if m.screen != ScreenConfirm {
		t.Errorf("screen = %v, want ScreenConfirm", m.screen)
	}
	if m.targetPath != "/global/path" {
		t.Errorf("targetPath = %q, want /global/path", m.targetPath)
	}
}

func TestEnterOnScopeLocalTransitionsToConfirm(t *testing.T) {
	m := newTestModel()
	m.cursor = 1
	next, _ := m.Update(keySpecial(tea.KeyEnter))
	m = next.(Model)
	if m.screen != ScreenConfirm {
		t.Errorf("screen = %v, want ScreenConfirm", m.screen)
	}
	if m.targetPath != "./local/path" {
		t.Errorf("targetPath = %q, want ./local/path", m.targetPath)
	}
}

func TestEscFromConfirmGoesBack(t *testing.T) {
	m := newTestModel()
	m.screen = ScreenConfirm
	m.targetPath = "/global/path"
	next, _ := m.Update(keySpecial(tea.KeyEsc))
	m = next.(Model)
	if m.screen != ScreenScope {
		t.Errorf("screen = %v, want ScreenScope", m.screen)
	}
}

func TestEnterOnConfirmCallsRegisterAndTransitionsToDone(t *testing.T) {
	called := false
	m := NewModel("opencode",
		func(path string) error { called = true; return nil },
		func() string { return "/global" },
		func() string { return "./local" },
	)
	m.screen = ScreenConfirm
	m.targetPath = "/global"
	next, _ := m.Update(keySpecial(tea.KeyEnter))
	m = next.(Model)
	if !called {
		t.Error("register func not called")
	}
	if m.screen != ScreenDone {
		t.Errorf("screen = %v, want ScreenDone", m.screen)
	}
	if m.err != nil {
		t.Errorf("err should be nil on success, got %v", m.err)
	}
}

func TestEnterOnConfirmRegisterError(t *testing.T) {
	boom := errors.New("permission denied")
	m := NewModel("opencode",
		func(path string) error { return boom },
		func() string { return "/global" },
		func() string { return "./local" },
	)
	m.screen = ScreenConfirm
	m.targetPath = "/global"
	next, _ := m.Update(keySpecial(tea.KeyEnter))
	m = next.(Model)
	if m.screen != ScreenDone {
		t.Errorf("screen = %v, want ScreenDone", m.screen)
	}
	if m.err != boom {
		t.Errorf("err = %v, want %v", m.err, boom)
	}
}

// TRIANGULATE: arrow key variants exercise a separate code path from j/k
func TestCursorDownArrowKey(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(keySpecial(tea.KeyDown))
	m = next.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor after ↓ = %d, want 1", m.cursor)
	}
}

func TestCursorUpArrowKey(t *testing.T) {
	m := newTestModel()
	m.cursor = 1
	next, _ := m.Update(keySpecial(tea.KeyUp))
	m = next.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor after ↑ = %d, want 0", m.cursor)
	}
}

// TRIANGULATE: ctrl+c also quits (separate key path)
func TestCtrlCQuits(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(keySpecial(tea.KeyCtrlC))
	m = next.(Model)
	if !m.quitting {
		t.Error("quitting should be true after ctrl+c")
	}
}

func TestQuitSetsQuitting(t *testing.T) {
	m := newTestModel()
	next, _ := m.Update(key("q"))
	m = next.(Model)
	if !m.quitting {
		t.Error("quitting should be true after q")
	}
}

func TestAnyKeyOnDoneQuits(t *testing.T) {
	m := newTestModel()
	m.screen = ScreenDone
	next, _ := m.Update(key("x"))
	m = next.(Model)
	if !m.quitting {
		t.Error("quitting should be true after any key on ScreenDone")
	}
}
