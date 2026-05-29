// Package workspace generates the folder structure for a DBeaver workspace.
package workspace

import (
	"path"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
)

// dirs are the four directories that every DBeaver workspace must contain.
var dirs = []string{
	".dbeaver",
	"Scripts",
	"Diagrams",
	"Bookmarks",
}

// Generate returns a slice of marker files that represent the DBeaver workspace
// directory structure rooted at name. Empty-content .gitkeep files are used so
// that the directories can be tracked by git without requiring actual files.
//
// The name parameter may contain path separators to express nesting
// (e.g. "corp-ws/ProjectA").
func Generate(name string, _ []string) []output.File {
	files := make([]output.File, 0, len(dirs))
	for _, dir := range dirs {
		files = append(files, output.File{
			Path:    path.Join(name, dir, ".gitkeep"),
			Content: []byte{},
		})
	}
	return files
}

// GenerateProject returns marker files for a DBeaver project nested inside an
// existing workspace. The output paths are <wsName>/<projectName>/<dir>/.gitkeep,
// producing a directory structure compatible with DBeaver's multi-project workspaces.
func GenerateProject(wsName, projectName string) []output.File {
	return Generate(path.Join(wsName, projectName), nil)
}
