package main

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestInstallCommandRegistered verifies that the install command is wired into rootCmd.
func TestInstallCommandRegistered(t *testing.T) {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "install" {
			return
		}
	}
	t.Error("'install' command not found in rootCmd — check that install.go init() registers it")
}

// TestInstallHelp verifies the install command has the expected flags registered.
func TestInstallHelp(t *testing.T) {
	var install *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "install" {
			install = cmd
			break
		}
	}
	if install == nil {
		t.Fatal("install command not found")
	}
	// Verify key flags are present
	expectedFlags := []string{"output", "workspace-path", "overwrite", "dry-run", "force", "backup"}
	for _, flag := range expectedFlags {
		if f := install.Flags().Lookup(flag); f == nil {
			t.Errorf("expected flag --%s to be registered on install command", flag)
		}
	}
}

// TestInstallMissingSourceReturnsError verifies that running install with a
// non-existent --output directory returns a "source not found" error.
func TestInstallMissingSourceReturnsError(t *testing.T) {
	// Reset persistent flags to avoid contamination from prior tests.
	_ = rootCmd.PersistentFlags().Set("dry-run", "false")

	rootCmd.SetArgs([]string{
		"install", "TestWs",
		"--dry-run",
		"--output", `C:\nonexistent-path-for-test`,
		"--workspace-path", t.TempDir(), // valid target so workspace check passes
	})
	err := rootCmd.Execute()

	// Reset args
	rootCmd.SetArgs([]string{})

	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
	if !strings.Contains(err.Error(), "source not found") {
		t.Errorf("expected 'source not found' in error, got: %q", err.Error())
	}
}
