package presetlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samueltorres/idasenctl/internal/config"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type presetItem struct {
	name   string
	height float32
}

func (i presetItem) Title() string {
	return i.name
}

func (i presetItem) Description() string {
	return fmt.Sprintf("Height: %.2f m", i.height)
}

func (i presetItem) FilterValue() string {
	return i.name
}

type presetListModel struct {
	list list.Model
}

func (m presetListModel) Init() tea.Cmd {
	return nil
}

func (m presetListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m presetListModel) View() string {
	return docStyle.Render(m.list.View())
}

type PresetListProgram struct {
	teaProgram *tea.Program
}

func NewProgram(configManager *config.ConfigManager, deskName string) *PresetListProgram {
	desk, err := configManager.GetDesk(deskName)

	var items []list.Item
	var title string

	if err != nil {
		items = append(items, presetItem{
			name:   "Error loading desk",
			height: 0,
		})
		title = "Error"
	} else {
		title = fmt.Sprintf("Presets for %s", deskName)

		for _, preset := range desk.Presets {
			items = append(items, presetItem{
				name:   preset.Name,
				height: preset.Height,
			})
		}

		if len(items) == 0 {
			items = append(items, presetItem{
				name:   "No presets configured",
				height: 0,
			})
		}
	}

	m := presetListModel{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = title
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("q", "quit")),
		}
	}

	return &PresetListProgram{
		teaProgram: tea.NewProgram(m),
	}
}

func (p *PresetListProgram) Run() error {
	_, err := p.teaProgram.Run()
	return err
}
