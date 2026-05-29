package workspace_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
	"github.com/jinkp/dbeaver-go-mcp/internal/workspace"
)

// collectPaths extracts all Path values from a file slice and sorts them.
func collectPaths(files []output.File) []string {
	paths := make([]string, 0, len(files))
	for _, f := range files {
		paths = append(paths, filepath.ToSlash(f.Path))
	}
	sort.Strings(paths)
	return paths
}

// hasPath returns true if any File in files matches path (slash-normalised).
func hasPath(files []output.File, path string) bool {
	for _, f := range files {
		if filepath.ToSlash(f.Path) == path {
			return true
		}
	}
	return false
}

// TestGenerate_FourDirs verifies that Generate creates the four required
// DBeaver workspace directories represented as .gitkeep marker files.
func TestGenerate_FourDirs(t *testing.T) {
	files := workspace.Generate("myws", nil)

	required := []string{
		"myws/.dbeaver/.gitkeep",
		"myws/Scripts/.gitkeep",
		"myws/Diagrams/.gitkeep",
		"myws/Bookmarks/.gitkeep",
	}

	if len(files) < len(required) {
		t.Fatalf("Generate returned %d files, want at least %d", len(files), len(required))
	}

	for _, want := range required {
		if !hasPath(files, want) {
			t.Errorf("missing required path %q in output; got: %v", want, collectPaths(files))
		}
	}
}

// TestGenerate_DifferentName verifies paths use the supplied workspace name.
func TestGenerate_DifferentName(t *testing.T) {
	files := workspace.Generate("project-alpha", nil)

	if !hasPath(files, "project-alpha/.dbeaver/.gitkeep") {
		t.Errorf("name not applied: expected 'project-alpha/.dbeaver/.gitkeep' in %v", collectPaths(files))
	}
	// Verify no old name leaks
	for _, f := range files {
		if strings.HasPrefix(filepath.ToSlash(f.Path), "myws/") {
			t.Errorf("unexpected path with old name 'myws': %s", f.Path)
		}
	}
}

// TestGenerate_GitkeepContentIsEmpty verifies that marker files have empty content.
func TestGenerate_GitkeepContentIsEmpty(t *testing.T) {
	files := workspace.Generate("ws", nil)
	for _, f := range files {
		if len(f.Content) != 0 {
			t.Errorf("expected empty content for %s, got %d bytes", f.Path, len(f.Content))
		}
	}
}

// TestGenerate_NestedProject verifies that when a workspace name contains a
// project suffix the structure is correctly nested.
func TestGenerate_NestedProject(t *testing.T) {
	// Simulate workspace/project nesting by calling Generate with a nested path.
	files := workspace.Generate("corp-ws/ProjectA", nil)

	required := []string{
		"corp-ws/ProjectA/.dbeaver/.gitkeep",
		"corp-ws/ProjectA/Scripts/.gitkeep",
	}
	for _, want := range required {
		if !hasPath(files, want) {
			t.Errorf("nested project: missing %q; got: %v", want, collectPaths(files))
		}
	}
}
