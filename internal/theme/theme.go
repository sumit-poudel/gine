// Package theme picks accent and muted colors for the terminal background.
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

func New(dark bool, colors config.Config) Theme {
	pick := lipgloss.LightDark(dark)
	return Theme{
		accent: pick(lipgloss.Color(colors.AccentLight), lipgloss.Color(colors.AccentDark)),
		muted:  pick(lipgloss.Color(colors.MutedLight), lipgloss.Color(colors.MutedDark)),
	}
}

func (theme Theme) Accent() color.Color { return theme.accent }

func (theme Theme) Muted() color.Color { return theme.muted }

const midBrightness = 0x7FFF

func (theme Theme) Contrast() color.Color {
	red, green, blue, _ := theme.accent.RGBA()
	brightness := (299*int(red) + 587*int(green) + 114*int(blue)) / 1000
	if brightness > midBrightness {
		return lipgloss.Color("#1C1917")
	}
	return lipgloss.Color("#FAFAF9")
}
