package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	hash    string
	taglist string
	tags    []string

	// ProblemWithConfigFile indicates whether or not there was a problem
	// loading the config
	ProblemWithConfigFile bool

	//	Create our map of replacement tokens
	tokens = make(map[string]string)
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "plexbot",
	Short: "Simple app to help organize tv shows and movies into the Plex naming format",
	Long: `Plexbot will determine the series and episode information using the passed filename
and move it into the correct naming format for Plex.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is plexbot.yaml)")
	RootCmd.PersistentFlags().StringVar(&hash, "hash", "", "Torrent hash used to identify the torrent")
	RootCmd.PersistentFlags().StringVar(&taglist, "tags", "", "List of tags (comma-seperated)")

	//	Parse the string list of tags to a slice of tags:
	tags = strings.Split(taglist, ",")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//	Set our defaults
	viper.SetDefault("plex.tvpath", "d:\\tv")
	viper.SetDefault("plex.errorpath", "d:\\errors")
	viper.SetDefault("preprocess.command", []string{})
	viper.SetDefault("postprocess.command", []string{})

	viper.SetConfigName("plexbot") // name of config file (without extension)
	viper.AddConfigPath(".")       // adding current directory as search path
	viper.AddConfigPath("$HOME")   // adding home directory as search path
	viper.AutomaticEnv()           // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	//	Set the hash token
	tokens["{hash}"] = hash

	// If a config file is found, read it in
	// otherwise, make note that there was a problem
	if err := viper.ReadInConfig(); err != nil {
		ProblemWithConfigFile = true
	}
}

// properTitle returns the proper title case for a given string
func properTitle(input string) string {
	words := strings.Fields(input)

	for index, word := range words {
		words[index] = strings.Title(word)
	}
	return strings.Join(words, " ")
}
