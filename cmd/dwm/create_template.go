package main

import (
	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
	"github.com/spf13/cobra"
)

var createTemplateCmd = &cobra.Command{
	Use:   "create-template <name>",
	Short: "Scaffold a new DBeaver project template",
	Long: `Creates a new template definition under <output>/templates/<name>/template.yaml
with placeholder content ready for customization.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateTemplate,
}

func init() {
	createTemplateCmd.Flags().String("from", "",
		"scaffold from an existing template (optional)")
	rootCmd.AddCommand(createTemplateCmd)
}

func runCreateTemplate(cmd *cobra.Command, args []string) error {
	if err := validateFlags(cmd); err != nil {
		return err
	}

	name := args[0]
	outputDir, _ := cmd.Flags().GetString("output")

	return templates.Scaffold(name, outputDir)
}
