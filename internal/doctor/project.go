package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brandonbloom/agent-hooks/internal/detect"
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
				Message: "No specific technologies detected (generic project)",
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
			Message: fmt.Sprintf("Detected technologies: %v", techNames),
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

func checkProjectTool(req ToolRequirement, verbose bool) CheckResult {
	result := CheckResult{Name: fmt.Sprintf("%s (%s)", req.Tool, req.Technology)}

	if !isCommandAvailable(req.Tool) {
		if req.Required {
			result.Status = CheckFailed
			result.Message = fmt.Sprintf("%s command not found (required for %s)", req.Tool, req.Technology)
		} else {
			result.Status = CheckWarning
			result.Message = fmt.Sprintf("%s command not found (optional for %s)", req.Tool, req.Technology)
		}
		return result
	}

	if verbose {
		version := getProjectToolVersion(req.Tool)
		if version != "" {
			result.Message = fmt.Sprintf("%s is installed (%s) - %s", req.Tool, version, req.Desc)
		} else {
			result.Message = fmt.Sprintf("%s is installed - %s", req.Tool, req.Desc)
		}
	}

	result.Status = CheckPassed
	return result
}

func getProjectToolVersion(command string) string {
	versionArgs := map[string][]string{
		"go":         {"version"},
		"node":       {"--version"},
		"npm":        {"--version"},
		"python":     {"--version"},
		"pip":        {"--version"},
		"cargo":      {"--version"},
		"rustc":      {"--version"},
		"ruby":       {"--version"},
		"gem":        {"--version"},
		"java":       {"-version"},
		"javac":      {"-version"},
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
