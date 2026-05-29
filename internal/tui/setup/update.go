package setup

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}

		switch m.screen {
		case ScreenScope:
			return m.updateScope(msg)
		case ScreenConfirm:
			return m.updateConfirm(msg)
		case ScreenDone:
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateScope(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 1 {
			m.cursor++
		}
	case "enter":
		if m.cursor == 0 {
			m.targetPath = m.globalPath()
		} else {
			m.targetPath = m.localPath()
		}
		m.screen = ScreenConfirm
	}
	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.err = m.register(m.targetPath)
		m.screen = ScreenDone
	case "esc":
		m.screen = ScreenScope
	}
	return m, nil
}
