package main

import (
	"os"

	mcpserver "github.com/jinkp/dbeaver-go-mcp/internal/mcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "mcp",
		Short:         "Start MCP stdio server (exposes dwm tools to AI assistants)",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// CRITICAL: redirect ALL output to stderr before ServeStdio.
			// MCP transport owns stdout from this point forward.
			log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
			return mcpserver.StartServer()
		},
	}
}

func init() {
	rootCmd.AddCommand(newMCPCmd())
}
