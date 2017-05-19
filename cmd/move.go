package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	sourceDirectory string
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Moves and renames files in a given directory",
	Long: `This command moves and renames files into the Plex naming format

For example:
Source directory contains: 'Once.Upon.a.Time.S03E01.720p.HDTV.X264-DIMENSION.mkv'
Plex base TV directory: 'D:\TV'

Then the file will get moved and renamed to:
D:\TV\Once Upon a Time\Season 3\s3e01.mkv`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("move called")
	},
}

func init() {
	RootCmd.AddCommand(moveCmd)

	// Flags for this command
	RootCmd.PersistentFlags().StringVar(&sourceDirectory, "sourceDirectory", "", "The source directory to look for files")

}
