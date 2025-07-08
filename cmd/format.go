package cmd

import (
	"fmt"
	"os"

	"github.com/brandonbloom/agent-hooks/internal/format"
	"github.com/brandonbloom/agent-hooks/internal/git"
	"github.com/brandonbloom/agent-hooks/internal/vcs"
	"github.com/spf13/cobra"
)

var (
	allFiles      bool
	formatVerbose bool
	dryRun        bool
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format code in the current project",
	Long: `Formats code in the current project using appropriate tools.
Only formats changed files by default. Use --all-files to format all tracked files.
Use --dry-run to preview what would be formatted without making changes.
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
			if formatVerbose {
				fmt.Println("No files to format")
			}
			return nil
		}

		opts := format.Options{
			DryRun:  dryRun,
			Verbose: formatVerbose,
		}

		result := format.FormatFilesWithOptions(filesToFormat, opts)

		if formatVerbose {
			if len(result.FormattedFiles) > 0 {
				action := "Formatted"
				if dryRun {
					action = "Would format"
				}
				for _, file := range result.FormattedFiles {
					fmt.Printf("%s: %s\n", action, file)
				}
			}

			if len(result.SkippedFiles) > 0 {
				for _, file := range result.SkippedFiles {
					fmt.Printf("Skipped: %s (no formatter available)\n", file)
				}
			}
		} else if dryRun {
			for _, file := range result.FormattedFiles {
				fmt.Printf("Would format: %s\n", file)
			}
		}

		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}

		if len(result.Errors) > 0 {
			for _, errMsg := range result.Errors {
				fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
			}
			cmd.SilenceUsage = true
			return fmt.Errorf("formatting failed")
		}

		return nil
	},
}

func init() {
	formatCmd.Flags().BoolVar(&allFiles, "all-files", false, "Format all tracked files instead of just changed files")
	formatCmd.Flags().BoolVarP(&formatVerbose, "verbose", "v", false, "Show detailed output about formatting operations")
	formatCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview what would be formatted without making changes")
}
