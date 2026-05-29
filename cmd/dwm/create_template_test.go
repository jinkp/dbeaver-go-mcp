package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreateTemplateCommandRegistered verifies the create-template command is registered.
func TestCreateTemplateCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "create-template <name>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("create-template command not registered on rootCmd")
	}
}

// TestCreateTemplateScaffoldsYAML verifies that create-template produces a
// template.yaml scaffold in the output/templates/<name>/ directory.
func TestCreateTemplateScaffoldsYAML(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"create-template", "my-template",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-template failed: %v", err)
	}

	yamlPath := filepath.Join(tmpDir, "templates", "my-template", "template.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Errorf("expected %s to exist: %v", yamlPath, err)
	}
}

// TestCreateTemplateYAMLHasName verifies the scaffold YAML contains the template name.
func TestCreateTemplateYAMLHasName(t *testing.T) {
	tmpDir := t.TempDir()

	_ = rootCmd.PersistentFlags().Set("dry-run", "false")
	_ = rootCmd.PersistentFlags().Set("overwrite", "false")

	rootCmd.SetArgs([]string{
		"create-template", "named-tpl",
		"--output", tmpDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-template failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "templates", "named-tpl", "template.yaml"))
	if err != nil {
		t.Fatalf("read template.yaml: %v", err)
	}

	content := string(data)
	if !containsSubstr(content, "named-tpl") {
		t.Errorf("template.yaml does not contain the template name 'named-tpl':\n%s", content)
	}
}

// containsSubstr is a helper to check if s contains substr.
func containsSubstr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
