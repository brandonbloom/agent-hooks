package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getVersionString())
	},
}

func getVersionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	version := info.Main.Version
	if version == "(devel)" {
		version = "dev"
	}

	var revision, modified string
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value
		case "vcs.modified":
			if kv.Value == "true" {
				modified = "-dirty"
			}
		}
	}

	if revision != "" {
		if len(revision) > 7 {
			revision = revision[:7]
		}
		return fmt.Sprintf("%s (%s%s)", version, revision, modified)
	}

	return version
}
