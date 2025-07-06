package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type FileStatus struct {
	Path   string
	Status string
}

func GetChangedFiles() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	var files []FileStatus
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}

		status := strings.TrimSpace(line[:2])
		path := strings.TrimSpace(line[3:])

		files = append(files, FileStatus{
			Path:   path,
			Status: status,
		})
	}

	return files, scanner.Err()
}

func GetAllTrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tracked files: %w", err)
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			files = append(files, line)
		}
	}

	return files, scanner.Err()
}
