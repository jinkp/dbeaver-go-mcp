package main

import (
	"strings"

	"github.com/jinkp/dbeaver-go-mcp/internal/config"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var generateScriptsCmd = &cobra.Command{
	Use:   "generate-scripts <config.yaml>",
	Short: "Generate SQL scripts from embedded templates for each engine",
	Long: `Reads a YAML configuration and generates SQL scripts for each engine
defined in the config. Scripts are written to Scripts/scripts/<engine>/<name>.sql.`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerateScripts,
}

func init() {
	generateScriptsCmd.Flags().String("scripts", "",
		"comma-separated list of scripts to generate (default: all)")
	rootCmd.AddCommand(generateScriptsCmd)
}

func runGenerateScripts(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	cfg, err := config.Load(args[0])
	if err != nil {
		return err
	}

	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	scriptsFlag, _ := cmd.Flags().GetString("scripts")

	spec := scripts.ScriptSpec{
		Project: cfg.Project,
		Engines: cfg.Engines,
		Scripts: parseScriptsFlag(scriptsFlag),
	}

	outFiles, err := scripts.Generate(spec)
	if err != nil {
		return err
	}

	files := convertScriptFiles(outFiles)

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}

// parseScriptsFlag splits the comma-separated scripts flag into a slice.
// Returns nil (meaning "all") when the flag is empty.
func parseScriptsFlag(flag string) []string {
	if flag == "" {
		return nil
	}
	return strings.Split(flag, ",")
}

// convertScriptFiles converts []scripts.OutFile to []output.File.
func convertScriptFiles(src []scripts.OutFile) []output.File {
	dst := make([]output.File, len(src))
	for i, f := range src {
		dst[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return dst
}
