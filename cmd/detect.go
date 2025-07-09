package cmd

import (
	"fmt"
	"os"

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

	evidence, err := detector.DetectWithEvidence(cwd)
	if err != nil {
		return fmt.Errorf("failed to detect technologies: %w", err)
	}

	if detectVerbose {
		fmt.Printf("Checking detection rules in %s:\n", cwd)

		for _, result := range evidence {
			if result.Found {
				fmt.Printf("✓ %s: %s\n", result.Technology, result.FormatEvidence())
			} else {
				fmt.Printf("✗ %s: not detected\n", result.Technology)
			}
		}
		return nil
	}

	// Non-verbose mode: only show detected technologies
	var detected []detect.Technology
	for _, result := range evidence {
		if result.Found {
			detected = append(detected, result.Technology)
		}
	}

	if len(detected) == 0 {
		return nil
	}

	for _, tech := range detected {
		fmt.Printf("%s\n", tech)
	}

	return nil
}
