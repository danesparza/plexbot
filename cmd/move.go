package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/danesparza/dlshow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	//	Emit our plex tv directory
	log.Printf("[INFO] Plex TV library path: %s\n", viper.GetString("plex.tvpath"))

	//	Make sure we were called with a directory
	if len(args) < 1 {
		fmt.Println(moveNoFile)
		return
	}
	log.Printf("[INFO] Looking for files in: %v...", args[0])

	//	See if the source directory exists
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		log.Printf("[ERROR] The directory doesn't exist: %v", args[0])
		return
	}

	//	If it does, see what movie files it contains:
	filesToMove := filesWithExtension([]string{".mp4", ".mkv", ".avi"}, args[0])

	for _, file := range filesToMove {
		log.Printf("[INFO] - Found file %v...", file)

		//	Parse show information:
		if showInfo, err := dlshow.GetEpisodeInfo(file); err == nil {
			newFile := fmt.Sprintf("s%de%02d%v", showInfo.SeasonNumber, showInfo.EpisodeNumber, filepath.Ext(file))
			newFile = filepath.Join(showInfo.ShowName, newFile)
			log.Printf("[INFO] -- Moving to %v", newFile)
		}

	}
}

func init() {
	RootCmd.AddCommand(moveCmd)
}

// filesWithExtension returns a list of files that contain the given
// extensions in the requested baseDirectory
func filesWithExtension(exts []string, baseDirectory string) []string {
	//	Sanity check that the source directory seems to exist
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		log.Panic("The directory doesn't exist " + baseDirectory)
	}

	var files []string
	filepath.Walk(baseDirectory, func(path string, f os.FileInfo, _ error) error {
		//	If it's a file...
		if !f.IsDir() {
			//	See if its extension matches one we're looking for...
			if contains(exts, filepath.Ext(f.Name())) {
				//	If it does, Add it to the pile of file results
				files = append(files, path)
			}
		}
		return nil
	})

	//	Return the list of files found
	return files
}

// contains returns true if the target slice contains the item 'e'
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
