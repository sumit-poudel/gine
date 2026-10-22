package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/sumit-poudel/gine/internal/theme"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.width = msg.Width
		m.window.height = msg.Height
		m.prompt.SetWidth(msg.Width - 10)
	case tea.BackgroundColorMsg:
		m.styles = theme.New(msg.IsDark())
	case tea.PasteMsg:
		m.prompt, cmd = m.prompt.Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			m.help.show = !m.help.show
		default:
			m.prompt, cmd = m.prompt.Update(msg)
		}
	}

	return m, cmd
}
