#!/bin/bash
# Simple test runner for agent-hooks using transcript

set -e

# Get the directory where this script lives
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Build agent-hooks
echo "Building agent-hooks..."
cd "$PROJECT_ROOT"
go build -o agent-hooks main.go

# Find all test.cmdt files
echo "Running tests..."
FAILED=0
PASSED=0

while IFS= read -r test_file; do
    test_dir=$(dirname "$test_file")
    test_name=$(basename "$test_dir")
    echo -n "  $test_name: "
    
    # Create temp directory for this test
    TEMP_DIR=$(mktemp -d)
    
    # Copy all files from the test directory to temp directory
    cp -r "$test_dir"/* "$TEMP_DIR/"
    
    # Run the test with transcript in the temp directory
    if cd "$TEMP_DIR" && PATH="$PROJECT_ROOT:$PATH" go run github.com/deref/transcript@latest check test.cmdt &> /dev/null; then
        echo -e "${GREEN}PASS${NC}"
        ((PASSED++))
    else
        echo -e "${RED}FAIL${NC}"
        ((FAILED++))
        # Show the failure details
        cd "$TEMP_DIR" && PATH="$PROJECT_ROOT:$PATH" go run github.com/deref/transcript@latest check test.cmdt 2>&1 | sed 's/^/    /'
    fi
    
    # Clean up temp directory
    rm -rf "$TEMP_DIR"
done < <(find "$SCRIPT_DIR" -name "test.cmdt" -type f)

echo
echo "Summary: $PASSED passed, $FAILED failed"

if [ $FAILED -gt 0 ]; then
    exit 1
fi