package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/config"
)

// TestLoadValidYAML verifies that Load correctly parses a valid YAML config file.
func TestLoadValidYAML(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
project: my-project
environments:
  - DEV
  - QA
connections:
  - name: primary-db
    type: postgresql
    host: localhost
    port: 5432
    database: mydb
    folder: Databases
engines:
  - postgresql
  - mysql
`
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatalf("setup: write config: %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("Load returned unexpected error: %v", err)
	}

	if cfg.Project != "my-project" {
		t.Errorf("expected Project 'my-project', got %q", cfg.Project)
	}
	if len(cfg.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(cfg.Environments))
	}
	if cfg.Environments[0] != "DEV" {
		t.Errorf("expected first env 'DEV', got %q", cfg.Environments[0])
	}
	if len(cfg.Connections) != 1 {
		t.Errorf("expected 1 connection, got %d", len(cfg.Connections))
	}
	conn := cfg.Connections[0]
	if conn.Name != "primary-db" {
		t.Errorf("expected connection name 'primary-db', got %q", conn.Name)
	}
	if conn.Type != "postgresql" {
		t.Errorf("expected connection type 'postgresql', got %q", conn.Type)
	}
	if conn.Host != "localhost" {
		t.Errorf("expected host 'localhost', got %q", conn.Host)
	}
	if conn.Port != 5432 {
		t.Errorf("expected port 5432, got %d", conn.Port)
	}
	if conn.Database != "mydb" {
		t.Errorf("expected database 'mydb', got %q", conn.Database)
	}
	if conn.Folder != "Databases" {
		t.Errorf("expected folder 'Databases', got %q", conn.Folder)
	}
	if len(cfg.Engines) != 2 {
		t.Errorf("expected 2 engines, got %d", len(cfg.Engines))
	}
}

// TestLoadMissingFile verifies that Load returns an error when the file does not exist.
func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

// TestLoadInvalidYAML verifies that Load returns an error for malformed YAML.
func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	badYAML := `
project: [broken yaml
  - this is not valid
    nested: wrong
`
	cfgPath := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(cfgPath, []byte(badYAML), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := config.Load(cfgPath)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
