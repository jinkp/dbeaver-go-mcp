package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCreateConnectionsCommandRegistered verifies create-connections is registered.
func TestCreateConnectionsCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "create-connections <config.yaml>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("create-connections command not registered on rootCmd")
	}
}

// TestCreateConnectionsGeneratesFiles verifies that create-connections parses a YAML
// config and writes data-sources.json + data-sources-config.json to the output dir.
func TestCreateConnectionsGeneratesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Reset persistent flags.
	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-connections failed: %v", err)
	}

	dsPath := filepath.Join(tmpDir, "data-sources.json")
	if _, err := os.Stat(dsPath); err != nil {
		t.Fatalf("data-sources.json not found: %v", err)
	}

	dscPath := filepath.Join(tmpDir, "data-sources-config.json")
	if _, err := os.Stat(dscPath); err != nil {
		t.Fatalf("data-sources-config.json not found: %v", err)
	}

	// Parse data-sources.json and verify no passwords are saved.
	data, err := os.ReadFile(dsPath)
	if err != nil {
		t.Fatalf("read data-sources.json: %v", err)
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse data-sources.json: %v", err)
	}

	conns, ok := doc["connections"].(map[string]interface{})
	if !ok {
		t.Fatal("data-sources.json missing 'connections' map")
	}
	if len(conns) == 0 {
		t.Fatal("data-sources.json has no connections")
	}
	for id, raw := range conns {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("connection %s is not an object", id)
		}
		if savePass, _ := entry["save-password"].(bool); savePass {
			t.Errorf("connection %s has save-password=true, expected false", id)
		}
	}
}

// TestCreateConnectionsFullPipelineFields verifies the full create-connections
// pipeline generates valid JSON with all DBeaver 26.0 required fields.
func TestCreateConnectionsFullPipelineFields(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-connections failed: %v", err)
	}

	dsPath := filepath.Join(tmpDir, "data-sources.json")
	data, err := os.ReadFile(dsPath)
	if err != nil {
		t.Fatalf("read data-sources.json: %v", err)
	}

	// Verify top-level JSON is valid.
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("data-sources.json is not valid JSON: %v", err)
	}

	// Verify required top-level keys: connections, folders.
	requiredKeys := []string{"connections", "folders"}
	for _, key := range requiredKeys {
		if _, ok := doc[key]; !ok {
			t.Errorf("data-sources.json missing required key %q", key)
		}
	}

	// Verify each connection has required DBeaver fields and save-password=false.
	conns, ok := doc["connections"].(map[string]interface{})
	if !ok || len(conns) == 0 {
		t.Fatal("data-sources.json has no connections")
	}
	requiredConnFields := []string{"provider", "driver", "name", "save-password", "folder", "configuration"}
	for id, raw := range conns {
		conn, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("connection %s is not an object", id)
		}
		for _, field := range requiredConnFields {
			if _, ok := conn[field]; !ok {
				t.Errorf("connection %s missing required field %q", id, field)
			}
		}
		if savePass, _ := conn["save-password"].(bool); savePass {
			t.Errorf("connection %s: save-password must be false", id)
		}
	}
}

// TestCreateConnectionsSplitByEnv verifies that --split-by-env produces one
// data-sources-<env>.json file per environment in addition to the global data-sources.json.
func TestCreateConnectionsSplitByEnv(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")
	_ = rootCmd.PersistentFlags().Set("merge", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--split-by-env",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-connections --split-by-env failed: %v", err)
	}

	// Global merged file must still exist.
	if _, err := os.Stat(filepath.Join(tmpDir, "data-sources.json")); err != nil {
		t.Errorf("data-sources.json not found: %v", err)
	}
	// Per-env files must exist for each environment in test-config.yaml (DEV, PROD).
	for _, env := range []string{"dev", "prod"} {
		envFile := filepath.Join(tmpDir, "data-sources-"+env+".json")
		if _, err := os.Stat(envFile); err != nil {
			t.Errorf("expected per-env file %q to exist: %v", "data-sources-"+env+".json", err)
		}
	}
}

// TestCreateConnectionsSplitByEnvContent verifies that each per-env file contains
// only connections whose folder matches that environment.
func TestCreateConnectionsSplitByEnvContent(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")
	_ = rootCmd.PersistentFlags().Set("merge", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--split-by-env",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-connections --split-by-env failed: %v", err)
	}

	// Each per-env file should only have connections whose folder = that env.
	for _, env := range []string{"DEV", "PROD"} {
		envFile := filepath.Join(tmpDir, "data-sources-"+strings.ToLower(env)+".json")
		data, err := os.ReadFile(envFile)
		if err != nil {
			t.Fatalf("read %s: %v", envFile, err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatalf("parse %s: %v", envFile, err)
		}
		conns, ok := doc["connections"].(map[string]interface{})
		if !ok || len(conns) == 0 {
			t.Fatalf("per-env file %s has no connections", envFile)
		}
		for id, raw := range conns {
			conn, ok := raw.(map[string]interface{})
			if !ok {
				t.Fatalf("connection %s in %s is not an object", id, envFile)
			}
			folder, _ := conn["folder"].(string)
			if folder != env {
				t.Errorf("connection %s in %s has folder=%q, want %q", id, envFile, folder, env)
			}
		}
	}
}

// TestCreateConnectionsDryRun verifies --dry-run writes no files.
func TestCreateConnectionsDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-connections --dry-run failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d files, expected 0", len(entries))
	}
}
