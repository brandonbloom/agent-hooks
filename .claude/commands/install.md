---
description: Install the latest version of agent-hooks globally after checking git status
allowed-tools: Bash(git:*), Bash(go install:*), Bash(which:*), Bash(agent-hooks:*)
---

# Install agent-hooks globally

Please check if all local changes have been pushed and then install the latest version of agent-hooks globally.

## Current git status
!`git status --porcelain`

## Check if branch is up to date with remote
!`git status -uno`

## Your task

1. **First, verify git status:**
   - Check if there are any uncommitted changes (working directory should be clean)
   - Check if the current branch has any unpushed commits
   - If either condition fails, bail out with an appropriate error message

2. **If git status is clean:**
   - Run `go install .` to install the latest version globally
   - Verify the installation was successful by:
     - Checking that `agent-hooks` command is available
     - Running `agent-hooks --version` to show the installed version
     - Showing the path where it was installed with `which agent-hooks`

3. **Report the results:**
   - Success message with version and installation path
   - Or appropriate error message if anything failed