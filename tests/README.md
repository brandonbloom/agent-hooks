# Agent-Hooks Test Suite

A simple test harness for testing the agent-hooks CLI tool.

## Architecture

The test harness consists of:

- **Test Runner** (`runner/main.go`) - Go program that executes test scenarios
- **Test Scenarios** - Directories containing `test.json` configuration files
- **Runner Script** (`run-tests.sh`) - Bash script to discover and run all tests

## Test Scenarios

Each test scenario is a directory containing:

- `test.json` - Test configuration (see format below)
- Optional additional files needed for the test

### Test Configuration Format

```json
{
  "name": "Human readable test name",
  "description": "Description of what this test validates",
  "setup": [
    "command1",
    "command2"
  ],
  "commands": [
    {
      "command": "which-vcs",
      "expectedExit": 0,
      "expectedOut": "git",
      "expectedErr": ""
    }
  ],
  "files": {
    "filename.txt": "file content",
    "another.md": "# Header\nContent"
  }
}
```

**Fields:**
- `setup` - Shell commands to run before testing (e.g., `git init`)
- `commands` - Array of agent-hooks commands to test
- `files` - Files to create in the test directory
- `expectedExit` - Expected exit code (0 = success, non-zero = error)
- `expectedOut` - Expected stdout content (exact match)
- `expectedErr` - Expected stderr content (exact match)

## Running Tests

### Run All Tests
```bash
cd tests
./run-tests.sh
```

### Run Individual Test
```bash
go run tests/runner/main.go tests/vcs-detection-git
```

## Current Test Scenarios

### VCS Detection Tests

- **vcs-detection-git** - Tests VCS detection in a git repository
  - Sets up a git repo with `git init` 
  - Runs `which-vcs` command
  - Expects output "git" with exit code 0

- **vcs-detection-none** - Tests VCS detection outside git repository  
  - No VCS setup
  - Runs `which-vcs` command
  - Expects error with exit code 1

## Design Principles

Following the agent-hooks design philosophy:

- **Simple and lightweight** - Minimal test framework overhead
- **Isolated environments** - Each test runs in a temporary directory
- **Real CLI testing** - Tests the actual `agent-hooks` commands via `go run`
- **Clear failure reporting** - Specific error messages for debugging

## Adding New Tests

1. Create a new directory under `tests/`
2. Add a `test.json` configuration file
3. Run the test suite to verify it works

Example:
```bash
mkdir tests/my-new-test
cat > tests/my-new-test/test.json << 'EOF'
{
  "name": "My New Test",
  "description": "Tests something important",
  "setup": [],
  "commands": [
    {
      "command": "version",
      "expectedExit": 0,
      "expectedOut": "",
      "expectedErr": ""
    }
  ],
  "files": {}
}
EOF
```