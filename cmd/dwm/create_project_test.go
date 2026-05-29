package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreateProjectCommandRegistered verifies the create-project command is registered.
func TestCreateProjectCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "create-project <name>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("create-project command not registered on rootCmd")
	}
}

// TestCreateProjectCreatesNestedStructure verifies that running create-project
// produces the expected nested workspace/project structure.
func TestCreateProjectCreatesNestedStructure(t *testing.T) {
	tmpDir := t.TempDir()

	// Reset persistent flags.
	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"create-project", "MyApp",
		"--workspace", "CorpWS",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-project failed: %v", err)
	}

	// Expect nested paths: <output>/CorpWS/MyApp/<dir>/.gitkeep
	expectedPaths := []string{
		filepath.Join(tmpDir, "CorpWS", "MyApp", ".dbeaver", ".gitkeep"),
		filepath.Join(tmpDir, "CorpWS", "MyApp", "Scripts", ".gitkeep"),
		filepath.Join(tmpDir, "CorpWS", "MyApp", "Diagrams", ".gitkeep"),
		filepath.Join(tmpDir, "CorpWS", "MyApp", "Bookmarks", ".gitkeep"),
	}
	for _, p := range expectedPaths {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected %s to exist: %v", p, err)
		}
	}
}

// TestCreateProjectDryRun verifies --dry-run writes no files.
func TestCreateProjectDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"create-project", "DryApp",
		"--workspace", "DryWS",
		"--output", tmpDir,
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-project --dry-run failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries, expected 0", len(entries))
	}
}
