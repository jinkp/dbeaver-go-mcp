package main

import (
	"strings"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var applyTemplateCmd = &cobra.Command{
	Use:   "apply-template <name>",
	Short: "Apply a DBeaver project template",
	Long: `Loads a named template (from local templates/, ~/.dwm/templates/, or built-in) and
renders it with the provided variable overrides. Generates all configured artifacts.`,
	Args: cobra.ExactArgs(1),
	RunE: runApplyTemplate,
}

func init() {
	applyTemplateCmd.Flags().StringArray("set", nil,
		"variable override in Key=Value format (repeatable)")
	rootCmd.AddCommand(applyTemplateCmd)
}

func runApplyTemplate(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	name := args[0]
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	setFlags, _ := cmd.Flags().GetStringArray("set")

	tmpl, err := templates.Load(name)
	if err != nil {
		return err
	}

	vars := parseSetFlags(setFlags)

	outFiles, err := templates.Apply(tmpl, vars)
	if err != nil {
		return err
	}

	files := convertTemplateFiles(outFiles)

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}

// parseSetFlags converts ["Key=Value", ...] into map[string]any.
func parseSetFlags(flags []string) map[string]any {
	vars := make(map[string]any, len(flags))
	for _, kv := range flags {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}

// convertTemplateFiles converts []templates.OutFile to []output.File.
func convertTemplateFiles(src []templates.OutFile) []output.File {
	dst := make([]output.File, len(src))
	for i, f := range src {
		dst[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return dst
}
