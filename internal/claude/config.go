// Package claude provides helpers to read and write Claude Code configuration files.
package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GlobalPath returns the path to the global Claude Code config (~/.claude.json).
func GlobalPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude.json")
}

// LocalPath returns the path to the local Claude Code settings (./.claude/settings.json).
func LocalPath() string {
	return filepath.Join(".claude", "settings.json")
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

// Register merges the dwm MCP entry into the Claude Code config at path.
// It preserves all existing keys in the file.
// Returns an error if the file contains malformed JSON (file is not modified).
func Register(path string) error {
	data, err := Load(path)
	if err != nil {
		return err
	}

	// Parse existing mcpServers section (if any).
	mcpServers := map[string]json.RawMessage{}
	if existing, ok := data["mcpServers"]; ok {
		_ = json.Unmarshal(existing, &mcpServers)
	}

	// Build the dwm MCP entry.
	entry, err := json.Marshal(map[string]any{
		"command": "dwm",
		"args":    []string{"mcp"},
	})
	if err != nil {
		return err
	}
	mcpServers["dwm"] = json.RawMessage(entry)

	serversBytes, err := json.Marshal(mcpServers)
	if err != nil {
		return err
	}
	data["mcpServers"] = json.RawMessage(serversBytes)

	return Save(path, data)
}
