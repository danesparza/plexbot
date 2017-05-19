package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	sourceDirectory string
	moveNoFile      = `You didn't pass anything to move.  

Move requires a given directory to move from

Example:
plexbot move c:\source\dir`
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
		fmt.Println(moveNoFile)
		return
	}
	log.Printf("[INFO] Looking for files in: %v", args[0])

	//	See if the source directory exists
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		log.Printf("The directory doesn't exist: %v", args[0])
		return
	}

	//	If it does, see what movie files it contains:
	filesToMove := filesWithExtension([]string{"mp4", "mkv", "avi"}, "")

	for _, file := range filesToMove {
		log.Printf("[INFO] Moving file %v...", file)
	}
}

func init() {
	RootCmd.AddCommand(moveCmd)
}

func filesWithExtension(exts []string, baseDirectory string) []string {
	//	Sanity check that the source directory seems to exist
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		log.Panic("The directory doesn't exist " + baseDirectory)
	}

	var files []string
	filepath.Walk(baseDirectory, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			//	For each of the extensions passed...
			for _, ext := range exts {

				//	If the file seems to match...
				r, err := regexp.MatchString(ext, f.Name())
				if err == nil && r {
					//	Add it to the pile of file results
					files = append(files, f.Name())
				}
			}
		}
		return nil
	})

	//	Return the list of files found
	return files
}
