package cmd

import (
	"fmt"

	"github.com/brandonbloom/agent-hooks/internal/config"
	"github.com/spf13/cobra"
)

var postToolUseCmd = &cobra.Command{
	Use:   "post-tool-use",
	Short: "Hook command for Claude Code PostToolUse events",
	Long: `This command is designed to be used as a Claude Code hook for PostToolUse events.
It checks the .agenthooks configuration file for the disable setting and only runs
formatting if hooks are not disabled. This command should be used in Claude Code
hooks instead of calling 'format' directly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration to check if hooks are disabled
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// If hooks are disabled, exit silently
		if cfg.Disable {
			return nil
		}

		// Delegate to format command with same arguments
		return formatCmd.RunE(cmd, args)
	},
}
