// Package app starts gine.
package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/sumit-poudel/gine/internal/config"
	"github.com/sumit-poudel/gine/internal/gui"
	"github.com/sumit-poudel/gine/internal/mimetype"
)

func Start() error {
	colors, err := config.Load()
	if err != nil {
		return fmt.Errorf("theme: %w", err)
	}
	mimeTypes, err := mimetype.Types()
	if err != nil {
		return fmt.Errorf("mime types: %w", err)
	}
	_, err = tea.NewProgram(gui.New(colors, mimeTypes, mimetype.LoadCatalog(), mimetype.LoadAssociations())).Run()
	return err
}
