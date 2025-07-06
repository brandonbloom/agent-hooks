package format

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result struct {
	FormattedFiles []string
	Warnings       []string
	Errors         []string
}

func FormatFiles(files []string) *Result {
	result := &Result{}

	goFiles := filterFilesByExtension(files, ".go")
	if len(goFiles) > 0 {
		if err := formatGoFiles(goFiles, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Go formatting failed: %v", err))
		}
	}

	unsupportedFiles := filterUnsupportedFiles(files)
	for _, file := range unsupportedFiles {
		result.Warnings = append(result.Warnings, fmt.Sprintf("No formatter available for: %s", file))
	}

	return result
}

func formatGoFiles(files []string, result *Result) error {
	if !isCommandAvailable("go") {
		return fmt.Errorf("go command not found - please install Go")
	}

	for _, file := range files {
		cmd := exec.Command("go", "fmt", file)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to format %s: %w", file, err)
		}
		result.FormattedFiles = append(result.FormattedFiles, file)
	}

	return nil
}

func filterFilesByExtension(files []string, ext string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(file, ext) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func filterUnsupportedFiles(files []string) []string {
	supportedExtensions := map[string]bool{
		".go": true,
	}

	var unsupported []string
	for _, file := range files {
		ext := filepath.Ext(file)
		if ext != "" && !supportedExtensions[ext] {
			if info, err := os.Stat(file); err == nil && !info.IsDir() {
				unsupported = append(unsupported, file)
			}
		}
	}
	return unsupported
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
