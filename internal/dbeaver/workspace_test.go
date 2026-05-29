package dbeaver

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWorkspacePath(t *testing.T) {
	t.Run("Windows: uses APPDATA env var", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Windows-only test")
		}
		// Create a temp dir to simulate the workspace existing
		tmp := t.TempDir()
		wsPath := filepath.Join(tmp, "DBeaverData", "workspace6")
		if err := os.MkdirAll(wsPath, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("APPDATA", tmp)

		got, err := WorkspacePath()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != wsPath {
			t.Errorf("got %q, want %q", got, wsPath)
		}
	})

	t.Run("returns error if workspace path does not exist", func(t *testing.T) {
		tmp := t.TempDir()
		// Point APPDATA/XDG_DATA_HOME at a dir that has no DBeaverData/workspace6 under it
		switch runtime.GOOS {
		case "windows":
			t.Setenv("APPDATA", tmp)
		case "darwin":
			// WorkspacePath on darwin uses UserHomeDir — not easily mockable via env
			// We test by verifying the error string when the path truly won't exist
			t.Skip("darwin home dir not easily overridable; covered by manual integration")
		default:
			t.Setenv("XDG_DATA_HOME", tmp)
		}

		_, err := WorkspacePath()
		if err == nil {
			t.Error("expected error when workspace does not exist, got nil")
		}
	})

	t.Run("Linux: uses XDG_DATA_HOME when set", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Linux-only test")
		}
		tmp := t.TempDir()
		wsPath := filepath.Join(tmp, "DBeaverData", "workspace6")
		if err := os.MkdirAll(wsPath, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("XDG_DATA_HOME", tmp)

		got, err := WorkspacePath()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != wsPath {
			t.Errorf("got %q, want %q", got, wsPath)
		}
	})
}

func TestListProjects(t *testing.T) {
	t.Run("returns only dirs with .dbeaver/ subfolder", func(t *testing.T) {
		root := t.TempDir()

		// Create project with .dbeaver/
		proj1 := filepath.Join(root, "MyProject")
		if err := os.MkdirAll(filepath.Join(proj1, ".dbeaver"), 0o755); err != nil {
			t.Fatal(err)
		}

		// Create dir without .dbeaver/ — should NOT be listed
		other := filepath.Join(root, "NotAProject")
		if err := os.MkdirAll(other, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create a regular file — should NOT be listed
		if err := os.WriteFile(filepath.Join(root, "somefile.txt"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}

		projects, err := ListProjects(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 1 {
			t.Errorf("expected 1 project, got %d: %v", len(projects), projects)
		}
		if projects[0] != "MyProject" {
			t.Errorf("expected MyProject, got %q", projects[0])
		}
	})

	t.Run("returns empty slice for empty directory", func(t *testing.T) {
		root := t.TempDir()
		projects, err := ListProjects(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 0 {
			t.Errorf("expected 0 projects, got %d", len(projects))
		}
	})

	t.Run("multiple projects returned", func(t *testing.T) {
		root := t.TempDir()
		for _, name := range []string{"Alpha", "Beta", "Gamma"} {
			if err := os.MkdirAll(filepath.Join(root, name, ".dbeaver"), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		// One dir without .dbeaver
		_ = os.MkdirAll(filepath.Join(root, "NotProject"), 0o755)

		projects, err := ListProjects(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 3 {
			t.Errorf("expected 3 projects, got %d: %v", len(projects), projects)
		}
	})
}
