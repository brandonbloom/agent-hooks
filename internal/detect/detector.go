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

type TechDetector interface {
	Detect(dir string) ([]Technology, error)
	GetRules() []DetectionRule
	GetToolRequirements(tech Technology) []ToolRequirement
	CheckRule(dir string, rule DetectionRule) (bool, error)
}

type ToolRequirement struct {
	Technology Technology
	Tool       string
	// Required determines how the doctor command reports missing tools:
	// - true: Missing tool shows as ERROR (project can't function without it)
	// - false: Missing tool shows as WARNING (optional, alternatives may exist)
	Required bool
	Desc     string
}

type DefaultDetector struct{}

var detectionRules = []DetectionRule{
	{Technology: Git, Files: []string{".git"}, Desc: "Git repository"},
	{Technology: Go, Files: []string{"go.mod", "*.go"}, Desc: "Go module or Go files"},
	{Technology: NodeJS, Files: []string{"package.json"}, Desc: "Node.js package"},
	{Technology: React, Files: []string{"*.jsx", "*.tsx"}, Desc: "React project"},
	{Technology: Vue, Files: []string{"*.vue", "vue.config.js", "vue.config.ts"}, Desc: "Vue.js project"},
	{Technology: Svelte, Files: []string{"*.svelte", "svelte.config.js", "vite.config.js"}, Desc: "Svelte project"},
	{Technology: NextJS, Files: []string{"next.config.js", "next.config.mjs", "next.config.ts"}, Desc: "Next.js project"},
	{Technology: Nuxt, Files: []string{"nuxt.config.js", "nuxt.config.ts"}, Desc: "Nuxt.js project"},
	{Technology: Angular, Files: []string{"angular.json", "*.component.ts"}, Desc: "Angular project"},
	{Technology: Python, Files: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}, Desc: "Python project"},
	{Technology: Rust, Files: []string{"Cargo.toml"}, Desc: "Rust project"},
	{Technology: Ruby, Files: []string{"Gemfile"}, Desc: "Ruby project"},
	{Technology: Java, Files: []string{"pom.xml", "build.gradle"}, Desc: "Java project"},
	{Technology: Clojure, Files: []string{"project.clj", "deps.edn", "shadow-cljs.edn", "bb.edn"}, Desc: "Clojure project"},
	{Technology: Transcript, Files: []string{"*.cmdt"}, Desc: "Transcript test files"},
	{Technology: Hurl, Files: []string{"*.hurl"}, Desc: "Hurl HTTP test files"},
	{Technology: C, Files: []string{"*.c", "*.h"}, Desc: "C source files"},
	{Technology: Cpp, Files: []string{"*.cpp", "*.cc", "*.cxx", "*.hpp", "*.hh", "*.hxx"}, Desc: "C++ source files"},
	{Technology: TypeScript, Files: []string{"*.ts", "*.dts"}, Desc: "TypeScript source files"},
	{Technology: JavaScript, Files: []string{"*.js", "*.mjs", "*.cjs"}, Desc: "JavaScript source files"},
	{Technology: Swift, Files: []string{"*.swift"}, Desc: "Swift source files"},
	{Technology: Kotlin, Files: []string{"*.kt", "*.kts"}, Desc: "Kotlin source files"},
	{Technology: CSharp, Files: []string{"*.cs"}, Desc: "C# source files"},
	{Technology: PHP, Files: []string{"*.php"}, Desc: "PHP source files"},
	{Technology: Dart, Files: []string{"*.dart"}, Desc: "Dart source files"},
	{Technology: Haskell, Files: []string{"*.hs"}, Desc: "Haskell source files"},
	{Technology: Perl, Files: []string{"*.pl"}, Desc: "Perl source files"},
	{Technology: Lua, Files: []string{"*.lua"}, Desc: "Lua source files"},
	{Technology: Shell, Files: []string{"*.sh", "*.bash", "*.zsh", "*.fish"}, Desc: "Shell scripts"},
	{Technology: Markdown, Files: []string{"*.md", "*.markdown"}, Desc: "Markdown files"},
	{Technology: HTML, Files: []string{"*.html", "*.htm"}, Desc: "HTML files"},
	{Technology: CSS, Files: []string{"*.css", "*.scss", "*.sass", "*.less"}, Desc: "CSS and preprocessor files"},
	{Technology: OCaml, Files: []string{"*.ml", "*.mli"}, Desc: "OCaml source files"},
	{Technology: VimScript, Files: []string{"*.vim"}, Desc: "Vim script files"},
	{Technology: Assembly, Files: []string{"*.asm", "*.s", "*.S"}, Desc: "Assembly source files"},
	{Technology: CoffeeScript, Files: []string{"*.coffee"}, Desc: "CoffeeScript source files"},
	{Technology: Elixir, Files: []string{"*.ex", "*.exs"}, Desc: "Elixir source files"},
	{Technology: Erlang, Files: []string{"*.erl"}, Desc: "Erlang source files"},
	{Technology: Fortran, Files: []string{"*.f90"}, Desc: "Fortran source files"},
	{Technology: Zig, Files: []string{"*.zig"}, Desc: "Zig source files"},
	{Technology: JSON, Files: []string{"*.json"}, Desc: "JSON files"},
	{Technology: YAML, Files: []string{"*.yaml", "*.yml"}, Desc: "YAML files"},
	{Technology: XML, Files: []string{"*.xml"}, Desc: "XML files"},
	{Technology: SQL, Files: []string{"*.sql"}, Desc: "SQL files"},
	{Technology: ProtocolBuffers, Files: []string{"*.proto"}, Desc: "Protocol Buffer files"},
	{Technology: GraphQL, Files: []string{"*.graphql", "*.gql"}, Desc: "GraphQL files"},
	{Technology: R, Files: []string{"*.r", "*.R"}, Desc: "R source files"},
	{Technology: TOML, Files: []string{"*.toml"}, Desc: "TOML files"},
	{Technology: INI, Files: []string{"*.ini", "*.cfg", "*.conf"}, Desc: "INI/Config files"},
	{Technology: LaTeX, Files: []string{"*.tex"}, Desc: "LaTeX files"},
	{Technology: PowerShell, Files: []string{"*.ps1"}, Desc: "PowerShell scripts"},
	{Technology: Batch, Files: []string{"*.bat", "*.cmd"}, Desc: "Batch files"},
	{Technology: Make, Files: []string{"Makefile", "makefile", "*.mk"}, Desc: "Makefiles"},
}

