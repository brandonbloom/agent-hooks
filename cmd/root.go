package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent-hooks",
	Short: "Zero-config convenience commands for humans and AI agents",
	Long: `agent-hooks provides zero-config convenience commands that can be used
by both humans and AI agents. It respects existing project configurations
when available but requires no special setup to be useful.`,
	Version:       getVersionString(),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(aboutCmd)
	rootCmd.AddCommand(detectCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(whichVcsCmd)
}
