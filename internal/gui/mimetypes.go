package gui

import (
	"charm.land/lipgloss/v2"

	"github.com/sumit-poudel/gine/internal/gui/style"
)

type rowStyle func(width int) lipgloss.Style

type TypeList struct {
	types     []string
	selected  int
	offset    int
	width     int
	height    int
	rowHeight int
}

func newTypeList(types []string) TypeList {
	return TypeList{types: types}
}

func (list *TypeList) setItems(items []string) {
	list.types = items
	list.selected = -1
	list.offset = 0
}

func (list TypeList) Empty() bool {
	return len(list.types) == 0
}

func (list *TypeList) Resize(width, height int, styles style.Styles) {
	list.width = max(0, width)
	list.height = max(0, height)
	list.rowHeight = 1
	if list.width > 0 {
		sample := styles.MimeRow(list.width).Render(" ")
		list.rowHeight = max(1, lipgloss.Height(sample))
	}
	list.slide()
}

func (list *TypeList) Previous() {
	if list.selected > 0 {
		list.selected--
	}
	list.slide()
}

func (list TypeList) Selected() string {
	if list.selected < 0 || list.selected >= len(list.types) {
		return ""
	}
	return list.types[list.selected]
}

func (list *TypeList) Next() {
	if list.selected < len(list.types)-1 {
		list.selected++
	}
	list.slide()
}

func (list TypeList) View(styles style.Styles, focused bool, row rowStyle) string {
	capacity := list.capacity()
	if capacity == 0 || list.width == 0 || len(list.types) == 0 {
		return ""
	}

	end := min(len(list.types), list.offset+capacity)
	rows := make([]string, 0, end-list.offset)
	for index := list.offset; index < end; index++ {
		line := row(list.width)
		if focused && index == list.selected {
			line = styles.SelectedTypeRow(list.width)
		}
		rows = append(rows, line.Render(list.types[index]))
	}

	return lipgloss.NewStyle().
		MaxWidth(list.width).
		MaxHeight(list.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, rows...))
}

func (list TypeList) capacity() int {
	if list.rowHeight < 1 || list.height < 1 {
		return 0
	}
	return list.height / list.rowHeight
}

func (list *TypeList) slide() {
	if list.selected < 0 {
		list.offset = 0
		return
	}
	capacity := list.capacity()
	if capacity < 1 {
		return
	}
	if list.selected < list.offset {
		list.offset = list.selected
	}
	if list.selected >= list.offset+capacity {
		list.offset = list.selected - capacity + 1
	}
	if list.offset < 0 {
		list.offset = 0
	}
}
