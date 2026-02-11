package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/samueltorres/idasenctl/internal/config"
	"github.com/samueltorres/idasenctl/internal/ui/schedulelist"
	"github.com/spf13/cobra"
)

var (
	scheduleTime     string
	scheduleDeskName string
	schedulePreset   string
	scheduleEnabled  bool
	scheduleDays     []string
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage desk movement schedules",
	Long:  `Create, list, and manage automated desk movement schedules.`,
}

var scheduleAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new schedule",
	Long:  `Add a new schedule for automated desk movements.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scheduleName := args[0]
		
		deskName := scheduleDeskName
		if deskName == "" {
			deskName = configManager.GetDefaultDesk()
		}
		
		days, err := parseDays(scheduleDays)
		if err != nil {
			log.Fatal(err)
		}
		
		schedule := config.Schedule{
			Name:       scheduleName,
			Time:       scheduleTime,
			DeskName:   deskName,
			PresetName: schedulePreset,
			Enabled:    scheduleEnabled,
			Days:       days,
		}
		
		err = configManager.AddSchedule(schedule)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("Schedule '%s' added successfully\n", scheduleName)
	},
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all schedules",
	Long:  `List all configured schedules.`,
	Run: func(cmd *cobra.Command, args []string) {
		program := schedulelist.NewProgram(configManager)
		err := program.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var scheduleRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a schedule",
	Long:  `Remove a schedule by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scheduleName := args[0]
		
		err := configManager.RemoveSchedule(scheduleName)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("Schedule '%s' removed successfully\n", scheduleName)
	},
}

func parseDays(dayStrings []string) ([]int, error) {
	dayMap := map[string]int{
		"sunday":    0,
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sun":       0,
		"mon":       1,
		"tue":       2,
		"wed":       3,
		"thu":       4,
		"fri":       5,
		"sat":       6,
	}
	
	var days []int
	for _, dayStr := range dayStrings {
		dayStr = strings.ToLower(dayStr)
		
		if dayNum, err := strconv.Atoi(dayStr); err == nil {
			if dayNum >= 0 && dayNum <= 6 {
				days = append(days, dayNum)
				continue
			}
		}
		
		if dayNum, ok := dayMap[dayStr]; ok {
			days = append(days, dayNum)
		} else {
			return nil, fmt.Errorf("invalid day: %s", dayStr)
		}
	}
	
	return days, nil
}

func formatDays(days []int) string {
	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var dayStrings []string
	
	for _, day := range days {
		if day >= 0 && day <= 6 {
			dayStrings = append(dayStrings, dayNames[day])
		}
	}
	
	return strings.Join(dayStrings, ",")
}

func init() {
	scheduleAddCmd.Flags().StringVarP(&scheduleTime, "time", "t", "", "Time in HH:MM format (required)")
	scheduleAddCmd.Flags().StringVarP(&scheduleDeskName, "desk", "d", "", "Desk name (defaults to default desk)")
	scheduleAddCmd.Flags().StringVarP(&schedulePreset, "preset", "p", "", "Preset name (required)")
	scheduleAddCmd.Flags().BoolVarP(&scheduleEnabled, "enabled", "e", true, "Enable the schedule")
	scheduleAddCmd.Flags().StringSliceVar(&scheduleDays, "days", []string{}, "Days of the week (e.g., monday,tuesday or 1,2 or mon,tue)")
	
	scheduleAddCmd.MarkFlagRequired("time")
	scheduleAddCmd.MarkFlagRequired("preset")
	scheduleAddCmd.MarkFlagRequired("days")
	
	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	scheduleCmd.AddCommand(scheduleRemoveCmd)
	
	rootCmd.AddCommand(scheduleCmd)
}