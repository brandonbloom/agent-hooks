package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/brandonbloom/agent-hooks/internal/git"
	"github.com/brandonbloom/agent-hooks/internal/vcs"
)

type DetectionRule struct {
	Technology Technology
	Files      []string
	Desc       string
	URL        string
}

type Detector struct {
	VCSType      vcs.VCS
	TrackedFiles []string
	fileIndex    map[string]bool // basename -> exists
	Verbose      bool
}

func (d *Detector) Detect(dir string) ([]Technology, error) {
	var detected []Technology

	// Phase 1: Do VCS detection and file listing once (if not already set)
	start := time.Now()
	if d.VCSType == "" {
		d.VCSType, _ = vcs.DetectVCS()
	}
	vcsTime := time.Since(start)

	start = time.Now()
	if d.VCSType == vcs.Git && d.TrackedFiles == nil {
		var err error
		d.TrackedFiles, err = git.GetAllTrackedFiles()
		if err != nil {
			return nil, err
		}

		// Build file index for fast lookups
		d.fileIndex = make(map[string]bool)
		for _, file := range d.TrackedFiles {
			d.fileIndex[filepath.Base(file)] = true
			d.fileIndex[file] = true // Also index full path
		}
	}
	gitTime := time.Since(start)

	// Phase 2: Check each rule with pre-computed information
	start = time.Now()
	for _, rule := range detectionRules {
		found, err := d.CheckRule(dir, rule)
		if err != nil {
			continue
		}
		if found {
			detected = append(detected, rule.Technology)
		}
	}
	rulesTime := time.Since(start)

	if d.Verbose {
		fmt.Printf("VCS: %v, Git: %v, Rules: %v\n", vcsTime, gitTime, rulesTime)
	}

	return detected, nil
}

func (d *Detector) GetRules() []DetectionRule {
	return detectionRules
}

func (d *Detector) CheckRule(dir string, rule DetectionRule) (bool, error) {
	// Special case: Git detection should use VCS walking logic
	if rule.Technology == Git {
		return d.VCSType == vcs.Git, nil
	}

	// For other technologies, try VCS-aware detection first
	if d.VCSType == vcs.Git {
		if found, err := d.checkRuleWithTrackedFiles(rule, d.TrackedFiles); err == nil {
			return found, nil
		}
		// If VCS detection fails, fall back to directory-only approach
	}

	// Fallback to current directory-only approach
	return d.checkRuleDirectoryOnly(dir, rule)
}

func (d *Detector) checkRuleDirectoryOnly(dir string, rule DetectionRule) (bool, error) {
	for _, file := range rule.Files {
		if containsWildcard(file) {
			matches, err := filepath.Glob(filepath.Join(dir, file))
			if err != nil {
				continue
			}
			if len(matches) > 0 {
				return true, nil
			}
		} else {
			path := filepath.Join(dir, file)
			if _, err := os.Stat(path); err == nil {
				return true, nil
			}
		}
	}
	return false, nil
}

func (d *Detector) checkRuleWithTrackedFiles(rule DetectionRule, trackedFiles []string) (bool, error) {
	for _, file := range rule.Files {
		if containsWildcard(file) {
			// For wildcards, check patterns against index keys
			for filename := range d.fileIndex {
				if matched, _ := filepath.Match(file, filename); matched {
					return true, nil
				}
			}
		} else {
			// Fast exact match lookup
			if d.fileIndex[file] {
				return true, nil
			}
		}
	}
	return false, nil
}

func containsWildcard(path string) bool {
	base := filepath.Base(path)
	return len(base) > 0 && (base[0] == '*' || base[len(base)-1] == '*')
}

func DetectInDirectory(dir string) ([]Technology, error) {
	detector := &Detector{}
	return detector.Detect(dir)
}

func DetectInCurrentDirectory() ([]Technology, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	return DetectInDirectory(cwd)
}
