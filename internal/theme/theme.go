package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/sumit-poudel/gine/internal/config"
)

type Theme struct {
	accent color.Color
	muted  color.Color
}

func New(dark bool) Theme {
	cfg := config.GetConfig()
	lightDark := lipgloss.LightDark(dark)
	return Theme{
		accent: lightDark(lipgloss.Color(cfg.AccentLight), lipgloss.Color(cfg.AccentDark)),
		muted:  lightDark(lipgloss.Color(cfg.MutedLight), lipgloss.Color(cfg.MutedDark)),
	}
}

func (t Theme) Screen(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center)
}

func (t Theme) Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.accent).
		Bold(true)
}

func (t Theme) Footer(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(t.muted)
}

func (t Theme) Input() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.accent).
		Bold(true)
}

func (t Theme) Dialog() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.accent).
		Padding(1, 2)
}

func (t Theme) DialogKey(gap int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.accent).
		Bold(true).
		PaddingRight(gap)
}

func (t Theme) DialogCell() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (t Theme) Flex(width, height int, panels ...string) string {
	if len(panels) == 0 {
		return ""
	}

	base, extra := width/len(panels), width%len(panels)
	blocks := make([]string, len(panels))
	for i, panel := range panels {
		w := base
		if i < extra {
			w++
		}
		blocks[i] = lipgloss.NewStyle().Width(w).Height(height).Render(panel)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, blocks...)
}
