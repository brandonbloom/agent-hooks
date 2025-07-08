package vcs

import (
	"fmt"
	"os"
	"path/filepath"
)

type VCS string

const (
	Git     VCS = "git"
	Unknown VCS = "unknown"
)

func DetectVCS() (VCS, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Unknown, fmt.Errorf("failed to get current directory: %w", err)
	}

	if isGitRepository(cwd) {
		return Git, nil
	}

	return Unknown, fmt.Errorf("unsupported or unknown version control system")
}

func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	root, found := findGitRoot(cwd)
	if !found {
		return "", fmt.Errorf("not in a git repository")
	}
	return root, nil
}

func findGitRoot(dir string) (string, bool) {
	current := dir
	for {
		gitDir := filepath.Join(current, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			if info.IsDir() || info.Mode().IsRegular() {
				return current, true
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", false
}

func isGitRepository(dir string) bool {
	_, found := findGitRoot(dir)
	return found
}
