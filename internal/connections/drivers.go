// Package connections provides the DBeaver driver registry and connection utilities.
package connections

import (
	"fmt"
	"sort"
)

// DriverEntry holds the DBeaver provider and JDBC driver identifiers for an engine.
// Values are verified against DBeaver Community 26.0.
type DriverEntry struct {
	Provider string
	Driver   string
}

// registry maps engine names to their DBeaver driver entries.
// All values verified against DBeaver Community 26.0.
var registry = map[string]DriverEntry{
	"postgresql": {Provider: "postgresql", Driver: "postgres-jdbc"},
	"sqlserver":  {Provider: "sqlserver", Driver: "sqlserver"},
	"mysql":      {Provider: "mysql", Driver: "mysql8"},
	"oracle":     {Provider: "oracle", Driver: "oracle_thin"},
	"mongodb":    {Provider: "mongodb", Driver: "mongodb"},
	"dynamodb":   {Provider: "dynamodb", Driver: "dynamodb"},
}

// Lookup returns the DriverEntry for the given engine name.
// It returns an error if the engine is not supported.
func Lookup(engine string) (DriverEntry, error) {
	entry, ok := registry[engine]
	if !ok {
		return DriverEntry{}, fmt.Errorf("connections.Lookup: unsupported engine %q; supported: %v", engine, SupportedEngines())
	}
	return entry, nil
}

// SupportedEngines returns a sorted list of all registered engine names.
func SupportedEngines() []string {
	engines := make([]string, 0, len(registry))
	for k := range registry {
		engines = append(engines, k)
	}
	sort.Strings(engines)
	return engines
}
