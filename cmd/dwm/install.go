package main

import (
	"fmt"

	"github.com/jinkp/dbeaver-go-mcp/internal/dbeaver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var (
		outputDir     string
		workspacePath string
		overwrite     bool
		dryRun        bool
		force         bool
		backup        bool
	)

	cmd := &cobra.Command{
		Use:   "install <workspace-name>",
		Short: "Install generated workspace directly into DBeaver",
		Long: `Install copies files from the generated output directory into your DBeaver workspace.

Safe by default:
  - Protected files (.dbeaver/credentials-config.json, .dbeaver/.settings/, .metadata/, .project) are NEVER modified.
  - data-sources.json is merged by connection ID (your connections are preserved).
  - Existing files are skipped unless --overwrite is passed.
  - DBeaver must be closed before installing (use --force to bypass this check).

Use --dry-run to preview the install plan without writing any files.
Use --backup to create a timestamped backup of the target workspace before writing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsName := args[0]

			// Resolve target workspace root.
			targetRoot := workspacePath
			if targetRoot == "" {
				var err error
				targetRoot, err = dbeaver.WorkspacePath()
				if err != nil {
					return err
				}
			}

			// Process safety check: abort if DBeaver is running (unless --force).
			if dbeaver.IsRunning() {
				if !force {
					return fmt.Errorf("DBeaver is running. Close it before installing, or use --force")
				}
				log.Warn().Msg("WARNING: DBeaver is running. Changes may be overwritten when it closes.")
				fmt.Println("WARNING: DBeaver is running. Changes may be overwritten when it closes.")
			}

			spec := dbeaver.InstallSpec{
				WorkspaceName: wsName,
				SourceDir:     outputDir,
				TargetRoot:    targetRoot,
				Overwrite:     overwrite,
				DryRun:        dryRun,
				Force:         force,
				Backup:        backup,
			}

			actions, err := dbeaver.Execute(spec)
			if err != nil {
				return err
			}

			// Print action summary.
			for _, a := range actions {
				switch a.Action {
				case "copy":
					fmt.Printf("  copy    %s\n", a.RelPath)
				case "merge":
					fmt.Printf("  merge   %s\n", a.RelPath)
				case "skip":
					fmt.Printf("  skip    %s (%s)\n", a.RelPath, a.Reason)
				case "skip-protected":
					fmt.Printf("  protect %s\n", a.RelPath)
				}
			}

			if dryRun {
				fmt.Println("\n[dry-run] No files written.")
			} else {
				fmt.Printf("\n✓ Installed %s into DBeaver workspace.\n", wsName)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "output", "./output", "Source output directory containing the workspace folder")
	cmd.Flags().StringVar(&workspacePath, "workspace-path", "", "Override the DBeaver workspace path (auto-detected by default)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Replace existing files instead of merging/skipping")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show install plan without writing any files")
	cmd.Flags().BoolVar(&force, "force", false, "Bypass the DBeaver-is-running safety check")
	cmd.Flags().BoolVar(&backup, "backup", false, "Create a timestamped backup of the target workspace before writing")

	return cmd
}

func init() {
	rootCmd.AddCommand(newInstallCmd())
}
