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

func isGitRepository(dir string) bool {
	current := dir
	for {
		gitDir := filepath.Join(current, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			return info.IsDir() || info.Mode().IsRegular()
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return false
}