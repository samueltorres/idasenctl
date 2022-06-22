package cmd

import (
	"context"
	"log"

	"github.com/samueltorres/idasenctl/internal/idasen"
	"github.com/samueltorres/idasenctl/internal/ui/deskmove"
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

		currentHeight, err := controller.GetCurrentHeight()
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		updates := make(chan float32)
		deskMoveProgram := deskmove.NewProgram(preset.Height, currentHeight, updates)

		go func() {
			err = controller.MoveTo(ctx, preset.Height, updates)
			if err != nil {
				log.Fatal(err)
			}
		}()

		err = deskMoveProgram.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	rootCmd.AddCommand(setCmd)
}
