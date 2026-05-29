package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/claude"
	"github.com/jinkp/dbeaver-go-mcp/internal/opencode"
	"github.com/spf13/cobra"
)

// findSubCommand looks for a direct subcommand matching the given use name.
func findSubCommand(cmd *cobra.Command, use string) *cobra.Command {
	for _, sub := range cmd.Commands() {
		if sub.Use == use || (len(sub.Use) >= len(use) && sub.Use[:len(use)] == use) {
			return sub
		}
	}
	return nil
}

// TestMCPCommandRegistered verifies that "mcp" is registered on rootCmd.
func TestMCPCommandRegistered(t *testing.T) {
	cmd := findSubCommand(rootCmd, "mcp")
	if cmd == nil {
		t.Fatal("expected 'mcp' subcommand to be registered on rootCmd")
	}
}

// TestSetupCommandRegistered verifies that "setup" is registered on rootCmd.
func TestSetupCommandRegistered(t *testing.T) {
	cmd := findSubCommand(rootCmd, "setup")
	if cmd == nil {
		t.Fatal("expected 'setup' subcommand to be registered on rootCmd")
	}
}

// TestSetupOpenCodeCommandRegistered verifies that "setup opencode" subcommand exists.
func TestSetupOpenCodeCommandRegistered(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "opencode")
	if sub == nil {
		t.Fatal("expected 'opencode' subcommand under 'setup'")
	}
}

// TestSetupClaudeCommandRegistered verifies that "setup claude" subcommand exists.
func TestSetupClaudeCommandRegistered(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "claude")
	if sub == nil {
		t.Fatal("expected 'claude' subcommand under 'setup'")
	}
}

// TestSetupOpenCodeGlobalFlag verifies that "setup opencode" exposes --global flag.
func TestSetupOpenCodeGlobalFlag(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "opencode")
	if sub == nil {
		t.Fatal("expected 'opencode' subcommand under 'setup'")
	}
	f := sub.Flags().Lookup("global")
	if f == nil {
		t.Fatal("expected --global flag on 'setup opencode'")
	}
}

// TestSetupClaudeGlobalFlag verifies that "setup claude" exposes --global flag.
func TestSetupClaudeGlobalFlag(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "claude")
	if sub == nil {
		t.Fatal("expected 'claude' subcommand under 'setup'")
	}
	f := sub.Flags().Lookup("global")
	if f == nil {
		t.Fatal("expected --global flag on 'setup claude'")
	}
}

// TestSetupOpenCodeRunsWithoutError verifies that opencode.Register writes
// the expected JSON structure to a temp file — same logic as the cobra command.
func TestSetupOpenCodeRunsWithoutError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opencode.json")

	if err := opencode.Register(path); err != nil {
		t.Fatalf("opencode.Register returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse JSON: %v", err)
	}

	if _, ok := m["mcp"]; !ok {
		t.Error("expected 'mcp' key in opencode.json after Register")
	}
}

// TestSetupClaudeRunsWithoutError verifies that claude.Register writes the
// expected JSON structure to a temp file — same logic as the cobra command.
func TestSetupClaudeRunsWithoutError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".claude", "settings.json")

	if err := claude.Register(path); err != nil {
		t.Fatalf("claude.Register returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse JSON: %v", err)
	}

	if _, ok := m["mcpServers"]; !ok {
		t.Error("expected 'mcpServers' key in settings.json after Register")
	}
}

// TestMCPCommandSilencesErrors verifies mcp command has SilenceErrors set,
// which prevents cobra from printing errors to stderr (MCP owns stdout/stderr).
func TestMCPCommandSilencesErrors(t *testing.T) {
	cmd := findSubCommand(rootCmd, "mcp")
	if cmd == nil {
		t.Fatal("expected 'mcp' subcommand on rootCmd")
	}
	if !cmd.SilenceErrors {
		t.Error("expected SilenceErrors=true on mcp command (MCP transport owns stderr)")
	}
	if !cmd.SilenceUsage {
		t.Error("expected SilenceUsage=true on mcp command")
	}
}

// TestSetupOpenCodeLocalFlag verifies that "setup opencode" exposes --local flag.
func TestSetupOpenCodeLocalFlag(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "opencode")
	if sub == nil {
		t.Fatal("expected 'opencode' subcommand under 'setup'")
	}
	f := sub.Flags().Lookup("local")
	if f == nil {
		t.Fatal("expected --local flag on 'setup opencode'")
	}
}

// TestSetupClaudeLocalFlag verifies that "setup claude" exposes --local flag.
func TestSetupClaudeLocalFlag(t *testing.T) {
	setupCmd := findSubCommand(rootCmd, "setup")
	if setupCmd == nil {
		t.Fatal("expected 'setup' subcommand on rootCmd")
	}
	sub := findSubCommand(setupCmd, "claude")
	if sub == nil {
		t.Fatal("expected 'claude' subcommand under 'setup'")
	}
	f := sub.Flags().Lookup("local")
	if f == nil {
		t.Fatal("expected --local flag on 'setup claude'")
	}
}

// TestNewSetupCmdStructure verifies that newSetupCmd() returns a command
// with exactly two subcommands: "opencode" and "claude".
func TestNewSetupCmdStructure(t *testing.T) {
	cmd := newSetupCmd()
	if cmd == nil {
		t.Fatal("newSetupCmd() returned nil")
	}
	subcmds := cmd.Commands()
	if len(subcmds) != 2 {
		t.Errorf("expected 2 subcommands under setup, got %d", len(subcmds))
	}
	names := make(map[string]bool)
	for _, sub := range subcmds {
		names[sub.Use] = true
	}
	if !names["opencode"] {
		t.Error("expected 'opencode' subcommand under setup")
	}
	if !names["claude"] {
		t.Error("expected 'claude' subcommand under setup")
	}
}
