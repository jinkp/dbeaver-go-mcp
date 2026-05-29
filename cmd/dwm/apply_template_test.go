package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestApplyTemplateCommandRegistered verifies the apply-template command is registered.
func TestApplyTemplateCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "apply-template <name>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("apply-template command not registered on rootCmd")
	}
}

// TestApplyTemplateBuiltin verifies that applying the built-in "multi-env" template
// produces workspace marker files.
func TestApplyTemplateBuiltin(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"apply-template", "multi-env",
		"--output", tmpDir,
		"--set", "Project=TestCorp",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("apply-template failed: %v", err)
	}

	// The multi-env template generates workspace structure with Project as root name.
	// With Project=TestCorp, expect .dbeaver/.gitkeep under TestCorp.
	expected := filepath.Join(tmpDir, "TestCorp", ".dbeaver", ".gitkeep")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected %s to exist: %v", expected, err)
	}
}

// TestApplyTemplateDryRun verifies --dry-run writes no files.
func TestApplyTemplateDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"apply-template", "multi-env",
		"--output", tmpDir,
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("apply-template --dry-run failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries, expected 0", len(entries))
	}
}

// TestApplyTemplateSetFlag verifies --set Key=Value overrides template variables.
func TestApplyTemplateSetFlag(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"apply-template", "multi-env",
		"--output", tmpDir,
		"--set", "Project=AlphaProject",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("apply-template --set failed: %v", err)
	}

	// Verify that the project name "AlphaProject" was applied (path contains it).
	expected := filepath.Join(tmpDir, "AlphaProject", ".dbeaver", ".gitkeep")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected AlphaProject path to exist: %v", err)
	}
}
