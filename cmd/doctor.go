package cmd

import (
	"fmt"
	"os"

	"github.com/brandonbloom/agent-hooks/internal/doctor"
	"github.com/spf13/cobra"
)

var verbose bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check development environment and Claude Code setup",
	Long: `Diagnoses common issues with development tools and Claude Code integration.
By default, only shows problems. Use --verbose to see all checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var allResults []doctor.CheckResult

		toolResults := doctor.RunToolChecks(verbose)
		allResults = append(allResults, toolResults...)

		projectResults := doctor.RunProjectChecks(verbose)
		allResults = append(allResults, projectResults...)

		claudeResults := doctor.RunClaudeChecks(verbose)
		allResults = append(allResults, claudeResults...)

		hasProblems := false

		for _, result := range allResults {
			switch result.Status {
			case doctor.CheckPassed:
				if verbose {
					fmt.Printf("✓ %s", result.Name)
					if result.Message != "" {
						fmt.Printf(": %s", result.Message)
					}
					fmt.Println()
				}
			case doctor.CheckWarning:
				hasProblems = true
				fmt.Fprintf(os.Stderr, "Warning: %s", result.Message)
				fmt.Fprintln(os.Stderr)
			case doctor.CheckFailed:
				hasProblems = true
				fmt.Fprintf(os.Stderr, "Error: %s", result.Message)
				fmt.Fprintln(os.Stderr)
			}
		}

		if verbose && len(allResults) > 0 {
			passed := 0
			warnings := 0
			failed := 0

			for _, result := range allResults {
				switch result.Status {
				case doctor.CheckPassed:
					passed++
				case doctor.CheckWarning:
					warnings++
				case doctor.CheckFailed:
					failed++
				}
			}

			fmt.Printf("\nSummary: %d passed", passed)
			if failed > 0 {
				fmt.Printf(", %d failed", failed)
			}
			if warnings > 0 {
				fmt.Printf(", %d warnings", warnings)
			}
			fmt.Println()
		}

		if hasProblems {
			return fmt.Errorf("found issues in environment setup")
		}

		return nil
	},
}

func init() {
	doctorCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show all checks, not just problems")
}
