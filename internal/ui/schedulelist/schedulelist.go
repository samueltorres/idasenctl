package schedulelist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samueltorres/idasenctl/internal/config"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type scheduleItem struct {
	name       string
	time       string
	deskName   string
	presetName string
	enabled    bool
	days       []int
}

func (i scheduleItem) Title() string {
	title := i.name
	if !i.enabled {
		title += " (disabled)"
	}
	return title
}

func (i scheduleItem) Description() string {
	daysStr := formatDays(i.days)
	return fmt.Sprintf("Time: %s | Desk: %s | Preset: %s | Days: %s", 
		i.time, i.deskName, i.presetName, daysStr)
}

func (i scheduleItem) FilterValue() string {
	return i.name
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

type scheduleListModel struct {
	list list.Model
}

func (m scheduleListModel) Init() tea.Cmd {
	return nil
}

func (m scheduleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m scheduleListModel) View() string {
	return docStyle.Render(m.list.View())
}

type ScheduleListProgram struct {
	teaProgram *tea.Program
}

func NewProgram(configManager *config.ConfigManager) *ScheduleListProgram {
	schedules := configManager.GetSchedules()

	var items []list.Item
	for _, schedule := range schedules {
		items = append(items, scheduleItem{
			name:       schedule.Name,
			time:       schedule.Time,
			deskName:   schedule.DeskName,
			presetName: schedule.PresetName,
			enabled:    schedule.Enabled,
			days:       schedule.Days,
		})
	}

	if len(items) == 0 {
		items = append(items, scheduleItem{
			name:       "No schedules configured",
			time:       "Run 'idasenctl schedule add' to add a schedule",
			deskName:   "",
			presetName: "",
			enabled:    false,
			days:       []int{},
		})
	}

	m := scheduleListModel{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Configured Schedules"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("q", "quit")),
		}
	}

	return &ScheduleListProgram{
		teaProgram: tea.NewProgram(m),
	}
}

func (p *ScheduleListProgram) Run() error {
	_, err := p.teaProgram.Run()
	return err
}