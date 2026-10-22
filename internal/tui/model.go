package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/sumit-poudel/gine/internal/theme"
)

type Model struct {
	window struct {
		width  int
		height int
	}
	prompt textinput.Model
	help   helpDialog
	styles theme.Theme
}

func New() Model {
	prompt := textinput.New()
	prompt.Placeholder = "Search Mime Type"
	prompt.SetVirtualCursor(false)
	prompt.Focus()
	prompt.CharLimit = 156

	return Model{
		prompt: prompt,
		styles: theme.New(true),
		help: helpDialog{
			table: helpTable{
				bindings: []binding{
					{"?", "Toggle this help menu"},
					{"q", "Quit the program"},
				},
				padding: 2,
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, tea.RequestBackgroundColor, textinput.Blink)
}
