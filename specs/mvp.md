# Clancy - TUI para Stream JSON do Claude Code

## Overview
TUI em Go + Bubbletea que faz tail -f em arquivo .jsonl gerado pelo Claude Code com `--output-format stream-json`. Roda em terminal separado.

## Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│  CLANCY TUI                                             │
├─────────────────────────────────────────────────────────┤
│  [Status Bar] watching: output.jsonl | 142 events      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ▶ message_start                                        │
│    model: claude-sonnet-4-5 | input: 25 tokens         │
│                                                         │
│  ▶ content_block [text]                                 │
│    "Let me help you with that..."                       │
│                                                         │
│  ▶ content_block [tool_use] Read                        │
│    file_path: "/src/main.go"                            │
│                                                         │
│  ▶ content_block [thinking]                             │
│    "I need to analyze the code structure..."            │
│                                                         │
│  ▶ message_delta                                        │
│    stop_reason: end_turn | output: 156 tokens          │
│                                                         │
├─────────────────────────────────────────────────────────┤
│  [Help] q:quit  ↑↓:scroll  g/G:top/bottom              │
└─────────────────────────────────────────────────────────┘
```

## Estrutura do Projeto

```
clancy/
├── main.go              # entry point, CLI args
├── watcher/
│   └── watcher.go       # tail -f do arquivo jsonl
├── parser/
│   └── parser.go        # parse eventos NDJSON
├── model/
│   └── events.go        # structs dos eventos stream-json
├── ui/
│   ├── app.go           # bubbletea model principal
│   ├── styles.go        # lipgloss styles
│   └── view.go          # renderização
└── go.mod
```

## Componentes

### 1. Watcher (tail -f)
- Usa `fsnotify` ou polling p/ detectar mudanças
- Lê novas linhas incrementalmente
- Emite linhas via channel para o parser

### 2. Parser
- Recebe linha JSON, faz unmarshal
- Identifica tipo do evento pelo campo `type`
- Acumula `partial_json` de tool_use até `content_block_stop`
- Acumula `text_delta` para mostrar texto completo do bloco

### 3. Model (eventos)
```go
type Event struct {
    Type    string          // message_start, content_block_delta, etc
    Index   int             // índice do content_block
    Raw     json.RawMessage // evento original
    Parsed  interface{}     // struct específica do tipo
}

type MessageStart struct {
    Message struct {
        ID        string
        Model     string
        Role      string
        Usage     Usage
    }
}

type ContentBlock struct {
    Type  string // text, tool_use, thinking
    Text  string // acumulado dos deltas
    Tool  *ToolUse
}

type ToolUse struct {
    ID    string
    Name  string
    Input string // JSON acumulado
}
```

### 4. UI (Bubbletea)
- **Model**: lista de eventos parseados, viewport position
- **Update**: recebe novos eventos do watcher via tea.Cmd
- **View**: renderiza lista scrollável com lipgloss

## Fluxo de Dados

```
arquivo.jsonl → Watcher → Parser → UI Model → View
     ↑              ↓          ↓
  [append]     channel     channel
                          tea.Cmd
```

## Uso

```bash
# Terminal 1: rodar claude
cat PROMPT.md | claude -p --output-format stream-json > output.jsonl

# Terminal 2: rodar clancy
clancy output.jsonl
# ou
clancy --file output.jsonl
# ou (default: procura *.jsonl no pwd)
clancy
```

## Implementação - Ordem

1. **Scaffolding**: go mod init, estrutura de pastas
2. **Events model**: structs para todos os tipos de evento
3. **Parser**: unmarshal + acumulação de deltas
4. **Watcher**: tail -f com fsnotify
5. **UI base**: bubbletea app com viewport
6. **Styles**: lipgloss para colorir por tipo
7. **Integração**: watcher → parser → ui via tea.Cmd

## Deps

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - styling
- `github.com/charmbracelet/bubbles/viewport` - scrolling
- `github.com/fsnotify/fsnotify` - file watching

## Verificação

1. Gerar arquivo de teste com eventos mock
2. Rodar `clancy test.jsonl` e verificar renderização
3. Testar append em tempo real: `echo '{"type":"ping"}' >> test.jsonl`
4. Rodar com claude real: `cat PROMPT.md | claude -p --output-format stream-json > out.jsonl`

---

## Decisões

- **Auto-scroll**: ON por default, tecla `f` toggle follow mode
- **Cores**: tool=azul, text=branco, thinking=dim/cinza, error=vermelho, ping=dim
- **Nome**: Clancy

## Keybindings

| Tecla | Ação |
|-------|------|
| `q` | quit |
| `↑/k` | scroll up |
| `↓/j` | scroll down |
| `g` | ir pro topo |
| `G` | ir pro fim |
| `f` | toggle follow mode (auto-scroll) |
| `pgup/pgdn` | página acima/abaixo |
