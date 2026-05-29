package connections_test

import (
	"testing"

	"github.com/jinkp/dbeaver-go-mcp/internal/connections"
)

// TestLookupAllSixEngines is a table-driven test that pins the Provider and Driver
// values for every supported engine — verified against DBeaver Community 26.0.
func TestLookupAllSixEngines(t *testing.T) {
	cases := []struct {
		engine           string
		expectedProvider string
		expectedDriver   string
	}{
		{"postgresql", "postgresql", "postgres-jdbc"},
		{"sqlserver", "sqlserver", "sqlserver"},
		{"mysql", "mysql", "mysql8"},
		{"oracle", "oracle", "oracle_thin"},
		{"mongodb", "mongodb", "mongodb"},
		{"dynamodb", "dynamodb", "dynamodb"},
	}

	for _, tc := range cases {
		t.Run(tc.engine, func(t *testing.T) {
			entry, err := connections.Lookup(tc.engine)
			if err != nil {
				t.Fatalf("Lookup(%q) returned unexpected error: %v", tc.engine, err)
			}
			if entry.Provider != tc.expectedProvider {
				t.Errorf("Lookup(%q).Provider = %q; want %q", tc.engine, entry.Provider, tc.expectedProvider)
			}
			if entry.Driver != tc.expectedDriver {
				t.Errorf("Lookup(%q).Driver = %q; want %q", tc.engine, entry.Driver, tc.expectedDriver)
			}
		})
	}
}

// TestLookupUnknownEngineReturnsError verifies that Lookup returns an error for
// an engine name not in the registry.
func TestLookupUnknownEngineReturnsError(t *testing.T) {
	unknownEngines := []string{"redis", "cassandra", "", "POSTGRESQL"}

	for _, engine := range unknownEngines {
		t.Run(engine, func(t *testing.T) {
			_, err := connections.Lookup(engine)
			if err == nil {
				t.Errorf("Lookup(%q) expected error, got nil", engine)
			}
		})
	}
}

// TestSupportedEnginesReturnsSortedList verifies that SupportedEngines returns
// all 6 supported engines in sorted (alphabetical) order.
func TestSupportedEnginesReturnsSortedList(t *testing.T) {
	engines := connections.SupportedEngines()

	if len(engines) != 6 {
		t.Fatalf("expected 6 supported engines, got %d: %v", len(engines), engines)
	}

	expected := []string{"dynamodb", "mongodb", "mysql", "oracle", "postgresql", "sqlserver"}
	for i, want := range expected {
		if engines[i] != want {
			t.Errorf("SupportedEngines()[%d] = %q; want %q", i, engines[i], want)
		}
	}
}
