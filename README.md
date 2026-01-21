# Clancy

![Ralph and Chief Wiggum](https://static0.cbrimages.com/wordpress/wp-content/uploads/2021/01/The-Simpsons-Ralph-Chief-Wiggum-2.jpg?q=50&fit=crop&w=1232&h=693&dpr=1.5)

TUI for viewing Claude Code sessions with live reload.

## Installation

Requires Go 1.21+

```bash
git clone https://github.com/aquila/clancy.git
cd clancy
go build -o clancy
sudo mv clancy /usr/local/bin/  # optional
```

## Usage

```bash
# Inside any repo: automatically opens the most recent session
clancy

# Or specify a file
clancy file.jsonl
```

When run without arguments, Clancy searches for sessions in order:

1. `~/.claude/projects/<current-repo>/` - saved Claude Code sessions
2. `*.jsonl` in current directory

## Keybindings

- `↑/↓` or `j/k` - Navigate messages
- `q` or `Ctrl+C` - Quit
