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
preprocess:
  - somecommand.exe
postprocess:
  - someothercommand.exe
	- somethingelse.exe
`)

var jsonDefault = []byte(`{
  "server" : {
    "port": "3000"
  },
  "datastore" : {
    "type": "boltdb",
    "database": "config.db"
  }
}`)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("defaults called")
	},
}

func init() {
	RootCmd.AddCommand(defaultsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// defaultsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// defaultsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
