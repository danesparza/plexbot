package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonConfig bool
	yamlConfig bool
)

var yamlDefault = []byte(`
plex:
	tvpath: d:\tv
	errorpath: d:\errors

# Token replacement for preprocess, postprocess and postprocessall sections:
# {oldfilepath} - Replaced with full path of existing file in source directory
# {newfilepath} - Replaced with full path of moved file in destination directory

# To have a process run before the 'move' process, 
# uncomment this section and add it here:
# preprocess:
#  - somecommand.exe {filepath}

# To have a process run after the 'move' process, 
# uncomment this section and add it here:
postprocess:
	- qbittorrentremove.exe -file "{oldfilepath}

# To have a process run after all of the 'move' processes
# uncomment this section and add it here
# postprocessall:
#  - someothercommand.exe {filepath}"
`)

var jsonDefault = []byte(`{
  "plex": {
		"tvpath": "d:\\tv",
		"errorpath": "d:\\errors"
  },
	/*
	Token replacement for preprocess, postprocess, and postprocessall sections:
	{oldfilepath} - Replaced with full path of existing file in source directory
	{newfilepath} - Replaced with full path of moved file in destination directory
	*/
  "postprocess": [
    "qbittorrentremove.exe -file \"{oldfilepath}\""
  ]
}`)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Prints default plexbot configuration files",
	Long: `Use this to create a default configuration file for plexbot. 

Example:
plexbot defaults > plexbot.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultsCmd)

	defaultsCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	defaultsCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
