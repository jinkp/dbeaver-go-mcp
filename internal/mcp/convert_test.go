package mcp_test

import (
	"testing"

	internalmcp "github.com/jinkp/dbeaver-go-mcp/internal/mcp"
	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
)

// ── conflictModeFromString ────────────────────────────────────────────────────

func TestConflictModeFromString_Error(t *testing.T) {
	got := internalmcp.ConflictModeFromString("error")
	if got != output.ConflictError {
		t.Errorf("expected ConflictError, got %v", got)
	}
}

func TestConflictModeFromString_Overwrite(t *testing.T) {
	got := internalmcp.ConflictModeFromString("overwrite")
	if got != output.ConflictOverwrite {
		t.Errorf("expected ConflictOverwrite, got %v", got)
	}
}

func TestConflictModeFromString_Skip(t *testing.T) {
	got := internalmcp.ConflictModeFromString("skip")
	if got != output.ConflictSkip {
		t.Errorf("expected ConflictSkip, got %v", got)
	}
}

func TestConflictModeFromString_Empty(t *testing.T) {
	// Empty string defaults to ConflictError
	got := internalmcp.ConflictModeFromString("")
	if got != output.ConflictError {
		t.Errorf("expected ConflictError for empty string, got %v", got)
	}
}

func TestConflictModeFromString_Unknown(t *testing.T) {
	// Unknown string defaults to ConflictError
	got := internalmcp.ConflictModeFromString("unknown-value")
	if got != output.ConflictError {
		t.Errorf("expected ConflictError for unknown string, got %v", got)
	}
}

// ── connectionFilesToOutput ───────────────────────────────────────────────────

func TestConnectionFilesToOutput_PreservesPathAndContent(t *testing.T) {
	in := []connections.OutFile{
		{Path: "data-sources.json", Content: []byte(`{"connections":{}}`)},
		{Path: "data-sources-config.json", Content: []byte("{}")},
	}
	got := internalmcp.ConnectionFilesToOutput(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 files, got %d", len(got))
	}
	if got[0].Path != "data-sources.json" {
		t.Errorf("path[0]: want 'data-sources.json', got %q", got[0].Path)
	}
	if string(got[0].Content) != `{"connections":{}}` {
		t.Errorf("content[0]: want '{\"connections\":{}}', got %q", got[0].Content)
	}
	if got[1].Path != "data-sources-config.json" {
		t.Errorf("path[1]: want 'data-sources-config.json', got %q", got[1].Path)
	}
}

func TestConnectionFilesToOutput_Empty(t *testing.T) {
	got := internalmcp.ConnectionFilesToOutput(nil)
	if len(got) != 0 {
		t.Errorf("expected 0 files for nil input, got %d", len(got))
	}
}

// ── scriptFilesToOutput ───────────────────────────────────────────────────────

func TestScriptFilesToOutput_PreservesPathAndContent(t *testing.T) {
	in := []scripts.OutFile{
		{Path: "Scripts/scripts/postgresql/health-check.sql", Content: []byte("SELECT 1;")},
	}
	got := internalmcp.ScriptFilesToOutput(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 file, got %d", len(got))
	}
	if got[0].Path != "Scripts/scripts/postgresql/health-check.sql" {
		t.Errorf("path: want 'Scripts/scripts/postgresql/health-check.sql', got %q", got[0].Path)
	}
	if string(got[0].Content) != "SELECT 1;" {
		t.Errorf("content: want 'SELECT 1;', got %q", got[0].Content)
	}
}

// ── snippetFilesToOutput ──────────────────────────────────────────────────────

func TestSnippetFilesToOutput_PreservesPathAndContent(t *testing.T) {
	in := []snippets.OutFile{
		{Path: "Scripts/snippets/postgresql/list-tables.sql", Content: []byte("SELECT tablename FROM pg_tables;")},
	}
	got := internalmcp.SnippetFilesToOutput(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 file, got %d", len(got))
	}
	if got[0].Path != "Scripts/snippets/postgresql/list-tables.sql" {
		t.Errorf("path: want correct path, got %q", got[0].Path)
	}
	if string(got[0].Content) != "SELECT tablename FROM pg_tables;" {
		t.Errorf("content mismatch, got %q", got[0].Content)
	}
}

// ── templateFilesToOutput ─────────────────────────────────────────────────────

func TestTemplateFilesToOutput_PreservesPathAndContent(t *testing.T) {
	in := []templates.OutFile{
		{Path: "templates/my-template/template.yaml", Content: []byte("name: my-template\n")},
	}
	got := internalmcp.TemplateFilesToOutput(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 file, got %d", len(got))
	}
	if got[0].Path != "templates/my-template/template.yaml" {
		t.Errorf("path: want 'templates/my-template/template.yaml', got %q", got[0].Path)
	}
	if string(got[0].Content) != "name: my-template\n" {
		t.Errorf("content mismatch, got %q", got[0].Content)
	}
}
