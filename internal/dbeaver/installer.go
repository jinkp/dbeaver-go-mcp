package dbeaver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// InstallSpec holds all parameters for a workspace install operation.
type InstallSpec struct {
	WorkspaceName string
	SourceDir     string // directory containing the workspace folder (default: "./output")
	TargetRoot    string // DBeaver workspace6/ path (from WorkspacePath() or --workspace-path)
	Overwrite     bool   // replace existing files instead of merging/skipping
	DryRun        bool   // plan only; do not write any files
	Force         bool   // bypass DBeaver-is-running check
	Backup        bool   // back up target workspace before writing
}

// InstallAction describes a single file operation planned or executed by Plan/Execute.
type InstallAction struct {
	Action  string // "copy" | "merge" | "skip" | "skip-protected" | "backup"
	RelPath string // slash-normalized relative path from workspace root
	Reason  string // human-readable explanation (populated for skip/skip-protected)
}

// Plan returns the list of InstallActions that would be performed for spec,
// without actually writing anything to disk.
func Plan(spec InstallSpec) ([]InstallAction, error) {
	sourceRoot := filepath.Join(spec.SourceDir, spec.WorkspaceName)
	if _, err := os.Stat(sourceRoot); os.IsNotExist(err) {
		return nil, fmt.Errorf("source not found: %s", sourceRoot)
	}
	targetRoot := filepath.Join(spec.TargetRoot, spec.WorkspaceName)

	var actions []InstallAction
	err := filepath.WalkDir(sourceRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		relPath, _ := filepath.Rel(sourceRoot, path)
		relSlash := filepath.ToSlash(relPath)

		// Protected files are always skipped.
		if IsProtected(relSlash) {
			actions = append(actions, InstallAction{
				Action:  "skip-protected",
				RelPath: relSlash,
				Reason:  "protected file",
			})
			return nil
		}

		targetPath := filepath.Join(targetRoot, relPath)
		_, statErr := os.Stat(targetPath)
		targetExists := statErr == nil

		// data-sources.json: merge by default unless --overwrite or target missing.
		if relSlash == ".dbeaver/data-sources.json" && targetExists && !spec.Overwrite {
			actions = append(actions, InstallAction{
				Action:  "merge",
				RelPath: relSlash,
				Reason:  "merge connections by ID",
			})
			return nil
		}

		// Other existing files: skip unless --overwrite.
		if targetExists && !spec.Overwrite {
			actions = append(actions, InstallAction{
				Action:  "skip",
				RelPath: relSlash,
				Reason:  "already exists",
			})
			return nil
		}

		actions = append(actions, InstallAction{
			Action:  "copy",
			RelPath: relSlash,
		})
		return nil
	})
	return actions, err
}

// Execute runs the install plan described by spec.
// If spec.DryRun is true, it returns the plan without writing anything.
// If spec.Backup is true, it backs up the target workspace first.
func Execute(spec InstallSpec) ([]InstallAction, error) {
	if spec.Backup {
		if err := BackupWorkspace(spec.TargetRoot, spec.WorkspaceName); err != nil {
			return nil, fmt.Errorf("backup failed: %w", err)
		}
	}

	actions, err := Plan(spec)
	if err != nil {
		return nil, err
	}
	if spec.DryRun {
		return actions, nil
	}

	sourceRoot := filepath.Join(spec.SourceDir, spec.WorkspaceName)
	targetRoot := filepath.Join(spec.TargetRoot, spec.WorkspaceName)

	for _, a := range actions {
		src := filepath.Join(sourceRoot, filepath.FromSlash(a.RelPath))
		dst := filepath.Join(targetRoot, filepath.FromSlash(a.RelPath))

		switch a.Action {
		case "copy":
			if err := copyFile(src, dst); err != nil {
				return actions, err
			}
		case "merge":
			existingBytes, err := os.ReadFile(dst)
			if err != nil {
				return actions, err
			}
			incomingBytes, err := os.ReadFile(src)
			if err != nil {
				return actions, err
			}
			merged, err := MergeDataSources(existingBytes, incomingBytes)
			if err != nil {
				return actions, err
			}
			if err := os.WriteFile(dst, merged, 0o644); err != nil {
				return actions, err
			}
		// "skip" and "skip-protected": nothing to do.
		}
	}
	return actions, nil
}

// BackupWorkspace copies the workspace directory at filepath.Join(root, name)
// to a timestamped sibling directory (e.g., <name>-backup-20060102-150405).
func BackupWorkspace(root, name string) error {
	src := filepath.Join(root, name)
	ts := time.Now().Format("20060102-150405")
	dst := filepath.Join(root, name+"-backup-"+ts)
	return copyDir(src, dst)
}

// ── internal helpers ──────────────────────────────────────────────────────────

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}
