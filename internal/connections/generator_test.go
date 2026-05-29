package connections_test

import (
	"encoding/json"
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
)

// TestGenerate_StableConnectionID verifies that the same input always
// produces the same connection ID (deterministic SHA-256 based).
func TestGenerate_StableConnectionID(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "MyApp",
		Environments: []string{"DEV"},
		Connections: []connections.ConnConfig{
			{Name: "{{.Project}}-{{.Env}}", Type: "postgresql", Host: "localhost", Port: 5432, Database: "mydb"},
		},
	}

	files1, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	files2, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate returned error on second call: %v", err)
	}

	if len(files1) == 0 || len(files2) == 0 {
		t.Fatal("Generate returned no files")
	}

	// Parse both data-sources.json and verify same connection IDs
	ds1 := parseDataSources(t, files1)
	ds2 := parseDataSources(t, files2)

	ids1 := connectionIDs(ds1)
	ids2 := connectionIDs(ds2)

	if len(ids1) != len(ids2) {
		t.Fatalf("different number of connection IDs: %v vs %v", ids1, ids2)
	}
	for i := range ids1 {
		if ids1[i] != ids2[i] {
			t.Errorf("ID mismatch at index %d: %q vs %q", i, ids1[i], ids2[i])
		}
	}
}

// TestGenerate_SavePasswordAlwaysFalse verifies the security constraint.
func TestGenerate_SavePasswordAlwaysFalse(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "SecureApp",
		Environments: []string{"PROD"},
		Connections: []connections.ConnConfig{
			{Name: "prod-db", Type: "sqlserver", Host: "sqlserver.prod", Port: 1433, Database: "mydb"},
		},
	}

	files, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	ds := parseDataSources(t, files)
	conns, ok := ds["connections"].(map[string]interface{})
	if !ok || len(conns) == 0 {
		t.Fatal("no connections in data-sources.json")
	}

	for id, raw := range conns {
		conn, ok := raw.(map[string]interface{})
		if !ok {
			t.Errorf("connection %q is not a map", id)
			continue
		}
		savePassword, exists := conn["save-password"]
		if !exists {
			t.Errorf("connection %q missing 'save-password' field", id)
			continue
		}
		if savePassword != false {
			t.Errorf("connection %q has save-password=%v, want false", id, savePassword)
		}
	}
}

// TestGenerate_DataSourcesConfigAlwaysEmitted verifies data-sources-config.json
// is always present and contains exactly `{}`.
func TestGenerate_DataSourcesConfigAlwaysEmitted(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "AnyProject",
		Environments: []string{"DEV"},
		Connections: []connections.ConnConfig{
			{Name: "dev-pg", Type: "postgresql", Host: "localhost", Port: 5432, Database: "mydb"},
		},
	}

	files, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var configFile *connections.OutFile
	for i := range files {
		if files[i].Path == "data-sources-config.json" {
			configFile = &files[i]
			break
		}
	}

	if configFile == nil {
		t.Fatal("data-sources-config.json not found in output")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(configFile.Content, &parsed); err != nil {
		t.Fatalf("data-sources-config.json is not valid JSON: %v", err)
	}
	if len(parsed) != 0 {
		t.Errorf("data-sources-config.json should be empty object {}, got: %v", parsed)
	}
}

// TestGenerate_UnknownEngineReturnsError verifies that unknown engine types
// cause Generate to return a descriptive error.
func TestGenerate_UnknownEngineReturnsError(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "BadProject",
		Environments: []string{"DEV"},
		Connections: []connections.ConnConfig{
			{Name: "bad-conn", Type: "oracle-xe-unknown", Host: "localhost", Port: 1521, Database: "xe"},
		},
	}

	_, err := connections.Generate(spec)
	if err == nil {
		t.Fatal("expected error for unknown engine, got nil")
	}
}