var toolRequirements = []ToolRequirement{
	{Technology: Git, Tool: "git", Required: true, Desc: "Git version control system"},
	{Technology: Go, Tool: "go", Required: true, Desc: "Go compiler and toolchain"},
	{Technology: Go, Tool: "gofmt", Required: false, Desc: "Go code formatter"},
	{Technology: NodeJS, Tool: "node", Required: true, Desc: "Node.js runtime"},
	{Technology: NodeJS, Tool: "npm", Required: false, Desc: "Node.js package manager"},
	{Technology: React, Tool: "node", Required: true, Desc: "Node.js runtime for React"},
	{Technology: React, Tool: "npm", Required: false, Desc: "Package manager for React"},
	{Technology: Vue, Tool: "node", Required: true, Desc: "Node.js runtime for Vue.js"},
	{Technology: Vue, Tool: "npm", Required: false, Desc: "Package manager for Vue.js"},
	{Technology: Svelte, Tool: "node", Required: true, Desc: "Node.js runtime for Svelte"},
	{Technology: Svelte, Tool: "npm", Required: false, Desc: "Package manager for Svelte"},
	{Technology: NextJS, Tool: "node", Required: true, Desc: "Node.js runtime for Next.js"},
	{Technology: NextJS, Tool: "npm", Required: false, Desc: "Package manager for Next.js"},
	{Technology: Nuxt, Tool: "node", Required: true, Desc: "Node.js runtime for Nuxt.js"},
	{Technology: Nuxt, Tool: "npm", Required: false, Desc: "Package manager for Nuxt.js"},
	{Technology: Angular, Tool: "node", Required: true, Desc: "Node.js runtime for Angular"},
	{Technology: Angular, Tool: "npm", Required: false, Desc: "Package manager for Angular"},
	{Technology: Angular, Tool: "ng", Required: false, Desc: "Angular CLI"},
	{Technology: Python, Tool: "python", Required: true, Desc: "Python interpreter"},
	{Technology: Python, Tool: "pip", Required: false, Desc: "Python package manager"},
	{Technology: Rust, Tool: "cargo", Required: true, Desc: "Rust package manager"},
	{Technology: Rust, Tool: "rustc", Required: true, Desc: "Rust compiler"},
	{Technology: Ruby, Tool: "ruby", Required: true, Desc: "Ruby interpreter"},
	{Technology: Ruby, Tool: "gem", Required: false, Desc: "Ruby package manager"},
	{Technology: Java, Tool: "java", Required: true, Desc: "Java runtime"},
	{Technology: Java, Tool: "javac", Required: true, Desc: "Java compiler"},
	{Technology: Clojure, Tool: "clojure", Required: true, Desc: "Clojure CLI tool"},
	{Technology: Clojure, Tool: "lein", Required: false, Desc: "Leiningen build tool"},
	{Technology: Transcript, Tool: "transcript", Required: true, Desc: "Transcript testing tool"},
	{Technology: Hurl, Tool: "hurl", Required: true, Desc: "Hurl HTTP testing tool"},
}

func NewDetector() TechDetector {
	return &DefaultDetector{}
}

func (d *DefaultDetector) Detect(dir string) ([]Technology, error) {
	var detected []Technology

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

func (d *DefaultDetector) CheckRule(dir string, rule DetectionRule) (bool, error) {
	// Special case: Git detection should use VCS walking logic
	if rule.Technology == Git {
		vcsType, err := vcs.DetectVCS()
		return err == nil && vcsType == vcs.Git, nil
	}

	// For other technologies, try VCS-aware detection first
	if vcsType, err := vcs.DetectVCS(); err == nil && vcsType == vcs.Git {
		if found, err := d.checkRuleWithTrackedFiles(rule); err == nil {
			return found, nil
		}
		// If VCS detection fails, fall back to directory-only approach
	}

	// Fallback to current directory-only approach
	return d.checkRuleDirectoryOnly(dir, rule)
}

func (d *DefaultDetector) checkRuleDirectoryOnly(dir string, rule DetectionRule) (bool, error) {
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

func (d *DefaultDetector) checkRuleWithTrackedFiles(rule DetectionRule) (bool, error) {
	trackedFiles, err := git.GetAllTrackedFiles()
	if err != nil {
		return false, err
	}

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
