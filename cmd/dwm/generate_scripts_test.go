package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateScriptsCommandRegistered verifies the generate-scripts command is registered.
func TestGenerateScriptsCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "generate-scripts <config.yaml>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("generate-scripts command not registered on rootCmd")
	}
}

// TestGenerateScriptsCreatesFiles verifies that generate-scripts writes SQL files
// to Scripts/scripts/<engine>/<name>.sql inside the output directory.
func TestGenerateScriptsCreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"generate-scripts", "testdata/test-config.yaml",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("generate-scripts failed: %v", err)
	}

	// test-config.yaml has engines: [postgresql]. Expect at least health-check.sql.
	expected := filepath.Join(tmpDir, "Scripts", "scripts", "postgresql", "health-check.sql")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected %s to exist: %v", expected, err)
	}
}

// TestGenerateScriptsDryRun verifies --dry-run writes no files.
func TestGenerateScriptsDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"generate-scripts", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("generate-scripts --dry-run failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries, expected 0", len(entries))
	}
}

// TestGenerateScriptsSubset verifies --scripts flag limits generation.
func TestGenerateScriptsSubset(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"generate-scripts", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--scripts", "health-check",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("generate-scripts --scripts failed: %v", err)
	}

	// Only health-check.sql should exist for postgresql.
	expected := filepath.Join(tmpDir, "Scripts", "scripts", "postgresql", "health-check.sql")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected health-check.sql to exist: %v", err)
	}

	// deadlocks-check.sql should NOT exist (not in --scripts subset).
	unexpected := filepath.Join(tmpDir, "Scripts", "scripts", "postgresql", "deadlocks-check.sql")
	if _, err := os.Stat(unexpected); err == nil {
		t.Errorf("deadlocks-check.sql should not exist when --scripts=health-check")
	}
}
