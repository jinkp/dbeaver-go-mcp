package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/output"
)

// TestConflictError verifies that Write returns an error when a file already exists
// and ConflictMode is ConflictError.
func TestConflictError(t *testing.T) {
	baseDir := t.TempDir()

	// Pre-create the file so a conflict exists
	existingPath := filepath.Join(baseDir, "some", "file.txt")
	if err := os.MkdirAll(filepath.Dir(existingPath), 0o755); err != nil {
		t.Fatalf("setup: mkdir: %v", err)
	}
	if err := os.WriteFile(existingPath, []byte("original"), 0o644); err != nil {
		t.Fatalf("setup: write existing file: %v", err)
	}

	files := []output.File{
		{Path: "some/file.txt", Content: []byte("new content")},
	}

	err := output.Write(files, baseDir, output.ConflictError)

	if err == nil {
		t.Fatal("expected error for ConflictError mode, got nil")
	}
	if !strings.Contains(err.Error(), "some/file.txt") {
		t.Errorf("expected error to mention path 'some/file.txt', got: %v", err)
	}

	// Original file must be untouched
	got, _ := os.ReadFile(existingPath)
	if string(got) != "original" {
		t.Errorf("ConflictError must not modify existing file, but got: %s", got)
	}
}

// TestConflictOverwrite verifies that Write replaces the existing file when
// ConflictMode is ConflictOverwrite.
func TestConflictOverwrite(t *testing.T) {
	baseDir := t.TempDir()

	existingPath := filepath.Join(baseDir, "replace_me.txt")
	if err := os.WriteFile(existingPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files := []output.File{
		{Path: "replace_me.txt", Content: []byte("replaced")},
	}

	if err := output.Write(files, baseDir, output.ConflictOverwrite); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("read replaced file: %v", err)
	}
	if string(got) != "replaced" {
		t.Errorf("expected 'replaced', got: %s", got)
	}
}

// TestConflictSkip verifies that Write skips an existing file (without error) and
// still writes non-conflicting files.
func TestConflictSkip(t *testing.T) {
	baseDir := t.TempDir()

	existingPath := filepath.Join(baseDir, "keep_me.txt")
	if err := os.WriteFile(existingPath, []byte("keep"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files := []output.File{
		{Path: "keep_me.txt", Content: []byte("would overwrite")},
		{Path: "new_file.txt", Content: []byte("brand new")},
	}

	if err := output.Write(files, baseDir, output.ConflictSkip); err != nil {
		t.Fatalf("expected no error for ConflictSkip, got: %v", err)
	}

	// Existing file must be untouched
	kept, _ := os.ReadFile(existingPath)
	if string(kept) != "keep" {
		t.Errorf("ConflictSkip must not change existing file, got: %s", kept)
	}

	// New file must be written
	newFile := filepath.Join(baseDir, "new_file.txt")
	written, err := os.ReadFile(newFile)
	if err != nil {
		t.Fatalf("new_file.txt was not written: %v", err)
	}
	if string(written) != "brand new" {
		t.Errorf("expected 'brand new', got: %s", written)
	}
}

// TestCreatesParentDirectories verifies that Write creates intermediate directories
// even when they do not exist.
func TestCreatesParentDirectories(t *testing.T) {
	baseDir := t.TempDir()

	files := []output.File{
		{Path: "a/b/c/deep.txt", Content: []byte("deep content")},
	}

	if err := output.Write(files, baseDir, output.ConflictError); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(baseDir, "a", "b", "c", "deep.txt"))
	if err != nil {
		t.Fatalf("deep file not found: %v", err)
	}
	if string(got) != "deep content" {
		t.Errorf("expected 'deep content', got: %s", got)
	}
}
