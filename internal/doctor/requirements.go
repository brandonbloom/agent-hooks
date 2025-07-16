package doctor

import (
	"github.com/brandonbloom/agent-hooks/internal/detect"
)

// FormattingToolSupport defines which tools can format which file extensions
type FormattingToolSupport struct {
	Extensions []string // File extensions this applies to
	Tools      []string // Available tools in preference order (first = most preferred)
}

// ToolRequirement represents an association between a technology and a tool.
// It specifies that when a particular technology is detected in a project,
// the associated tool should be checked with the given required/optional status.
type ToolRequirement struct {
	Technology detect.Technology
	Tool       string // Name of tool from AllTools registry
	// Required determines how the doctor command reports missing tools:
	// - true: Missing tool shows as ERROR (project can't function without it)
	// - false: Missing tool shows as WARNING (optional, alternatives may exist)
	Required bool
}

// Tool requirements are sorted alphabetically by technology to minimize merge conflicts
// when adding new requirements. Please maintain this order.
var toolRequirements = []ToolRequirement{
	{Technology: detect.Angular, Tool: "node", Required: true},
	{Technology: detect.Angular, Tool: "npm", Required: false},
	{Technology: detect.Angular, Tool: "ng", Required: false},
	{Technology: detect.Clojure, Tool: "clojure", Required: true},
	{Technology: detect.Clojure, Tool: "lein", Required: false},
	{Technology: detect.Git, Tool: "git", Required: true},
	{Technology: detect.Go, Tool: "go", Required: true},
	{Technology: detect.Go, Tool: "gofmt", Required: false},
	{Technology: detect.Hurl, Tool: "hurl", Required: true},
	{Technology: detect.Java, Tool: "java", Required: true},
	{Technology: detect.Java, Tool: "javac", Required: true},
	{Technology: detect.NextJS, Tool: "node", Required: true},
	{Technology: detect.NextJS, Tool: "npm", Required: false},
	{Technology: detect.NodeJS, Tool: "node", Required: true},
	{Technology: detect.NodeJS, Tool: "npm", Required: false},
	{Technology: detect.Nuxt, Tool: "node", Required: true},
	{Technology: detect.Nuxt, Tool: "npm", Required: false},
	{Technology: detect.Procfile, Tool: "procfile-runner", Required: false},
	{Technology: detect.Python, Tool: "python", Required: true},
	{Technology: detect.Python, Tool: "pip", Required: false},
	{Technology: detect.React, Tool: "node", Required: true},
	{Technology: detect.React, Tool: "npm", Required: false},
	{Technology: detect.Ruby, Tool: "ruby", Required: true},
	{Technology: detect.Ruby, Tool: "gem", Required: false},
	{Technology: detect.Rust, Tool: "cargo", Required: true},
	{Technology: detect.Rust, Tool: "rustc", Required: true},
	{Technology: detect.Svelte, Tool: "node", Required: true},
	{Technology: detect.Svelte, Tool: "npm", Required: false},
	{Technology: detect.Swift, Tool: "swift", Required: true},
	{Technology: detect.Transcript, Tool: "transcript", Required: true},
	{Technology: detect.Vue, Tool: "node", Required: true},
	{Technology: detect.Vue, Tool: "npm", Required: false},
}

func GetToolRequirements(tech detect.Technology) []ToolRequirement {
	var requirements []ToolRequirement
	for _, req := range toolRequirements {
		if req.Technology == tech {
			requirements = append(requirements, req)
		}
	}
	return requirements
}

// Formatting tool support - tools are listed in preference order
var formattingSupport = []FormattingToolSupport{
	{
		Extensions: []string{".go"},
		Tools:      []string{"goimports", "gofmt"}, // goimports preferred over gofmt
	},
	{
		Extensions: []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".mts", ".cts"},
		Tools:      []string{"biome", "prettier"}, // biome preferred over prettier
	},
}

func GetFormattingSupport() []FormattingToolSupport {
	return formattingSupport
}
