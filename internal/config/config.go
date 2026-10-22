package config

import (
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

func GetConfig() Config {
	conf := defaults
	if path, err := getPath(); err == nil {
		_, _ = toml.DecodeFile(path, &conf)
	}
	return conf
}

func getPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gine", "theme.toml"), nil
}
