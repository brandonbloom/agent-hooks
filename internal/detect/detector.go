package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

type DetectionEvidence struct {
	Technology    Technology
	Found         bool
	MatchedFiles  []string
	PatternCounts map[string]int
	Method        string
}

type Detector struct {
	VCSType      vcs.VCS
	TrackedFiles []string
	fileIndex    map[string]bool // basename -> exists
	Verbose      bool
}

func (d *Detector) Detect(dir string) ([]Technology, error) {
	evidence, err := d.DetectWithEvidence(dir)
	if err != nil {
		return nil, err
	}

	var detected []Technology
	for _, ev := range evidence {
		if ev.Found {
			detected = append(detected, ev.Technology)
		}
	}
	return detected, nil
}

func (d *Detector) DetectWithEvidence(dir string) ([]DetectionEvidence, error) {
	var evidence []DetectionEvidence

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
		ev := d.CheckRuleWithEvidence(dir, rule)
		evidence = append(evidence, ev)
	}
	rulesTime := time.Since(start)

	if d.Verbose {
		fmt.Printf("VCS: %v, Git: %v, Rules: %v\n", vcsTime, gitTime, rulesTime)
	}

	return evidence, nil
}

func (d *Detector) GetRules() []DetectionRule {
	return detectionRules
}

func (d *Detector) CheckRule(dir string, rule DetectionRule) (bool, error) {
	evidence := d.CheckRuleWithEvidence(dir, rule)
	return evidence.Found, nil
}

func (d *Detector) CheckRuleWithEvidence(dir string, rule DetectionRule) DetectionEvidence {
	evidence := DetectionEvidence{
		Technology:    rule.Technology,
		Found:         false,
		MatchedFiles:  []string{},
		PatternCounts: make(map[string]int),
		Method:        "",
	}

	// Special case: Git detection should use VCS walking logic
	if rule.Technology == Git {
		evidence.Found = d.VCSType == vcs.Git
		evidence.Method = "vcs-detection"
		if evidence.Found {
			evidence.MatchedFiles = []string{".git"}
		}
		return evidence
	}

	// For other technologies, try VCS-aware detection first
	if d.VCSType == vcs.Git {
		if vcsEvidence := d.checkRuleWithTrackedFiles(rule, d.TrackedFiles); vcsEvidence.Found {
			return vcsEvidence
		}
		// If VCS detection fails, fall back to directory-only approach
	}

	// Fallback to current directory-only approach
	return d.checkRuleDirectoryOnly(dir, rule)
}

func (d *Detector) checkRuleDirectoryOnly(dir string, rule DetectionRule) DetectionEvidence {
	evidence := DetectionEvidence{
		Technology:    rule.Technology,
		Found:         false,
		MatchedFiles:  []string{},
		PatternCounts: make(map[string]int),
		Method:        "directory-scan",
	}

	for _, file := range rule.Files {
		if containsWildcard(file) {
			matches, err := filepath.Glob(filepath.Join(dir, file))
			if err != nil {
				continue
			}
			if len(matches) > 0 {
				evidence.Found = true
				evidence.PatternCounts[file] = len(matches)
				for _, match := range matches {
					evidence.MatchedFiles = append(evidence.MatchedFiles, filepath.Base(match))
				}
			}
		} else {
			path := filepath.Join(dir, file)
			if _, err := os.Stat(path); err == nil {
				evidence.Found = true
				evidence.MatchedFiles = append(evidence.MatchedFiles, file)
			}
		}
	}
	return evidence
}

func (d *Detector) checkRuleWithTrackedFiles(rule DetectionRule, trackedFiles []string) DetectionEvidence {
	evidence := DetectionEvidence{
		Technology:    rule.Technology,
		Found:         false,
		MatchedFiles:  []string{},
		PatternCounts: make(map[string]int),
		Method:        "git-tracked",
	}

	for _, file := range rule.Files {
		if containsWildcard(file) {
			// For wildcards, check patterns against tracked files
			matchedFiles := []string{}
			for _, trackedFile := range trackedFiles {
				if matched, _ := filepath.Match(file, filepath.Base(trackedFile)); matched {
					matchedFiles = append(matchedFiles, trackedFile)
				}
			}
			if len(matchedFiles) > 0 {
				evidence.Found = true
				evidence.PatternCounts[file] = len(matchedFiles)
				evidence.MatchedFiles = append(evidence.MatchedFiles, matchedFiles...)
			}
		} else {
			// Fast exact match lookup
			if d.fileIndex[file] {
				evidence.Found = true
				evidence.MatchedFiles = append(evidence.MatchedFiles, file)
			}
		}
	}
	return evidence
}

func containsWildcard(path string) bool {
	base := filepath.Base(path)
	return len(base) > 0 && (base[0] == '*' || base[len(base)-1] == '*')
}

func (e DetectionEvidence) FormatEvidence() string {
	if !e.Found {
		return "not detected"
	}

	if len(e.MatchedFiles) == 0 {
		return "detected"
	}

	// If we have pattern counts, format them appropriately
	if len(e.PatternCounts) > 0 {
		var parts []string
		for pattern, count := range e.PatternCounts {
			if count <= 3 {
				// Show individual files for small counts
				matchedForPattern := []string{}
				for _, file := range e.MatchedFiles {
					if matched, _ := filepath.Match(pattern, file); matched {
						matchedForPattern = append(matchedForPattern, file)
					}
				}
				if len(matchedForPattern) > 0 {
					return joinFiles(matchedForPattern)
				}
			} else {
				// Show count for large numbers
				parts = append(parts, fmt.Sprintf("%q (%d files)", pattern, count))
			}
		}
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// For exact matches, show the files
	if len(e.MatchedFiles) <= 3 {
		return joinFiles(e.MatchedFiles)
	}

	return fmt.Sprintf("%d files", len(e.MatchedFiles))
}

func joinFiles(files []string) string {
	if len(files) == 1 {
		return fmt.Sprintf("%q", files[0])
	}
	quotedFiles := make([]string, len(files))
	for i, file := range files {
		quotedFiles[i] = fmt.Sprintf("%q", file)
	}
	return strings.Join(quotedFiles, ", ")
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
