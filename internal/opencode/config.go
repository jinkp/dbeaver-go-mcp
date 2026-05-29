// Package opencode provides helpers to read and write opencode.json configuration files.
package opencode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// GlobalPath returns the platform-aware path to the global opencode.json.
// Windows: %APPDATA%\opencode\opencode.json
// Linux/macOS: ~/.config/opencode/opencode.json
func GlobalPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "opencode", "opencode.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

// LocalPath returns the path to the local opencode.json in the current directory.
func LocalPath() string {
	return "opencode.json"
}

// Load reads the JSON file at path into a raw map.
// Returns an empty map if the file does not exist.
// Returns an error if the file exists but contains malformed JSON.
func Load(path string) (map[string]json.RawMessage, error) {
	data := map[string]json.RawMessage{}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// Save writes data to path as indented JSON, creating parent directories as needed.
func Save(path string, data map[string]json.RawMessage) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

// Register merges the dwm MCP entry into the opencode.json at path.
// It preserves all existing keys in the file.
// Returns an error if the file contains malformed JSON (file is not modified).
func Register(path string) error {
	data, err := Load(path)
	if err != nil {
		return err
	}

	// Parse existing mcp section (if any).
	mcpSection := map[string]json.RawMessage{}
	if existing, ok := data["mcp"]; ok {
		_ = json.Unmarshal(existing, &mcpSection)
	}

	// Build the dwm MCP entry.
	entry, err := json.Marshal(map[string]any{
		"type":    "local",
		"command": "dwm",
		"args":    []string{"mcp"},
	})
	if err != nil {
		return err
	}
	mcpSection["dwm"] = json.RawMessage(entry)

	mcpBytes, err := json.Marshal(mcpSection)
	if err != nil {
		return err
	}
	data["mcp"] = json.RawMessage(mcpBytes)

	return Save(path, data)
}
