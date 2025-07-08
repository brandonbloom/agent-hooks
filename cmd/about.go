package cmd

import (
	"fmt"
	"strings"

	"github.com/brandonbloom/agent-hooks/internal/detect"
	"github.com/brandonbloom/agent-hooks/internal/doctor"
	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about <name>",
	Short: "Show information about a technology or tool",
	Long: `Show detailed information about a technology or tool that agent-hooks knows about.
This includes the name, description, file patterns/extensions, and reference URL.

Examples:
  agent-hooks about go
  agent-hooks about git
  agent-hooks about typescript`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.ToLower(args[0])

		// First check if it's a technology
		detector := &detect.Detector{}
		rules := detector.GetRules()

		for _, rule := range rules {
			if strings.ToLower(string(rule.Technology)) == name {
				fmt.Printf("Name: %s\n", rule.Technology)
				fmt.Printf("Type: Technology\n")
				fmt.Printf("Description: %s\n", rule.Desc)
				fmt.Printf("File patterns: %s\n", strings.Join(rule.Files, ", "))
				fmt.Printf("URL: %s\n", rule.URL)
				return nil
			}
		}

		// Then check if it's a tool
		for _, tool := range doctor.AllTools {
			if strings.ToLower(tool.Name) == name {
				fmt.Printf("Name: %s\n", tool.Name)
				fmt.Printf("Type: Tool\n")
				if tool.Command != "" {
					fmt.Printf("Command: %s\n", tool.Command)
				} else {
					fmt.Printf("Command: (meta-tool)\n")
				}
				if tool.Required {
					fmt.Printf("Required: Yes\n")
				} else {
					fmt.Printf("Required: No\n")
				}
				fmt.Printf("URL: %s\n", tool.URL)
				return nil
			}
		}

		return fmt.Errorf("unknown technology or tool: %s", name)
	},
}
