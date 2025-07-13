#!/bin/bash

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Running agent-hooks test suite..."
echo "Project root: $PROJECT_ROOT"
echo

# Build the test runner
echo "Building test runner..."
go build -o "$SCRIPT_DIR/test-runner" "$SCRIPT_DIR/runner/main.go"

# Find all test directories (directories containing test.json)
test_dirs=()
while IFS= read -r -d '' test_dir; do
    test_dirs+=("$test_dir")
done < <(find "$SCRIPT_DIR" -name "test.json" -type f -print0 | sed -z 's|/test.json||g')

if [ ${#test_dirs[@]} -eq 0 ]; then
    echo "No test directories found!"
    exit 1
fi

echo "Found ${#test_dirs[@]} test(s):"
for test_dir in "${test_dirs[@]}"; do
    echo "  - $(basename "$test_dir")"
done
echo

# Run each test
failed_tests=()
passed_tests=()

for test_dir in "${test_dirs[@]}"; do
    test_name=$(basename "$test_dir")
    echo "Running test: $test_name"
    echo "----------------------------------------"
    
    if cd "$PROJECT_ROOT" && "$SCRIPT_DIR/test-runner" "$test_dir"; then
        echo "✓ PASSED: $test_name"
        passed_tests+=("$test_name")
    else
        echo "✗ FAILED: $test_name"
        failed_tests+=("$test_name")
    fi
    echo
done

# Summary
echo "========================================"
echo "Test Summary:"
echo "  Passed: ${#passed_tests[@]}"
echo "  Failed: ${#failed_tests[@]}"
echo "  Total:  ${#test_dirs[@]}"

if [ ${#failed_tests[@]} -gt 0 ]; then
    echo
    echo "Failed tests:"
    for test in "${failed_tests[@]}"; do
        echo "  - $test"
    done
    exit 1
fi

echo
echo "All tests passed! 🎉"