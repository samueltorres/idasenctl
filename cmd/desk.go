package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/samueltorres/idasenctl/internal/config"
	"github.com/samueltorres/idasenctl/internal/idasen"
	"github.com/samueltorres/idasenctl/internal/ui/desklist"
	"github.com/samueltorres/idasenctl/internal/ui/deskselect"

	"github.com/spf13/cobra"
)

var deskCmd = &cobra.Command{
	Use:   "desk",
	Short: "manages idasen desks",
	Long:  ``,
}

var deskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		controller, err := idasen.NewScanner()
		if err != nil {
			log.Fatal(err)
		}
		deskScans := make(chan idasen.DeviceInfo)
		go controller.Scan(ctx, deskScans)

		deskSelectProgram := deskselect.NewProgram(deskScans)
		err = deskSelectProgram.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}

		selectedDesk := deskSelectProgram.GetSelectedDesk()
		if selectedDesk == nil {
			return
		}

		err = configManager.SetDesk(config.Desk{
			Name:    selectedDesk.Name,
			Address: selectedDesk.Address,
		})
		if err != nil {
			log.Fatal(err)
		}
		err = configManager.SetDefaultDesk(selectedDesk.Name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Selected desk:", selectedDesk.Name)
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

var deskListCmd = &cobra.Command{
	Use:   "list",
	Short: "list configured desks",
	Long:  `List all configured desks with their addresses and preset counts.`,
	Run: func(cmd *cobra.Command, args []string) {
		program := desklist.NewProgram(configManager)
		err := program.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	deskCmd.AddCommand(deskAddCmd)
	deskCmd.AddCommand(deskDefaultCmd)
	deskCmd.AddCommand(deskListCmd)

	rootCmd.AddCommand(deskCmd)
}
