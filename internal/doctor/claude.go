package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ClaudeSettings struct {
	Hooks ClaudeHooks `json:"hooks"`
}

type ClaudeHooks struct {
	PostToolUse []ClaudeHook `json:"PostToolUse"`
	PreToolUse  []ClaudeHook `json:"PreToolUse"`
}

type ClaudeHook struct {
	Matcher string       `json:"matcher"`
	Hooks   []HookConfig `json:"hooks"`
}

type HookConfig struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

func RunClaudeChecks(verbose bool) []CheckResult {
	var results []CheckResult
	
	settingsPath := getClaudeSettingsPath()
	
	settingsResult := checkClaudeSettingsFile(settingsPath, verbose)
	if !verbose && settingsResult.Status == CheckPassed {
		// Only continue to hook checks if settings file exists
	} else {
		results = append(results, settingsResult)
	}
	
	if settingsResult.Status != CheckFailed {
		hookResult := checkClaudeHookConfiguration(settingsPath, verbose)
		if !verbose && hookResult.Status == CheckPassed {
			// Don't add passed checks in non-verbose mode
		} else {
			results = append(results, hookResult)
		}
	}
	
	return results
}

func checkClaudeSettingsFile(settingsPath string, verbose bool) CheckResult {
	result := CheckResult{Name: "Claude settings file"}
	
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		result.Status = CheckFailed
		result.Message = fmt.Sprintf("Claude settings file missing at %s", settingsPath)
		return result
	}
	
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		result.Status = CheckFailed
		result.Message = fmt.Sprintf("Cannot read Claude settings file: %v", err)
		return result
	}
	
	var settings ClaudeSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		result.Status = CheckFailed
		result.Message = fmt.Sprintf("Invalid JSON in Claude settings file: %v", err)
		return result
	}
	
	result.Status = CheckPassed
	if verbose {
		result.Message = fmt.Sprintf("Claude settings file found at %s", settingsPath)
	}
	
	return result
}

func checkClaudeHookConfiguration(settingsPath string, verbose bool) CheckResult {
	result := CheckResult{Name: "Claude hook configuration"}
	
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		result.Status = CheckFailed
		result.Message = fmt.Sprintf("Cannot read settings file: %v", err)
		return result
	}
	
	var settings ClaudeSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		result.Status = CheckFailed
		result.Message = fmt.Sprintf("Cannot parse settings file: %v", err)
		return result
	}
	
	agentHooksFound := false
	correctMatcher := false
	
	for _, hook := range settings.Hooks.PostToolUse {
		for _, config := range hook.Hooks {
			if strings.Contains(config.Command, "agent-hooks format") {
				agentHooksFound = true
				if strings.Contains(hook.Matcher, "Write") && 
				   strings.Contains(hook.Matcher, "Edit") && 
				   strings.Contains(hook.Matcher, "MultiEdit") {
					correctMatcher = true
				}
				break
			}
		}
	}
	
	if !agentHooksFound {
		result.Status = CheckWarning
		result.Message = "agent-hooks not configured in Claude hooks"
		return result
	}
	
	if !correctMatcher {
		result.Status = CheckWarning
		result.Message = "agent-hooks hook matcher should include 'Write|Edit|MultiEdit'"
		return result
	}
	
	result.Status = CheckPassed
	if verbose {
		result.Message = "agent-hooks properly configured in Claude PostToolUse hooks"
	}
	
	return result
}

func getClaudeSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".claude", "settings.json")
}