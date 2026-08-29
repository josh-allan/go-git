package differ

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josh-allan/go-git/internal/ui"
)

type pager struct {
	viewport viewport.Model
	content  string
}

func newPager(content string, width, height int) pager {
	vp := viewport.New(width, height-1)
	vp.SetContent(content)

	return pager{
		viewport: vp,
		content:  content,
	}
}

func (p pager) Init() tea.Cmd {
	return nil
}

func (p pager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		}
	case tea.WindowSizeMsg:
		p.viewport.Width = msg.Width
		p.viewport.Height = msg.Height - 1
	}

	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p pager) View() string {
	return p.viewport.View() + "\n" + p.statusLine()
}

func (p pager) statusLine() string {
	pct := p.viewport.ScrollPercent() * 100

	info := fmt.Sprintf(" %3.0f%% ", pct)
	help := " q quit · ↑↓/j/k scroll · pgup/pgdn page "

	gap := p.viewport.Width - len(info) - len(help)
	if gap < 0 {
		gap = 0
	}

	style := lipgloss.NewStyle().
		Foreground(ui.StatusFg).
		Background(ui.StatusBg)

	return style.Render(help + strings.Repeat(" ", gap) + info)
}

func Page(content string, width, height int) error {
	lines := strings.Count(content, "\n")
	if lines <= height {
		fmt.Print(content)
		return nil
	}

	p := newPager(content, width, height)
	_, err := tea.NewProgram(p, tea.WithAltScreen()).Run()
	if err != nil {
		return fmt.Errorf("running pager: %w", err)
	}
	return nil
}
