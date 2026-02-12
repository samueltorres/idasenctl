package daemon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samueltorres/idasenctl/internal/config"
	"github.com/samueltorres/idasenctl/internal/idasen"
	"github.com/samueltorres/idasenctl/internal/notification"
)

type Daemon struct {
	configManager *config.ConfigManager
	notifier      *notification.Notifier
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewDaemon(configManager *config.ConfigManager) *Daemon {
	ctx, cancel := context.WithCancel(context.Background())
	return &Daemon{
		configManager: configManager,
		notifier:      notification.NewNotifier(),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (d *Daemon) Start() error {
	log.Println("Starting idasenctl daemon...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go d.runScheduler()

	select {
	case <-sigChan:
		log.Println("Received termination signal, shutting down...")
		d.cancel()
		return nil
	case <-d.ctx.Done():
		return nil
	}
}

func (d *Daemon) runScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.checkSchedules()
		}
	}
}

func (d *Daemon) checkSchedules() {
	schedules := d.configManager.GetSchedules()
	now := time.Now()

	for _, schedule := range schedules {
		if !schedule.Enabled {
			continue
		}

		if !d.shouldRunToday(schedule, now) {
			continue
		}

		scheduledTime, err := d.parseScheduleTime(schedule.Time, now)
		if err != nil {
			log.Printf("Error parsing schedule time for %s: %v", schedule.Name, err)
			continue
		}

		timeDiff := scheduledTime.Sub(now)

		if timeDiff > 0 && timeDiff <= 10*time.Second {
			d.sendNotification(schedule)
		}

		if timeDiff > -30*time.Second && timeDiff <= 0 {
			d.executeSchedule(schedule)
		}
	}
}

func (d *Daemon) shouldRunToday(schedule config.Schedule, now time.Time) bool {
	weekday := int(now.Weekday())

	for _, day := range schedule.Days {
		if day == weekday {
			return true
		}
	}

	return false
}

func (d *Daemon) parseScheduleTime(timeStr string, now time.Time) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	year, month, day := now.Date()
	return time.Date(year, month, day, parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location()), nil
}

func (d *Daemon) sendNotification(schedule config.Schedule) {
	message := fmt.Sprintf("Your desk will move to preset '%s' in 10 seconds", schedule.PresetName)
	title := "Desk Movement Scheduled"

	err := d.notifier.SendNotification(title, message)
	if err != nil {
		log.Printf("Error sending notification for schedule %s: %v", schedule.Name, err)
	}
}

func (d *Daemon) executeSchedule(schedule config.Schedule) {
	log.Printf("Executing schedule: %s", schedule.Name)

	desk, err := d.configManager.GetDesk(schedule.DeskName)
	if err != nil {
		log.Printf("Error getting desk %s: %v", schedule.DeskName, err)
		return
	}

	preset, ok := desk.Presets[schedule.PresetName]
	if !ok {
		log.Printf("Error: preset %s not found for desk %s", schedule.PresetName, schedule.DeskName)
		return
	}

	controller, err := idasen.NewController(desk.Address)
	if err != nil {
		log.Printf("Error creating controller for desk %s: %v", schedule.DeskName, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err = controller.MoveTo(ctx, preset.Height, nil)
	if err != nil {
		log.Printf("Error moving desk %s to preset %s: %v", schedule.DeskName, schedule.PresetName, err)
		return
	}

	log.Printf("Successfully moved desk %s to preset %s (height: %.2f)", schedule.DeskName, schedule.PresetName, preset.Height)
}
