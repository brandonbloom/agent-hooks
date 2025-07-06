package cmd

import (
	"fmt"

	"github.com/brandonbloom/agent-hooks/internal/vcs"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format code in the current project",
	Long: `Formats code in the current project using appropriate tools.
Currently requires a Git repository to operate.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		detectedVcs, err := vcs.DetectVCS()
		if err != nil {
			return fmt.Errorf("cannot format: %w", err)
		}
		
		if detectedVcs != vcs.Git {
			return fmt.Errorf("formatting is only supported in Git repositories, detected: %s", detectedVcs)
		}
		
		fmt.Println("Format command is not yet implemented, but Git repository detected successfully.")
		return nil
	},
}