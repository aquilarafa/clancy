# Clancy

Go TUI for viewing Claude Code Ralph sessions with live reload. Parses `--output-format stream-json` output.

## Tech Stack

- Go 1.21+
- Bubble Tea (TUI framework)
- Lipgloss (styling)
- fsnotify (file watching)

## Project Structure

```
main.go      - Entry point, CLI args
ui/          - Bubble Tea model and view
model/       - Data structures
parser/      - JSONL parsing
watcher/     - File change detection
```

## Commands

```bash
# Build
go build -o clancy

# Run
./clancy file.jsonl
./clancy              # auto-finds *.jsonl in cwd

# Test
go test ./...
```

## Conventions

- Follow Bubble Tea patterns: Model, Update, View
- Use Lipgloss for all styling
- Keep Update pure, side effects via Cmd
