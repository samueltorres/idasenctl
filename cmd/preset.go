package cmd

import (
	"log"

	"github.com/samueltorres/idasenctl/internal/idasen"
	"github.com/samueltorres/idasenctl/internal/ui/presetlist"
	"github.com/spf13/cobra"
)

var (
	deskFlag          string
	deskPresetHeight  float32
	deskPresetCurrent bool
)

// presetCmd represents the preset command
var presetCmd = &cobra.Command{
	Use:   "preset",
	Short: "manages desk presets",
}

var presetAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "adds desk presets",
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

		height := deskPresetHeight
		if deskPresetCurrent {
			controller, err := idasen.NewController(desk.Address)
			if err != nil {
				log.Fatal(err)
			}

			height, err = controller.GetCurrentHeight()
			if err != nil {
				log.Fatal(err)
			}
		}

		err = configManager.SetDeskPreset(deskName, presetName, height)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var presetDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "delete desk presets",
	Run: func(cmd *cobra.Command, args []string) {
		presetName := args[0]
		deskName := deskFlag
		if deskName == "" {
			deskName = configManager.GetDefaultDesk()
		}

		err := configManager.DeleteDeskPreset(deskName, presetName)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var presetListCmd = &cobra.Command{
	Use:   "list",
	Short: "list desk presets",
	Run: func(cmd *cobra.Command, args []string) {
		deskName := deskFlag
		if deskName == "" {
			deskName = configManager.GetDefaultDesk()
		}

		program := presetlist.NewProgram(configManager, deskName)
		err := program.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	presetAddCmd.Flags().StringVarP(&deskFlag, "desk", "d", "", "The name of the desk")
	presetAddCmd.Flags().Float32VarP(&deskPresetHeight, "height", "", 0, "The height of the desk on the preset")
	presetAddCmd.Flags().BoolVarP(&deskPresetCurrent, "current", "c", false, "The height of the desk on the preset")

	presetListCmd.Flags().StringVarP(&deskFlag, "desk", "d", "", "The name of the desk")
	presetDeleteCmd.Flags().StringVarP(&deskFlag, "desk", "d", "", "The name of the desk")

	presetCmd.AddCommand(presetAddCmd)
	presetCmd.AddCommand(presetDeleteCmd)
	presetCmd.AddCommand(presetListCmd)
	rootCmd.AddCommand(presetCmd)
}
