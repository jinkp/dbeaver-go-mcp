// Package mcp provides the MCP server transport layer for dwm.
package mcp

import (
	"encoding/json"

	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
)

// ToolResult is the JSON shape returned by every MCP tool call.
type ToolResult struct {
	FilesWritten []string `json:"files_written"`
	OutputDir    string   `json:"output_dir"`
}

// MarshalToolResult serializes a ToolResult to JSON string, ignoring marshal errors
// (impossible for this simple struct).
func MarshalToolResult(r ToolResult) string {
	b, _ := json.Marshal(r)
	return string(b)
}

// ConflictModeFromString converts a string conflict mode to output.ConflictMode.
// Valid values: "overwrite", "skip". Anything else (including "", "error")
// maps to output.ConflictError.
func ConflictModeFromString(s string) output.ConflictMode {
	switch s {
	case "overwrite":
		return output.ConflictOverwrite
	case "skip":
		return output.ConflictSkip
	default:
		return output.ConflictError
	}
}

// ConnectionFilesToOutput converts connections.OutFile slice to output.File slice.
func ConnectionFilesToOutput(in []connections.OutFile) []output.File {
	out := make([]output.File, len(in))
	for i, f := range in {
		out[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return out
}

// ScriptFilesToOutput converts scripts.OutFile slice to output.File slice.
func ScriptFilesToOutput(in []scripts.OutFile) []output.File {
	out := make([]output.File, len(in))
	for i, f := range in {
		out[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return out
}

// SnippetFilesToOutput converts snippets.OutFile slice to output.File slice.
func SnippetFilesToOutput(in []snippets.OutFile) []output.File {
	out := make([]output.File, len(in))
	for i, f := range in {
		out[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return out
}

// TemplateFilesToOutput converts templates.OutFile slice to output.File slice.
func TemplateFilesToOutput(in []templates.OutFile) []output.File {
	out := make([]output.File, len(in))
	for i, f := range in {
		out[i] = output.File{Path: f.Path, Content: f.Content}
	}
	return out
}
