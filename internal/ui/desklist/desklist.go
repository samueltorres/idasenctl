package desklist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samueltorres/idasenctl/internal/config"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type deskItem struct {
	name         string
	address      string
	presetCount  int
	isDefault    bool
}

func (i deskItem) Title() string {
	title := i.name
	if i.isDefault {
		title += " (default)"
	}
	return title
}

func (i deskItem) Description() string {
	return fmt.Sprintf("Address: %s | Presets: %d", i.address, i.presetCount)
}

func (i deskItem) FilterValue() string {
	return i.name
}

type deskListModel struct {
	list list.Model
}

func (m deskListModel) Init() tea.Cmd {
	return nil
}

func (m deskListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m deskListModel) View() string {
	return docStyle.Render(m.list.View())
}

type DeskListProgram struct {
	teaProgram *tea.Program
}

func NewProgram(configManager *config.ConfigManager) *DeskListProgram {
	desks := configManager.GetAllDesks()
	defaultDesk := configManager.GetDefaultDesk()

	var items []list.Item
	for _, desk := range desks {
		items = append(items, deskItem{
			name:        desk.Name,
			address:     desk.Address,
			presetCount: len(desk.Presets),
			isDefault:   desk.Name == defaultDesk,
		})
	}

	if len(items) == 0 {
		items = append(items, deskItem{
			name:        "No desks configured",
			address:     "Run 'idasenctl desk add' to add a desk",
			presetCount: 0,
			isDefault:   false,
		})
	}

	m := deskListModel{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Configured Desks"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("q", "quit")),
		}
	}

	return &DeskListProgram{
		teaProgram: tea.NewProgram(m),
	}
}

func (p *DeskListProgram) Run() error {
	_, err := p.teaProgram.Run()
	return err
}