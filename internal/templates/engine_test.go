package templates_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/templates"
)

// TestLoad_BuiltinMultiEnv verifies that the built-in "multi-env" template
// can be loaded and has the expected metadata.
func TestLoad_BuiltinMultiEnv(t *testing.T) {
	tmpl, err := templates.Load("multi-env")
	if err != nil {
		t.Fatalf("Load('multi-env') error: %v", err)
	}
	if tmpl == nil {
		t.Fatal("Load returned nil template")
	}
	if tmpl.Name != "multi-env" {
		t.Errorf("expected Name='multi-env', got %q", tmpl.Name)
	}
	if len(tmpl.Variables) == 0 {
		t.Error("expected at least one variable, got none")
	}
	if len(tmpl.Generates) == 0 {
		t.Error("expected at least one generate spec, got none")
	}
}

// TestLoad_UnknownTemplateReturnsError verifies that requesting a non-existent
// template returns a descriptive error.
func TestLoad_UnknownTemplateReturnsError(t *testing.T) {
	_, err := templates.Load("does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown template, got nil")
	}
}

// TestScaffold_CreatesTemplateYAML verifies that Scaffold creates a
// template.yaml file in the target directory.
func TestScaffold_CreatesTemplateYAML(t *testing.T) {
	dir := t.TempDir()

	err := templates.Scaffold("my-custom-tmpl", dir)
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}

	yamlPath := filepath.Join(dir, "templates", "my-custom-tmpl", "template.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("template.yaml not created at %s: %v", yamlPath, err)
	}
	if len(data) == 0 {
		t.Error("template.yaml is empty")
	}
}

// TestScaffold_DifferentName verifies that Scaffold uses the provided name
// in the created YAML (triangulation: different name, different directory).
func TestScaffold_DifferentName(t *testing.T) {
	dir := t.TempDir()

	err := templates.Scaffold("analytics-platform", dir)
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}

	yamlPath := filepath.Join(dir, "templates", "analytics-platform", "template.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Fatalf("expected template.yaml at %s, got: %v", yamlPath, err)
	}
}

// TestApply_WithVarsProducesFiles verifies that Apply renders template variables
// and returns output files from the template's generates spec.
func TestApply_WithVarsProducesFiles(t *testing.T) {
	tmpl, err := templates.Load("multi-env")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	vars := map[string]any{
		"Project":      "PaymentService",
		"Environments": []string{"DEV", "PROD"},
		"DBEngine":     "postgresql",
	}

	files, err := templates.Apply(tmpl, vars)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("Apply returned no files")
	}
}

// TestApply_WorkspaceFilesPresent verifies that Apply produces workspace
// marker files when the template includes a workspace generate spec.
func TestApply_WorkspaceFilesPresent(t *testing.T) {
	tmpl, err := templates.Load("multi-env")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	vars := map[string]any{
		"Project":  "BankingCore",
		"DBEngine": "sqlserver",
	}

	files, err := templates.Apply(tmpl, vars)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}

	// At least workspace files should be present (.gitkeep markers).
	hasGitkeep := false
	for _, f := range files {
		if filepath.Ext(f.Path) == "" {
			// Check by looking for .gitkeep files
		}
		_ = f
	}

	// We need at least some files
	if len(files) == 0 {
		t.Fatal("Apply produced no files")
	}

	// Verify at least one .gitkeep file is present (workspace structure)
	for _, f := range files {
		if filepath.Base(f.Path) == ".gitkeep" {
			hasGitkeep = true
			break
		}
	}
	if !hasGitkeep {
		t.Error("expected at least one .gitkeep workspace marker file")
	}
}
