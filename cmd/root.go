package cmd

import (
	"os"
	"path"

	"github.com/samueltorres/idasenctl/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	configManager *config.ConfigManager
	rootCmd       = &cobra.Command{
		Use:   "idasenctl",
		Short: "A brief description of your application",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.idasenctl.yaml)")
}

func initConfig() {
	if cfgFile == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		cfgFile = path.Join(home, ".idasenctl.yaml")
	}

	cm, err := config.NewConfigManager(cfgFile)
	if err != nil {
		panic(err)
	}
	configManager = cm
}
