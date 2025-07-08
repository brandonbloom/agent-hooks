package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brandonbloom/agent-hooks/internal/detect"
	"github.com/brandonbloom/agent-hooks/internal/vcs"
)

func RunProjectChecks(verbose bool) []CheckResult {
	var results []CheckResult

	cwd, err := os.Getwd()
	if err != nil {
		results = append(results, CheckResult{
			Name:    "Project Detection",
			Status:  CheckWarning,
			Message: fmt.Sprintf("Failed to get current directory: %v", err),
		})
		return results
	}

	// Get project root info once
	var projectInfo string
	projectRoot, err := vcs.FindProjectRoot()
	if err != nil {
		projectInfo = fmt.Sprintf("current directory: %s", cwd)
	} else {
		projectInfo = fmt.Sprintf("project root: %s", projectRoot)
	}

	detector := &detect.Detector{}
	technologies, err := detector.Detect(cwd)
	if err != nil {
		results = append(results, CheckResult{
			Name:    "Project Detection",
			Status:  CheckWarning,
			Message: fmt.Sprintf("Failed to detect technologies: %v", err),
		})
		return results
	}

	if len(technologies) == 0 {
		if verbose {
			results = append(results, CheckResult{
				Name:    "Project Detection",
				Status:  CheckPassed,
				Message: fmt.Sprintf("No specific technologies detected (generic project) - %s", projectInfo),
			})
		}
		return results
	}

	if verbose {
		techNames := make([]string, len(technologies))
		for i, tech := range technologies {
			techNames[i] = string(tech)
		}
		results = append(results, CheckResult{
			Name:    "Project Detection",
			Status:  CheckPassed,
			Message: fmt.Sprintf("Detected technologies: %v - %s", techNames, projectInfo),
		})
	}

	for _, tech := range technologies {
		requirements := GetToolRequirements(tech)
		for _, req := range requirements {
			result := checkProjectTool(req, verbose)
			results = append(results, result)
		}
	}

	return results
}

// checkProjectTool validates a tool for a specific project technology.
// It looks up the tool definition from AllTools and delegates to the unified
// tool checking system, using the requirement's required flag rather than
// the tool's default required setting.
func checkProjectTool(req ToolRequirement, verbose bool) CheckResult {
	// Look up the tool definition
	tool, exists := GetToolByName(req.Tool)
	if !exists {
		return CheckResult{
			Name:    fmt.Sprintf("%s (%s)", req.Tool, req.Technology),
			Status:  CheckFailed,
			Message: fmt.Sprintf("Tool definition not found: %s", req.Tool),
		}
	}

	// Use the unified tool checking with the requirement's required flag
	result := checkTool(tool, req.Required, verbose)

	// Update the result name to include technology context
	result.Name = fmt.Sprintf("%s (%s)", req.Tool, req.Technology)

	return result
}

func getProjectToolVersion(command string) string {
	// Version arguments are sorted alphabetically to minimize merge conflicts
	// when adding new tools. Please maintain this order.
	versionArgs := map[string][]string{
		"cargo":      {"--version"},
		"gem":        {"--version"},
		"go":         {"version"},
		"java":       {"-version"},
		"javac":      {"-version"},
		"node":       {"--version"},
		"npm":        {"--version"},
		"pip":        {"--version"},
		"python":     {"--version"},
		"ruby":       {"--version"},
		"rustc":      {"--version"},
		"transcript": {"--version"},
	}

	args, exists := versionArgs[command]
	if !exists {
		return ""
	}

	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
