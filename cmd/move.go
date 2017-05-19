package cmd

import (
	"log"

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
plexbot move c:\source\dir

Source directory contains: 'Once.Upon.a.Time.S03E01.720p.HDTV.X264-DIMENSION.mkv'
Plex base TV directory: 'D:\TV'

Then the file will get moved and renamed to:
D:\TV\Once Upon a Time\Season 3\s3e01.mkv`,
	Run: parseAndMove,
}

func parseAndMove(cmd *cobra.Command, args []string) {
	//	Make sure we were called with a directory
	if len(args) < 1 {
		log.Println("Move requires a given directory to move from")
		return
	}
	log.Printf("[INFO] Looking for files in: %v", args[0])

	//	See if the source directory exists
	//	If it does, see if the directory contains movie files

}

func init() {
	RootCmd.AddCommand(moveCmd)
}
