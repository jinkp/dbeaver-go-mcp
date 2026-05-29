package dbeaver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// buildSourceTree creates a test source tree under root/<workspaceName>/.
// files is a map of relPath → content.
func buildSourceTree(t *testing.T, root, workspaceName string, files map[string][]byte) {
	t.Helper()
	for relPath, content := range files {
		full := filepath.Join(root, workspaceName, relPath)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// makeSpec returns an InstallSpec wired to temp source + target dirs.
func makeSpec(sourceRoot, targetRoot, wsName string, overwrite, dryRun bool) InstallSpec {
	return InstallSpec{
		WorkspaceName: wsName,
		SourceDir:     sourceRoot,
		TargetRoot:    targetRoot,
		Overwrite:     overwrite,
		DryRun:        dryRun,
	}
}

// ── Plan tests ────────────────────────────────────────────────────────────────

func TestPlan_ProtectedFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	buildSourceTree(t, src, "WS", map[string][]byte{
		".dbeaver/credentials-config.json": []byte(`{}`),
	})
	actions, err := Plan(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Action != "skip-protected" {
		t.Errorf("expected skip-protected, got %q", actions[0].Action)
	}
	if actions[0].RelPath != ".dbeaver/credentials-config.json" {
		t.Errorf("unexpected relPath: %q", actions[0].RelPath)
	}
}

func TestPlan_DataSourcesMerge(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	// Create data-sources.json in both source and target
	buildSourceTree(t, src, "WS", map[string][]byte{
		".dbeaver/data-sources.json": []byte(`{"connections":{},"folders":{}}`),
	})
	// Target workspace .dbeaver already has data-sources.json
	targetDS := filepath.Join(dst, "WS", ".dbeaver")
	_ = os.MkdirAll(targetDS, 0o755)
	_ = os.WriteFile(filepath.Join(targetDS, "data-sources.json"), []byte(`{"connections":{},"folders":{}}`), 0o644)

	actions, err := Plan(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d: %v", len(actions), actions)
	}
	if actions[0].Action != "merge" {
		t.Errorf("expected merge action, got %q", actions[0].Action)
	}
}

func TestPlan_NewSQLFile_Copy(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	buildSourceTree(t, src, "WS", map[string][]byte{
		"Scripts/health-check.sql": []byte(`SELECT 1;`),
	})
	actions, err := Plan(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Action != "copy" {
		t.Errorf("expected copy action, got %q", actions[0].Action)
	}
}

func TestPlan_ExistingFile_Skip(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	buildSourceTree(t, src, "WS", map[string][]byte{
		"Scripts/existing.sql": []byte(`SELECT 1;`),
	})
	// Pre-create the target file
	targetPath := filepath.Join(dst, "WS", "Scripts", "existing.sql")
	_ = os.MkdirAll(filepath.Dir(targetPath), 0o755)
	_ = os.WriteFile(targetPath, []byte(`SELECT 2;`), 0o644)

	actions, err := Plan(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Action != "skip" {
		t.Errorf("expected skip action, got %q", actions[0].Action)
	}
}

func TestPlan_SourceMissing_Error(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	// Do NOT create the workspace source dir
	_, err := Plan(makeSpec(src, dst, "NonExistent", false, false))
	if err == nil {
		t.Error("expected error for missing source, got nil")
	}
}

// ── Execute tests ─────────────────────────────────────────────────────────────

func TestExecute_FileCopied(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	content := []byte(`SELECT 1;`)
	buildSourceTree(t, src, "WS", map[string][]byte{
		"Scripts/check.sql": content,
	})

	_, err := Execute(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dst, "WS", "Scripts", "check.sql"))
	if err != nil {
		t.Fatalf("file was not created: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestExecute_DataSourcesMerged(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	existingDS := `{"connections":{"abc123":{"name":"user-conn"}},"folders":{}}`
	incomingDS := `{"connections":{"def456":{"name":"new-conn"}},"folders":{}}`

	buildSourceTree(t, src, "WS", map[string][]byte{
		".dbeaver/data-sources.json": []byte(incomingDS),
	})
	// Pre-create existing data-sources.json in target
	targetDS := filepath.Join(dst, "WS", ".dbeaver")
	_ = os.MkdirAll(targetDS, 0o755)
	_ = os.WriteFile(filepath.Join(targetDS, "data-sources.json"), []byte(existingDS), 0o644)

	_, err := Execute(makeSpec(src, dst, "WS", false, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := os.ReadFile(filepath.Join(targetDS, "data-sources.json"))
	if err != nil {
		t.Fatalf("data-sources.json not found: %v", err)
	}

	var ds dataSources
	if err := json.Unmarshal(result, &ds); err != nil {
		t.Fatalf("invalid JSON in merged result: %v", err)
	}
	if _, ok := ds.Connections["abc123"]; !ok {
		t.Error("abc123 (user conn) was lost after merge")
	}
	if _, ok := ds.Connections["def456"]; !ok {
		t.Error("def456 (new conn) was not added")
	}
}

func TestExecute_Idempotency(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	existingDS := `{"connections":{"abc123":{"name":"user-conn"}},"folders":{}}`
	incomingDS := `{"connections":{"def456":{"name":"new-conn"}},"folders":{}}`

	buildSourceTree(t, src, "WS", map[string][]byte{
		".dbeaver/data-sources.json": []byte(incomingDS),
	})
	targetDS := filepath.Join(dst, "WS", ".dbeaver")
	_ = os.MkdirAll(targetDS, 0o755)
	_ = os.WriteFile(filepath.Join(targetDS, "data-sources.json"), []byte(existingDS), 0o644)

	spec := makeSpec(src, dst, "WS", false, false)

	// First execute
	_, err := Execute(spec)
	if err != nil {
		t.Fatalf("first execute error: %v", err)
	}
	first, _ := os.ReadFile(filepath.Join(targetDS, "data-sources.json"))

	// Second execute (idempotency check)
	_, err = Execute(spec)
	if err != nil {
		t.Fatalf("second execute error: %v", err)
	}
	second, _ := os.ReadFile(filepath.Join(targetDS, "data-sources.json"))

	// Parse both and compare connection keys
	var ds1, ds2 dataSources
	_ = json.Unmarshal(first, &ds1)
	_ = json.Unmarshal(second, &ds2)

	if len(ds1.Connections) != len(ds2.Connections) {
		t.Errorf("idempotency violated: first=%d conns, second=%d conns", len(ds1.Connections), len(ds2.Connections))
	}
}

func TestExecute_DryRun_NoWrites(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	buildSourceTree(t, src, "WS", map[string][]byte{
		"Scripts/check.sql": []byte(`SELECT 1;`),
	})

	spec := makeSpec(src, dst, "WS", false, true) // dryRun = true
	actions, err := Execute(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Actions should be planned
	if len(actions) == 0 {
		t.Error("expected at least one action in dry-run plan")
	}

	// But the file should NOT be written to target
	targetFile := filepath.Join(dst, "WS", "Scripts", "check.sql")
	if _, err := os.Stat(targetFile); err == nil {
		t.Error("dry-run should not have written any files, but target file exists")
	}
}

func TestBackupWorkspace(t *testing.T) {
	src := t.TempDir()
	wsName := "MyWorkspace"
	files := map[string][]byte{
		".dbeaver/data-sources.json":  []byte(`{}`),
		"Scripts/check.sql":           []byte(`SELECT 1;`),
		"Scripts/nested/advanced.sql": []byte(`SELECT 2;`),
	}
	buildSourceTree(t, src, wsName, files)

	err := BackupWorkspace(src, wsName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find the backup dir (it has a timestamp suffix)
	entries, _ := os.ReadDir(src)
	var backupDir string
	for _, e := range entries {
		if e.IsDir() && e.Name() != wsName {
			backupDir = filepath.Join(src, e.Name())
			break
		}
	}
	if backupDir == "" {
		t.Fatal("no backup directory was created")
	}

	// Verify all files are present in the backup
	for relPath := range files {
		backupFile := filepath.Join(backupDir, relPath)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			t.Errorf("backup missing file: %s", relPath)
		}
	}
}
