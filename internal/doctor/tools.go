package doctor

import (
	"fmt"
	"os/exec"
	"strings"
)

type CheckResult struct {
	Name    string
	Status  CheckStatus
	Message string
}

type CheckStatus int

const (
	CheckPassed CheckStatus = iota
	CheckWarning
	CheckFailed
)

type ToolCheck struct {
	Name      string
	Command   string
	Validator func() error
	Required  bool
}

var DefaultTools = []ToolCheck{
	{Name: "git", Command: "git", Required: true},
	{Name: "go", Command: "go", Required: false},
	{Name: "goimports", Command: "goimports", Required: false},
	{Name: "agent-hooks", Command: "agent-hooks", Required: false},
}

func RunToolChecks(verbose bool) []CheckResult {
	var results []CheckResult

	for _, tool := range DefaultTools {
		result := checkTool(tool, verbose)

		if !verbose && result.Status == CheckPassed {
			continue
		}

		results = append(results, result)
	}

	return results
}

func checkTool(tool ToolCheck, verbose bool) CheckResult {
	result := CheckResult{Name: tool.Name}

	if !isCommandAvailable(tool.Command) {
		result.Status = CheckFailed
		if tool.Required {
			result.Message = fmt.Sprintf("%s command not found", tool.Command)
		} else {
			result.Status = CheckWarning
			result.Message = fmt.Sprintf("%s command not found (optional)", tool.Command)
		}
		return result
	}

	if verbose {
		version := getToolVersion(tool.Command)
		if version != "" {
			result.Message = fmt.Sprintf("%s is installed (%s)", tool.Command, version)
		} else {
			result.Message = fmt.Sprintf("%s is installed", tool.Command)
		}
	}

	if tool.Validator != nil {
		if err := tool.Validator(); err != nil {
			result.Status = CheckWarning
			result.Message = fmt.Sprintf("%s configuration issue: %v", tool.Name, err)
			return result
		}
	}

	result.Status = CheckPassed
	return result
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func getToolVersion(command string) string {
	versionArgs := map[string][]string{
		"git":         {"--version"},
		"go":          {"version"},
		"goimports":   {"--help"},
		"agent-hooks": {"--version"},
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

	version := strings.TrimSpace(string(output))

	switch command {
	case "git":
		if strings.HasPrefix(version, "git version ") {
			return strings.TrimPrefix(version, "git version ")
		}
	case "go":
		if strings.HasPrefix(version, "go version ") {
			parts := strings.Fields(version)
			if len(parts) >= 3 {
				return parts[2]
			}
		}
	}

	return version
}
