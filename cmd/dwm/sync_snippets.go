package main

import (
	"strings"

	"github.com/jinkp/dbeaver-go-mcp/internal/config"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var syncSnippetsCmd = &cobra.Command{
	Use:   "sync-snippets <config.yaml>",
	Short: "Sync SQL snippet files from embedded templates",
	Long: `Reads a YAML configuration and generates SQL snippet files for each engine.
Snippets are written to Scripts/snippets/<engine>/<name>.sql.`,
	Args: cobra.ExactArgs(1),
	RunE: runSyncSnippets,
}

func init() {
	syncSnippetsCmd.Flags().String("engines", "",
		"comma-separated list of engines to generate snippets for (default: from config)")
	rootCmd.AddCommand(syncSnippetsCmd)
}

func runSyncSnippets(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	cfg, err := config.Load(args[0])
	if err != nil {
		return err
	}

	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	enginesFlag, _ := cmd.Flags().GetString("engines")

	// --engines flag overrides config engines.
	engines := cfg.Engines
	if enginesFlag != "" {
		engines = strings.Split(enginesFlag, ",")
	}

	spec := snippets.SnippetSpec{Engines: engines}

	outFiles, err := snippets.Generate(spec)
	if err != nil {
		return err
	}

	files := convertSnippetFiles(outFiles)

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}

// convertSnippetFiles converts []snippets.OutFile to []output.File.
func convertSnippetFiles(src []snippets.OutFile) []output.File {
	dst := make([]output.File, len(src))
	for i, f := range src {
		dst[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return dst
}
