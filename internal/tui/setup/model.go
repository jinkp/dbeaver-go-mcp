package setup

import tea "github.com/charmbracelet/bubbletea"

// Screen represents the current wizard screen.
type Screen int

const (
	ScreenScope   Screen = iota // 0: global vs local selection
	ScreenConfirm               // 1: show path, ask confirm
	ScreenDone                  // 2: success or error
)

// Model is the Bubbletea model for the dwm setup wizard.
type Model struct {
	screen     Screen
	client     string // "opencode" | "claude" (display only)
	cursor     int    // 0=global, 1=local
	targetPath string // resolved after scope chosen
	err        error  // nil=success on ScreenDone
	quitting   bool

	// Injected for testability — no real FS calls in unit tests
	register   func(path string) error
	globalPath func() string
	localPath  func() string
}

// NewModel creates a new wizard model for the given client.
// register, globalPath, localPath are injected to allow testing without FS.
func NewModel(client string, register func(string) error, globalPath, localPath func() string) Model {
	return Model{
		screen:     ScreenScope,
		client:     client,
		cursor:     0,
		register:   register,
		globalPath: globalPath,
		localPath:  localPath,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }
