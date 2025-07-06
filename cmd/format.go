package cmd

import (
	"fmt"
	"os"

	"github.com/brandonbloom/agent-hooks/internal/format"
	"github.com/brandonbloom/agent-hooks/internal/git"
	"github.com/brandonbloom/agent-hooks/internal/vcs"
	"github.com/spf13/cobra"
)

var allFiles bool

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format code in the current project",
	Long: `Formats code in the current project using appropriate tools.
Only formats changed files by default. Use --all-files to format all tracked files.
Currently requires a Git repository and supports Go files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		detectedVcs, err := vcs.DetectVCS()
		if err != nil {
			return fmt.Errorf("cannot format: %w", err)
		}

		if detectedVcs != vcs.Git {
			return fmt.Errorf("formatting is only supported in Git repositories, detected: %s", detectedVcs)
		}

		var filesToFormat []string

		if allFiles {
			trackedFiles, err := git.GetAllTrackedFiles()
			if err != nil {
				return fmt.Errorf("failed to get tracked files: %w", err)
			}
			filesToFormat = trackedFiles
		} else {
			changedFiles, err := git.GetChangedFiles()
			if err != nil {
				return fmt.Errorf("failed to get changed files: %w", err)
			}

			for _, file := range changedFiles {
				filesToFormat = append(filesToFormat, file.Path)
			}
		}

		if len(filesToFormat) == 0 {
			return nil
		}

		result := format.FormatFiles(filesToFormat)

		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}

		if len(result.Errors) > 0 {
			for _, errMsg := range result.Errors {
				fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
			}
			return fmt.Errorf("formatting failed")
		}

		return nil
	},
}

func init() {
	formatCmd.Flags().BoolVar(&allFiles, "all-files", false, "Format all tracked files instead of just changed files")
}
