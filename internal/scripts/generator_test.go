package scripts_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
)

// allScripts lists the 7 recognised script names per the spec.
var allScripts = []string{
	"health-check",
	"migration-validation",
	"compare-row-counts",
	"deadlocks-check",
	"long-running-queries",
	"table-size-analysis",
	"failed-jobs",
}

// TestGenerate_AllScriptNamesRecognised verifies every known script name
// produces output files for a given engine.
func TestGenerate_AllScriptNamesRecognised(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql"},
		Scripts: allScripts,
	}

	files, err := scripts.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("Generate returned no files")
	}
	if len(files) != len(allScripts) {
		t.Errorf("expected %d files, got %d", len(allScripts), len(files))
	}
}

// TestGenerate_OutputPathFormat verifies that output paths match
// `Scripts/scripts/<engine>/<script-name>.sql`.
func TestGenerate_OutputPathFormat(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql"},
		Scripts: []string{"health-check"},
	}

	files, err := scripts.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no files generated")
	}

	expected := "Scripts/scripts/postgresql/health-check.sql"
	got := filepath.ToSlash(files[0].Path)
	if got != expected {
		t.Errorf("expected path %q, got %q", expected, got)
	}
}

// TestGenerate_ContentIsNonEmpty verifies that every generated SQL file
// has actual content.
func TestGenerate_ContentIsNonEmpty(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql"},
		Scripts: []string{"health-check", "deadlocks-check"},
	}

	files, err := scripts.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	for _, f := range files {
		if len(f.Content) == 0 {
			t.Errorf("file %q has empty content", f.Path)
		}
	}
}

// TestGenerate_UnknownScriptReturnsError verifies that requesting a script
// that doesn't exist returns an error.
func TestGenerate_UnknownScriptReturnsError(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql"},
		Scripts: []string{"nonexistent-script"},
	}

	_, err := scripts.Generate(spec)
	if err == nil {
		t.Fatal("expected error for unknown script, got nil")
	}
}

// TestGenerate_UnknownEngineReturnsError verifies that requesting an unknown
// engine returns an error.
func TestGenerate_UnknownEngineReturnsError(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"teradata"},
		Scripts: []string{"health-check"},
	}

	_, err := scripts.Generate(spec)
	if err == nil {
		t.Fatal("expected error for unknown engine, got nil")
	}
}

// TestGenerate_MultipleEngines verifies that multiple engines each produce
// their own set of files.
func TestGenerate_MultipleEngines(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql", "sqlserver"},
		Scripts: []string{"health-check"},
	}

	files, err := scripts.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files (one per engine), got %d", len(files))
	}

	hasPg := false
	hasSql := false
	for _, f := range files {
		p := filepath.ToSlash(f.Path)
		if strings.Contains(p, "/postgresql/") {
			hasPg = true
		}
		if strings.Contains(p, "/sqlserver/") {
			hasSql = true
		}
	}
	if !hasPg {
		t.Error("missing postgresql output")
	}
	if !hasSql {
		t.Error("missing sqlserver output")
	}
}

// TestGenerate_EmptyScriptsGeneratesAll verifies that an empty Scripts slice
// means "generate all known scripts" for each engine.
func TestGenerate_EmptyScriptsGeneratesAll(t *testing.T) {
	spec := scripts.ScriptSpec{
		Project: "Test",
		Engines: []string{"postgresql"},
		Scripts: nil,
	}

	files, err := scripts.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// Should produce all 7 scripts for postgresql
	if len(files) < len(allScripts) {
		t.Errorf("expected at least %d files for all scripts, got %d", len(allScripts), len(files))
	}
}
