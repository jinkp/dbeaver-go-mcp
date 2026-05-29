package main

import (
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/workspace"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "create-project <name>",
	Short: "Create a DBeaver project nested inside an existing workspace",
	Long: `Creates a DBeaver-compatible project folder structure nested under <workspace>/<name>.
Produces: .dbeaver/, Scripts/, Diagrams/, Bookmarks/ with .gitkeep markers.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateProject,
}

func init() {
	createProjectCmd.Flags().String("workspace", "", "parent workspace name (required)")
	if err := createProjectCmd.MarkFlagRequired("workspace"); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(createProjectCmd)
}

func runCreateProject(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	projectName := args[0]
	wsName, _ := cmd.Flags().GetString("workspace")
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	files := workspace.GenerateProject(wsName, projectName)

	if dryRun {
		for _, f := range files {
			log.Info().Str("path", f.Path).Int("bytes", len(f.Content)).Msg("[dry-run] would write")
		}
		return nil
	}

	return output.Write(files, outputDir, conflictMode(cmd))
}
