package doctor

import (
	"fmt"
	"os"
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
	URL       string
}

// Default tools are sorted alphabetically to minimize merge conflicts
// when adding new tools. Please maintain this order.
var DefaultTools = []ToolCheck{
	{Name: "agent-hooks", Command: "agent-hooks", Required: false, URL: "https://github.com/brandonbloom/agent-hooks"},
	{Name: "direnv", Command: "direnv", Required: false, Validator: validateDirenvSetup, URL: "https://direnv.net"},
	{Name: "git", Command: "git", Required: true, URL: "https://git-scm.com"},
	{Name: "go", Command: "go", Required: false, URL: "https://golang.org"},
	{Name: "goimports", Command: "goimports", Required: false, URL: "https://pkg.go.dev/golang.org/x/tools/cmd/goimports"},
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
	// Version arguments are sorted alphabetically to minimize merge conflicts
	// when adding new tools. Please maintain this order.
	versionArgs := map[string][]string{
		"agent-hooks": {"--version"},
		"direnv":      {"version"},
		"git":         {"--version"},
		"go":          {"version"},
		"goimports":   {"--help"},
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

func validateDirenvSetup() error {
	// Only validate direnv setup if we're in a project that uses direnv
	if !hasDirenvFile() {
		return nil // No .envrc file, so no additional validation needed
	}

	// Check if direnv status command works (indicates shell integration)
	cmd := exec.Command("direnv", "status")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("direnv status failed - shell integration may not be setup")
	}

	status := string(output)

	// Check if current directory's .envrc is allowed
	if strings.Contains(status, "Found RC allowed false") {
		return fmt.Errorf(".envrc file found but not allowed - run 'direnv allow' to trust it")
	}

	return nil
}

func hasDirenvFile() bool {
	_, err := os.Stat(".envrc")
	return err == nil
}
