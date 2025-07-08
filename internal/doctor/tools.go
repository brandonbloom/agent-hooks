package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brandonbloom/agent-hooks/internal/detect"
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
	URL       string
}

// AllTools contains all known tools, sorted alphabetically to minimize merge conflicts
// when adding new tools. Please maintain this order.
//
// This is the master registry of all tools that can be checked by the doctor command.
// Tools can be:
//   - Regular tools: Have a Command field and are checked for availability in PATH
//   - Meta-tools: Have no Command field (empty string) and use only a Validator function
//     to check for "one of" requirements (e.g., procfile-runner checks for foreman OR hivemind OR overmind)
var AllTools = []ToolCheck{
	{Name: "agent-hooks", Command: "agent-hooks", URL: "https://github.com/brandonbloom/agent-hooks"},
	{Name: "cargo", Command: "cargo", URL: "https://doc.rust-lang.org/cargo/"},
	{Name: "clojure", Command: "clojure", URL: "https://clojure.org"},
	{Name: "direnv", Command: "direnv", Validator: validateDirenvSetup, URL: "https://direnv.net"},
	{Name: "foreman", Command: "foreman", URL: "https://github.com/ddollar/foreman"},
	{Name: "gem", Command: "gem", URL: "https://rubygems.org"},
	{Name: "git", Command: "git", URL: "https://git-scm.com"},
	{Name: "go", Command: "go", URL: "https://golang.org"},
	{Name: "gofmt", Command: "gofmt", URL: "https://golang.org"},
	{Name: "goimports", Command: "goimports", URL: "https://pkg.go.dev/golang.org/x/tools/cmd/goimports"},
	{Name: "hivemind", Command: "hivemind", URL: "https://github.com/DarthSim/hivemind"},
	{Name: "hurl", Command: "hurl", URL: "https://hurl.dev"},
	{Name: "java", Command: "java", URL: "https://www.oracle.com/java/"},
	{Name: "javac", Command: "javac", URL: "https://www.oracle.com/java/"},
	{Name: "lein", Command: "lein", URL: "https://leiningen.org"},
	{Name: "ng", Command: "ng", URL: "https://angular.io/cli"},
	{Name: "node", Command: "node", URL: "https://nodejs.org"},
	{Name: "npm", Command: "npm", URL: "https://www.npmjs.com"},
	{Name: "overmind", Command: "overmind", URL: "https://github.com/DarthSim/overmind"},
	{Name: "pip", Command: "pip", URL: "https://pip.pypa.io"},
	{Name: "procfile-runner", Validator: validateProcfileRunner, URL: "https://devcenter.heroku.com/articles/procfile"},
	{Name: "python", Command: "python", URL: "https://www.python.org"},
	{Name: "ruby", Command: "ruby", URL: "https://www.ruby-lang.org"},
	{Name: "rustc", Command: "rustc", URL: "https://www.rust-lang.org"},
	{Name: "transcript", Command: "transcript", URL: "https://github.com/jspahrsummers/transcript"},
}

// DefaultTools are the core tools checked for ALL projects.
// These tools are checked globally regardless of what technologies are detected in the project.
// All default tools are required by definition.
var DefaultTools = []string{
	"agent-hooks",
	"git",
	"go",
}

// GetToolByName returns the ToolCheck for the given tool name from the AllTools registry.
// This is used to look up tool definitions when processing requirements or handling the about command.
func GetToolByName(name string) (ToolCheck, bool) {
	for _, tool := range AllTools {
		if tool.Name == name {
			return tool, true
		}
	}
	return ToolCheck{}, false
}

func RunToolChecks(verbose bool) []CheckResult {
	var results []CheckResult

	for _, toolName := range DefaultTools {
		tool, exists := GetToolByName(toolName)
		if !exists {
			results = append(results, CheckResult{
				Name:    toolName,
				Status:  CheckFailed,
				Message: fmt.Sprintf("Tool definition not found: %s", toolName),
			})
			continue
		}

		required := true // All default tools are required
		result := checkTool(tool, required, verbose)

		if !verbose && result.Status == CheckPassed {
			continue
		}

		results = append(results, result)
	}

	return results
}

// checkTool validates a tool with the specified required flag.
// This function implements the unified tool checking logic where the required status
// is determined by the context (global default tools vs. project-specific requirements)
// rather than being an intrinsic property of the tool itself.
func checkTool(tool ToolCheck, required bool, verbose bool) CheckResult {
	result := CheckResult{Name: tool.Name}

	commandAvailable := tool.Command != "" && isCommandAvailable(tool.Command)

	if !commandAvailable && tool.Command != "" {
		if required {
			result.Status = CheckFailed
			result.Message = fmt.Sprintf("%s command not found", tool.Command)
		} else {
			result.Status = CheckWarning
			result.Message = fmt.Sprintf("%s command not found (optional)", tool.Command)
		}
		return result
	}

	if commandAvailable && verbose {
		version := getToolVersion(tool.Command)
		if version != "" {
			result.Message = fmt.Sprintf("%s is installed (%s)", tool.Command, version)
		} else {
			result.Message = fmt.Sprintf("%s is installed", tool.Command)
		}
	}

	if tool.Validator != nil {
		if err := tool.Validator(); err != nil {
			if required {
				result.Status = CheckFailed
				result.Message = fmt.Sprintf("%s: %v", tool.Name, err)
			} else {
				result.Status = CheckWarning
				result.Message = fmt.Sprintf("%s: %v", tool.Name, err)
			}
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

func validateProcfileRunner() error {
	runners := []string{"foreman", "hivemind", "overmind"}
	result := CheckForOneToolOf(detect.Procfile, runners, "procfile-runner", false, false)
	if result.Status == CheckPassed {
		return nil
	}
	return fmt.Errorf("no procfile runner found, install one of: %v", runners)
}

// CheckForOneToolOf checks if at least one tool from the given list is available
// for the specified technology. Returns a CheckResult indicating success or failure.
//
// This function implements "one of" requirements where a project needs ANY of several
// alternative tools to function (e.g., npm OR yarn OR pnpm for Node.js projects).
// It's used by meta-tool validators like validateProcfileRunner.
func CheckForOneToolOf(tech detect.Technology, tools []string, groupName string, required bool, verbose bool) CheckResult {
	result := CheckResult{Name: fmt.Sprintf("%s (%s)", groupName, tech)}

	var availableTools []string
	for _, tool := range tools {
		if isCommandAvailable(tool) {
			availableTools = append(availableTools, tool)
		}
	}

	if len(availableTools) == 0 {
		if required {
			result.Status = CheckFailed
			result.Message = fmt.Sprintf("No %s tools found (required for %s). Available options: %v", groupName, tech, tools)
		} else {
			result.Status = CheckWarning
			result.Message = fmt.Sprintf("No %s tools found (optional for %s). Available options: %v", groupName, tech, tools)
		}
		return result
	}

	result.Status = CheckPassed
	if verbose {
		result.Message = fmt.Sprintf("%s satisfied by: %v", groupName, availableTools)
	}

	return result
}
