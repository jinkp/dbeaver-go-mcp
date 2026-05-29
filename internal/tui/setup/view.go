package setup

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	pathStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	switch m.screen {
	case ScreenScope:
		return m.viewScope()
	case ScreenConfirm:
		return m.viewConfirm()
	case ScreenDone:
		return m.viewDone()
	}
	return ""
}

func (m Model) viewScope() string {
	clientName := clientDisplayName(m.client)
	s := titleStyle.Render(fmt.Sprintf("Register dwm MCP server in %s", clientName)) + "\n\n"

	options := []struct {
		label string
		path  func() string
	}{
		{"Global", m.globalPath},
		{"Local ", m.localPath},
	}

	for i, opt := range options {
		cursor := "  "
		radio := "○"
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			radio = "●"
			style = selectedStyle
		}
		s += fmt.Sprintf("%s%s %s  %s\n",
			cursor,
			radio,
			style.Render(opt.label),
			pathStyle.Render(opt.path()),
		)
	}

	s += "\n" + footerStyle.Render("↑/↓ · Enter to select · q to quit")
	return s
}

func (m Model) viewConfirm() string {
	clientName := clientDisplayName(m.client)
	s := titleStyle.Render(fmt.Sprintf("Register dwm in %s", clientName)) + "\n\n"
	s += "Target: " + pathStyle.Render(m.targetPath) + "\n"
	s += normalStyle.Render("Existing keys will be preserved.") + "\n\n"
	s += footerStyle.Render("Enter to confirm · Esc to go back · q to quit")
	return s
}

func (m Model) viewDone() string {
	if m.err != nil {
		s := errorStyle.Render("✗ Registration failed:") + "\n"
		s += "  " + m.err.Error() + "\n\n"
		s += footerStyle.Render("Press any key to exit.")
		return s
	}
	s := successStyle.Render("✓ dwm registered in "+m.targetPath) + "\n"
	s += normalStyle.Render("Ensure 'dwm' is on your PATH.") + "\n\n"
	s += footerStyle.Render("Press any key to exit.")
	return s
}

func clientDisplayName(client string) string {
	switch client {
	case "opencode":
		return "OpenCode"
	case "claude":
		return "Claude Code"
	default:
		return client
	}
}
