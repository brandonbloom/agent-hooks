package doctor

import "github.com/brandonbloom/agent-hooks/internal/detect"

type ToolRequirement struct {
	Technology detect.Technology
	Tool       string
	// Required determines how the doctor command reports missing tools:
	// - true: Missing tool shows as ERROR (project can't function without it)
	// - false: Missing tool shows as WARNING (optional, alternatives may exist)
	Required bool
	Desc     string
}

// Tool requirements are sorted alphabetically by technology to minimize merge conflicts
// when adding new requirements. Please maintain this order.
var toolRequirements = []ToolRequirement{
	{Technology: detect.Angular, Tool: "node", Required: true, Desc: "Node.js runtime for Angular"},
	{Technology: detect.Angular, Tool: "npm", Required: false, Desc: "Package manager for Angular"},
	{Technology: detect.Angular, Tool: "ng", Required: false, Desc: "Angular CLI"},
	{Technology: detect.Clojure, Tool: "clojure", Required: true, Desc: "Clojure CLI tool"},
	{Technology: detect.Clojure, Tool: "lein", Required: false, Desc: "Leiningen build tool"},
	{Technology: detect.Git, Tool: "git", Required: true, Desc: "Git version control system"},
	{Technology: detect.Go, Tool: "go", Required: true, Desc: "Go compiler and toolchain"},
	{Technology: detect.Go, Tool: "gofmt", Required: false, Desc: "Go code formatter"},
	{Technology: detect.Hurl, Tool: "hurl", Required: true, Desc: "Hurl HTTP testing tool"},
	{Technology: detect.Java, Tool: "java", Required: true, Desc: "Java runtime"},
	{Technology: detect.Java, Tool: "javac", Required: true, Desc: "Java compiler"},
	{Technology: detect.NextJS, Tool: "node", Required: true, Desc: "Node.js runtime for Next.js"},
	{Technology: detect.NextJS, Tool: "npm", Required: false, Desc: "Package manager for Next.js"},
	{Technology: detect.NodeJS, Tool: "node", Required: true, Desc: "Node.js runtime"},
	{Technology: detect.NodeJS, Tool: "npm", Required: false, Desc: "Node.js package manager"},
	{Technology: detect.Nuxt, Tool: "node", Required: true, Desc: "Node.js runtime for Nuxt.js"},
	{Technology: detect.Nuxt, Tool: "npm", Required: false, Desc: "Package manager for Nuxt.js"},
	{Technology: detect.Python, Tool: "python", Required: true, Desc: "Python interpreter"},
	{Technology: detect.Python, Tool: "pip", Required: false, Desc: "Python package manager"},
	{Technology: detect.React, Tool: "node", Required: true, Desc: "Node.js runtime for React"},
	{Technology: detect.React, Tool: "npm", Required: false, Desc: "Package manager for React"},
	{Technology: detect.Ruby, Tool: "ruby", Required: true, Desc: "Ruby interpreter"},
	{Technology: detect.Ruby, Tool: "gem", Required: false, Desc: "Ruby package manager"},
	{Technology: detect.Rust, Tool: "cargo", Required: true, Desc: "Rust package manager"},
	{Technology: detect.Rust, Tool: "rustc", Required: true, Desc: "Rust compiler"},
	{Technology: detect.Svelte, Tool: "node", Required: true, Desc: "Node.js runtime for Svelte"},
	{Technology: detect.Svelte, Tool: "npm", Required: false, Desc: "Package manager for Svelte"},
	{Technology: detect.Transcript, Tool: "transcript", Required: true, Desc: "Transcript testing tool"},
	{Technology: detect.Vue, Tool: "node", Required: true, Desc: "Node.js runtime for Vue.js"},
	{Technology: detect.Vue, Tool: "npm", Required: false, Desc: "Package manager for Vue.js"},
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
