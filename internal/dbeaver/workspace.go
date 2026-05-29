package dbeaver

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// WorkspacePath returns the OS-specific DBeaver Community workspace6 directory path.
// Returns an error if the environment variable is missing or the path does not exist.
//
// Platform paths:
//   - Windows: %APPDATA%\DBeaverData\workspace6\
//   - macOS:   ~/Library/DBeaverData/workspace6/
//   - Linux:   $XDG_DATA_HOME/DBeaverData/workspace6/ (fallback: ~/.local/share/DBeaverData/workspace6/)
func WorkspacePath() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", fmt.Errorf("APPDATA env var not set")
		}
		base = filepath.Join(appdata, "DBeaverData", "workspace6")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		base = filepath.Join(home, "Library", "DBeaverData", "workspace6")
	default: // linux
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot determine home directory: %w", err)
			}
			xdg = filepath.Join(home, ".local", "share")
		}
		base = filepath.Join(xdg, "DBeaverData", "workspace6")
	}

	if _, err := os.Stat(base); os.IsNotExist(err) {
		return "", fmt.Errorf("DBeaver workspace not found at %s", base)
	}
	return base, nil
}

// ListProjects returns the names of DBeaver projects (subdirectories containing a .dbeaver/ subfolder)
// found directly under root. Non-project directories are silently skipped.
func ListProjects(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var projects []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dbeaverDir := filepath.Join(root, e.Name(), ".dbeaver")
		if _, err := os.Stat(dbeaverDir); err == nil {
			projects = append(projects, e.Name())
		}
	}
	return projects, nil
}
