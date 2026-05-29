package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jinkp/dbeaver-go-mcp/internal/claude"
	"github.com/jinkp/dbeaver-go-mcp/internal/opencode"
	tuisetup "github.com/jinkp/dbeaver-go-mcp/internal/tui/setup"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newSetupCmd() *cobra.Command {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Register dwm as an MCP server in AI clients",
	}
	setupCmd.AddCommand(newSetupOpenCodeCmd())
	setupCmd.AddCommand(newSetupClaudeCmd())
	return setupCmd
}

func newSetupOpenCodeCmd() *cobra.Command {
	var global bool
	var local bool
	cmd := &cobra.Command{
		Use:   "opencode",
		Short: "Register dwm MCP server in OpenCode",
		RunE: func(cmd *cobra.Command, args []string) error {
			globalChanged := cmd.Flags().Changed("global")
			localChanged := cmd.Flags().Changed("local")

			if globalChanged || localChanged {
				// bypass: text mode (backward compatible + new --local flag)
				path := opencode.LocalPath()
				if global {
					path = opencode.GlobalPath()
				}
				if err := opencode.Register(path); err != nil {
					return err
				}
				fmt.Printf("✓ dwm registered in %s\n", path)
				fmt.Println("  Ensure 'dwm' is on your PATH.")
				return nil
			}

			// non-TTY fallback: auto-register locally
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				path := opencode.LocalPath()
				if err := opencode.Register(path); err != nil {
					return err
				}
				fmt.Printf("✓ dwm registered in %s\n", path)
				return nil
			}

			// TUI wizard
			m := tuisetup.NewModel("opencode", opencode.Register, opencode.GlobalPath, opencode.LocalPath)
			_, err := tea.NewProgram(m).Run()
			return err
		},
	}
	cmd.Flags().BoolVar(&global, "global", false, "Write to global OpenCode config")
	cmd.Flags().BoolVar(&local, "local", false, "Write to local config (./opencode.json)")
	return cmd
}

func newSetupClaudeCmd() *cobra.Command {
	var global bool
	var local bool
	cmd := &cobra.Command{
		Use:   "claude",
		Short: "Register dwm MCP server in Claude Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			globalChanged := cmd.Flags().Changed("global")
			localChanged := cmd.Flags().Changed("local")

			if globalChanged || localChanged {
				// bypass: text mode (backward compatible + new --local flag)
				path := claude.LocalPath()
				if global {
					path = claude.GlobalPath()
				}
				if err := claude.Register(path); err != nil {
					return err
				}
				fmt.Printf("✓ dwm registered in %s\n", path)
				fmt.Println("  Ensure 'dwm' is on your PATH.")
				return nil
			}

			// non-TTY fallback: auto-register locally
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				path := claude.LocalPath()
				if err := claude.Register(path); err != nil {
					return err
				}
				fmt.Printf("✓ dwm registered in %s\n", path)
				return nil
			}

			// TUI wizard
			m := tuisetup.NewModel("claude", claude.Register, claude.GlobalPath, claude.LocalPath)
			_, err := tea.NewProgram(m).Run()
			return err
		},
	}
	cmd.Flags().BoolVar(&global, "global", false, "Write to global Claude Code config")
	cmd.Flags().BoolVar(&local, "local", false, "Write to local config (./.claude/settings.json)")
	return cmd
}

func init() {
	rootCmd.AddCommand(newSetupCmd())
}
