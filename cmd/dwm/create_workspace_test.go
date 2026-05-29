package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreateWorkspaceCommandRegistered verifies the create-workspace command is registered.
func TestCreateWorkspaceCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "create-workspace <name>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("create-workspace command not registered on rootCmd")
	}
}

// TestCreateWorkspaceCreatesStructure verifies that running create-workspace produces
// the expected DBeaver folder structure in a temp directory.
func TestCreateWorkspaceCreatesStructure(t *testing.T) {
	tmpDir := t.TempDir()

	rootCmd.SetArgs([]string{"create-workspace", "myproject", "--output", tmpDir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-workspace failed: %v", err)
	}

	expectedPaths := []string{
		filepath.Join(tmpDir, "myproject", ".dbeaver", ".gitkeep"),
		filepath.Join(tmpDir, "myproject", "Scripts", ".gitkeep"),
		filepath.Join(tmpDir, "myproject", "Diagrams", ".gitkeep"),
		filepath.Join(tmpDir, "myproject", "Bookmarks", ".gitkeep"),
	}
	for _, p := range expectedPaths {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist: %v", p, err)
		}
	}
}

// TestCreateWorkspaceDryRun verifies that --dry-run does NOT write any files.
func TestCreateWorkspaceDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	rootCmd.SetArgs([]string{"create-workspace", "ws-dry", "--output", tmpDir, "--dry-run"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-workspace --dry-run failed: %v", err)
	}

	// Nothing should have been written.
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries, expected 0", len(entries))
	}
}

// TestCreateWorkspaceCustomEnvironments verifies --environments flag is accepted
// and the workspace is still written (not dry-run).
func TestCreateWorkspaceCustomEnvironments(t *testing.T) {
	tmpDir := t.TempDir()

	// Reset persistent flags to default values before this test to avoid state
	// leaking from TestCreateWorkspaceDryRun (which sets --dry-run).
	if err := rootCmd.PersistentFlags().Set("dry-run", "false"); err != nil {
		t.Fatalf("reset dry-run: %v", err)
	}

	rootCmd.SetArgs([]string{
		"create-workspace", "envws",
		"--output", tmpDir,
		"--environments", "STAGING,UAT",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-workspace --environments failed: %v", err)
	}

	// Workspace structure should still be created (environments flag is accepted and passed through).
	if _, err := os.Stat(filepath.Join(tmpDir, "envws", ".dbeaver", ".gitkeep")); err != nil {
		t.Errorf("expected .dbeaver/.gitkeep to exist: %v", err)
	}
}
