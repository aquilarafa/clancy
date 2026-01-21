package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindClaudeSessionFile(t *testing.T) {
	// Use actual home directory to avoid os.UserHomeDir() caching issues
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home dir")
	}

	// Create a unique temp project directory
	tmpProject := t.TempDir()

	// Resolve symlinks (macOS: /var -> /private/var)
	tmpProject, err = filepath.EvalSymlinks(tmpProject)
	if err != nil {
		t.Fatal(err)
	}

	// Expected Claude path conversion
	expectedClaudeDir := convertPathToClaudeFormat(tmpProject)
	claudeProjectPath := filepath.Join(homeDir, ".claude", "projects", expectedClaudeDir)

	// Clean up after test
	defer os.RemoveAll(claudeProjectPath)

	if err := os.MkdirAll(claudeProjectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test jsonl files with different modification times
	oldFile := filepath.Join(claudeProjectPath, "old-session.jsonl")
	newFile := filepath.Join(claudeProjectPath, "new-session.jsonl")

	if err := os.WriteFile(oldFile, []byte(`{"type":"test"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte(`{"type":"test"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Make old file actually older
	oldTime := mustParseTime("2020-01-01T00:00:00Z")
	os.Chtimes(oldFile, oldTime, oldTime)

	// Change to project directory
	originalWd, _ := os.Getwd()
	os.Chdir(tmpProject)
	defer os.Chdir(originalWd)

	// Test: should find the newer session file
	result := findClaudeSessionFile()

	if result == "" {
		t.Fatal("expected to find a session file, got empty string")
	}

	if result != newFile {
		t.Errorf("expected newest file %s, got %s", newFile, result)
	}
}

func TestFindClaudeSessionFileNoClaudeDir(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	// Create project dir but NO .claude directory
	projectPath := filepath.Join(tmpHome, "Projects", "nocloud")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	originalWd, _ := os.Getwd()
	os.Chdir(projectPath)
	defer os.Chdir(originalWd)

	// Should return empty when no Claude sessions exist
	result := findClaudeSessionFile()
	if result != "" {
		t.Errorf("expected empty string when no Claude dir, got %s", result)
	}
}

func TestConvertPathToClaudeFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/Users/aquila/Projects/clarifica", "-Users-aquila-Projects-clarifica"},
		{"/home/user/code/app", "-home-user-code-app"},
		{"/tmp/test", "-tmp-test"},
		{"relative/path", "-relative-path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertPathToClaudeFormat(tt.input)
			if result != tt.expected {
				t.Errorf("convertPathToClaudeFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper to convert path to Claude format (extracted for testing)
func convertPathToClaudeFormat(path string) string {
	result := ""
	for _, c := range path {
		if c == '/' {
			result += "-"
		} else {
			result += string(c)
		}
	}
	if len(result) > 0 && result[0] != '-' {
		result = "-" + result
	}
	return result
}

func mustParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
