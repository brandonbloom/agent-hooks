package detect

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brandonbloom/agent-hooks/internal/git"
	"github.com/brandonbloom/agent-hooks/internal/vcs"
)

type Technology string

const (
	Git             Technology = "git"
	Go              Technology = "go"
	NodeJS          Technology = "nodejs"
	React           Technology = "react"
	Vue             Technology = "vue"
	Svelte          Technology = "svelte"
	NextJS          Technology = "nextjs"
	Nuxt            Technology = "nuxt"
	Angular         Technology = "angular"
	Python          Technology = "python"
	Rust            Technology = "rust"
	Ruby            Technology = "ruby"
	Java            Technology = "java"
	Clojure         Technology = "clojure"
	Transcript      Technology = "transcript"
	Hurl            Technology = "hurl"
	C               Technology = "c"
	Cpp             Technology = "cpp"
	TypeScript      Technology = "typescript"
	JavaScript      Technology = "javascript"
	Swift           Technology = "swift"
	Kotlin          Technology = "kotlin"
	CSharp          Technology = "csharp"
	PHP             Technology = "php"
	Dart            Technology = "dart"
	Haskell         Technology = "haskell"
	Perl            Technology = "perl"
	Lua             Technology = "lua"
	Shell           Technology = "shell"
	Markdown        Technology = "markdown"
	HTML            Technology = "html"
	CSS             Technology = "css"
	OCaml           Technology = "ocaml"
	VimScript       Technology = "vimscript"
	Assembly        Technology = "assembly"
	CoffeeScript    Technology = "coffeescript"
	Elixir          Technology = "elixir"
	Erlang          Technology = "erlang"
	Fortran         Technology = "fortran"
	Zig             Technology = "zig"
	JSON            Technology = "json"
	YAML            Technology = "yaml"
	XML             Technology = "xml"
	SQL             Technology = "sql"
	ProtocolBuffers Technology = "protobuf"
	GraphQL         Technology = "graphql"
	R               Technology = "r"
	TOML            Technology = "toml"
	INI             Technology = "ini"
	LaTeX           Technology = "latex"
	PowerShell      Technology = "powershell"
	Batch           Technology = "batch"
	Make            Technology = "make"
)

type DetectionRule struct {
	Technology Technology
	Files      []string
	Desc       string
}

type Detector struct {
	VCSType      vcs.VCS
	TrackedFiles []string
}

func (d *Detector) Detect(dir string) ([]Technology, error) {
	var detected []Technology

	// Phase 1: Do VCS detection and file listing once (if not already set)
	if d.VCSType == "" {
		d.VCSType, _ = vcs.DetectVCS()
	}
	if d.VCSType == vcs.Git && d.TrackedFiles == nil {
		var err error
		d.TrackedFiles, err = git.GetAllTrackedFiles()
		if err != nil {
			return nil, err
		}
	}

	// Phase 2: Check each rule with pre-computed information
	for _, rule := range detectionRules {
		found, err := d.CheckRule(dir, rule)
		if err != nil {
			continue
		}
		if found {
			detected = append(detected, rule.Technology)
		}
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
			// Check pattern against all tracked files
			for _, tracked := range trackedFiles {
				if matched, _ := filepath.Match(file, filepath.Base(tracked)); matched {
					return true, nil
				}
			}
		} else {
			// Check exact match in tracked files
			for _, tracked := range trackedFiles {
				if tracked == file || filepath.Base(tracked) == file {
					return true, nil
				}
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
