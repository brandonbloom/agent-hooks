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
	SkippedFiles   []string
}

type Options struct {
	DryRun  bool
	Verbose bool
}

func FormatFiles(files []string) *Result {
	return FormatFilesWithOptions(files, Options{})
}

func FormatFilesWithOptions(files []string, opts Options) *Result {
	result := &Result{}

	goFiles := filterFilesByExtension(files, ".go")
	if len(goFiles) > 0 {
		if err := formatGoFiles(goFiles, result, opts); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Go formatting failed: %v", err))
		}
	}

	unsupportedFiles := filterUnsupportedFiles(files)
	for _, file := range unsupportedFiles {
		if opts.Verbose {
			result.Warnings = append(result.Warnings, fmt.Sprintf("No formatter available for: %s", file))
		} else {
			result.SkippedFiles = append(result.SkippedFiles, file)
		}
	}

	return result
}

func formatGoFiles(files []string, result *Result, opts Options) error {
	if !isCommandAvailable("goimports") {
		return fmt.Errorf("goimports command not found - please install with: go install golang.org/x/tools/cmd/goimports@latest")
	}

	for _, file := range files {
		if opts.DryRun {
			result.FormattedFiles = append(result.FormattedFiles, file)
			continue
		}

		cmd := exec.Command("goimports", "-w", file)
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
