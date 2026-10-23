package gui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/sumit-poudel/gine/internal/gui/style"
)

const shortcutColumnGap = 2

type Help struct {
	visible bool
}

func newHelp() Help {
	return Help{}
}

func (help *Help) Toggle() {
	help.visible = !help.visible
}

func (help Help) Visible() bool {
	return help.visible
}

func (help Help) render(styles style.Styles, vimMode bool) string {
	grid := table.New().
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderRow(false).
		StyleFunc(func(_, column int) lipgloss.Style {
			if column == 0 {
				return styles.Shortcut(shortcutColumnGap)
			}
			return styles.Description()
		})

	for _, binding := range contextFor(vimMode).shortcuts {
		grid.Row(binding.label, binding.description)
	}
	return styles.HelpFrame().Render(grid.Render())
}

func (model Model) helpOverlay() string {
	dialog := model.help.render(model.styles, model.vimMode)
	if lipgloss.Width(dialog) > model.window.width || lipgloss.Height(dialog) > model.window.height {
		return ""
	}
	return dialog
}
