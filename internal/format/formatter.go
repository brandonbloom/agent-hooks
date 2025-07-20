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
	// Get project-aware tool preference
	preferredTools := getProjectAwareToolPreference(support.Tools, detectedTechs)

	// Find the first available tool from the project-aware preference list
	for _, toolName := range preferredTools {
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
		// Biome can format JS/TS files if the command is available globally
		return isCommandAvailable(toolName)
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

// getProjectAwareToolPreference reorders tools based on project configuration.
// If exactly one tool has project config, it gets priority.
// Otherwise, use the original preference order.
func getProjectAwareToolPreference(tools []string, detectedTechs []detect.Technology) []string {
	var configuredTools []string
	var unconfiguredTools []string

	for _, tool := range tools {
		hasConfig := false
		switch tool {
		case "prettier":
			hasConfig = containsTechnology(detectedTechs, detect.Prettier)
		case "biome":
			hasConfig = containsTechnology(detectedTechs, detect.Biome)
		default:
			// Unknown tool, treat as unconfigured
			hasConfig = false
		}

		if hasConfig {
			configuredTools = append(configuredTools, tool)
		} else {
			unconfiguredTools = append(unconfiguredTools, tool)
		}
	}

	// If exactly one tool is configured, prefer it
	if len(configuredTools) == 1 {
		result := make([]string, 0, len(tools))
		result = append(result, configuredTools...)
		result = append(result, unconfiguredTools...)
		return result
	}

	// If zero or multiple tools configured, use original order
	return tools
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

// formatterCommand encapsulates the parameters needed for formatting with availability checking
type formatterCommand struct {
	command      string   // command to check availability for
	errorMessage string   // error message if command not available
	toolName     string   // name of the tool for error messages
	cmdArgs      []string // the command arguments
	files        []string // files to format
	result       *Result  // result structure to populate
	opts         Options  // options for formatting
}

// Run executes the formatter command with availability checking and error handling
func (fc *formatterCommand) Run() error {
	// Check availability
	if !isCommandAvailable(fc.command) {
		return fmt.Errorf(fc.errorMessage)
	}

	// Execute formatting
	for _, file := range fc.files {
		if fc.opts.DryRun {
			fc.result.FormattedFiles = append(fc.result.FormattedFiles, file)
			continue
		}

		// Build command with file appended
		fullArgs := append(fc.cmdArgs, file)
		cmd := exec.Command(fullArgs[0], fullArgs[1:]...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to format %s with %s: %w\nOutput: %s", file, fc.toolName, err, string(output))
		}
		fc.result.FormattedFiles = append(fc.result.FormattedFiles, file)
	}

	return nil
}

func formatWithGoimports(files []string, result *Result, opts Options) error {
	return (&formatterCommand{
		command:      "goimports",
		errorMessage: "goimports command not found - install with: go install golang.org/x/tools/cmd/goimports@latest",
		toolName:     "goimports",
		cmdArgs:      []string{"goimports", "-w"},
		files:        files,
		result:       result,
		opts:         opts,
	}).Run()
}

func formatWithGofmt(files []string, result *Result, opts Options) error {
	return (&formatterCommand{
		command:      "gofmt",
		errorMessage: "gofmt command not found",
		toolName:     "gofmt",
		cmdArgs:      []string{"gofmt", "-w"},
		files:        files,
		result:       result,
		opts:         opts,
	}).Run()
}

func formatWithBiome(files []string, result *Result, opts Options) error {
	return (&formatterCommand{
		command:      "biome",
		errorMessage: "biome command not found - install with: npm install -g @biomejs/biome",
		toolName:     "biome",
		cmdArgs:      []string{"biome", "format", "--write"},
		files:        files,
		result:       result,
		opts:         opts,
	}).Run()
}

func formatWithPrettier(files []string, result *Result, opts Options) error {
	return (&formatterCommand{
		command:      "npx",
		errorMessage: "npx command not found - install Node.js to get npx",
		toolName:     "prettier",
		cmdArgs:      []string{"npx", "prettier", "--write"},
		files:        files,
		result:       result,
		opts:         opts,
	}).Run()
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
