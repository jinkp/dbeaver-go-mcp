package main

import (
	"errors"
	"os"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dwm",
	Short: "DBeaver Workspace Manager — generate DBeaver workspace artifacts from YAML",
	Long: `dwm is a CLI tool that generates DBeaver-compatible workspace structures,
connection files, SQL scripts, and snippets from a declarative YAML configuration.`,
}

// Execute runs the root cobra command. Called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Configure zerolog to write to stderr.
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Persistent flags available to all sub-commands.
	rootCmd.PersistentFlags().String("output", "./output", "base output directory")
	rootCmd.PersistentFlags().Bool("overwrite", false, "delete and regenerate existing files")
	rootCmd.PersistentFlags().Bool("skip", false, "skip existing files (log warning, continue)")
	rootCmd.PersistentFlags().Bool("dry-run", false, "print plan to stdout, write nothing")
	rootCmd.PersistentFlags().Bool("merge", false, "merge into existing files (not supported in v1)")
}

// validateFlags checks for unsupported flags and returns an error if any are set.
// Call this at the top of every command's RunE before doing any real work.
func validateFlags(cmd *cobra.Command) error {
	merge, _ := cmd.Flags().GetBool("merge")
	if merge {
		return errors.New("merge not supported in v1")
	}
	return nil
}

// conflictMode reads the --overwrite and --skip flags from cmd and returns the
// corresponding ConflictMode. --overwrite takes precedence over --skip.
func conflictMode(cmd *cobra.Command) output.ConflictMode {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	if overwrite {
		return output.ConflictOverwrite
	}
	skip, _ := cmd.Flags().GetBool("skip")
	if skip {
		return output.ConflictSkip
	}
	return output.ConflictError
}
