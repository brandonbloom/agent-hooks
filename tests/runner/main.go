package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TestConfig defines a test scenario
type TestConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Setup       []string          `json:"setup"`       // Commands to run for setup (e.g., "git init")
	Commands    []TestCommand     `json:"commands"`    // Commands to test
	Files       map[string]string `json:"files"`       // Files to create: filename -> content
}

// TestCommand defines a command to run and its expected result
type TestCommand struct {
	Command      string `json:"command"`      // Command to run (e.g., "which-vcs")
	ExpectedExit int    `json:"expectedExit"` // Expected exit code (0 = success)
	ExpectedOut  string `json:"expectedOut"`  // Expected stdout (exact match)
	ExpectedErr  string `json:"expectedErr"`  // Expected stderr (exact match)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <test-directory>\n", os.Args[0])
		os.Exit(1)
	}

	testDir := os.Args[1]
	if err := runTest(testDir); err != nil {
		fmt.Fprintf(os.Stderr, "Test failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")
}

func runTest(testDir string) error {
	// Read test configuration
	configPath := filepath.Join(testDir, "test.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read test config: %w", err)
	}

	var config TestConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse test config: %w", err)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "agent-hooks-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy test directory to temp location
	if err := copyDir(testDir, tempDir); err != nil {
		return fmt.Errorf("failed to copy test dir: %w", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current dir: %w", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		return fmt.Errorf("failed to change to temp dir: %w", err)
	}

	// Create test files
	for filename, content := range config.Files {
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
	}

	// Run setup commands
	for _, setupCmd := range config.Setup {
		if err := runCommand(setupCmd); err != nil {
			return fmt.Errorf("setup command failed '%s': %w", setupCmd, err)
		}
	}

	// Build agent-hooks binary
	agentHooksPath := filepath.Join(oldDir, "main.go")
	buildCmd := exec.Command("go", "run", agentHooksPath)
	buildCmd.Dir = tempDir

	// Run test commands
	for i, testCmd := range config.Commands {
		fmt.Printf("Running test command %d: %s\n", i+1, testCmd.Command)
		
		// Prepare command
		parts := strings.Fields(testCmd.Command)
		cmd := exec.Command("go", append([]string{"run", agentHooksPath}, parts...)...)
		cmd.Dir = tempDir

		// Capture output
		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Run command
		err := cmd.Run()
		actualExit := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				actualExit = exitError.ExitCode()
			} else {
				return fmt.Errorf("command failed to execute: %w", err)
			}
		}

		// Check results
		actualOut := strings.TrimSpace(stdout.String())
		actualErr := strings.TrimSpace(stderr.String())

		if actualExit != testCmd.ExpectedExit {
			return fmt.Errorf("command %d: expected exit code %d, got %d", i+1, testCmd.ExpectedExit, actualExit)
		}

		if testCmd.ExpectedOut != "" && actualOut != testCmd.ExpectedOut {
			return fmt.Errorf("command %d: expected stdout '%s', got '%s'", i+1, testCmd.ExpectedOut, actualOut)
		}

		if testCmd.ExpectedErr != "" && actualErr != testCmd.ExpectedErr {
			return fmt.Errorf("command %d: expected stderr '%s', got '%s'", i+1, testCmd.ExpectedErr, actualErr)
		}

		fmt.Printf("  ✓ Exit code: %d\n", actualExit)
		if actualOut != "" {
			fmt.Printf("  ✓ Stdout: %s\n", actualOut)
		}
		if actualErr != "" {
			fmt.Printf("  ✓ Stderr: %s\n", actualErr)
		}
	}

	return nil
}

func runCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the temp directory itself if it's already there
		if path == src {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}