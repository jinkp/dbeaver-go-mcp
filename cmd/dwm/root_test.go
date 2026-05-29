package main

import (
	"strings"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/spf13/cobra"
)

// TestConflictModeDefault verifies that conflictMode returns ConflictError
// when neither --overwrite nor --skip is set.
func TestConflictModeDefault(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("overwrite", false, "")
	cmd.Flags().Bool("skip", false, "")
	// Neither flag set — default values (both false)

	got := conflictMode(cmd)
	if got != output.ConflictError {
		t.Errorf("expected ConflictError (%d), got %d", output.ConflictError, got)
	}
}

// TestConflictModeOverwrite verifies that conflictMode returns ConflictOverwrite
// when --overwrite is set.
func TestConflictModeOverwrite(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("overwrite", false, "")
	cmd.Flags().Bool("skip", false, "")
	if err := cmd.Flags().Set("overwrite", "true"); err != nil {
		t.Fatalf("set --overwrite: %v", err)
	}

	got := conflictMode(cmd)
	if got != output.ConflictOverwrite {
		t.Errorf("expected ConflictOverwrite (%d), got %d", output.ConflictOverwrite, got)
	}
}

// TestConflictModeSkip verifies that conflictMode returns ConflictSkip
// when --skip is set (and --overwrite is not).
func TestConflictModeSkip(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("overwrite", false, "")
	cmd.Flags().Bool("skip", false, "")
	if err := cmd.Flags().Set("skip", "true"); err != nil {
		t.Fatalf("set --skip: %v", err)
	}

	got := conflictMode(cmd)
	if got != output.ConflictSkip {
		t.Errorf("expected ConflictSkip (%d), got %d", output.ConflictSkip, got)
	}
}

// TestConflictModeOverwriteTakesPrecedence verifies that --overwrite takes
// precedence over --skip when both are set.
func TestConflictModeOverwriteTakesPrecedence(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("overwrite", false, "")
	cmd.Flags().Bool("skip", false, "")
	if err := cmd.Flags().Set("overwrite", "true"); err != nil {
		t.Fatalf("set --overwrite: %v", err)
	}
	if err := cmd.Flags().Set("skip", "true"); err != nil {
		t.Fatalf("set --skip: %v", err)
	}

	got := conflictMode(cmd)
	if got != output.ConflictOverwrite {
		t.Errorf("expected ConflictOverwrite (%d) when both flags set, got %d", output.ConflictOverwrite, got)
	}
}

// TestValidateFlagsMergeReturnsError verifies that validateFlags returns an error
// with the "merge not supported in v1" message when --merge is passed.
func TestValidateFlagsMergeReturnsError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("merge", false, "")
	if err := cmd.Flags().Set("merge", "true"); err != nil {
		t.Fatalf("set --merge: %v", err)
	}

	err := validateFlags(cmd)
	if err == nil {
		t.Fatal("expected error when --merge is set, got nil")
	}
	if !strings.Contains(err.Error(), "merge not supported in v1") {
		t.Errorf("expected error message to contain 'merge not supported in v1', got: %q", err.Error())
	}
}

// TestValidateFlagsNoMergeReturnsNil verifies that validateFlags returns nil
// when --merge is NOT set (normal operation).
func TestValidateFlagsNoMergeReturnsNil(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("merge", false, "")
	// --merge not set — default false

	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("expected nil error when --merge is not set, got: %v", err)
	}
}

// TestCreateConnectionsMergeReturnsError verifies that passing --merge to
// create-connections returns an error containing "merge not supported in v1".
func TestCreateConnectionsMergeReturnsError(t *testing.T) {
	// Reset persistent flags before test.
	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")
	_ = rootCmd.PersistentFlags().Set("merge", "false")

	rootCmd.SetArgs([]string{
		"create-connections", "testdata/test-config.yaml",
		"--merge",
	})
	err := rootCmd.Execute()

	// Always reset --merge after the test so subsequent tests are not affected.
	_ = rootCmd.PersistentFlags().Set("merge", "false")

	if err == nil {
		t.Fatal("expected error when --merge is passed, got nil")
	}
	if !strings.Contains(err.Error(), "merge not supported in v1") {
		t.Errorf("expected 'merge not supported in v1' in error, got: %q", err.Error())
	}
}
