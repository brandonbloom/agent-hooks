package format

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brandonbloom/agent-hooks/internal/detect"
	"github.com/brandonbloom/agent-hooks/internal/doctor"
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

	// Detect available formatting technologies once
	detectedTechs, err := detect.DetectInCurrentDirectory()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to detect technologies: %v", err))
		return result
	}

	// Get formatting support configuration
	supportConfigs := doctor.GetFormattingSupport()

	// Group files by their formatting support
	for _, config := range supportConfigs {
		matchingFiles := filterFilesByExtensions(files, config.Extensions)
		if len(matchingFiles) > 0 {
			if err := formatFilesBySupport(matchingFiles, config, detectedTechs, result, opts); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Formatting failed for %v: %v", config.Extensions, err))
			}
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

func formatFilesBySupport(files []string, support doctor.FormattingToolSupport, detectedTechs []detect.Technology, result *Result, opts Options) error {
	// Find the first available tool from the preference-ordered list
	for _, toolName := range support.Tools {
		if canUseFormatter(toolName, detectedTechs) {
			return formatWithTool(toolName, files, result, opts)
		}
	}

	// No formatter available
	return fmt.Errorf("no formatter available for extensions %v - available tools: %v", support.Extensions, support.Tools)
}

func canUseFormatter(toolName string, detectedTechs []detect.Technology) bool {
	switch toolName {
	case "goimports", "gofmt":
		// Go tools are always available if the command exists
		return isCommandAvailable(toolName)
	case "biome":
		// JS/TS tools require configuration detection
		return containsTechnology(detectedTechs, detect.Biome) && isCommandAvailable(toolName)
	case "prettier":
		return containsTechnology(detectedTechs, detect.Prettier) && isCommandAvailable("npx")
	default:
		return false
	}
}

func containsTechnology(techs []detect.Technology, target detect.Technology) bool {
	for _, tech := range techs {
		if tech == target {
			return true
		}
	}
	return false
}

func formatWithTool(toolName string, files []string, result *Result, opts Options) error {
	switch toolName {
	case "goimports":
		return formatWithGoimports(files, result, opts)
	case "gofmt":
		return formatWithGofmt(files, result, opts)
	case "biome":
		return formatWithBiome(files, result, opts)
	case "prettier":
		return formatWithPrettier(files, result, opts)
	default:
		return fmt.Errorf("unsupported formatter: %s", toolName)
	}
}

// formatWithCommandCheck checks command availability and executes formatting with common error handling
func formatWithCommandCheck(command string, errorMessage string, toolName string, cmdArgs []string, files []string, result *Result, opts Options) error {
	if !isCommandAvailable(command) {
		return fmt.Errorf(errorMessage)
	}
	
	return runFormatterCommand(toolName, cmdArgs, files, result, opts)
}

// runFormatterCommand executes a formatter command on a set of files with common error handling
func runFormatterCommand(toolName string, cmdArgs []string, files []string, result *Result, opts Options) error {
	for _, file := range files {
		if opts.DryRun {
			result.FormattedFiles = append(result.FormattedFiles, file)
			continue
		}

		// Build command with file appended
		fullArgs := append(cmdArgs, file)
		cmd := exec.Command(fullArgs[0], fullArgs[1:]...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to format %s with %s: %w\nOutput: %s", file, toolName, err, string(output))
		}
		result.FormattedFiles = append(result.FormattedFiles, file)
	}

	return nil
}

func formatWithGoimports(files []string, result *Result, opts Options) error {
	return formatWithCommandCheck("goimports", "goimports command not found - install with: go install golang.org/x/tools/cmd/goimports@latest", "goimports", []string{"goimports", "-w"}, files, result, opts)
}

func formatWithGofmt(files []string, result *Result, opts Options) error {
	return formatWithCommandCheck("gofmt", "gofmt command not found", "gofmt", []string{"gofmt", "-w"}, files, result, opts)
}

func formatWithBiome(files []string, result *Result, opts Options) error {
	return formatWithCommandCheck("biome", "biome command not found - install with: npm install -g @biomejs/biome", "biome", []string{"biome", "format", "--write"}, files, result, opts)
}

func formatWithPrettier(files []string, result *Result, opts Options) error {
	return formatWithCommandCheck("npx", "npx command not found - install Node.js to get npx", "prettier", []string{"npx", "prettier", "--write"}, files, result, opts)
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

func filterFilesByExtensions(files []string, extensions []string) []string {
	var filtered []string
	for _, file := range files {
		for _, ext := range extensions {
			if strings.HasSuffix(file, ext) {
				filtered = append(filtered, file)
				break
			}
		}
	}
	return filtered
}

func filterUnsupportedFiles(files []string) []string {
	// Build supported extensions map from formatting support
	supportedExtensions := make(map[string]bool)
	for _, config := range doctor.GetFormattingSupport() {
		for _, ext := range config.Extensions {
			supportedExtensions[ext] = true
		}
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
