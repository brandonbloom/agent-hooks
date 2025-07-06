package cmd

import (
	"fmt"

	"github.com/brandonbloom/agent-hooks/internal/vcs"
	"github.com/spf13/cobra"
)

var whichVcsCmd = &cobra.Command{
	Use:   "which-vcs",
	Short: "Detect which version control system is in use",
	Long:  `Detects which version control system is being used in the current directory or any parent directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		detectedVcs, err := vcs.DetectVCS()
		if err != nil {
			return err
		}
		
		fmt.Println(detectedVcs)
		return nil
	},
}