// TestGenerate_NameTemplateRendering verifies {{.Project}} and {{.Env}} substitution.
func TestGenerate_NameTemplateRendering(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "Fintech",
		Environments: []string{"QA"},
		Connections: []connections.ConnConfig{
			{Name: "{{.Project}}-{{.Env}}-pg", Type: "postgresql", Host: "pg.qa", Port: 5432, Database: "finance"},
		},
	}

	files, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	ds := parseDataSources(t, files)
	conns, ok := ds["connections"].(map[string]interface{})
	if !ok || len(conns) == 0 {
		t.Fatal("no connections in data-sources.json")
	}

	found := false
	for _, raw := range conns {
		conn, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := conn["name"].(string); ok && name == "Fintech-QA-pg" {
			found = true
			break
		}
	}
	if !found {
		t.Error("template rendering failed: expected name 'Fintech-QA-pg' not found in connections")
	}
}

// TestGenerate_FolderMatchesEnv verifies that each connection is placed in a
// folder named after its environment.
func TestGenerate_FolderMatchesEnv(t *testing.T) {
	spec := connections.ConnSpec{
		Project:      "BiApp",
		Environments: []string{"STAGE"},
		Connections: []connections.ConnConfig{
			{Name: "stage-db", Type: "mysql", Host: "mysql.stage", Port: 3306, Database: "bi"},
		},
	}

	files, err := connections.Generate(spec)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	ds := parseDataSources(t, files)
	conns, ok := ds["connections"].(map[string]interface{})
	if !ok || len(conns) == 0 {
		t.Fatal("no connections in data-sources.json")
	}

	for id, raw := range conns {
		conn, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		folder, _ := conn["folder"].(string)
		if folder != "STAGE" {
			t.Errorf("connection %q: expected folder='STAGE', got %q", id, folder)
		}
	}
}

// TestGenerate_DifferentInputsProduceDifferentIDs verifies that distinct
// connection specs (different name, env, or type) produce distinct connection IDs.
// The ID is derived from provider|driver|name|env, so changes in any of those fields
// must produce a different hash.
func TestGenerate_DifferentInputsProduceDifferentIDs(t *testing.T) {
	// Two specs with all four ID components different: different name, env, and type.
	specA := connections.ConnSpec{
		Project:      "App",
		Environments: []string{"DEV"},
		Connections: []connections.ConnConfig{
			{Name: "conn-A", Type: "postgresql", Host: "localhost", Port: 5432, Database: "db"},
		},
	}
	specB := connections.ConnSpec{
		Project:      "App",
		Environments: []string{"PROD"}, // different env → different ID
		Connections: []connections.ConnConfig{
			{Name: "conn-B", Type: "sqlserver", Host: "localhost", Port: 1433, Database: "db"}, // different name+type → different ID
		},
	}

	filesA, err := connections.Generate(specA)
	if err != nil {
		t.Fatalf("Generate(specA) error: %v", err)
	}
	filesB, err := connections.Generate(specB)
	if err != nil {
		t.Fatalf("Generate(specB) error: %v", err)
	}

	dsA := parseDataSources(t, filesA)
	dsB := parseDataSources(t, filesB)

	idsA := connectionIDs(dsA)
	idsB := connectionIDs(dsB)

	if len(idsA) == 0 || len(idsB) == 0 {
		t.Fatal("expected at least one connection ID in each spec")
	}
	for _, idA := range idsA {
		for _, idB := range idsB {
			if idA == idB {
				t.Errorf("distinct connection specs produced the same ID %q — IDs collide", idA)
			}
		}
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

// OutFile is a type alias exported by the connections package for tests.
// (We reference it via connections.OutFile to keep the type in one place.)

func parseDataSources(t *testing.T, files []connections.OutFile) map[string]interface{} {
	t.Helper()
	for _, f := range files {
		if f.Path == "data-sources.json" {
			var m map[string]interface{}
			if err := json.Unmarshal(f.Content, &m); err != nil {
				t.Fatalf("data-sources.json parse error: %v\ncontent: %s", err, f.Content)
			}
			return m
		}
	}
	t.Fatal("data-sources.json not found")
	return nil
}

func connectionIDs(ds map[string]interface{}) []string {
	conns, ok := ds["connections"].(map[string]interface{})
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(conns))
	for id := range conns {
		ids = append(ids, id)
	}
	return ids
}
