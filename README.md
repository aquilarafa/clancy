# Clancy

Go TUI for viewing Claude Code Ralph sessions with live reload. Parses `--output-format stream-json` output.

## Installation

### Build from source

Requires Go 1.21+

```bash
git clone https://github.com/aquila/clancy.git
cd clancy
go build -o clancy
```

Optionally, move to your PATH:

```bash
sudo mv clancy /usr/local/bin/
```

## Commands

```bash
# View a specific JSONL file
clancy session.jsonl

# Using --file flag
clancy --file session.jsonl
clancy -f session.jsonl

# Auto-detect *.jsonl in current directory (picks most recent)
clancy
```
