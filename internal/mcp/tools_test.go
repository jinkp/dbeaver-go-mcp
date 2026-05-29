package mcp_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"

	internalmcp "github.com/jinkp/dbeaver-go-mcp/internal/mcp"
)

// helper: build a CallToolRequest with given arguments map
func makeReq(args map[string]any) mcplib.CallToolRequest {
	return mcplib.CallToolRequest{
		Params: mcplib.CallToolParams{
			Arguments: args,
		},
	}
}

// helper: parse ToolResult from a *CallToolResult text content
func parseToolResult(t *testing.T, res *mcplib.CallToolResult) internalmcp.ToolResult {
	t.Helper()
	if res == nil {
		t.Fatal("result is nil")
	}
	if len(res.Content) == 0 {
		t.Fatal("result has no content")
	}
	text, ok := res.Content[0].(mcplib.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	var tr internalmcp.ToolResult
	if err := json.Unmarshal([]byte(text.Text), &tr); err != nil {
		t.Fatalf("parse ToolResult JSON: %v (raw: %s)", err, text.Text)
	}
	return tr
}

// TestToolCreateWorkspace_ValidArgs verifies create_workspace with a valid name
// returns a JSON result with files_written and output_dir.
func TestToolCreateWorkspace_ValidArgs(t *testing.T) {
	dir := t.TempDir()
	req := makeReq(map[string]any{
		"name":   "MyWorkspace",
		"output": dir,
	})

	res, err := internalmcp.ToolCreateWorkspace(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %v", res.Content)
	}

	tr := parseToolResult(t, res)
	if len(tr.FilesWritten) == 0 {
		t.Error("expected at least one file_written, got 0")
	}
	if tr.OutputDir != dir {
		t.Errorf("expected output_dir=%q, got %q", dir, tr.OutputDir)
	}
}

// TestToolCreateWorkspace_MissingName verifies that missing 'name' returns a ToolResultError.
func TestToolCreateWorkspace_MissingName(t *testing.T) {
	req := makeReq(map[string]any{})

	res, err := internalmcp.ToolCreateWorkspace(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError=true for missing name, got false")
	}
}

// TestToolCreateWorkspace_ConflictOverwrite verifies that conflict=overwrite
// succeeds even when files already exist.
func TestToolCreateWorkspace_ConflictOverwrite(t *testing.T) {
	dir := t.TempDir()

	// First call to create files
	req1 := makeReq(map[string]any{"name": "WS", "output": dir})
	if _, err := internalmcp.ToolCreateWorkspace(context.Background(), req1); err != nil {
		t.Fatalf("first call: %v", err)
	}

	// Second call with same dir but conflict=overwrite — should succeed
	req2 := makeReq(map[string]any{"name": "WS", "output": dir, "conflict": "overwrite"})
	res, err := internalmcp.ToolCreateWorkspace(context.Background(), req2)
	if err != nil {
		t.Fatalf("second call unexpected Go error: %v", err)
	}
	if res.IsError {
		t.Fatalf("expected success with conflict=overwrite, got error: %v", res.Content)
	}
}

// TestToolCreateConnections_InvalidConfigPath verifies that an invalid config path
// returns a ToolResultError (no panic, no Go error).
func TestToolCreateConnections_InvalidConfigPath(t *testing.T) {
	req := makeReq(map[string]any{
		"config": "/nonexistent/path/config.yaml",
	})

	res, err := internalmcp.ToolCreateConnections(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError=true for bad config path, got false")
	}
}

// TestToolApplyTemplate_UnknownTemplate verifies that an unknown template name
// returns a ToolResultError (no panic, no Go error).
func TestToolApplyTemplate_UnknownTemplate(t *testing.T) {
	req := makeReq(map[string]any{
		"name": "nonexistent-template-xyz",
	})

	res, err := internalmcp.ToolApplyTemplate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError=true for unknown template, got false")
	}
}

// TestToolCreateProject_ValidArgs verifies create_project returns files_written.
func TestToolCreateProject_ValidArgs(t *testing.T) {
	dir := t.TempDir()
	req := makeReq(map[string]any{
		"name":      "ProjectA",
		"workspace": "MyWS",
		"output":    dir,
	})

	res, err := internalmcp.ToolCreateProject(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %v", res.Content)
	}

	tr := parseToolResult(t, res)
	if len(tr.FilesWritten) == 0 {
		t.Error("expected files_written to be non-empty")
	}
}

// TestToolCreateTemplate_ValidArgs verifies create_template creates a template scaffold.
func TestToolCreateTemplate_ValidArgs(t *testing.T) {
	dir := t.TempDir()
	req := makeReq(map[string]any{
		"name":   "my-tmpl",
		"output": dir,
	})

	res, err := internalmcp.ToolCreateTemplate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %v", res.Content)
	}

	// Verify the template.yaml was created
	expectedPath := filepath.Join(dir, "templates", "my-tmpl", "template.yaml")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected template file at %s, got: %v", expectedPath, err)
	}
}

// TestToolGenerateScripts_InvalidConfigPath verifies bad config → ToolResultError.
func TestToolGenerateScripts_InvalidConfigPath(t *testing.T) {
	req := makeReq(map[string]any{
		"config": "/nonexistent/scripts.yaml",
	})

	res, err := internalmcp.ToolGenerateScripts(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError=true for bad config path")
	}
}

// TestToolSyncSnippets_InvalidConfigPath verifies bad config → ToolResultError.
func TestToolSyncSnippets_InvalidConfigPath(t *testing.T) {
	req := makeReq(map[string]any{
		"config": "/nonexistent/snippets.yaml",
	})

	res, err := internalmcp.ToolSyncSnippets(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError=true for bad config path")
	}
}
