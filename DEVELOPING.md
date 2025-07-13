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
- **Verbose instrumentation**: Verbose modes should instrument the same code 
  rather than doing extra work. Verbosity is about output formatting, not 
  logic or behavior changes

## Project Structure

```
agent-hooks/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Root command setup
│   ├── about.go           # Technology and tool introspection subcommand
│   ├── detect.go          # Technology detection subcommand
│   ├── doctor.go          # Environment diagnostics subcommand
│   ├── format.go          # Format subcommand
│   ├── version.go         # Version information subcommand
│   └── which_vcs.go       # VCS detection subcommand
├── internal/
│   ├── detect/
│   │   ├── detector.go     # Main detection engine
│   │   ├── technologies.go # Technology constants (alphabetical)
│   │   └── rules.go        # Detection rules (alphabetical)
│   ├── vcs/
│   │   └── detector.go     # VCS detection logic
│   ├── git/
│   │   └── status.go       # Git operations
│   ├── format/
│   │   └── formatter.go    # Code formatting logic
│   └── doctor/
│       ├── tools.go        # Development tool checks (alphabetical)
│       ├── claude.go       # Claude Code setup validation
│       ├── requirements.go # Environment requirements
│       └── project.go      # Project-specific checks
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
- **Dry-run support**: Preview mode that shows what would be formatted without changes
- **Verbose output**: Detailed reporting of formatting operations and skipped files

### Diagnostics System

The `internal/doctor` package provides environment and setup validation:

- **Generalized tool checking**: Configurable list of development tools to verify
- **Claude Code integration validation**: Comprehensive Claude settings and hook validation
- **Silence is golden**: Only shows problems by default, verbose mode shows all checks
- **Actionable feedback**: Specific error messages with guidance on fixing issues
- **Reference URLs**: Each tool includes official documentation URL for introspection

### Technology Detection System

The `internal/detect` package provides extensible technology detection:

- **Technology constants** (`technologies.go`): Alphabetically sorted technology identifiers
- **Detection rules** (`rules.go`): File patterns, descriptions, and reference URLs for each technology  
- **VCS-aware detection**: Prioritizes git-tracked files for performance
- **Fallback scanning**: Falls back to directory traversal when VCS unavailable
- **Alphabetical ordering**: All technology lists maintain strict alphabetical order to minimize merge conflicts
- **Reference URLs**: Each technology includes official documentation URL for introspection

#### Key files:
- `internal/detect/technologies.go` - Technology constant definitions
- `internal/detect/rules.go` - Detection rules mapping technologies to file patterns
- `internal/detect/detector.go` - Main detection engine

## Building and Testing

### Building

```bash
go build -o agent-hooks main.go
```

### Running Tests

Unit tests:
```bash
go test ./...
```

CLI integration tests use [transcript](https://github.com/deref/transcript):
```bash
cd tests && ./run-tests.sh
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

## Adding New Technologies

To add support for detecting a new technology:

1. **Add technology constant** to `internal/detect/technologies.go`:
   - Insert in alphabetical order to minimize merge conflicts
   - Use lowercase naming convention

2. **Add detection rule** to `internal/detect/rules.go`:
   - Insert in alphabetical order by technology name
   - Specify file patterns that indicate the technology's presence
   - Provide descriptive text for user-facing output
   - Include official documentation URL for reference

3. **Add tool check** (optional) to `internal/doctor/tools.go`:
   - Add to `DefaultTools` slice in alphabetical order
   - Specify command name and whether it's required
   - Add version detection to `versionArgs` map if needed
   - Include official documentation URL for reference

### Example: Adding Direnv Support

```go
// In internal/detect/technologies.go
Direnv          Technology = "direnv"

// In internal/detect/rules.go  
{Technology: Direnv, Files: []string{".envrc"}, Desc: "Direnv environment configuration", URL: "https://direnv.net"},

// In internal/doctor/tools.go
{Name: "direnv", Command: "direnv", Required: false, URL: "https://direnv.net"},
```

### Important Notes

- **Alphabetical ordering is critical** - All technology collections must maintain strict alphabetical order
- **Detection uses VCS-aware scanning** - Files must be git-tracked to be detected (fallback to directory scan when not in git repo)
- **Test both detection and doctor commands** after adding new technologies

## Testing Approach

### Manual Testing Scenarios

#### Format Command
1. **No changed files**: `agent-hooks format` should exit silently
2. **Changed Go files**: Should format only changed `.go` files
3. **All files**: `agent-hooks format --all-files` should format all tracked Go files
4. **Mixed file types**: Should format supported files and warn about unsupported
5. **Missing tools**: Should fail with helpful error message
6. **Non-Git repository**: Should fail with VCS error

#### Detect Command  
1. **Technology detection**: `agent-hooks detect` should identify all technologies in project
2. **Git-tracked files**: Only git-tracked files are considered for detection
3. **Untracked files**: Create untracked technology files, verify they're not detected until tracked

#### Doctor Command
1. **Tool availability**: `agent-hooks doctor` should check all configured tools
2. **Verbose output**: `agent-hooks doctor --verbose` should show all tool versions
3. **Missing tools**: Remove tools from PATH, verify appropriate warnings

#### About Command
1. **Technology lookup**: `agent-hooks about go` should show technology details
2. **Tool lookup**: `agent-hooks about goimports` should show tool details
3. **Case insensitive**: Should work with any casing (go, Go, GO)
4. **Unknown items**: Should provide clear error for unknown technologies/tools

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