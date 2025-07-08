package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brandonbloom/agent-hooks/internal/detect"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect technologies and tools in the current project",
	Long: `Detect technologies and tools used in the current project directory.
By default, only shows detected technologies. Use --verbose to see all detection attempts.`,
	RunE: runDetect,
}

var detectVerbose bool

func init() {
	detectCmd.Flags().BoolVarP(&detectVerbose, "verbose", "v", false, "Show all detection attempts, not just detected technologies")
	rootCmd.AddCommand(detectCmd)
}

func runDetect(cmd *cobra.Command, args []string) error {
	detector := &detect.Detector{Verbose: detectVerbose}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if detectVerbose {
		technologies, err := detector.Detect(cwd)
		if err != nil {
			return fmt.Errorf("failed to detect technologies: %w", err)
		}

		rules := detector.GetRules()
		fmt.Printf("Checking detection rules in %s:\n", cwd)

		techSet := make(map[detect.Technology]bool)
		for _, tech := range technologies {
			techSet[tech] = true
		}

		for _, rule := range rules {
			if techSet[rule.Technology] {
				fmt.Printf("✓ %s: %s\n", rule.Technology, rule.Desc)
			} else {
				fmt.Printf("✗ %s: not detected\n", rule.Technology)
			}
		}
		return nil
	}

	technologies, err := detector.Detect(cwd)
	if err != nil {
		return fmt.Errorf("failed to detect technologies: %w", err)
	}

	if len(technologies) == 0 {
		return nil
	}

	for _, tech := range technologies {
		fmt.Printf("%s\n", tech)
	}

	return nil
}

func containsWildcard(path string) bool {
	return filepath.Base(path) != path && (filepath.Base(path)[0] == '*' || filepath.Base(path)[len(filepath.Base(path))-1] == '*')
}
