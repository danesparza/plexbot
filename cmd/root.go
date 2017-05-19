package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	plexTVDirectory string
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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/plexbot.yaml)")
	RootCmd.PersistentFlags().StringVar(&plexTVDirectory, "tvdir", "", "Base Plex TV directory")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("plexbot") // name of config file (without extension)
	viper.AddConfigPath(".")       // adding current directory as search path
	viper.AddConfigPath("$HOME")   // adding home directory as search path
	viper.AutomaticEnv()           // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
