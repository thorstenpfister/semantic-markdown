package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the semantic-markdown version (set during build)
	Version = "0.1.0-dev"
	// GitCommit is the git commit hash (set during build)
	GitCommit = "unknown"
	// BuildDate is the build date (set during build)
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information for semantic-md",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("semantic-md version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
