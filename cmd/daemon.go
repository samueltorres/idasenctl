package cmd

import (
	"log"

	"github.com/samueltorres/idasenctl/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run idasenctl as a daemon with scheduled desk movements",
	Long:  `Run idasenctl as a background daemon that will automatically move your desk based on configured schedules.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := daemon.NewDaemon(configManager)
		err := d.Start()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
