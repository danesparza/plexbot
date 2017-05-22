package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Emit our plex tv directory
	log.Printf("[INFO] Plex TV library path: %s\n", viper.GetString("plex.tvpath"))

	//	Make sure we were called with a directory
	if len(args) < 1 {
		fmt.Println(moveNoFile)
		return
	}
	log.Printf("[INFO] Looking for files in: %v...", args[0])

	//	See if the source directory exists
	sourceBaseDir := args[0]
	if _, err := os.Stat(sourceBaseDir); os.IsNotExist(err) {
		log.Printf("[ERROR] The directory doesn't exist: %v", sourceBaseDir)
		return
	}

	//	See if the destination directory exists
	destBaseDir := viper.GetString("plex.tvpath")
	if _, err := os.Stat(destBaseDir); err != nil {
		log.Printf("[ERROR] The plex TV directory doesn't exist: %v", destBaseDir)
		return
	}

	//	If it does, see what movie files it contains:
	filesToMove := filesWithExtension([]string{".mp4", ".mkv", ".avi"}, sourceBaseDir)
	log.Printf("[INFO] Found %d file(s) to process", len(filesToMove))

	for _, file := range filesToMove {
		log.Printf("[INFO] - Found file %v...", file)

		//	Create our map of replacement tokens and add
		//	our first token:
		tokens := make(map[string]string)
		tokens["{oldfilepath}"] = file

		//	Perform preprocessing like this:
		//	http://stackoverflow.com/a/20438245/19020
		if viper.InConfig("preprocess") {
			preProcessItems := viper.GetStringSlice("preprocess")
			for _, item := range preProcessItems {
				item = formatTokenizedString(item, tokens)
				log.Printf("[INFO] -- Executing %v", item)

				// splitting head => g++ parts => rest of the command
				parts := strings.Fields(item)
				head := parts[0]
				parts = parts[1:len(parts)]

				exec.Command(head, parts...).Output()
			}
		}

		//	Parse show information:
		if showInfo, err := dlshow.GetEpisodeInfo(file); err == nil {

			//	Format the new filepath:
			seasonDir := fmt.Sprintf("Season %d", showInfo.SeasonNumber)
			newPath := filepath.Join(destBaseDir, showInfo.ShowName, seasonDir)
			newFileName := fmt.Sprintf("s%de%02d%v", showInfo.SeasonNumber, showInfo.EpisodeNumber, filepath.Ext(file))
			newFile := filepath.Join(destBaseDir, showInfo.ShowName, seasonDir, newFileName)

			//	Add to our replacement tokens:
			tokens["{newfilepath}"] = newFile

			//	Make sure the new path exists:
			os.MkdirAll(newPath, os.ModePerm)

			//	Move the file
			log.Printf("[INFO] -- Moving to %v", newFile)
			if err := CopyFile(file, newFile, os.ModePerm); err != nil {
				log.Printf("[ERROR] %v", err)
			}

			//	Perform postprocessing like this:
			//	http://stackoverflow.com/a/20438245/19020
			if viper.InConfig("postprocess") {
				postProcessItems := viper.GetStringSlice("postprocess")
				for _, item := range postProcessItems {
					item = formatTokenizedString(item, tokens)
					log.Printf("[INFO] -- Executing %v", item)

					// splitting head => g++ parts => rest of the command
					parts := strings.Fields(item)
					head := parts[0]
					parts = parts[1:len(parts)]
					exec.Command(head, parts...).Output()
				}
			}

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

// formatTokenizedString will format a string containing tokens by replacing
// the tokens with their actual values and returning the new string
func formatTokenizedString(originalString string, tokens map[string]string) string {
	retval := originalString

	//	Replace each token with its value:
	for token, value := range tokens {
		retval = strings.Replace(retval, token, value, -1)
	}

	return retval
}

// CopyFile copies the contents from src to dst using io.Copy.
// If dst does not exist, CopyFile creates it with permissions perm;
// otherwise CopyFile truncates it before writing.
func CopyFile(src, dst string, perm os.FileMode) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	return
}
