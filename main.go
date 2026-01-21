package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aquila/clancy/ui"
	"github.com/aquila/clancy/watcher"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	filename := findFile()
	if filename == "" {
		fmt.Fprintln(os.Stderr, "Usage: clancy [file.jsonl]")
		fmt.Fprintln(os.Stderr, "       clancy --file file.jsonl")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "If no file specified, looks for *.jsonl in current directory")
		os.Exit(1)
	}

	// Check file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File not found: %s\n", filename)
		os.Exit(1)
	}

	// Create watcher
	w := watcher.New(filename)
	if err := w.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting watcher: %v\n", err)
		os.Exit(1)
	}

	// Create and run UI
	model := ui.New(filename, w)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func findFile() string {
	args := os.Args[1:]

	// Parse args
	for i, arg := range args {
		if arg == "--file" || arg == "-f" {
			if i+1 < len(args) {
				return args[i+1]
			}
			return ""
		}
		if arg == "--help" || arg == "-h" {
			return ""
		}
		// First non-flag argument is the file
		if arg[0] != '-' {
			return arg
		}
	}

	// No file specified, first check Claude sessions directory
	if file := findClaudeSessionFile(); file != "" {
		return file
	}

	// Fall back to looking for *.jsonl in current directory
	matches, err := filepath.Glob("*.jsonl")
	if err != nil || len(matches) == 0 {
		return ""
	}

	return findNewestFile(matches)
}

// findClaudeSessionFile looks for the most recent .jsonl in ~/.claude/projects/<project-dir>/
func findClaudeSessionFile() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Convert path to Claude format: /Users/aquila/Projects/foo -> -Users-aquila-Projects-foo
	projectDir := strings.ReplaceAll(cwd, "/", "-")
	if !strings.HasPrefix(projectDir, "-") {
		projectDir = "-" + projectDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	claudeProjectPath := filepath.Join(homeDir, ".claude", "projects", projectDir)

	// Check if directory exists
	if _, err := os.Stat(claudeProjectPath); os.IsNotExist(err) {
		return ""
	}

	// Find all .jsonl files in this directory
	matches, err := filepath.Glob(filepath.Join(claudeProjectPath, "*.jsonl"))
	if err != nil || len(matches) == 0 {
		return ""
	}

	return findNewestFile(matches)
}

// findNewestFile returns the most recently modified file from a list
func findNewestFile(files []string) string {
	if len(files) == 1 {
		return files[0]
	}

	var newest string
	var newestTime int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.ModTime().Unix() > newestTime {
			newestTime = info.ModTime().Unix()
			newest = file
		}
	}
	return newest
}
