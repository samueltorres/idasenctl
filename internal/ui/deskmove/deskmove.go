package deskmove

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type progressMsg float64

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type progressModel struct {
	progress      progress.Model
	desiredHeight float32
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progressMsg:
		var cmds []tea.Cmd
		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequentially(finalPause(), tea.Quit))
		}
		cmds = append(cmds, m.progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (e progressModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + fmt.Sprintf("Setting height to %.2f m \n\n", e.desiredHeight) +
		pad + e.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

type DeskMoveProgram struct {
	initialHeight float32
	desiredHeight float32
	updates       chan float32
	progressModel *progressModel
	teaProgram    *tea.Program
}

func NewProgram(desiredHeight float32, initialHeight float32, updates chan float32) *DeskMoveProgram {
	m := &progressModel{
		progress:      progress.New(progress.WithDefaultGradient()),
		desiredHeight: desiredHeight,
	}

	return &DeskMoveProgram{
		initialHeight: initialHeight,
		desiredHeight: desiredHeight,
		updates:       updates,
		progressModel: m,
		teaProgram:    tea.NewProgram(m),
	}
}

func (p *DeskMoveProgram) Run(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case currentHeight := <-p.updates:
				totalDistance := math.Abs(float64(p.desiredHeight - p.initialHeight))
				distanceCovered := math.Abs(float64(p.initialHeight - currentHeight))

				progress := progressMsg(distanceCovered / totalDistance)
				if math.Abs(float64(p.desiredHeight-currentHeight)) < 0.005 {
					progress = progressMsg(1)
				}
				p.teaProgram.Send(progress)
			}
		}
	}()

	err := p.teaProgram.Start()
	if err != nil {
		return err
	}
	return nil
}
