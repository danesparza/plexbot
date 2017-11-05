package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/danesparza/dlshow"
	"github.com/danesparza/plexbot/files"
	"github.com/danesparza/plexbot/plugin"
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

	//	Add the tv path to the list of tokens
	tokens["{tvpath}"] = viper.GetString("plex.tvpath")

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
	filesToMove := files.FindWithExtension([]string{".mp4", ".mkv", ".avi"}, sourceBaseDir)
	log.Printf("[INFO] Found %d file(s) to process", len(filesToMove))

	for _, file := range filesToMove {
		log.Printf("[INFO] - Found file %v...", file)
		tokens["{oldfilepath}"] = file

		//	Perform preprocessing
		if viper.InConfig("preprocess") {
			preProcessItems := viper.GetStringSlice("preprocess")
			for _, item := range preProcessItems {
				item = plugin.FormatTokenizedString(item, tokens)
				log.Printf("[INFO] -- Executing %v", item)
				plugin.ExecutePlugin(item)
			}
		}

		//	Parse show information:
		if showInfo, err := dlshow.GetEpisodeInfo(file); err == nil {

			//	If we can't parse the filename,
			//	just go to the next filename
			if showInfo.ParseType == 0 {
				continue
			}

			//	Add our showinfo tokens:
			tokens["{showname}"] = showInfo.ShowName
			tokens["{showseasonnumber}"] = strconv.Itoa(showInfo.SeasonNumber)
			tokens["{showepisodenumber}"] = strconv.Itoa(showInfo.EpisodeNumber)

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
			if err := files.Copy(file, newFile, os.ModePerm); err != nil {
				log.Printf("[ERROR] %v", err)
			}

			//	Perform 'postprocess each' items
			if viper.InConfig("postprocess") {
				postProcessItems := viper.GetStringSlice("postprocess")
				for _, item := range postProcessItems {
					item = plugin.FormatTokenizedString(item, tokens)
					log.Printf("[INFO] -- Executing %v", item)
					plugin.ExecutePlugin(item)
				}
			}

		}

	}

	//	Perform 'postprocess all' items
	if viper.InConfig("postprocessall") {
		postProcessItems := viper.GetStringSlice("postprocessall")
		for _, item := range postProcessItems {
			item = plugin.FormatTokenizedString(item, tokens)
			log.Printf("[INFO] -- Executing %v", item)
			plugin.ExecutePlugin(item)
		}
	}
}

func init() {
	RootCmd.AddCommand(moveCmd)
}
