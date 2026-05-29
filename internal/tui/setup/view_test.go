package setup

import (
	"errors"
	"strings"
	"testing"
)

func TestViewScopeContainsOptions(t *testing.T) {
	m := newTestModel()
	v := m.View()
	if !strings.Contains(v, "Global") {
		t.Error("scope view should contain 'Global'")
	}
	if !strings.Contains(v, "Local") {
		t.Error("scope view should contain 'Local'")
	}
	if !strings.Contains(v, "/global/path") {
		t.Error("scope view should contain global path")
	}
}

func TestViewConfirmContainsPath(t *testing.T) {
	m := newTestModel()
	m.screen = ScreenConfirm
	m.targetPath = "/some/path"
	v := m.View()
	if !strings.Contains(v, "/some/path") {
		t.Error("confirm view should contain target path")
	}
	if !strings.Contains(v, "Existing keys will be preserved") {
		t.Error("confirm view should mention key preservation")
	}
}

func TestViewDoneSuccess(t *testing.T) {
	m := newTestModel()
	m.screen = ScreenDone
	m.targetPath = "/registered/path"
	v := m.View()
	if !strings.Contains(v, "✓") {
		t.Error("success view should contain ✓")
	}
	if !strings.Contains(v, "/registered/path") {
		t.Error("success view should contain registered path")
	}
	if !strings.Contains(v, "PATH") {
		t.Error("success view should mention PATH")
	}
}

func TestViewDoneError(t *testing.T) {
	m := newTestModel()
	m.screen = ScreenDone
	m.err = errors.New("something went wrong")
	v := m.View()
	if !strings.Contains(v, "✗") {
		t.Error("error view should contain ✗")
	}
	if !strings.Contains(v, "something went wrong") {
		t.Error("error view should contain error message")
	}
}

func TestViewQuittingIsEmpty(t *testing.T) {
	m := newTestModel()
	m.quitting = true
	if m.View() != "" {
		t.Error("quitting view should be empty string")
	}
}
