package cmd

import (
	"log"

	"github.com/samueltorres/idasenctl/internal/idasen"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [presetName]",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		presetName := args[0]
		deskName := deskFlag
		if deskName == "" {
			deskName = configManager.GetDefaultDesk()
		}
		desk, err := configManager.GetDesk(deskName)
		if err != nil {
			panic(err)
		}

		preset, ok := desk.Presets[presetName]
		if !ok {
			log.Fatal("Could not find that preset")
		}

		controller, err := idasen.NewController(desk.Address)
		if err != nil {
			log.Fatal(err)
		}

		err = controller.MoveTo(preset.Height)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	rootCmd.AddCommand(setCmd)
}
