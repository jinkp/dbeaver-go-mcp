package claude_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/claude"
)

// TestRegister_CreatesFileWhenAbsent verifies Register creates the config file
// with the correct mcpServers.dwm entry when the file doesn't exist.
func TestRegister_CreatesFileWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".claude.json")

	if err := claude.Register(path); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to be created at %s, got: %v", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse JSON: %v", err)
	}

	serversRaw, ok := m["mcpServers"]
	if !ok {
		t.Fatal("expected 'mcpServers' key in config")
	}

	var servers map[string]json.RawMessage
	if err := json.Unmarshal(serversRaw, &servers); err != nil {
		t.Fatalf("parse mcpServers: %v", err)
	}

	dwmRaw, ok := servers["dwm"]
	if !ok {
		t.Fatal("expected 'mcpServers.dwm' entry")
	}

	var dwmEntry map[string]any
	if err := json.Unmarshal(dwmRaw, &dwmEntry); err != nil {
		t.Fatalf("parse dwm entry: %v", err)
	}

	if dwmEntry["command"] != "dwm" {
		t.Errorf("expected command=dwm, got %v", dwmEntry["command"])
	}
}

// TestRegister_PreservesExistingServers verifies Register preserves existing
// mcpServers entries when adding the dwm entry.
func TestRegister_PreservesExistingServers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".claude.json")

	// Write a file with an existing mcpServer entry
	existing := `{"mcpServers": {"other-tool": {"command": "other", "args": ["run"]}}}`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := claude.Register(path); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	data, _ := os.ReadFile(path)
	var m map[string]json.RawMessage
	json.Unmarshal(data, &m)

	var servers map[string]json.RawMessage
	json.Unmarshal(m["mcpServers"], &servers)

	if _, ok := servers["other-tool"]; !ok {
		t.Error("expected 'other-tool' entry to be preserved")
	}
	if _, ok := servers["dwm"]; !ok {
		t.Error("expected 'dwm' entry to be added")
	}
}

// TestRegister_CreatesParentDirForLocal verifies that Register creates the parent
// directory if it does not exist (e.g., .claude/ dir for local settings.json).
func TestRegister_CreatesParentDirForLocal(t *testing.T) {
	dir := t.TempDir()
	// .claude/ subdir does NOT exist yet
	path := filepath.Join(dir, ".claude", "settings.json")

	if err := claude.Register(path); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s, got: %v", path, err)
	}
}

// TestLoad_ErrorOnMalformedJSON verifies Load returns an error for malformed JSON.
func TestLoad_ErrorOnMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".claude.json")

	if err := os.WriteFile(path, []byte("{bad json"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := claude.Load(path)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

// TestRegister_DoesNotModifyOnMalformedJSON verifies that Register returns
// an error and does not modify a file with malformed JSON.
func TestRegister_DoesNotModifyOnMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".claude.json")

	malformed := []byte("{bad json")
	if err := os.WriteFile(path, malformed, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := claude.Register(path); err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}

	got, _ := os.ReadFile(path)
	if string(got) != string(malformed) {
		t.Errorf("file must not be modified; got: %s", got)
	}
}
