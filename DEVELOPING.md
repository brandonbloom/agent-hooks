# Developing agent-hooks

Developer and contributor guide for the agent-hooks CLI tool.

## Design Philosophy

- **Zero-config by default**: Works without project-specific configuration
- **Respects existing config**: Will use project configurations when present
- **Evolution over magic**: Hard-coded knowledge of tools and practices, 
  not dynamic learning
- **Human and AI friendly**: Designed for both manual and automated use
- **Silence is golden**: Commands produce no output on success, only warnings 
  and errors when needed
- **Do no harm**: Prefer warnings over failures, only fail on critical errors

## Project Structure

```
agent-hooks/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Root command setup
│   ├── which_vcs.go       # VCS detection subcommand
│   ├── format.go          # Format subcommand
│   └── doctor.go          # Environment diagnostics subcommand
├── internal/
│   ├── vcs/
│   │   └── detector.go    # VCS detection logic
│   ├── git/
│   │   └── status.go      # Git operations
│   ├── format/
│   │   └── formatter.go   # Code formatting logic
│   └── doctor/
│       ├── tools.go       # Development tool checks
│       └── claude.go      # Claude Code setup validation
├── go.mod                 # Go module
├── go.sum                 # Dependencies
├── README.md              # User documentation
└── DEVELOPING.md          # This file
```

## Architecture

### Command Structure

The CLI uses the Cobra framework for command organization:

- **Root command** (`cmd/root.go`): Main entry point and command registration
- **Subcommands**: Each subcommand is in its own file under `cmd/`
- **Internal packages**: Business logic is separated into focused modules

### VCS Detection

The `internal/vcs` package provides version control system detection:

- **Extensible design**: Easy to add support for other VCS systems
- **Git support**: Searches parent directories for `.git` directory
- **Future-proof**: Returns typed VCS constants for type safety

### Git Operations

The `internal/git` package handles Git-specific operations:

- **File status parsing**: Uses `git status --porcelain` for changed files
- **Tracked file listing**: Uses `git ls-files` for all tracked files
- **Clean output parsing**: Robust handling of Git command output

### Formatting System

The `internal/format` package provides the extensible formatting system:

- **Language-agnostic interface**: Easy to add new formatters
- **Result aggregation**: Collects formatted files, warnings, and errors
- **Tool availability checking**: Verifies required tools are installed
- **Graceful degradation**: Warns about unsupported files instead of failing

### Diagnostics System

The `internal/doctor` package provides environment and setup validation:

- **Generalized tool checking**: Configurable list of development tools to verify
- **Claude Code integration validation**: Comprehensive Claude settings and hook validation
- **Silence is golden**: Only shows problems by default, verbose mode shows all checks
- **Actionable feedback**: Specific error messages with guidance on fixing issues

## Building and Testing

### Building

```bash
go build -o agent-hooks main.go
```

### Running Tests

```bash
go test ./...
```

### Local Development

Run commands directly with:

```bash
go run main.go [command]
```

## Adding New Commands

1. Create a new file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Define the command using Cobra patterns
3. Add the command to `cmd/root.go` in the `init()` function
4. Create supporting packages in `internal/` as needed

### Example Command Template

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCommand = &cobra.Command{
    Use:   "new-command",
    Short: "Brief description",
    Long:  `Longer description of what this command does.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation here
        return nil
    },
}
```

## Adding New Formatters

1. Add the file extension to `supportedExtensions` in `internal/format/formatter.go`
2. Create a new formatter function following the pattern of `formatGoFiles`
3. Add the formatter call to `FormatFiles` function
4. Implement tool availability checking

### Formatter Requirements

- Check if the required tool is available before attempting to format
- Return descriptive errors that help users understand what went wrong
- Follow the "silence is golden" principle: no output on success
- Handle individual file failures gracefully

## Testing Approach

### Manual Testing Scenarios

1. **No changed files**: `agent-hooks format` should exit silently
2. **Changed Go files**: Should format only changed `.go` files
3. **All files**: `agent-hooks format --all-files` should format all tracked Go files
4. **Mixed file types**: Should format supported files and warn about unsupported
5. **Missing tools**: Should fail with helpful error message
6. **Non-Git repository**: Should fail with VCS error

### Test Data Setup

Create test scenarios by:
- Making changes to Go files
- Adding unsupported file types
- Testing in non-Git directories

## Contributing Guidelines

This is a personal tool, but contributions are welcome:

1. Follow the existing code style and patterns
2. Maintain the "silence is golden" philosophy
3. Add appropriate error handling and user-friendly messages
4. Test both success and failure scenarios
5. Update documentation for new features

## Claude Code Integration

When developing for Claude Code integration:

- Commands should be fast and lightweight
- Error messages should be actionable
- Consider both human and AI agent usage patterns
- Test hook integration scenarios

### Hook Development Tips

- Test hooks with both successful and failing scenarios
- Ensure hook failures don't break Claude Code workflows
- Consider timing implications for large codebases
- Document expected hook behavior clearly