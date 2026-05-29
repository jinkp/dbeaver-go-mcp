// Package output provides file writing with conflict resolution.
package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// File represents a generated output artifact: a relative path and its content.
type File struct {
	Path    string
	Content []byte
}

// ConflictMode controls what happens when a target file already exists.
type ConflictMode int

const (
	// ConflictError returns an error immediately if any target file already exists.
	ConflictError ConflictMode = iota
	// ConflictOverwrite deletes the existing file and writes the new content.
	ConflictOverwrite
	// ConflictSkip logs a warning and skips the conflicting file; other files are still written.
	ConflictSkip
)

// Write writes each file in files to <baseDir>/<file.Path>.
// Parent directories are always created automatically.
// The conflict parameter governs behaviour when a target file already exists.
func Write(files []File, baseDir string, mode ConflictMode) error {
	for _, f := range files {
		target := filepath.Join(baseDir, filepath.FromSlash(f.Path))

		// Always ensure parent directories exist.
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("output.Write: mkdir %s: %w", filepath.Dir(target), err)
		}

		// Check for existing file.
		if _, err := os.Stat(target); err == nil {
			// File exists — apply conflict mode.
			switch mode {
			case ConflictError:
				return fmt.Errorf("output.Write: file already exists: %s", f.Path)

			case ConflictOverwrite:
				if err := os.Remove(target); err != nil {
					return fmt.Errorf("output.Write: remove %s: %w", f.Path, err)
				}

			case ConflictSkip:
				log.Warn().Str("path", f.Path).Msg("output.Write: skipping existing file")
				continue
			}
		}

		if err := os.WriteFile(target, f.Content, 0o644); err != nil {
			return fmt.Errorf("output.Write: write %s: %w", f.Path, err)
		}
	}
	return nil
}
