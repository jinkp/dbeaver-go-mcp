package opencode_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/opencode"
)

// TestRegister_CreatesFileWhenAbsent verifies that Register creates the config file
// when it does not exist yet, with the correct mcp.dwm entry.
func TestRegister_CreatesFileWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	if err := opencode.Register(path); err != nil {
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

	mcpRaw, ok := m["mcp"]
	if !ok {
		t.Fatal("expected 'mcp' key in config")
	}

	var mcpSection map[string]json.RawMessage
	if err := json.Unmarshal(mcpRaw, &mcpSection); err != nil {
		t.Fatalf("parse mcp section: %v", err)
	}

	dwmRaw, ok := mcpSection["dwm"]
	if !ok {
		t.Fatal("expected 'mcp.dwm' entry")
	}

	var dwmEntry map[string]any
	if err := json.Unmarshal(dwmRaw, &dwmEntry); err != nil {
		t.Fatalf("parse dwm entry: %v", err)
	}

	if dwmEntry["type"] != "local" {
		t.Errorf("expected type=local, got %v", dwmEntry["type"])
	}
	if dwmEntry["command"] != "dwm" {
		t.Errorf("expected command=dwm, got %v", dwmEntry["command"])
	}
}

// TestRegister_PreservesExistingKeys verifies that Register preserves other
// keys in the config file when adding the mcp.dwm entry.
func TestRegister_PreservesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	// Write a file with existing keys
	existing := `{"theme": "dark", "keybinds": {"submit": "enter"}}`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := opencode.Register(path); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse JSON: %v", err)
	}

	// Existing keys must still be present
	if _, ok := m["theme"]; !ok {
		t.Error("expected 'theme' key to be preserved")
	}
	if _, ok := m["keybinds"]; !ok {
		t.Error("expected 'keybinds' key to be preserved")
	}

	// New mcp entry must exist
	if _, ok := m["mcp"]; !ok {
		t.Error("expected 'mcp' key to be added")
	}
}

// TestLoad_ErrorOnMalformedJSON verifies that Load returns an error
// when the file contains malformed JSON.
func TestLoad_ErrorOnMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	if err := os.WriteFile(path, []byte("{not valid json"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := opencode.Load(path)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

// TestRegister_DoesNotModifyFileOnLoadError verifies that Register does not
// overwrite a file with malformed JSON (it returns an error instead).
func TestRegister_DoesNotModifyFileOnLoadError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	malformedContent := []byte("{not valid json")
	if err := os.WriteFile(path, malformedContent, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := opencode.Register(path)
	if err == nil {
		t.Fatal("expected Register to return error on malformed JSON, got nil")
	}

	// File must be untouched
	got, _ := os.ReadFile(path)
	if string(got) != string(malformedContent) {
		t.Errorf("Register must not modify file on load error; got: %s", got)
	}
}

// TestRegister_Idempotent verifies that calling Register twice does not
// duplicate entries.
func TestRegister_Idempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	if err := opencode.Register(path); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	if err := opencode.Register(path); err != nil {
		t.Fatalf("second Register: %v", err)
	}

	data, _ := os.ReadFile(path)
	var m map[string]json.RawMessage
	json.Unmarshal(data, &m)

	var mcpSection map[string]json.RawMessage
	json.Unmarshal(m["mcp"], &mcpSection)

	if len(mcpSection) != 1 {
		t.Errorf("expected exactly 1 entry in mcp section, got %d", len(mcpSection))
	}
}
