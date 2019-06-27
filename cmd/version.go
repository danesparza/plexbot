package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildVersion = "Unknown"
var commitID string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version information",
	Run: func(cmd *cobra.Command, args []string) {
		//	Show the version number
		fmt.Printf("\nPlexbot version %s", buildVersion)

		//	Show the commitID if available:
		if commitID != "" {
			fmt.Printf(" (%s)", commitID[:7])
		}

		//	Trailing space and newline
		fmt.Println(" ")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
