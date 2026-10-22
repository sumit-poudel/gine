package tui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/sumit-poudel/gine/internal/theme"
)

type binding struct {
	key  string
	desc string
}

type helpDialog struct {
	show  bool
	table helpTable
}

type helpTable struct {
	bindings []binding
	padding  int
}

func (h helpTable) render(t theme.Theme) string {
	tbl := table.New().
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderRow(false).
		StyleFunc(func(_, col int) lipgloss.Style {
			if col == 0 {
				return t.DialogKey(h.padding)
			}
			return t.DialogCell()
		})

	for _, b := range h.bindings {
		tbl.Row(b.key, b.desc)
	}

	return t.Dialog().Render(tbl.Render())
}

func (h helpDialog) layer(t theme.Theme, width, height int) *lipgloss.Layer {
	dialog := h.table.render(t)
	w, dialogHeight := lipgloss.Width(dialog), lipgloss.Height(dialog)
	if w > width || dialogHeight > height {
		return nil
	}

	return lipgloss.NewLayer(dialog).
		X((width - w) / 2).
		Y((height - dialogHeight) / 2).
		Z(10)
}
