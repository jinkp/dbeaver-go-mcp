package main

import (
	"strings"

	"github.com/jinkp/dbeaver-go-mcp/internal/config"
	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var createConnectionsCmd = &cobra.Command{
	Use:   "create-connections <config.yaml>",
	Short: "Generate DBeaver connection files from a YAML config",
	Long: `Reads a YAML configuration file and generates data-sources.json and
data-sources-config.json compatible with DBeaver 26.0. Passwords are never saved.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateConnections,
}

func init() {
	createConnectionsCmd.Flags().Bool("split-by-env", false,
		"generate separate connection files per environment")
	rootCmd.AddCommand(createConnectionsCmd)
}

func runCreateConnections(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	cfg, err := config.Load(args[0])
	if err != nil {
		return err
	}

	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	splitByEnv, _ := cmd.Flags().GetBool("split-by-env")

	// Build the full (merged) spec used both for the global file and as the base for per-env files.
	fullSpec := connections.ConnSpec{
		Project:      cfg.Project,
		Environments: cfg.Environments,
		Connections:  mapConnConfigs(cfg.Connections),
	}

	// Generate the global merged data-sources.json and data-sources-config.json.
	outFiles, err := connections.Generate(fullSpec)
	if err != nil {
		return err
	}
	files := convertConnectionFiles(outFiles)

	// If --split-by-env is set, also generate one data-sources-<env>.json per environment.
	if splitByEnv {
		for _, env := range cfg.Environments {
			envSpec := connections.ConnSpec{
				Project:      cfg.Project,
				Environments: []string{env},
				Connections:  fullSpec.Connections,
			}
			envOutFiles, err := connections.Generate(envSpec)
			if err != nil {
				return err
			}
			// envOutFiles[0] is data-sources.json — rename it to data-sources-<env>.json.
			// envOutFiles[1] is data-sources-config.json — skip it (global already included).
			if len(envOutFiles) > 0 {
				renamed := output.File{
					Path:    "data-sources-" + strings.ToLower(env) + ".json",
					Content: envOutFiles[0].Content,
				}
				files = append(files, renamed)
			}
		}
	}

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}

// mapConnConfigs converts []config.ConnConfig to []connections.ConnConfig.
func mapConnConfigs(src []config.ConnConfig) []connections.ConnConfig {
	dst := make([]connections.ConnConfig, len(src))
	for i, c := range src {
		dst[i] = connections.ConnConfig{
			Name:     c.Name,
			Type:     c.Type,
			Host:     c.Host,
			Port:     c.Port,
			Database: c.Database,
		}
	}
	return dst
}

// convertConnectionFiles converts []connections.OutFile to []output.File.
func convertConnectionFiles(src []connections.OutFile) []output.File {
	dst := make([]output.File, len(src))
	for i, f := range src {
		dst[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return dst
}
