// Package style builds lipgloss styles from a color theme.
package style

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/sumit-poudel/gine/internal/theme"
)

type Styles struct {
	theme theme.Theme
}

func New(theme theme.Theme) Styles {
	return Styles{theme: theme}
}

func (styles Styles) Accent() color.Color { return styles.theme.Accent() }

func (styles Styles) Muted() color.Color { return styles.theme.Muted() }

func (styles Styles) Window(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Left)
}

func (styles Styles) Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(styles.theme.Accent()).
		Bold(true)
}

func (styles Styles) StatusLine(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width)
}

func (styles Styles) SearchBox() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.theme.Accent())
}

func (styles Styles) MimeRow(width int) lipgloss.Style {
	return styles.line(styles.theme.Muted(), width)
}

func (styles Styles) AppRow(width int) lipgloss.Style {
	return styles.line(styles.theme.Accent(), width)
}

func (styles Styles) EmptyRow(width int) lipgloss.Style {
	return styles.line(styles.theme.Muted(), width).Italic(true)
}

func (styles Styles) SelectedTypeRow(width int) lipgloss.Style {
	return styles.line(styles.theme.Contrast(), width).
		Background(styles.theme.Accent()).
		Bold(true)
}

func (styles Styles) HelpFrame() lipgloss.Style {
	return styles.accentFrame()
}

func (styles Styles) Shortcut(paddingRight int) lipgloss.Style {
	return styles.boldAccent(paddingRight)
}

func (styles Styles) Description() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.theme.Muted())
}

func (styles Styles) InfoFrame(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.theme.Accent()).
		Padding(0, 1)
}

func (styles Styles) InfoName() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.theme.Accent()).Bold(true)
}

func (styles Styles) InfoLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.theme.Muted())
}

func (styles Styles) InfoValue() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.theme.Accent())
}

func (styles Styles) InfoHint() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.theme.Muted()).Italic(true)
}

func (styles Styles) PaneDivider() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(styles.theme.Accent())
}

func (styles Styles) line(foreground color.Color, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(foreground).
		Width(width).
		MaxHeight(1)
}

func (styles Styles) boldAccent(paddingRight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(styles.theme.Accent()).
		Bold(true).
		PaddingRight(paddingRight)
}

func (styles Styles) accentFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.theme.Accent()).
		Padding(1, 2)
}
