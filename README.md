# agent-hooks

Zero-config convenience commands for humans and AI agents.

## Why This Tool Exists

AI agents like Claude are powerful but often inconsistent with development workflows—they might forget to format code, skip linting, or miss other routine tasks. This tool brings **deterministic execution** to AI-assisted development by using hooks to ensure critical steps happen reliably.

### Design Principles

- **Convention over configuration**: Sensible defaults that work across projects, eliminating the need for per-project config files
- **Unix-centric**: Built around standard Unix tools and philosophy rather than being tied to specific programming languages  
- **Get out of your way**: Works silently when successful, provides clear explanations when it fails
- **Low false positive rate**: Like `go fmt`, opinionated enough to maintain consistency but unopinionated enough to avoid conflicts

This tool was built primarily for personal use but may serve as a useful example for others tackling similar workflow automation challenges.

NOTE: I am aggressively "vibe-coding" this project as a learning experience.

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
agent-hooks format                    # Format changed files only
agent-hooks format --all-files       # Format all tracked files  
agent-hooks format --verbose         # Show what files are formatted
agent-hooks format --dry-run         # Preview what would be formatted
agent-hooks format --dry-run -v      # Preview with detailed output
```

### `post-tool-use`
Hook command for Claude Code PostToolUse events. Checks `.agenthooks` configuration and only runs formatting if hooks are not disabled. Use this command in Claude Code hooks instead of calling `format` directly.

```bash
agent-hooks post-tool-use             # For use in Claude Code hooks
```

### `doctor`
Checks development environment and Claude Code setup. Silent by default, shows all checks with `--verbose`.

```bash
agent-hooks doctor              # Only show problems
agent-hooks doctor --verbose    # Show all checks
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
            "command": "agent-hooks post-tool-use"
          }
        ]
      }
    ]
  }
}
```

## Configuration

### Disabling Hook Execution

You can disable automatic execution of agent-hooks when running as Claude Code hooks by creating a `.agenthooks` configuration file in your project root:

```yaml
disable: true
```

This setting only affects automatic execution via Claude Code hooks (when using `agent-hooks post-tool-use`). Manual invocation (running `agent-hooks format` directly in the terminal) will still work normally.

The configuration file is searched in the current directory and parent directories, allowing you to disable hooks at the project level or higher in the directory hierarchy.

## Contributing

See [DEVELOPING.md](DEVELOPING.md) for development setup and architecture details.
