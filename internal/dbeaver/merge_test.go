package dbeaver

import (
	"encoding/json"
	"testing"
)

func TestMergeDataSources(t *testing.T) {
	// Helper to build a data-sources JSON with given connections and folders.
	buildDS := func(conns map[string]string, folders map[string]string) []byte {
		c := map[string]json.RawMessage{}
		for k, v := range conns {
			c[k] = json.RawMessage(`{"name":"` + v + `"}`)
		}
		f := map[string]json.RawMessage{}
		for k, v := range folders {
			f[k] = json.RawMessage(`{"description":"` + v + `"}`)
		}
		b, _ := json.Marshal(map[string]interface{}{
			"connections": c,
			"folders":     f,
		})
		return b
	}

	// Helper to unmarshal result connections map.
	parseConns := func(b []byte) map[string]json.RawMessage {
		var ds dataSources
		_ = json.Unmarshal(b, &ds)
		return ds.Connections
	}

	parseFolders := func(b []byte) map[string]json.RawMessage {
		var ds dataSources
		_ = json.Unmarshal(b, &ds)
		return ds.Folders
	}

	t.Run("existing ID preserved when both have it (existing wins)", func(t *testing.T) {
		existing := buildDS(map[string]string{"abc123": "existing-name"}, nil)
		incoming := buildDS(map[string]string{"abc123": "incoming-name"}, nil)
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		conns := parseConns(result)
		var entry map[string]interface{}
		_ = json.Unmarshal(conns["abc123"], &entry)
		if entry["name"] != "existing-name" {
			t.Errorf("expected existing-name, got %v", entry["name"])
		}
	})

	t.Run("new incoming ID is added", func(t *testing.T) {
		existing := buildDS(map[string]string{"abc123": "conn-a"}, nil)
		incoming := buildDS(map[string]string{"def456": "conn-b"}, nil)
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		conns := parseConns(result)
		if _, ok := conns["def456"]; !ok {
			t.Error("expected def456 to be added from incoming")
		}
		if _, ok := conns["abc123"]; !ok {
			t.Error("expected abc123 to be preserved from existing")
		}
	})

	t.Run("existing-only ID NOT deleted", func(t *testing.T) {
		existing := buildDS(map[string]string{"abc123": "conn-a", "only-existing": "conn-x"}, nil)
		incoming := buildDS(map[string]string{"abc123": "conn-a"}, nil)
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		conns := parseConns(result)
		if _, ok := conns["only-existing"]; !ok {
			t.Error("expected only-existing to be preserved, but it was removed")
		}
	})

	t.Run("folder union additive", func(t *testing.T) {
		existing := buildDS(nil, map[string]string{"DEV": "dev env"})
		incoming := buildDS(nil, map[string]string{"PROD": "prod env"})
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		folders := parseFolders(result)
		if _, ok := folders["DEV"]; !ok {
			t.Error("expected DEV folder preserved")
		}
		if _, ok := folders["PROD"]; !ok {
			t.Error("expected PROD folder added")
		}
	})

	t.Run("no duplicate folders - existing wins", func(t *testing.T) {
		existing := buildDS(nil, map[string]string{"DEV": "existing-dev"})
		incoming := buildDS(nil, map[string]string{"DEV": "incoming-dev"})
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		folders := parseFolders(result)
		var f map[string]interface{}
		_ = json.Unmarshal(folders["DEV"], &f)
		if f["description"] != "existing-dev" {
			t.Errorf("expected existing-dev, got %v", f["description"])
		}
	})

	t.Run("idempotency: merge(merge(e,i), i) == merge(e,i)", func(t *testing.T) {
		existing := buildDS(map[string]string{"abc123": "conn-a"}, map[string]string{"DEV": "dev"})
		incoming := buildDS(map[string]string{"def456": "conn-b"}, map[string]string{"PROD": "prod"})

		first, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("first merge error: %v", err)
		}
		second, err := MergeDataSources(first, incoming)
		if err != nil {
			t.Fatalf("second merge error: %v", err)
		}

		// Compare connection keys
		firstConns := parseConns(first)
		secondConns := parseConns(second)
		if len(firstConns) != len(secondConns) {
			t.Errorf("idempotency violated: first has %d conns, second has %d", len(firstConns), len(secondConns))
		}
		for k := range firstConns {
			if _, ok := secondConns[k]; !ok {
				t.Errorf("key %q disappeared after second merge", k)
			}
		}
	})

	t.Run("empty existing → result equals incoming", func(t *testing.T) {
		existing := []byte(`{"connections":{},"folders":{}}`)
		incoming := buildDS(map[string]string{"def456": "conn-b"}, map[string]string{"PROD": "prod"})
		result, err := MergeDataSources(existing, incoming)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		conns := parseConns(result)
		if _, ok := conns["def456"]; !ok {
			t.Error("expected def456 from incoming")
		}
	})

	t.Run("malformed existing returns error", func(t *testing.T) {
		_, err := MergeDataSources([]byte(`not json`), buildDS(nil, nil))
		if err == nil {
			t.Error("expected error for malformed existing JSON")
		}
	})

	t.Run("malformed incoming returns error", func(t *testing.T) {
		existing := buildDS(nil, nil)
		_, err := MergeDataSources(existing, []byte(`not json`))
		if err == nil {
			t.Error("expected error for malformed incoming JSON")
		}
	})
}
