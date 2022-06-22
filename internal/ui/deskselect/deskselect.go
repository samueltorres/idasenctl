package deskselect

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samueltorres/idasenctl/internal/idasen"
)

var scanUIDocStyle = lipgloss.NewStyle().Margin(1, 2)

type scanItem struct {
	deskName, deskAddress string
}

func (i scanItem) Title() string       { return i.deskName }
func (i scanItem) Description() string { return i.deskAddress }
func (i scanItem) FilterValue() string { return i.deskName }

type scanModel struct {
	List         list.Model
	SelectedItem *scanItem
}

func (m *scanModel) Init() tea.Cmd {
	return nil
}

func (m *scanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", " ":
			item := m.List.SelectedItem().(scanItem)
			m.SelectedItem = &item
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := scanUIDocStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m scanModel) View() string {
	return scanUIDocStyle.Render(m.List.View())
}

type DeskSelectProgram struct {
	teaProgram *tea.Program
	scanModel  *scanModel
	deskScans  chan idasen.DeviceInfo
}

func NewProgram(deskScans chan idasen.DeviceInfo) *DeskSelectProgram {
	items := []list.Item{}
	m := &scanModel{List: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Scanning desks"
	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{key.NewBinding(key.WithHelp("enter", "select"))}
	}

	return &DeskSelectProgram{
		teaProgram: tea.NewProgram(m, tea.WithAltScreen()),
		scanModel:  m,
		deskScans:  deskScans,
	}
}

func (p *DeskSelectProgram) Run(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ds := <-p.deskScans:
				cmd := p.scanModel.List.InsertItem(0, scanItem{deskName: ds.Name, deskAddress: ds.Address})
				p.teaProgram.Send(cmd)
			}
		}
	}()

	_, err := p.teaProgram.StartReturningModel()
	if err != nil {
		fmt.Println("Error running desk selection:", err)
		os.Exit(1)
	}

	return nil
}

func (p *DeskSelectProgram) GetSelectedDesk() *idasen.DeviceInfo {
	if p.scanModel.SelectedItem == nil {
		return nil
	}
	return &idasen.DeviceInfo{Name: p.scanModel.SelectedItem.deskName, Address: p.scanModel.SelectedItem.deskAddress}
}
