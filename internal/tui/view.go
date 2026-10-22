package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	content := m.styles.Title().Render("GINE") + "\n\n"
	content += m.styles.Input().Render(m.prompt.View())

	main := m.styles.Screen(m.window.width, m.window.height).Render(content)

	v := tea.NewView(m.compose(main))
	v.WindowTitle = "GINE"
	v.AltScreen = true
	return v
}

func (m Model) compose(main string) string {
	layers := []*lipgloss.Layer{lipgloss.NewLayer(main)}

	if m.window.height > 0 {
		footer := m.styles.Footer(m.window.width).Render("Press q to quit.  ? for help")
		layers = append(layers, lipgloss.NewLayer(footer).Y(m.window.height-1).Z(1))
	}

	if m.help.show {
		if dialog := m.help.layer(m.styles, m.window.width, m.window.height); dialog != nil {
			layers = append(layers, dialog)
		}
	}

	return lipgloss.NewCompositor(layers...).Render()
}
