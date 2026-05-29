package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSyncSnippetsCommandRegistered verifies the sync-snippets command is registered.
func TestSyncSnippetsCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "sync-snippets <config.yaml>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("sync-snippets command not registered on rootCmd")
	}
}

// TestSyncSnippetsCreatesFiles verifies sync-snippets writes snippet SQL files.
func TestSyncSnippetsCreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"sync-snippets", "testdata/test-config.yaml",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("sync-snippets failed: %v", err)
	}

	// test-config.yaml has engines: [postgresql]. Expect list-tables.sql.
	expected := filepath.Join(tmpDir, "Scripts", "snippets", "postgresql", "list-tables.sql")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected %s to exist: %v", expected, err)
	}
}

// TestSyncSnippetsDryRun verifies --dry-run writes no files.
func TestSyncSnippetsDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"sync-snippets", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("sync-snippets --dry-run failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries, expected 0", len(entries))
	}
}

// TestSyncSnippetsEnginesOverride verifies --engines overrides config engines.
func TestSyncSnippetsEnginesOverride(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"sync-snippets", "testdata/test-config.yaml",
		"--output", tmpDir,
		"--engines", "mysql",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("sync-snippets --engines failed: %v", err)
	}

	// Should produce mysql snippets, not postgresql.
	mysqlSnippet := filepath.Join(tmpDir, "Scripts", "snippets", "mysql", "list-tables.sql")
	if _, err := os.Stat(mysqlSnippet); err != nil {
		t.Errorf("expected mysql list-tables.sql to exist: %v", err)
	}
	// postgresql snippet should NOT exist because --engines=mysql.
	pgSnippet := filepath.Join(tmpDir, "Scripts", "snippets", "postgresql", "list-tables.sql")
	if _, err := os.Stat(pgSnippet); err == nil {
		t.Errorf("postgresql snippet should not exist when --engines=mysql")
	}
}
