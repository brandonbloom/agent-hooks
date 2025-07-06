# agent-hooks

Zero-config convenience commands for humans and AI agents.

## Installation

```bash
go install github.com/brandonbloom/agent-hooks@latest
```

## Commands

### `which-vcs`
Detects which version control system is in use (currently supports Git).

```bash
agent-hooks which-vcs
```

### `format`
Formats changed files (or all files with `--all-files`). Supports Go files.

```bash
agent-hooks format              # Format changed files only
agent-hooks format --all-files  # Format all tracked files
```

## Claude Code Hooks

Automatically format code after file modifications by adding this to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
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

## Contributing

See [DEVELOPING.md](DEVELOPING.md) for development setup and architecture details.