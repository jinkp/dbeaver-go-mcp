// Package scripts generates SQL script files from embedded templates.
package scripts

import (
	"fmt"
	"io/fs"
	"path"

	dbeavergoMCP "github.com/jinkp/dbeaver-go-mcp"
)

// knownScripts lists all supported script names in canonical order.
var knownScripts = []string{
	"health-check",
	"migration-validation",
	"compare-row-counts",
	"deadlocks-check",
	"long-running-queries",
	"table-size-analysis",
	"failed-jobs",
}

// knownEngines lists supported engine names.
var knownEngines = map[string]bool{
	"postgresql": true,
	"sqlserver":  true,
	"mysql":      true,
}

// OutFile represents a generated output file.
type OutFile struct {
	Path    string
	Content []byte
}

// ScriptSpec describes what to generate.
type ScriptSpec struct {
	Project string
	Engines []string
	Scripts []string // subset; if empty, generate all
}

// Generate returns SQL script files from embedded templates.
// Output path format: Scripts/scripts/<engine>/<script-name>.sql
// Returns error if any engine or script name is not recognised.
func Generate(spec ScriptSpec) ([]OutFile, error) {
	// Resolve scripts list.
	scriptNames := spec.Scripts
	if len(scriptNames) == 0 {
		scriptNames = knownScripts
	}

	// Validate all requested engines.
	for _, eng := range spec.Engines {
		if !knownEngines[eng] {
			return nil, fmt.Errorf("scripts.Generate: unsupported engine %q", eng)
		}
	}

	// Validate all requested script names.
	knownSet := make(map[string]bool, len(knownScripts))
	for _, s := range knownScripts {
		knownSet[s] = true
	}
	for _, s := range scriptNames {
		if !knownSet[s] {
			return nil, fmt.Errorf("scripts.Generate: unknown script %q", s)
		}
	}

	var files []OutFile

	for _, eng := range spec.Engines {
		for _, scriptName := range scriptNames {
			templatePath := path.Join("templates", "sql", eng, scriptName+".sql")

			content, err := fs.ReadFile(dbeavergoMCP.SQLTemplates, templatePath)
			if err != nil {
				// Engine doesn't have this script — skip gracefully.
				continue
			}

			files = append(files, OutFile{
				Path:    path.Join("Scripts", "scripts", eng, scriptName+".sql"),
				Content: content,
			})
		}
	}

	// If requested specific scripts but none found (because engine supports none), that's fine.
	// But if we requested specific scripts and ALL are unknown, we already caught it above.

	return files, nil
}
