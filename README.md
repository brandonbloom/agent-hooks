# agent-hooks

Zero-config convenience commands for humans and AI agents.

## Overview

`agent-hooks` is a personal CLI tool designed to provide convenient, 
zero-configuration commands that both humans and AI agents can invoke as 
needed. The tool respects existing project configurations when available but 
requires no special setup to be useful.

## Installation

### From Source

```bash
git clone https://github.com/brandonbloom/agent-hooks.git
cd agent-hooks
go build -o agent-hooks main.go
```

Move the binary to a directory in your PATH:

```bash
mv agent-hooks /usr/local/bin/
```

### Direct Installation

```bash
go install github.com/brandonbloom/agent-hooks@latest
```

## Usage

### Available Commands

#### `which-vcs`

Detects which version control system is being used in the current directory 
or any parent directory.

```bash
agent-hooks which-vcs
# Output: git
```

Currently supports:
- Git (detects `.git` directory)

#### `format`

Formats code in the current project using appropriate tools. By default, only 
formats files that have changed according to Git status. Requires a Git 
repository to operate.

```bash
agent-hooks format              # Format only changed files
agent-hooks format --all-files  # Format all tracked files
```

**Supported Formatters:**
- **Go**: Uses `go fmt` for `.go` files

**Philosophy**: "Silence is golden" - produces no output on success, warnings 
to stderr for unsupported file types, and only fails on critical errors 
(e.g., missing formatter tools).

### Help

```bash
agent-hooks --help
agent-hooks [command] --help
```

## Claude Code Hooks Integration

This tool is designed to work seamlessly with Claude Code's hooks feature. 
Hooks allow you to run shell commands at various points in Claude Code's 
lifecycle.

### Example Hook Configuration

Add this to your `~/.claude/settings.json` to automatically format code 
before file edits:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Edit",
        "hooks": [
          {
            "type": "command", 
            "command": "agent-hooks format"
          }
        ]
      }
    ]
  }
}
```

### Hook Use Cases

- **Pre-commit formatting**: Run `agent-hooks format` before commits
- **Project detection**: Use `agent-hooks which-vcs` to determine project type
- **Automated workflows**: Chain multiple agent-hooks commands for complex 
  workflows

## Design Philosophy

- **Zero-config by default**: Works without project-specific configuration
- **Respects existing config**: Will use project configurations when present
- **Evolution over magic**: Hard-coded knowledge of tools and practices, 
  not dynamic learning
- **Human and AI friendly**: Designed for both manual and automated use
- **Silence is golden**: Commands produce no output on success, only warnings 
  and errors when needed
- **Do no harm**: Prefer warnings over failures, only fail on critical errors

## Development

### Project Structure

```
agent-hooks/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Root command setup
│   ├── which_vcs.go       # VCS detection subcommand
│   └── format.go          # Format subcommand
├── internal/
│   ├── vcs/
│   │   └── detector.go    # VCS detection logic
│   ├── git/
│   │   └── status.go      # Git operations
│   └── format/
│       └── formatter.go   # Code formatting logic
├── go.mod                 # Go module
├── go.sum                 # Dependencies
└── README.md              # This file
```

### Building

```bash
go build -o agent-hooks main.go
```

### Testing

```bash
go test ./...
```

## Contributing

This is a personal tool, but contributions are welcome. The tool is designed 
to evolve with hard-coded knowledge of specific tools and practices rather 
than attempting dynamic learning.

## License

MIT License - see the source repository for details.