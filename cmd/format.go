package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

// resolveFilePaths converts relative paths from Git to absolute paths
func resolveFilePaths(files []string) ([]string, error) {
	gitRoot, err := vcs.FindProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	resolvedFiles := make([]string, len(files))
	for i, file := range files {
		if filepath.IsAbs(file) {
			resolvedFiles[i] = file
		} else {
			resolvedFiles[i] = filepath.Join(gitRoot, file)
		}
	}
	return resolvedFiles, nil
}

var formatCmd = &cobra.Command{
	Use:   "format [files...]",
	Short: "Format code in the current project",
	Long: `Formats code in the current project using appropriate tools.
With no arguments, formats only changed files.
With file arguments, formats only those specific files.
Use --all-files to format all tracked files (mutually exclusive with file arguments).
Use --dry-run to preview what would be formatted without making changes.
Currently requires a Git repository and supports Go files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate mutually exclusive options
		if allFiles && len(args) > 0 {
			return fmt.Errorf("cannot use --all-files with specific file arguments")
		}

		detectedVcs, err := vcs.DetectVCS()
		if err != nil {
			return fmt.Errorf("cannot format: %w", err)
		}

		if detectedVcs != vcs.Git {
			return fmt.Errorf("formatting is only supported in Git repositories, detected: %s", detectedVcs)
		}

		var filesToFormat []string

		if len(args) > 0 {
			// Format specific files provided as arguments
			filesToFormat = args
		} else if allFiles {
			// Format all tracked files
			trackedFiles, err := git.GetAllTrackedFiles()
			if err != nil {
				return fmt.Errorf("failed to get tracked files: %w", err)
			}
			filesToFormat = trackedFiles
		} else {
			// Format only changed files (default behavior)
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

		// Resolve relative paths from Git to absolute paths
		resolvedFiles, err := resolveFilePaths(filesToFormat)
		if err != nil {
			return fmt.Errorf("failed to resolve file paths: %w", err)
		}
		filesToFormat = resolvedFiles

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
			return fmt.Errorf("%s", result.Errors[0])
		}

		return nil
	},
}

func init() {
	formatCmd.Flags().BoolVar(&allFiles, "all-files", false, "Format all tracked files instead of just changed files")
	formatCmd.Flags().BoolVarP(&formatVerbose, "verbose", "v", false, "Show detailed output about formatting operations")
	formatCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview what would be formatted without making changes")
}
