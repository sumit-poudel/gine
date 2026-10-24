// Package config loads ~/.config/gine/theme.toml.
package config

import (
	"errors"
	"io/fs"
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

func Load() (Config, error) {
	loaded := defaults
	path, err := themePath()
	if err == nil {
		if _, err := toml.DecodeFile(path, &loaded); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return Config{}, err
		}
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
