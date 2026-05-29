package snippets_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
)

// TestGenerate_SupportedEnginesProduceFiles verifies that each supported
// engine produces at least one snippet file.
func TestGenerate_SupportedEnginesProduceFiles(t *testing.T) {
	engines := []string{"postgresql", "sqlserver", "mysql", "shared"}

	for _, eng := range engines {
		t.Run(eng, func(t *testing.T) {
			spec := snippets.SnippetSpec{Engines: []string{eng}}
			files, err := snippets.Generate(spec)
			if err != nil {
				t.Fatalf("Generate(%q) error: %v", eng, err)
			}
			if len(files) == 0 {
				t.Errorf("Generate(%q) returned no files", eng)
			}
		})
	}
}

// TestGenerate_OutputPathFormat verifies output path format:
// Scripts/snippets/<engine>/<snippet-name>.sql
func TestGenerate_OutputPathFormat(t *testing.T) {
	spec := snippets.SnippetSpec{Engines: []string{"postgresql"}}
	files, err := snippets.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no files returned")
	}

	for _, f := range files {
		p := filepath.ToSlash(f.Path)
		if !strings.HasPrefix(p, "Scripts/snippets/postgresql/") {
			t.Errorf("unexpected path format: %q", p)
		}
		if !strings.HasSuffix(p, ".sql") {
			t.Errorf("expected .sql suffix: %q", p)
		}
	}
}

// TestGenerate_UnknownEngineReturnsError verifies that an unknown engine
// name returns an error.
func TestGenerate_UnknownEngineReturnsError(t *testing.T) {
	spec := snippets.SnippetSpec{Engines: []string{"neo4j"}}
	_, err := snippets.Generate(spec)
	if err == nil {
		t.Fatal("expected error for unknown engine, got nil")
	}
}

// TestGenerate_ContentIsNonEmpty verifies each snippet file has actual content.
func TestGenerate_ContentIsNonEmpty(t *testing.T) {
	spec := snippets.SnippetSpec{Engines: []string{"postgresql", "sqlserver"}}
	files, err := snippets.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	for _, f := range files {
		if len(f.Content) == 0 {
			t.Errorf("file %q has empty content", f.Path)
		}
	}
}

// TestGenerate_PostgreSQLMinimumSnippets verifies postgresql has at least
// the four required snippets: list-tables, list-indexes, table-row-count, active-connections.
func TestGenerate_PostgreSQLMinimumSnippets(t *testing.T) {
	spec := snippets.SnippetSpec{Engines: []string{"postgresql"}}
	files, err := snippets.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	required := []string{"list-tables.sql", "list-indexes.sql", "table-row-count.sql", "active-connections.sql"}
	foundNames := make(map[string]bool)
	for _, f := range files {
		parts := strings.Split(filepath.ToSlash(f.Path), "/")
		foundNames[parts[len(parts)-1]] = true
	}

	for _, req := range required {
		if !foundNames[req] {
			t.Errorf("missing required snippet: %q", req)
		}
	}
}

// TestGenerate_SharedSnippets verifies that shared snippets produce files
// under Scripts/snippets/shared/.
func TestGenerate_SharedSnippets(t *testing.T) {
	spec := snippets.SnippetSpec{Engines: []string{"shared"}}
	files, err := snippets.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no shared snippet files returned")
	}
	for _, f := range files {
		p := filepath.ToSlash(f.Path)
		if !strings.HasPrefix(p, "Scripts/snippets/shared/") {
			t.Errorf("unexpected shared path: %q", p)
		}
	}
}
