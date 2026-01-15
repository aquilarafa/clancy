package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	// No file specified, look for *.jsonl in current directory
	matches, err := filepath.Glob("*.jsonl")
	if err != nil || len(matches) == 0 {
		return ""
	}

	// If multiple files, pick the most recently modified
	if len(matches) == 1 {
		return matches[0]
	}

	var newest string
	var newestTime int64
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.ModTime().Unix() > newestTime {
			newestTime = info.ModTime().Unix()
			newest = match
		}
	}
	return newest
}
