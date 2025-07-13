# Agent-Hooks Test Suite Makefile

.PHONY: all test build clean help

# Default target
all: test

# Build the test runner
build:
	go build -o test-runner runner/main.go

# Run all tests using the shell script
test: build
	./run-tests.sh

# Run individual tests manually
test-vcs-git: build
	./test-runner vcs-detection-git

test-vcs-none: build
	./test-runner vcs-detection-none

# Clean up generated files
clean:
	rm -f test-runner

# Show help
help:
	@echo "Available targets:"
	@echo "  all         - Run all tests (default)"
	@echo "  test        - Run all tests"
	@echo "  build       - Build the test runner"
	@echo "  test-vcs-git   - Run VCS detection git test"
	@echo "  test-vcs-none  - Run VCS detection none test"
	@echo "  clean       - Clean up generated files"
	@echo "  help        - Show this help"