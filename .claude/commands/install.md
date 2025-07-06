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
   - Run `go install github.com/brandonbloom/agent-hooks@latest` to install the latest version globally
   - Verify the installation was successful by:
     - Checking that `agent-hooks` command is available with `which agent-hooks`
     - Running `agent-hooks version` to show the installed version (should include git commit SHA)
     - Confirm the version matches the latest commit SHA from this repository

3. **Report the results:**
   - Success message showing:
     - Installation path from `which agent-hooks`
     - Full version string from `agent-hooks version`
     - Confirmation that the git SHA in the version matches the current repository state
   - Or appropriate error message if anything failed