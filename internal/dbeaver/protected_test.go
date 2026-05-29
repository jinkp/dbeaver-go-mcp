package dbeaver

import "testing"

func TestIsProtected(t *testing.T) {
	tests := []struct {
		name     string
		relPath  string
		expected bool
	}{
		// Each of the 4 protected paths returns true
		{"credentials-config.json exact", ".dbeaver/credentials-config.json", true},
		{"settings exact", ".dbeaver/.settings", true},
		{"metadata exact", ".metadata", true},
		{"project exact", ".project", true},
		// Prefix match: sub-paths under protected dirs also return true
		{"settings sub-path", ".dbeaver/.settings/something.prefs", true},
		{"metadata sub-path", ".metadata/version.ini", true},
		// Non-protected paths return false
		{"data-sources.json not protected", ".dbeaver/data-sources.json", false},
		{"sql script not protected", "Scripts/health-check.sql", false},
		{"root file not protected", "somefile.txt", false},
		{"data-sources-config not protected", ".dbeaver/data-sources-config.json", false},
		// Windows-style backslash paths should also work (normalized)
		{"windows backslash settings", `.dbeaver\.settings`, true},
		{"windows backslash sub-path", `.dbeaver\.settings\prefs.ini`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsProtected(tt.relPath)
			if got != tt.expected {
				t.Errorf("IsProtected(%q) = %v, want %v", tt.relPath, got, tt.expected)
			}
		})
	}
}
