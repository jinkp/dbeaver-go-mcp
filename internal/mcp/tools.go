package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"

	"github.com/jinkp/dbeaver-go-mcp/internal/config"
	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
	"github.com/jinkp/dbeaver-go-mcp/internal/workspace"
)

// toolResult builds the standard CallToolResult JSON response for successful tool calls.
func toolResult(files []output.File, outputDir string) *mcplib.CallToolResult {
	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = f.Path
	}
	b, _ := json.Marshal(ToolResult{FilesWritten: paths, OutputDir: outputDir})
	return mcplib.NewToolResultText(string(b))
}

// errResult builds a ToolResultError. Never returns a Go error (keeps stdio transport clean).
func errResult(msg string) *mcplib.CallToolResult {
	return mcplib.NewToolResultError(msg)
}

// ToolCreateWorkspace implements the create_workspace tool.
// Exported for testability; registered via registerTools.
func ToolCreateWorkspace(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		return errResult("name is required"), nil
	}
	envStr := req.GetString("environments", "DEV,QA,STAGE,PROD")
	envs := strings.Split(envStr, ",")
	for i, e := range envs {
		envs[i] = strings.TrimSpace(e)
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	files := workspace.Generate(name, envs)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}

// ToolCreateProject implements the create_project tool.
func ToolCreateProject(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		return errResult("name is required"), nil
	}
	ws := req.GetString("workspace", "")
	if ws == "" {
		return errResult("workspace is required"), nil
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	files := workspace.GenerateProject(ws, name)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}

// ToolCreateConnections implements the create_connections tool.
func ToolCreateConnections(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	configPath := req.GetString("config", "")
	if configPath == "" {
		return errResult("config is required"), nil
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	cfg, err := config.Load(configPath)
	if err != nil {
		return errResult(fmt.Sprintf("load config: %s", err)), nil
	}

	// Map config.ConnConfig → connections.ConnConfig
	connConfs := make([]connections.ConnConfig, len(cfg.Connections))
	for i, c := range cfg.Connections {
		connConfs[i] = connections.ConnConfig{
			Name:     c.Name,
			Type:     c.Type,
			Host:     c.Host,
			Port:     c.Port,
			Database: c.Database,
		}
	}

	spec := connections.ConnSpec{
		Project:      cfg.Project,
		Environments: cfg.Environments,
		Connections:  connConfs,
	}

	outFiles, err := connections.Generate(spec)
	if err != nil {
		return errResult(fmt.Sprintf("generate connections: %s", err)), nil
	}

	files := ConnectionFilesToOutput(outFiles)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}

// ToolGenerateScripts implements the generate_scripts tool.
func ToolGenerateScripts(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	configPath := req.GetString("config", "")
	if configPath == "" {
		return errResult("config is required"), nil
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	cfg, err := config.Load(configPath)
	if err != nil {
		return errResult(fmt.Sprintf("load config: %s", err)), nil
	}

	// Parse optional scripts filter
	var scriptNames []string
	if s := req.GetString("scripts", ""); s != "" {
		for _, n := range strings.Split(s, ",") {
			scriptNames = append(scriptNames, strings.TrimSpace(n))
		}
	}

	spec := scripts.ScriptSpec{
		Project: cfg.Project,
		Engines: cfg.Engines,
		Scripts: scriptNames,
	}

	outFiles, err := scripts.Generate(spec)
	if err != nil {
		return errResult(fmt.Sprintf("generate scripts: %s", err)), nil
	}

	files := ScriptFilesToOutput(outFiles)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}

// ToolSyncSnippets implements the sync_snippets tool.
func ToolSyncSnippets(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	configPath := req.GetString("config", "")
	if configPath == "" {
		return errResult("config is required"), nil
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	cfg, err := config.Load(configPath)
	if err != nil {
		return errResult(fmt.Sprintf("load config: %s", err)), nil
	}

	// Use engines from config; allow override via "engines" param
	engines := cfg.Engines
	if e := req.GetString("engines", ""); e != "" {
		engines = nil
		for _, eng := range strings.Split(e, ",") {
			engines = append(engines, strings.TrimSpace(eng))
		}
	}

	spec := snippets.SnippetSpec{Engines: engines}
	outFiles, err := snippets.Generate(spec)
	if err != nil {
		return errResult(fmt.Sprintf("generate snippets: %s", err)), nil
	}

	files := SnippetFilesToOutput(outFiles)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}

// ToolCreateTemplate implements the create_template tool.
func ToolCreateTemplate(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		return errResult("name is required"), nil
	}
	outputDir := req.GetString("output", "./output")

	if err := templates.Scaffold(name, outputDir); err != nil {
		return errResult(fmt.Sprintf("scaffold template: %s", err)), nil
	}

	// Return the path of the created template.yaml
	templatePath := "templates/" + name + "/template.yaml"
	files := []output.File{{Path: templatePath}}
	return toolResult(files, outputDir), nil
}

// ToolApplyTemplate implements the apply_template tool.
func ToolApplyTemplate(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		return errResult("name is required"), nil
	}
	outputDir := req.GetString("output", "./output")
	mode := ConflictModeFromString(req.GetString("conflict", ""))

	// Parse optional vars object
	vars := map[string]any{}
	if raw, ok := req.GetArguments()["vars"]; ok {
		if m, ok := raw.(map[string]any); ok {
			vars = m
		}
	}

	tmpl, err := templates.Load(name)
	if err != nil {
		return errResult(fmt.Sprintf("load template %q: %s", name, err)), nil
	}

	outFiles, err := templates.Apply(tmpl, vars)
	if err != nil {
		return errResult(fmt.Sprintf("apply template: %s", err)), nil
	}

	files := TemplateFilesToOutput(outFiles)
	if err := output.Write(files, outputDir, mode); err != nil {
		return errResult(err.Error()), nil
	}
	return toolResult(files, outputDir), nil
}
