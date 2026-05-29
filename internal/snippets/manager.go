// Package snippets generates SQL snippet files from embedded templates.
package snippets

import (
	"fmt"
	"io/fs"
	"path"

	dbeavergoMCP "github.com/jinkp/dbeaver-go-mcp"
)

// knownEngines lists supported snippet engines including "shared".
var knownEngines = map[string]bool{
	"postgresql": true,
	"sqlserver":  true,
	"mysql":      true,
	"shared":     true,
}

// OutFile represents a generated output file.
type OutFile struct {
	Path    string
	Content []byte
}

// SnippetSpec describes which engines to generate snippets for.
type SnippetSpec struct {
	Engines []string // postgresql, sqlserver, mysql, shared
}

// Generate returns SQL snippet files from embedded templates.
// Output path format: Scripts/snippets/<engine>/<snippet-name>.sql
// Returns error if any engine is not recognised.
func Generate(spec SnippetSpec) ([]OutFile, error) {
	// Validate engines.
	for _, eng := range spec.Engines {
		if !knownEngines[eng] {
			return nil, fmt.Errorf("snippets.Generate: unsupported engine %q", eng)
		}
	}

	var files []OutFile

	for _, eng := range spec.Engines {
		dirPath := path.Join("templates", "snippets", eng)

		entries, err := fs.ReadDir(dbeavergoMCP.SnippetTemplates, dirPath)
		if err != nil {
			return nil, fmt.Errorf("snippets.Generate: read dir %q: %w", dirPath, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := path.Join(dirPath, entry.Name())
			content, err := fs.ReadFile(dbeavergoMCP.SnippetTemplates, filePath)
			if err != nil {
				return nil, fmt.Errorf("snippets.Generate: read %q: %w", filePath, err)
			}

			files = append(files, OutFile{
				Path:    path.Join("Scripts", "snippets", eng, entry.Name()),
				Content: content,
			})
		}
	}

	return files, nil
}
