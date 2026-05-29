// Package templates provides DBeaver workspace template loading and application.
package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	dbeavergoMCP "github.com/jinkp/dbeaver-go-mcp"
	"github.com/jinkp/dbeaver-go-mcp/internal/scripts"
	"github.com/jinkp/dbeaver-go-mcp/internal/snippets"
	"github.com/jinkp/dbeaver-go-mcp/internal/workspace"
	"gopkg.in/yaml.v3"
)

// OutFile represents a generated output file.
type OutFile struct {
	Path    string
	Content []byte
}

// Template describes a dwm project template.
type Template struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Version     string        `yaml:"version"`
	Variables   []Variable    `yaml:"variables"`
	Generates   []GenerateSpec `yaml:"generates"`
}

// Variable describes a template variable.
type Variable struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default"`
	Values      []string    `yaml:"values"`
}

// GenerateSpec describes a generation step in a template.
type GenerateSpec struct {
	Type   string   `yaml:"type"`   // workspace, connections, scripts, snippets
	Subset []string `yaml:"subset"` // optional filter
}

// Load resolves a template by name using this lookup order:
//  1. ./templates/<name>/template.yaml  (project-local override)
//  2. ~/.dwm/templates/<name>/template.yaml  (user-level override)
//  3. Embedded built-in templates (templates/projects/<name>/template.yaml)
//
// Returns an error if the template is not found in any location.
func Load(name string) (*Template, error) {
	// 1. Project-local templates directory.
	localPath := filepath.Join("templates", name, "template.yaml")
	if data, err := os.ReadFile(localPath); err == nil {
		return parseTemplate(data)
	}

	// 2. User home directory templates.
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, ".dwm", "templates", name, "template.yaml")
		if data, err := os.ReadFile(homePath); err == nil {
			return parseTemplate(data)
		}
	}

	// 3. Embedded built-in templates.
	embeddedPath := "templates/projects/" + name + "/template.yaml"
	data, err := fs.ReadFile(dbeavergoMCP.ProjectTemplates, embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("templates.Load: template %q not found (local, home, or built-in)", name)
	}
	return parseTemplate(data)
}

// Scaffold creates a new template definition at targetDir/templates/<name>/template.yaml
// with placeholder content.
func Scaffold(name, targetDir string) error {
	dir := filepath.Join(targetDir, "templates", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("templates.Scaffold: mkdir %s: %w", dir, err)
	}

	placeholder := Template{
		Name:        name,
		Description: "TODO: describe this template",
		Version:     "1.0.0",
		Variables: []Variable{
			{Name: "Project", Description: "Project name", Required: true},
		},
		Generates: []GenerateSpec{
			{Type: "workspace"},
			{Type: "connections"},
			{Type: "scripts"},
			{Type: "snippets"},
		},
	}

	data, err := yaml.Marshal(placeholder)
	if err != nil {
		return fmt.Errorf("templates.Scaffold: marshal: %w", err)
	}

	yamlPath := filepath.Join(dir, "template.yaml")
	if err := os.WriteFile(yamlPath, data, 0o644); err != nil {
		return fmt.Errorf("templates.Scaffold: write %s: %w", yamlPath, err)
	}

	return nil
}

// Apply renders a template with the given variables and returns the list of
// files to write. Variable defaults are applied for any missing keys.
func Apply(tmpl *Template, vars map[string]any) ([]OutFile, error) {
	// Apply defaults for missing variables.
	resolved := applyDefaults(tmpl, vars)

	var result []OutFile

	for _, gen := range tmpl.Generates {
		switch gen.Type {
		case "workspace":
			project, _ := resolved["Project"].(string)
			wsFiles := workspace.Generate(project, nil)
			for _, f := range wsFiles {
				result = append(result, OutFile{Path: f.Path, Content: f.Content})
			}

		case "scripts":
			engine := resolveEngine(resolved)
			subset := gen.Subset
			spec := scripts.ScriptSpec{
				Project: resolveProject(resolved),
				Engines: []string{engine},
				Scripts: subset,
			}
			sFiles, err := scripts.Generate(spec)
			if err != nil {
				return nil, fmt.Errorf("templates.Apply: scripts: %w", err)
			}
			for _, f := range sFiles {
				result = append(result, OutFile{Path: f.Path, Content: f.Content})
			}

		case "snippets":
			engine := resolveEngine(resolved)
			spec := snippets.SnippetSpec{Engines: []string{engine}}
			snFiles, err := snippets.Generate(spec)
			if err != nil {
				return nil, fmt.Errorf("templates.Apply: snippets: %w", err)
			}
			for _, f := range snFiles {
				result = append(result, OutFile{Path: f.Path, Content: f.Content})
			}

		case "connections":
			// Connections require a config file — generate an empty placeholder.
			// Full connection generation is handled by the CLI layer with a config file.
			// Skip silently; the CLI wires connections.Generate separately.

		default:
			return nil, fmt.Errorf("templates.Apply: unknown generate type %q", gen.Type)
		}
	}

	return result, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func parseTemplate(data []byte) (*Template, error) {
	var t Template
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("templates: parse template: %w", err)
	}
	return &t, nil
}

func applyDefaults(tmpl *Template, vars map[string]any) map[string]any {
	resolved := make(map[string]any, len(vars))
	for k, v := range vars {
		resolved[k] = v
	}
	for _, v := range tmpl.Variables {
		if _, ok := resolved[v.Name]; !ok && v.Default != nil {
			resolved[v.Name] = v.Default
		}
	}
	return resolved
}

func resolveProject(vars map[string]any) string {
	if p, ok := vars["Project"].(string); ok {
		return p
	}
	return "project"
}

func resolveEngine(vars map[string]any) string {
	if e, ok := vars["DBEngine"].(string); ok && e != "" {
		return e
	}
	return "postgresql"
}
