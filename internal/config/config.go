// Package config loads ~/.config/gine/theme.toml.
package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	AccentLight string `toml:"accent_light"`
	AccentDark  string `toml:"accent_dark"`
	MutedLight  string `toml:"muted_light"`
	MutedDark   string `toml:"muted_dark"`
}

var defaults = Config{
	AccentLight: "#6D28D9",
	AccentDark:  "#A78BFA",
	MutedLight:  "#6B7280",
	MutedDark:   "#9CA3AF",
}

// Load reads ~/.config/gine/theme.toml. A missing file uses the defaults.
func Load() (Config, error) {
	loaded := defaults
	path, err := themePath()
	if err != nil {
		return loaded, nil
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return loaded, nil
	} else if err != nil {
		return Config{}, err
	}
	if _, err := toml.DecodeFile(path, &loaded); err != nil {
		return Config{}, err
	}
	return loaded, nil
}

func themePath() (string, error) {
	directory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(directory, "gine", "theme.toml"), nil
}
