package mcp

import (
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// StartServer creates and starts the MCP stdio server.
// All tool registrations happen before ServeStdio takes ownership of stdout.
func StartServer() error {
	s := server.NewMCPServer("dwm", "2.0.0",
		server.WithToolCapabilities(false),
	)
	registerTools(s)
	return server.ServeStdio(s)
}

// registerTools attaches all 7 dwm tool handlers to the MCP server.
func registerTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("create_workspace",
		mcplib.WithDescription("Create a DBeaver workspace folder structure"),
		mcplib.WithString("name", mcplib.Required(), mcplib.Description("Workspace name")),
		mcplib.WithString("environments", mcplib.Description("Comma-separated environments (default: DEV,QA,STAGE,PROD)")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip (default: error)")),
	), ToolCreateWorkspace)

	s.AddTool(mcplib.NewTool("create_project",
		mcplib.WithDescription("Create a DBeaver project nested inside a workspace"),
		mcplib.WithString("name", mcplib.Required(), mcplib.Description("Project name")),
		mcplib.WithString("workspace", mcplib.Required(), mcplib.Description("Parent workspace name")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip")),
	), ToolCreateProject)

	s.AddTool(mcplib.NewTool("create_connections",
		mcplib.WithDescription("Generate DBeaver data-sources.json from a YAML config file"),
		mcplib.WithString("config", mcplib.Required(), mcplib.Description("Path to YAML config file")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithBoolean("split_by_env", mcplib.Description("Generate one file per environment")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip")),
	), ToolCreateConnections)

	s.AddTool(mcplib.NewTool("generate_scripts",
		mcplib.WithDescription("Generate SQL maintenance scripts from a YAML config file"),
		mcplib.WithString("config", mcplib.Required(), mcplib.Description("Path to YAML config file")),
		mcplib.WithString("scripts", mcplib.Description("Comma-separated script names (default: all)")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip")),
	), ToolGenerateScripts)

	s.AddTool(mcplib.NewTool("sync_snippets",
		mcplib.WithDescription("Sync SQL snippet files organized by database engine"),
		mcplib.WithString("config", mcplib.Required(), mcplib.Description("Path to YAML config file")),
		mcplib.WithString("engines", mcplib.Description("Comma-separated engines (default: from config)")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip")),
	), ToolSyncSnippets)

	s.AddTool(mcplib.NewTool("create_template",
		mcplib.WithDescription("Scaffold a new DBeaver project template definition"),
		mcplib.WithString("name", mcplib.Required(), mcplib.Description("Template name")),
		mcplib.WithString("output", mcplib.Description("Output directory for template files (default: ./output)")),
	), ToolCreateTemplate)

	s.AddTool(mcplib.NewTool("apply_template",
		mcplib.WithDescription("Apply a DBeaver project template with variable substitution"),
		mcplib.WithString("name", mcplib.Required(), mcplib.Description("Template name")),
		mcplib.WithObject("vars", mcplib.Description("Template variable overrides as key-value pairs")),
		mcplib.WithString("output", mcplib.Description("Output directory (default: ./output)")),
		mcplib.WithString("conflict", mcplib.Description("Conflict mode: error|overwrite|skip")),
	), ToolApplyTemplate)
}
