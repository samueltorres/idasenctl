package cmd

import (
	"log"

	"github.com/samueltorres/idasenctl/internal/config"
	"github.com/samueltorres/idasenctl/internal/idasen"

	"github.com/spf13/cobra"
)

var (
	deskAddressFlag string
)

var deskCmd = &cobra.Command{
	Use:   "desk",
	Short: "manages idasen desks",
	Long:  ``,
}

var deskScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans idasen desks",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		controller, err := idasen.NewScanner()
		if err != nil {
			log.Fatal(err)
		}
		controller.Scan()
	},
}

var deskAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "adds an idasen desk",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deskName := args[0]
		err := configManager.SetDesk(config.Desk{
			Name:    args[0],
			Address: deskAddressFlag,
		})
		if err != nil {
			log.Fatal(err)
		}
		err = configManager.SetDefaultDesk(deskName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deskDefaultCmd = &cobra.Command{
	Use:   "default [name]",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deskName := args[0]
		err := configManager.SetDefaultDesk(deskName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	deskAddCmd.Flags().StringVarP(&deskAddressFlag, "address", "a", "", "The address uuid of the desk you are adding")
	deskAddCmd.MarkFlagRequired("address")

	deskCmd.AddCommand(deskAddCmd)
	deskCmd.AddCommand(deskDefaultCmd)
	deskCmd.AddCommand(deskScanCmd)

	rootCmd.AddCommand(deskCmd)
}
