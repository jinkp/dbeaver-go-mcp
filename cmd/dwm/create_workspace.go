package main

import (
	"strings"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/workspace"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var createWorkspaceCmd = &cobra.Command{
	Use:   "create-workspace <name>",
	Short: "Create a DBeaver workspace folder structure",
	Long: `Creates a DBeaver-compatible workspace folder structure at <name> inside the output directory.
Produces: .dbeaver/, Scripts/, Diagrams/, Bookmarks/ with .gitkeep markers.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateWorkspace,
}

func init() {
	createWorkspaceCmd.Flags().String("environments", "DEV,QA,STAGE,PROD",
		"comma-separated list of environment names (passed to workspace generator)")
	rootCmd.AddCommand(createWorkspaceCmd)
}

func runCreateWorkspace(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	name := args[0]

	envsFlag, _ := cmd.Flags().GetString("environments")
	envs := strings.Split(envsFlag, ",")

	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	files := workspace.Generate(name, envs)

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}
