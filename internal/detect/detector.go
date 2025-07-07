package detect

import (
	"fmt"
	"os"
	"path/filepath"
)

type Technology string

const (
	Git        Technology = "git"
	Go         Technology = "go"
	NodeJS     Technology = "nodejs"
	Python     Technology = "python"
	Rust       Technology = "rust"
	Ruby       Technology = "ruby"
	Java       Technology = "java"
	Transcript Technology = "transcript"
	Hurl       Technology = "hurl"
)

type DetectionRule struct {
	Technology Technology
	Files      []string
	Desc       string
}

type TechDetector interface {
	Detect(dir string) ([]Technology, error)
	GetRules() []DetectionRule
	GetToolRequirements(tech Technology) []ToolRequirement
}

type ToolRequirement struct {
	Technology Technology
	Tool       string
	Required   bool
	Desc       string
}

type DefaultDetector struct{}

var detectionRules = []DetectionRule{
	{Technology: Git, Files: []string{".git"}, Desc: "Git repository"},
	{Technology: Go, Files: []string{"go.mod"}, Desc: "Go module"},
	{Technology: NodeJS, Files: []string{"package.json"}, Desc: "Node.js package"},
	{Technology: Python, Files: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}, Desc: "Python project"},
	{Technology: Rust, Files: []string{"Cargo.toml"}, Desc: "Rust project"},
	{Technology: Ruby, Files: []string{"Gemfile"}, Desc: "Ruby project"},
	{Technology: Java, Files: []string{"pom.xml", "build.gradle"}, Desc: "Java project"},
	{Technology: Transcript, Files: []string{"*.cmdt"}, Desc: "Transcript test files"},
	{Technology: Hurl, Files: []string{"*.hurl"}, Desc: "Hurl HTTP test files"},
}

var toolRequirements = []ToolRequirement{
	{Technology: Git, Tool: "git", Required: true, Desc: "Git version control system"},
	{Technology: Go, Tool: "go", Required: true, Desc: "Go compiler and toolchain"},
	{Technology: Go, Tool: "gofmt", Required: false, Desc: "Go code formatter"},
	{Technology: NodeJS, Tool: "node", Required: true, Desc: "Node.js runtime"},
	{Technology: NodeJS, Tool: "npm", Required: false, Desc: "Node.js package manager"},
	{Technology: Python, Tool: "python", Required: true, Desc: "Python interpreter"},
	{Technology: Python, Tool: "pip", Required: false, Desc: "Python package manager"},
	{Technology: Rust, Tool: "cargo", Required: true, Desc: "Rust package manager"},
	{Technology: Rust, Tool: "rustc", Required: true, Desc: "Rust compiler"},
	{Technology: Ruby, Tool: "ruby", Required: true, Desc: "Ruby interpreter"},
	{Technology: Ruby, Tool: "gem", Required: false, Desc: "Ruby package manager"},
	{Technology: Java, Tool: "java", Required: true, Desc: "Java runtime"},
	{Technology: Java, Tool: "javac", Required: true, Desc: "Java compiler"},
	{Technology: Transcript, Tool: "transcript", Required: true, Desc: "Transcript testing tool"},
	{Technology: Hurl, Tool: "hurl", Required: true, Desc: "Hurl HTTP testing tool"},
}

func NewDetector() TechDetector {
	return &DefaultDetector{}
}

func (d *DefaultDetector) Detect(dir string) ([]Technology, error) {
	var detected []Technology

	for _, rule := range detectionRules {
		found, err := d.checkRule(dir, rule)
		if err != nil {
			continue
		}
		if found {
			detected = append(detected, rule.Technology)
		}
	}

	return detected, nil
}

func (d *DefaultDetector) GetRules() []DetectionRule {
	return detectionRules
}

func (d *DefaultDetector) GetToolRequirements(tech Technology) []ToolRequirement {
	var requirements []ToolRequirement
	for _, req := range toolRequirements {
		if req.Technology == tech {
			requirements = append(requirements, req)
		}
	}
	return requirements
}

func (d *DefaultDetector) checkRule(dir string, rule DetectionRule) (bool, error) {
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

func containsWildcard(path string) bool {
	base := filepath.Base(path)
	return len(base) > 0 && (base[0] == '*' || base[len(base)-1] == '*')
}

func DetectInDirectory(dir string) ([]Technology, error) {
	detector := NewDetector()
	return detector.Detect(dir)
}

func DetectInCurrentDirectory() ([]Technology, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	return DetectInDirectory(cwd)
}
