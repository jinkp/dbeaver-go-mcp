// Package dbeaver provides domain functions for DBeaver workspace management.
package dbeaver

import (
	"path/filepath"
	"strings"
)

// protectedPaths lists file paths that must never be modified during installation.
// These contain sensitive or DBeaver-managed state.
var protectedPaths = []string{
	".dbeaver/credentials-config.json",
	".dbeaver/.settings",
	".metadata",
	".project",
}

// IsProtected returns true if relPath matches any protected path exactly or by prefix.
// Uses filepath.ToSlash normalization for cross-platform comparison.
func IsProtected(relPath string) bool {
	normalized := filepath.ToSlash(relPath)
	for _, p := range protectedPaths {
		if normalized == p || strings.HasPrefix(normalized, p+"/") {
			return true
		}
	}
	return false
}
